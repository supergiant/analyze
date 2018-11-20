import { NgModule }             from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'checks'
  },
  {
    path: 'checks',
    loadChildren: './checks/checks.module#ChecksModule'
  },
  {
    path: '**',
    redirectTo: 'checks'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})

export class AppRoutingModule {
}
