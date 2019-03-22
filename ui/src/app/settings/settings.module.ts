import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { SettingsRoutingModule } from './settings-routing.module';
import { SettingsComponent }     from './settings.component';
import { HttpClientModule }     from "@angular/common/http";
import { MatCardModule }        from "@angular/material";

@NgModule({
  declarations: [SettingsComponent],
  imports: [
    CommonModule,
    SettingsRoutingModule,
    HttpClientModule,
    MatCardModule,
  ]
})
export class SettingsModule { }
