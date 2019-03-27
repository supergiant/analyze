import { Component, AfterViewInit, ViewEncapsulation, ElementRef } from '@angular/core';
import { HttpClient }                           from "@angular/common/http";
import { map }                                  from "rxjs/operators";
import { Observable }                           from "rxjs";

import { Plugin } from 'src/app/models/plugin';
import { PluginsService } from 'src/app/shared/services/plugins.service';
import { CeRegisterService } from "../shared/services/ce-register.service";
import { CustomElementsService } from "../shared/services/custom-elements.service";
import { CeCacheService } from "../shared/services/ce-cache.service";

@Component({
  selector: 'app-checks',
  templateUrl: './checks.component.html',
  styleUrls: ['./checks.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class ChecksComponent implements AfterViewInit {

  private ceLoadedEvents$: Observable<CustomEvent>;
  private registeredCEs: Map<string, string>;

  constructor(
    private http: HttpClient,
    private pluginsService: PluginsService,
    private ceRegisterService: CeRegisterService,
    private customElService: CustomElementsService,
    private elRef: ElementRef,
    private ceCache: CeCacheService
  ) { this.registeredCEs = ceCache.getAllRegisteredCEs(); }

  ngAfterViewInit() {
    const container = this.elRef.nativeElement.tagName.toLowerCase()

    this.pluginsService.getAll().map((plugin: Plugin) => {
      const entrypoint = plugin.checkComponentEntryPoint;

      if (!this.registeredCEs.has(entrypoint)) {
        this.ceRegisterService.registerAndMountCe(entrypoint, container)
      } else {
        this.customElService.mountCustomElement(container, this.registeredCEs.get(entrypoint));
      }

    });
  }
}
