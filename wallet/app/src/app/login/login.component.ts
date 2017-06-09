import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AppSettingsService } from '../app.settings.service';

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css']
})

export class LoginComponent implements OnInit {
    private model: any = {};
    private loading = false;
    private error = '';

    constructor(
        private router: Router,
        private appSettingsService: AppSettingsService,
    ) {

    }

    ngOnInit() {

    }

    login() {
        this.loading = true;
    }
}
