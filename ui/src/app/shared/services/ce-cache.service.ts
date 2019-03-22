import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class CeCacheService {

  registeredCEs: Map<string, string>;

  constructor() { this.registeredCEs = new Map<string, string>(); }

  addRegisteredCE(entrypoint, selector) {
    this.registeredCEs.set(entrypoint, selector);
  }

  deleteRegisteredCE(entrypoint) {
    this.registeredCEs.delete(entrypoint);
  }

  getAllRegisteredCEs() {
    return this.registeredCEs;
  }
}
