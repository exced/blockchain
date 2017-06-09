import { Injectable } from '@angular/core';

@Injectable()
export class AppSettingsService {

    private URL: string = 'http://localhost:3000';

    constructor() {
    }

    public getURL(): string {
        return this.URL;
    }
}