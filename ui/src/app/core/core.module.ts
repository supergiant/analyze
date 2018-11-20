import { NgModule, Optional, SkipSelf }      from '@angular/core';
import { CommonModule }                      from '@angular/common';
import { HeaderComponent }                   from './header/header.component';
import { UserMenuComponent }                 from "./header/user-menu/user-menu.component";
import { MatDialogModule, MatToolbarModule } from "@angular/material";
import { FooterComponent }                   from './footer/footer.component';
import { RouterModule }                      from "@angular/router";
import { MenuModalComponent }                from "src/app/core/header/user-menu/menu-modal/menu-modal.component";

@NgModule({
  declarations: [
    //  nav component etc
    UserMenuComponent,
    HeaderComponent,
    FooterComponent,
    MenuModalComponent,
  ],
  imports: [
    CommonModule,
    MatDialogModule,
    MatToolbarModule,
    RouterModule,
  ],
  exports: [
    HeaderComponent,
    FooterComponent,
  ],
  entryComponents: [
    MenuModalComponent,
  ],
})
export class CoreModule {
  constructor(@Optional() @SkipSelf() parentModule: CoreModule) {
    throwIfAlreadyLoaded(parentModule, 'CoreModule');
  }
}

export function throwIfAlreadyLoaded(parentModule: any, moduleName: string) {
  if (parentModule) {
    throw new Error(`${moduleName} has already been loaded. Import Core modules in the AppModule only.`);
  }
}
