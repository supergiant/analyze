import { Injectable } from '@angular/core';
import { fromEvent, fromEventPattern, Observable } from "rxjs";
import { publish } from "rxjs/operators";
import { CELoadedEvent, EventType } from "../../models/events";
import { CustomElementsService } from "src/app/shared/services/custom-elements.service";
import { CeCacheService } from "src/app/shared/services/ce-cache.service";

@Injectable()
export class CeRegisterService {
  // <componentEntryPoint, componentRef>
  readonly registeredCEs: Map<string, string>;
  readonly bus: Element;
  readonly ceLoadedEvents$: Observable<CustomEvent>;

  constructor(private customElService: CustomElementsService, private ceCache: CeCacheService) {
    this.bus = document.querySelector<Element>('head');

    this.ceLoadedEvents$ = fromEventPattern(this.addHandler.bind(this), this.removeHandler.bind(this));
  }

  private addHandler(handler) {
    this.bus.addEventListener(EventType.CE_LOADED_EVENT, handler);
  }

  private removeHandler(handler) {
    this.bus.removeEventListener(EventType.CE_LOADED_EVENT, handler);
  }

  public registerCe(componentEntryPoint: string, container: string) {
    this.ceLoadedEvents$.subscribe((event: CustomEvent) => {
      const selector = event.detail.selector;
      const customEl = this.customElService.createCustomElement(selector)
      this.ceCache.addRegisteredCE(componentEntryPoint, selector);
      this.customElService.mountCustomElement(container, customEl);
    });

    const script = document.createElement('script');
    script.src = 'http://54.183.122.86:32291' + componentEntryPoint;

    this.bus.appendChild(script);
  }

  public getAllCeLoadedEvents() {
    return this.ceLoadedEvents$;
  }
}
