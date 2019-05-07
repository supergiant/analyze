import { Component, AfterViewInit, OnDestroy, ViewEncapsulation, ElementRef } from '@angular/core';
import { HttpClient }                           from "@angular/common/http";
import { map, takeUntil, mergeMap }                       from "rxjs/operators";
import { Observable, Subject, concat, forkJoin }                           from "rxjs";

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
export class ChecksComponent implements AfterViewInit, OnDestroy {

  private ceLoadedEvents$: Observable<CustomEvent>;
  private registeredCEs: Map<string, {[key: string]: any}>;
  private container: string;
  private getChecks$: Observable<Object>;
  private checks: {};
  private ngUnsubscribe = new Subject();
  private lateChecks$: any;

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
      this.getChecks$ = this.pluginsService.getChecks().pipe(
        map((checks: Check[]) => {
          return checks.reduce((obj, ck) => Object.assign({[ck.id]: ck}, obj), {})
        }),
        takeUntil(this.ngUnsubscribe)
      );
    }

  ngAfterViewInit() {
    this.container = this.elRef.nativeElement.tagName.toLowerCase();

    this.getChecks$.subscribe(
      checks => {
        // TODO: hack until we can figure out how to combine getChecks and ceLoadedEvents
        this.checks = checks;
        this.pluginsService.getAll().map((plugin: Plugin) => {
          const pluginId = plugin.id;
          const entrypoint = plugin.checkComponentEntryPoint;
          const checkData = checks[plugin.id];

          if (this.registeredCEs.get(pluginId)) {
            if (this.registeredCEs.get(pluginId).hasOwnProperty("check")) {
              this.customElService.mountCustomElement(this.container, this.registeredCEs.get(pluginId).check, "check-result", checkData);
            } else { this.ceRegisterService.registerCe(entrypoint); }
          } else {
            this.ceRegisterService.registerCe(entrypoint);
          }
        });
      },
      err => console.log(err)
    )

    this.ceLoadedEvents$
      .pipe(
        takeUntil(this.ngUnsubscribe)
      )
      .subscribe((event: CustomEvent) => {
        const ceSelector = event.detail.selector;
        // TODO: we need to swap this out for a pluginId property (here and in plugin)
        const pluginId = event.detail.pluginName;
        const checkData = this.checks[event.detail.pluginName]
        this.ceCache.addRegisteredCE(pluginId, ceSelector, "check");
        this.customElService.mountCustomElement(this.container, ceSelector, "check-result", checkData);
      }
    )
  }

  ngOnDestroy() {
    this.ngUnsubscribe.next();
    this.ngUnsubscribe.complete();
  }
}
