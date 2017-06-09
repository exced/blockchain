import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { Routes, RouterModule } from '@angular/router';

import { LoginComponent } from './login/login.component';
import { SigninComponent } from './signin/signin.component';
import { HomeComponent } from './home/home.component';
import { AuthGuard } from '../guards/auth.guard';

const appRoutes: Routes = [
  { path: 'login', component: LoginComponent },
  { path: 'signin', component: SigninComponent },
  { path: '', component: HomeComponent, canActivate: [AuthGuard] },
  { path: '**', redirectTo: '' } // otherwise redirect to home
];

@NgModule({
  imports: [
    BrowserModule,
    RouterModule.forRoot(appRoutes, { useHash: true }),
  ],
  providers: [],
  declarations: [],
  exports: [
    RouterModule
  ],
})

export class AppRoutingModule { }