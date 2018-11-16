import { NgModule }     from '@angular/core';
import { CommonModule } from '@angular/common';

import { ChecksRoutingModule }                                                             from './checks-routing.module';
import { ChecksComponent }                                                                 from './checks.component';
import { HttpClientModule }                                                                from "@angular/common/http";
import { MatCardModule, MatTabsModule, MatExpansionModule, MatIconModule, MatChipsModule } from "@angular/material";

@NgModule({
  declarations: [ChecksComponent],
  imports: [
    CommonModule,
    ChecksRoutingModule,
    HttpClientModule,
    MatCardModule,
    MatTabsModule,
    MatExpansionModule,
    MatIconModule,
    MatChipsModule,
  ],
})
export class ChecksModule {
}
