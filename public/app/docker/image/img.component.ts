///<reference path="../../../node_modules/@angular/core/src/metadata/lifecycle_hooks.d.ts"/>
/**
 * Created by liwei on 2016/6/24.
 */
import {Component, OnInit} from '@angular/core';
import {Router, RouteParams}            from '@angular/router-deprecated';

import {RepoService} from './img.service'

@Component({
    selector: 'my-repos',
    templateUrl: 'app/docker/image/img.component.html'
})
export class ReposComponent implements OnInit {

    repos:string[];
    error:any;

    constructor(private router:Router,
                private repoService:RepoService) {
    }

    getRepos() {
        return this.repoService.getRepositories()
            .then(repos => this.repos = repos)
            .catch(error => this.error = error)
    }

    ngOnInit():any {
        this.getRepos()
    }

    gotoTags(group:string, image:string) {
        if (group == "libs") {
            group = "";
        }
        let link = ['Tags', {group: group, name: image}];

        this.router.navigate(link)
    }

}


@Component({
    selector: 'my-tags',
    templateUrl: 'app/docker/image/tags.component.html',
})
export class TagsComponent implements OnInit {
    error:any;
    tags:{};
    group:string;
    image:string;

    constructor(private router:Router,
                private routeParams:RouteParams,
                private repoService:RepoService) {
    }


    getTags(group:string, name:string) {
        return this.repoService.getTags(group, name)
            .then(tags=>this.tags = tags)
            .catch(error =>this.error = error)
    }

    ngOnInit():any {
        let group = this.routeParams.get('group');
        let name = this.routeParams.get('name');
        this.group = group;
        this.image = name;
        this.getTags(group, name)

    }

    gotoTag(group:string, image:string, tag:string) {
        if (group == "libs") {
            group = "";
        }
        let link = ['Tag', {group: group, name: image, version: tag}];

        this.router.navigate(link)
    }
}

@Component({
    selector: 'my-tag',
    templateUrl: 'app/docker/image/tag-detail.component.html'
})
export class TagComponent implements OnInit {
    error:any;
    group:string;
    image:string;
    version:string;
    tag:Tag;

    constructor(private routeParams:RouteParams,
                private repoService:RepoService) {
    }


    //get tag from server

    getTag(group:string, name:string, version:string) {
        return this.repoService.getTag(group, name, version)
            .then(tag => {
                    this.tag = tag;
                    console.dir(this.tag)
                }
            )
            .catch(error =>this.error = error)
    }

    ngOnInit():any {
        let group = this.routeParams.get('group');
        let name = this.routeParams.get('name');
        let version = this.routeParams.get('version');
        this.group = group;
        this.image = name;
        this.version = version;
        this.getTag(group, name, version);
    }
}
