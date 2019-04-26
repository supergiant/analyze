import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class CustomElementsService {

  constructor() { }

  public mountCustomElement(containerSelector, CESelector, attr, data) {
    const customEl: HTMLElement = document.createElement(CESelector);
    // TODO: 'attr' is a hack until loading/registering is synchronous
    customEl.setAttribute(attr, JSON.stringify(data))
    const container = document.querySelector(containerSelector);
    container.appendChild(customEl);
  }
}
