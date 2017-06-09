import { Injectable } from '@angular/core';
import { AppSettingsService } from '../app/app.settings.service';
import { Http, Headers, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map'

const HAS_LOGGED_IN = 'BlockchainWalletHasLoggedIn';

@Injectable()
export class AuthService {

    private URL: string;

    constructor(private http: Http, private appSettingsService: AppSettingsService) {
        this.URL = appSettingsService.getURL();
    }

    public getUser() {
        return localStorage.getItem(HAS_LOGGED_IN);
    }

    /**
     * Login with user public key and password
     * @param key
     * @param password
     */
    public login(key: string, password: string): Observable<any> {
        return this.http.post(this.URL + '/login', JSON.stringify({ key: key, password: password }))
            .map((response: Response) => {
                // login successful if there's a jwt token in the response
                let user = response.json();
                if (user && user.token) {
                    // store user details and jwt token in local storage to keep user logged in between page refreshes
                    localStorage.setItem(HAS_LOGGED_IN, JSON.stringify(user));
                }
            });
    }

    /**
     * Signin with user public key and password
     * @param key
     * @param password
     */
    public signin(key: string, password: string): Observable<any> {
        return this.http.post('/signin', JSON.stringify({ key: key, password: password }))
            .map((response: Response) => {
                // login successful if there's a jwt token in the response
                let user = response.json();
                if (user && user.token) {
                    // store user details and jwt token in local storage to keep user logged in between page refreshes
                    localStorage.setItem(HAS_LOGGED_IN, JSON.stringify(user));
                }
            });
    }

    public logout() {
        // remove user from local storage to log user out
        localStorage.removeItem(HAS_LOGGED_IN);
    }
}