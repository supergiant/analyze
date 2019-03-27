import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class CustomElementsService {

  constructor() { }

  public mountCustomElement(containerSelector, CESelector) {
    const customEl: HTMLElement = document.createElement(CESelector);
    const container = document.querySelector(containerSelector);
    container.appendChild(customEl);
  }
}
