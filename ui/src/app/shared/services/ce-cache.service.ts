import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class CeCacheService {

  registeredCEs: Map<string, {[key: string]: any}>;

  constructor() { this.registeredCEs = new Map(); }

  addRegisteredCE(pluginId, selector, key) {
    let val;

    if (this.registeredCEs.get(pluginId)) {
      val = this.registeredCEs.get(pluginId);
      val[key] = selector;
    } else {
      val = {[key]: selector}
    }

    this.registeredCEs.set(pluginId, val);
  }

  deleteRegisteredCE(plugId) {
    this.registeredCEs.delete(plugId);
  }

  getAllRegisteredCEs() {
    return this.registeredCEs;
  }
}
