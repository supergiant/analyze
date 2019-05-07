import { Component, AfterViewInit, OnDestroy, ElementRef } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, Subject, of } from "rxjs";
import { takeUntil, switchMap, tap } from "rxjs/operators";

import { Plugin } from 'src/app/models/plugin';
import { CELoadedEvent, EventType } from "../models/events";
import { PluginsService } from 'src/app/shared/services/plugins.service';
import { CeRegisterService } from "../shared/services/ce-register.service";
import { CustomElementsService } from "../shared/services/custom-elements.service";
import { CeCacheService } from "../shared/services/ce-cache.service";

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss'],
})
export class SettingsComponent implements AfterViewInit, OnDestroy {

  private ceLoadedEvents$: Observable<CustomEvent>;
  private registeredCEs: Map<string, {[key: string]: any}>;
  private container: string;
  private ngUnsubscribe = new Subject()

  constructor(
    private pluginsService: PluginsService,
    private ceRegisterService: CeRegisterService,
    private customElService: CustomElementsService,
    private elRef: ElementRef,
    private ceCache: CeCacheService,
    private http: HttpClient
  ){
    this.registeredCEs = ceCache.getAllRegisteredCEs();
    this.ceLoadedEvents$ = this.ceRegisterService.getAllCeLoadedEvents();
   }

  getKeyFromValue(map, value) {
    for (let [k, v] of map) {
      if (v.settings === value)
        return k;
    }
  }

  ngAfterViewInit() {
    this.container = this.elRef.nativeElement.tagName.toLowerCase()
    this.elRef.nativeElement.addEventListener('ConfigUpdate',
      e => {
        // temporary
        const pluginId = this.getKeyFromValue(this.registeredCEs, e.target.tagName.toLowerCase());
        this.pluginsService.updateConfig(pluginId.replace("-settings", ""), e.detail).pipe(
          takeUntil(this.ngUnsubscribe),
          switchMap(res => this.pluginsService.getPluginConfig(pluginId.replace("-settings", "")))
        ).subscribe(
          config => {
            let el = document.querySelector(this.registeredCEs.get(pluginId).settings)
            el.setAttribute('plugin-config', JSON.stringify(config))
          }
        )
      });

    this.pluginsService.getAll().map((plugin: Plugin) => {
      const pluginId = plugin.id
      const entrypoint = plugin.settingsComponentEntryPoint;

      if (this.registeredCEs.get(pluginId)) {
        if (this.registeredCEs.get(pluginId).hasOwnProperty("settings")) {
          this.customElService.mountCustomElement(this.container, this.registeredCEs.get(pluginId).settings);
          this.pluginsService.getPluginConfig(plugin.id).pipe(
            takeUntil(this.ngUnsubscribe)
          )
          .subscribe(
            config => {
              let el = document.querySelector(this.registeredCEs.get(pluginId).settings);
              el.setAttribute('plugin-config', JSON.stringify(config));
            }
          )
        } else { this.ceRegisterService.registerCe(entrypoint) }
      } else {
        this.ceRegisterService.registerCe(entrypoint)
      }

    });

    this.ceLoadedEvents$
      .pipe(
        takeUntil(this.ngUnsubscribe),
        tap((event: CustomEvent) => {
          // TODO: we need to swap this out for a pluginId property (here and in plugin)
          const pluginName = event.detail.pluginName;
          this.ceCache.addRegisteredCE(pluginName, event.detail.selector, "settings");
          this.customElService.mountCustomElement(this.container, event.detail.selector);
        }),
        switchMap((event: CustomEvent) => {
          // TODO: trying to avoid having to subscribe in here
          this.pluginsService.getPluginConfig(event.detail.pluginName).toPromise().then(
            config => {
              let el = document.querySelector(event.detail.selector);
              el.setAttribute('plugin-config', JSON.stringify(config));
            }
          )
          return of(null)
        })
      ).subscribe()
  }

  ngOnDestroy() {
    this.ngUnsubscribe.next();
    this.ngUnsubscribe.complete();
  }
}
