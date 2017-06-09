import { Component, OnInit } from '@angular/core';

@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

    private values = [];

    constructor(
    ) {
        this.values = [
            {
                currency: "BTC",
                value: 10
            },
            {
                currency: "ETH",
                value: 100
            }, {
                currency: "EXC",
                value: 1000
            },
        ]
    }

    ngOnInit() {
    }

}
