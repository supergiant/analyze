package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/dre1080/recover"
	"github.com/go-openapi/loads"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/rest"

	"github.com/supergiant/analyze/asset"
	"github.com/supergiant/analyze/pkg/analyze"
	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
	"github.com/supergiant/analyze/pkg/api/operations"
	"github.com/supergiant/analyze/pkg/config"
	"github.com/supergiant/analyze/pkg/kube"
	"github.com/supergiant/analyze/pkg/logger"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/proxy"
	"github.com/supergiant/analyze/pkg/scheduler"
	"github.com/supergiant/analyze/pkg/storage/etcd"
)

var configFilePaths = flag.StringArrayP(
	"config",
	"c",
	[]string{"./config.yaml", "/etc/analyzed/config.yaml", "$HOME/.analyzed/config.yaml"},
	"config file path")

func main() {
	flag.Parse()

	cfg := &analyze.Config{}

	// configFileReadError is not critical due to possibility that configuration is done by environment variables
	configFileReadError := config.ReadFromFiles(cfg, *configFilePaths)

	if err := config.MergeEnv("AZ", cfg); err != nil {
		log.Fatalf("unable to merge env variables %v\n", err)
	}

	//TODO: try to unify APIs discovery which are hosted in k8s
	//TODO: and rewrite config population logic
	if etcdEndpoint := discoverETCDEndpoint(); etcdEndpoint != "" {
		cfg.ETCD.Endpoints = append(cfg.ETCD.Endpoints, discoverETCDEndpoint())
	}

	log := logger.NewLogger(cfg.Logging).WithField("app", "analyze-core")
	mainLogger := log.WithField("component", "main")

	mainLogger.Infof("config: %+v", cfg)
	mainLogger.Infof("config file name: %s", config.UsedFileName())
	if configFileReadError != nil {
		log.Warnf("unable to read config file, %v", configFileReadError)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validation error, err: %v", err)
	}

	kubeClient, err := kube.NewKubeClient(log.WithField("component", "kubeClient"))
	if err != nil {
		log.Fatalf("unable to create kube client, err: %v", err)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("can't get kube config, err: %v", err)
	}

	tr, err := rest.TransportFor(config)
	if err != nil {
		log.Fatalf("can't get kube client transport, err: %v", err)
	}
	proxySet := proxy.NewProxySet(tr, log.WithField("component", "proxySet"))

	etcdStorage, err := etcd.NewETCDStorage(cfg.ETCD, log.WithField("component", "etcdClient"))
	if err != nil {
		log.Fatalf("unable to create ETCD client, err: %v", err)
	}

	defer etcdStorage.Close()

	scheduler := scheduler.NewScheduler(log.WithField("component", "scheduler"))
	defer scheduler.Stop()

	watchChan := etcdStorage.WatchRange(context.Background(), models.PluginPrefix)
	log.Debug("etcd watch is started")
	pluginController := analyze.NewPluginController(
		watchChan,
		etcdStorage,
		kubeClient,
		scheduler,
		proxySet,
		log.WithField("component", "pluginController"),
	)
	defer pluginController.Stop()

	swaggerSpec, err := loads.Analyzed(api.SwaggerJSON, "2.0")
	if err != nil {
		log.Fatalf("unable to create spec analyzed document, err: %v", err)
	}

	//TODO: add request logging middleware
	//TODO: add metrics middleware
	analyzeAPI := operations.NewAnalyzeAPI(swaggerSpec)
	analyzeAPI.Logger = log.WithField("component", "analyzeApi").Errorf

	analyzeAPI.GetCheckResultsHandler = handlers.NewChecksResultsHandler(
		etcdStorage,
		log.WithField("handler", "CheckResultsHandler"),
	)
	analyzeAPI.GetPluginHandler = handlers.NewPluginHandler(
		etcdStorage,
		log.WithField("handler", "PluginHandler"),
	)
	analyzeAPI.GetPluginsHandler = handlers.NewPluginsHandler(
		etcdStorage,
		log.WithField("handler", "PluginsHandler"),
	)
	analyzeAPI.RegisterPluginHandler = handlers.NewRegisterPluginHandler(
		etcdStorage,
		log.WithField("handler", "RegisterPluginHandler"),
	)
	analyzeAPI.UnregisterPluginHandler = handlers.NewUnregisterPluginHandler(
		etcdStorage,
		log.WithField("handler", "UnregisterPluginHandler"),
	)
	analyzeAPI.GetPluginConfigHandler = handlers.NewPluginConfigHandler(
		etcdStorage,
		log.WithField("handler", "PluginConfigHandler"),
	)

	analyzeAPI.ReplacePluginConfigHandler = handlers.NewReplacePluginConfigHandler(
		etcdStorage,
		log.WithField("handler", "ReplacePluginConfigHandler"),
	)

	err = analyzeAPI.Validate()
	if err != nil {
		log.Fatalf("API configuration error, err: %v", err)
	}

	server := api.NewServer(analyzeAPI)
	server.Port = cfg.API.ServerPort
	server.Host = cfg.API.ServerHost
	server.ConfigureAPI()

	handlerWithRecovery := recover.New(&recover.Options{
		Log: logrus.Error,
	})

	//TODO fix CORS till release
	corsHandler := cors.New(cors.Options{
		Debug:          false,
		AllowedHeaders: []string{"*"},
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{},
		MaxAge:         1000,
	}).Handler

	handler := alice.New(
		handlerWithRecovery,
		corsHandler,
		swaggerMiddleware,
		newProxyMiddleware(proxySet, log.WithField("middleware", "proxy")),
		uiMiddleware,
	).Then(analyzeAPI.Serve(nil))

	server.SetHandler(handler)

	defer server.Shutdown()

	if servingError := server.Serve(); servingError != nil {
		log.Fatalf("unable to serve HTTP API, err: %v", servingError)
	}
}

func discoverETCDEndpoint() string {
	etcdHost, hostExists := os.LookupEnv("ETCD_SERVICE_HOST")
	etcdPort, portExists := os.LookupEnv("ETCD_SERVICE_PORT")
	if !hostExists || !portExists {
		return ""
	}
	return etcdHost + ":" + etcdPort
}

func swaggerMiddleware(handler http.Handler) http.Handler {
	var staticServer = http.FileServer(asset.Assets)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/api/v1/swagger-ui" || r.URL.Path == "/api/v1/help" {
			http.Redirect(w, r, "/api/v1/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.HasPrefix(r.URL.Path, "/api/v1/swagger-ui/") {
			url := strings.TrimPrefix(r.URL.Path, "/api/v1/swagger-ui/")
			r.URL.Path = "/swagger/" + url
			staticServer.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})

}

func uiMiddleware(handler http.Handler) http.Handler {
	var staticServer = http.FileServer(asset.Assets)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.HasPrefix(r.URL.Path, "/api/v1") {
			r.URL.Path = "/ui" + r.URL.Path
			staticServer.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func newProxyMiddleware(proxySet *proxy.Set, logger logrus.FieldLogger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if strings.HasPrefix(req.URL.Path, "/api/v1") {
				handler.ServeHTTP(res, req)
				return
			}

			var targetProxy *httputil.ReverseProxy
			for id, proxy := range proxySet.GetProxies() {
				if strings.Contains(req.URL.Path, id) {
					targetProxy = proxy
				}
			}
			if targetProxy == nil {
				handler.ServeHTTP(res, req)
				return
			}

			logger.Debugf("got proxy request at: %v, request: %+v", time.Now(), req.URL)
			defer logger.Debugf("proxy request processing finished at: %v, request: %+v", time.Now(), req.URL)

			// Update the headers to allow for SSL redirection
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

			// Note that ServeHttp is non blocking and uses a go routine under the hood
			targetProxy.ServeHTTP(res, req)
		})
	}
}
