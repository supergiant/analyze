import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Plugin } from 'src/app/models/plugin';

@Injectable()
export class PluginsService {
  private installedPlugins: Plugin[];

  constructor(
    private http: HttpClient,
  ) {
  }

  //TODO rewrite
  public async refreshAll() {
    await this.http.get<Plugin[]>('/api/v1/plugins').toPromise().then(
      (plugins: Plugin[]) => {
        console.log(plugins);
        this.installedPlugins = plugins;
      }, (reason: any) => {
        console.log('cant get plugins: ', reason)
      }
    );
  }

  public getAll(): Plugin[] {
    return this.installedPlugins
  }

  public getChecks() {
    return this.http.get('api/v1/checks')
  }

  public getPluginConfig(pluginId: string) {
    return this.http.get('api/v1/plugins/' + pluginId + '/config');
  }

  public updateConfig(pluginId, config) {
    return this.http.patch('api/v1/plugins/' + pluginId + '/config', config)
  }
}
