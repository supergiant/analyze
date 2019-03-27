import { CUSTOM_ELEMENTS_SCHEMA, NgModule, APP_INITIALIZER } from '@angular/core';
import { CommonModule } from '@angular/common';

import { AppRoutingModule }        from './app-routing.module';
import { AppComponent }            from './app.component';
import { CoreModule }              from './core/core.module';
import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { PluginsService } from "./shared/services/plugins.service";
import { CeRegisterService } from "./shared/services/ce-register.service";
import { BrowserModule } from "@angular/platform-browser";


export function startupServiceFactory(pluginsService: PluginsService): Function {
  return () => {
    return pluginsService.refreshAll();
  };
}

@NgModule({
  providers: [
    PluginsService,
    CeRegisterService,
    {
      provide: APP_INITIALIZER,
      useFactory: startupServiceFactory,
      multi: true,
      deps: [PluginsService]
    },
  ],
  declarations: [
    AppComponent,
  ],
  imports: [
    CommonModule,
    AppRoutingModule,
    CoreModule,
    BrowserAnimationsModule,
    HttpClientModule,
    BrowserModule,
  ],
  bootstrap: [AppComponent],
  schemas: [
    CUSTOM_ELEMENTS_SCHEMA
  ]
})
export class AppModule {
}
