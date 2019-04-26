import { Component, AfterViewInit, OnDestroy, ViewEncapsulation, ElementRef } from '@angular/core';
import { HttpClient }                           from "@angular/common/http";
import { map, takeUntil }                       from "rxjs/operators";
import { Observable, Subject }                  from "rxjs";

import { Plugin } from 'src/app/models/plugin';
import { Check } from 'src/app/models/Check';
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
  private registeredCEs: Map<string, string>;
  private readonly onDestroy = new Subject<void>();
  private DOMWatcher: MutationObserver;
  private container: string;

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
    }

  ngAfterViewInit() {
    this.container = this.elRef.nativeElement.tagName.toLowerCase();

    // TODO: how can we do this without DOM?
    this.DOMWatcher = new MutationObserver(this.populateChecks.bind(this));
    this.DOMWatcher.observe(document.querySelector(this.container), { attributes: false, childList: true, subtree: false })

    this.pluginsService.getAll().map((plugin: Plugin) => {
      const entrypoint = plugin.checkComponentEntryPoint;

      if (!this.registeredCEs.has(plugin.id)) {
        this.ceRegisterService.registerAndMountCe(entrypoint, plugin.id, this.container);
      } else {
        this.customElService.mountCustomElement(this.container, this.registeredCEs.get(plugin.id));
      }
    });
  }

  ngOnDestroy() {
    this.onDestroy.next();
    this.DOMWatcher.disconnect();
  }

  populateChecks(mutationsList, observer) {
    // TODO: can we make this pure?
    if (document.querySelector(this.container).childElementCount === this.pluginsService.getAll().length) {
      this.pluginsService.getChecks().subscribe(
        (checks: Check[]) => {
          checks.forEach(c => {
            const selector = this.registeredCEs.get(c.id);
            const el = document.querySelector(selector);
            el.setAttribute("check-result", JSON.stringify(c));
          })
        },
        err => console.log(err)
      )
    }
  }
}
