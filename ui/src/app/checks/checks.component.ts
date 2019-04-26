import { Component, AfterViewInit, ViewEncapsulation, ElementRef } from '@angular/core';
import { HttpClient }                           from "@angular/common/http";
import { map, takeUntil }                       from "rxjs/operators";
import { Observable }                           from "rxjs";

import { Plugin } from 'src/app/models/plugin';
import { Check } from 'src/app/models/check';
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
  private container: string;
  private checks$: Observable<Object>;

  constructor(
    private http: HttpClient,
    private pluginsService: PluginsService,
    private ceRegisterService: CeRegisterService,
    private customElService: CustomElementsService,
    private elRef: ElementRef,
    private ceCache: CeCacheService
  ) {
      this.registeredCEs = ceCache.getAllRegisteredCEs();
      this.ceLoadedEvents$ = this.ceRegisterService.getAllCeLoadedEvents();
      this.checks$ = this.pluginsService.getChecks().pipe(
        map((checks: Check[]) => {
          return checks.reduce((obj, ck) => Object.assign({[ck.id]: ck}, obj), {})
        })
      );
    }

  ngAfterViewInit() {
    this.container = this.elRef.nativeElement.tagName.toLowerCase();

    this.checks$.subscribe(
      checks => {
        this.pluginsService.getAll().map((plugin: Plugin) => {
          const entrypoint = plugin.checkComponentEntryPoint;
          const checkData = checks[plugin.id];

          if (!this.registeredCEs.has(plugin.id)) {
            this.ceRegisterService.registerAndMountCe(entrypoint, plugin.id, this.container, "check-result", checkData);
          } else {
            this.customElService.mountCustomElement(this.container, this.registeredCEs.get(plugin.id), "check-result", checkData);
          }
        });
      },
      err => console.log(err)
    )
  }
}
