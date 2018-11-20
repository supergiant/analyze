import { NgModule }             from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { PluginsComponent }     from "src/app/plugins/plugins.component";

const routes: Routes = [
  {
    path: '',
    component: PluginsComponent,
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class PluginsRoutingModule { }
