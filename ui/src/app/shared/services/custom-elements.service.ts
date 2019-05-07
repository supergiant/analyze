import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class CustomElementsService {

  constructor() { }

  public mountCustomElement(containerSelector, CESelector, attr?, data?) {
    const customEl: HTMLElement = document.createElement(CESelector);
    const container = document.querySelector(containerSelector);
    if (attr && data) {
      customEl.setAttribute(attr, JSON.stringify(data))
    }
    container.appendChild(customEl);
  }
}
