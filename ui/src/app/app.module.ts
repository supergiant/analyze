import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { AppRoutingModule }        from './app-routing.module';
import { AppComponent }            from 'src/app/app.component';
import { CoreModule }              from './core/core.module';
import { SharedModule }            from './shared/shared.module';
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { BrowserModule }           from "@angular/platform-browser";

@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    CommonModule,
    AppRoutingModule,
    CoreModule,
    SharedModule,
    BrowserAnimationsModule,

  ],
  bootstrap: [AppComponent],
})
export class AppModule {
}
