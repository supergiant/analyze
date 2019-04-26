import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class CeCacheService {

  registeredCEs: Map<string, string>;

  constructor() { this.registeredCEs = new Map<string, string>(); }

  addRegisteredCE(plugId, selector) {
    this.registeredCEs.set(plugId, selector);
  }

  deleteRegisteredCE(plugId) {
    this.registeredCEs.delete(plugId);
  }

  getAllRegisteredCEs() {
    return this.registeredCEs;
  }
}
