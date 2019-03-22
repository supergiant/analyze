import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class CustomElementsService {

  constructor() { }

  public createCustomElement(selector): HTMLElement {
    const customEl: HTMLElement = document.createElement(selector);
    return customEl
  }

  public mountCustomElement(containerSelector, customEl) {
    const container = document.querySelector(containerSelector);
    container.appendChild(customEl);
  }
}
