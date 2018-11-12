import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ChecksRoutingModule } from './checks-routing.module';
import { ChecksComponent } from './checks.component';

@NgModule({
  declarations: [ChecksComponent],
  imports: [
    CommonModule,
    ChecksRoutingModule
  ]
})
export class ChecksModule { }
