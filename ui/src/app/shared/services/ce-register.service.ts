import { Injectable } from '@angular/core';
import { fromEvent, fromEventPattern, Observable } from "rxjs";
import { publish, single } from "rxjs/operators";
import { CELoadedEvent, EventType } from "../../models/events";
import { CustomElementsService } from "src/app/shared/services/custom-elements.service";
import { CeCacheService } from "src/app/shared/services/ce-cache.service";
import { environment } from 'src/environments/environment';

@Injectable()
export class CeRegisterService {
  // <componentEntryPoint, componentRef>
  readonly registeredCEs: Map<string, string>;
  readonly bus: Element;
  readonly ceLoadedEvents$: Observable<CustomEvent>;

  constructor(private customElService: CustomElementsService, private ceCache: CeCacheService) {
    this.bus = document.querySelector<Element>('head');

    this.ceLoadedEvents$ = fromEventPattern(this.addHandler.bind(this), this.removeHandler.bind(this));

    // this.test();
  }

  async test() {
    await console.log(this.ceLoadedEvents$.toPromise);
  }

  private addHandler(handler) {
    this.bus.addEventListener(EventType.CE_LOADED_EVENT, handler);
  }

  private removeHandler(handler) {
    this.bus.removeEventListener(EventType.CE_LOADED_EVENT, handler);
  }

  public registerAndMountCe(componentEntryPoint: string, containerSelector: string) {
    // TODO: this fn currently has too many responsibilities (register w/ DOM, cache, create and mount ce)
    // need to extract these into a synchronous workflow
    this.ceLoadedEvents$.subscribe((event: CustomEvent) => {
      const ceSelector = event.detail.selector;
      this.ceCache.addRegisteredCE(componentEntryPoint, ceSelector)
      this.customElService.mountCustomElement(containerSelector, ceSelector)
    });

    const script = document.createElement('script');
    script.src = environment.hostUrl + componentEntryPoint;

    this.bus.appendChild(script);
  }

  public getAllCeLoadedEvents() {
    return this.ceLoadedEvents$;
  }
}
