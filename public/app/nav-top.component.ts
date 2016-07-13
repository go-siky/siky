import {Component} from '@angular/core';
import {RouteConfig, ROUTER_DIRECTIVES, ROUTER_PROVIDERS} from '@angular/router-deprecated';

@Component({
    selector: 'nav-top',
    templateUrl: 'app/nav-top.component.html'
})
export class NavTopComponent {
    name = "Docker Siky"
}
