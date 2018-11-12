import { NgModule }             from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ChecksComponent }      from "src/app/checks/checks.component";

const routes: Routes = [
  {
    path: '',
    component: ChecksComponent
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ChecksRoutingModule { }
