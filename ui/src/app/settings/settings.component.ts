import { Component, OnInit, ElementRef } from '@angular/core';
import { Observable } from "rxjs";
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
export class SettingsComponent {

  private ceLoadedEvents$: Observable<CustomEvent>;
  private registeredCEs: Map<string, string>;

  constructor(
    private pluginsService: PluginsService,
    private ceRegisterService: CeRegisterService,
    private customElService: CustomElementsService,
    private elRef: ElementRef,
    private ceCache: CeCacheService
  ){ this.registeredCEs = ceCache.getAllRegisteredCEs(); }

  ngAfterViewInit() {
    const container = this.elRef.nativeElement.tagName.toLowerCase()

    this.pluginsService.getAll().map((plugin: Plugin) => {
      const entrypoint = plugin.settingsComponentEntryPoint;
      if (!this.registeredCEs.has(entrypoint)) {
        this.ceRegisterService.registerCe(entrypoint, container);
      } else {
        const customEl = this.customElService.createCustomElement(this.registeredCEs.get(entrypoint))
        this.customElService.mountCustomElement(container, customEl)
      }
    });
  }
}
