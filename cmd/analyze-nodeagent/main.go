package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:          "analyze-nodeagent",
		Short:        "analyze-nodeagent deployed on each node to collect node related info for analyze",
		RunE:         runCommand,
		SilenceUsage: true,
	}

	command.PersistentFlags().StringP(
		"api-port",
		"p",
		"9292",
		"tcp port where node agent API is serving")

	if err := command.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}

func runCommand(cmd *cobra.Command, _ []string) error {

	apiPort, err := cmd.Flags().GetString("api-port")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag api-port")
	}

	logger := logrus.New()

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		logger.Fatal("unable to load SDK config, " + err.Error())
	}

	var ec2MetadataService = ec2metadata.New(cfg)

	if strings.ToLower(os.Getenv("AWS_EC2_METADATA_DISABLED")) == "true" {
		logger.Error("metadata is available")
	}

	if available := ec2MetadataService.Available(); !available {
		logger.Error("metadata is NOT available")
	}

	result, err := ec2MetadataService.GetMetadata("instance-id")
	if err != nil {
		logger.Error(err)
	}

	logger.Errorf("metadata is available, instance-id: %s", result)

	var httpServer = &http.Server{}
	var router = mux.NewRouter()
	awsAPI := router.PathPrefix("/aws").Subrouter()
	httpServer.Handler = awsAPI

	awsAPI.HandleFunc("/meta-data/{path}", func(ec2MetadataService *ec2metadata.EC2Metadata, logger logrus.FieldLogger) func(http.ResponseWriter, *http.Request) {
		return func(res http.ResponseWriter, req *http.Request) {
			vars := mux.Vars(req)

			path := vars["path"]
			result, err := ec2MetadataService.GetMetadata(path)
			if err != nil {
				logger.Error(err)
			}
			res.Write([]byte(result))
		}
	}(ec2MetadataService, logger)).Methods(http.MethodGet)

	awsAPI.HandleFunc("/dynamic/{path}", func(ec2MetadataService *ec2metadata.EC2Metadata, logger logrus.FieldLogger) func(http.ResponseWriter, *http.Request) {
		return func(res http.ResponseWriter, req *http.Request) {
			vars := mux.Vars(req)

			path := vars["path"]
			result, err := ec2MetadataService.GetDynamicData(path)
			if err != nil {
				logger.Error(err)
			}
			res.Write([]byte(result))
		}
	}(ec2MetadataService, logger)).Methods(http.MethodGet)

	awsAPI.HandleFunc("/user-data", func(ec2MetadataService *ec2metadata.EC2Metadata, logger logrus.FieldLogger) func(http.ResponseWriter, *http.Request) {
		return func(res http.ResponseWriter, req *http.Request) {
			result, err := ec2MetadataService.GetUserData()
			if err != nil {
				logger.Error(err)
			}
			res.Write([]byte(result))
		}
	}(ec2MetadataService, logger)).Methods(http.MethodGet)

	listener, err := net.Listen("tcp", ":"+apiPort)
	if err != nil {
		return err
	}

	var addr = listener.Addr().String()
	var addrParts = strings.Split(addr, ":")
	if len(addrParts) == 0 {
		return errors.Errorf("can't get non occupied port, addr %v", addr)
	}

	if err := httpServer.Serve(listener); err != nil {
		if err != nil {
			return err
		}
	}

	return nil
}
