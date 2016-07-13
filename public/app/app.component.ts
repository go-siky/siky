import {Component} from '@angular/core';
// import {ClickMeComponent}    from './click-me.component';
import {RouteConfig, ROUTER_DIRECTIVES, ROUTER_PROVIDERS} from '@angular/router-deprecated';

import {RepoService} from "./docker/image/img.service";
import {ReposComponent, TagsComponent, TagComponent} from "./docker/image/img.component";
import {WelcomeComponent} from "./welcome/welcome.component";

@Component({
    selector: 'my-app',
    templateUrl: 'app/app.component.html',
    directives: [
        ReposComponent, ROUTER_DIRECTIVES
    ], providers: [
        ROUTER_PROVIDERS,
        RepoService,
    ]
})
@RouteConfig([
    {path: '/', name: 'HOME', component: WelcomeComponent, useAsDefault: true},
    {path: '/repos', name: 'Repos', component: ReposComponent},
    {path: '/repos/:group/:name/tags', name: 'Tags', component: TagsComponent},
    {path: '/repos/:group/:name/tag/:version/', name: 'Tag', component: TagComponent},
])
export class AppComponent {
    title = "Siky X";
}
