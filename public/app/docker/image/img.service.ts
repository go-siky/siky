/**
 * Created by liwei on 2016/6/24.
 */
import {Injectable}    from '@angular/core';
import {Headers, Http} from '@angular/http';
import 'rxjs/add/operator/toPromise';
@Injectable()
export class RepoService {
    // private repoUrl = '/v2/_catalog';
    private repoUrl = '/api/v2/repositories';

    constructor(private http:Http) {
    }

    getRepositories():Promise<string[]> {
        // let headers = new Headers({
        //     'Authorization':'Basic ZHp4eTpNZlkxMjM='
        // })
        return this.http.get(this.repoUrl)
            .toPromise().then(response => response.json())
            .catch(this.handleError)
    }

    getTags(group:string, name:string):Promise<{}> {
        let tagsUrl = '/api/v2/tags?group=' + group + '&name=' + name;
        return this.http.get(tagsUrl)
            .toPromise().then(response => response.json())
            .catch(this.handleError)
    }

    getTag(group:string, name:string, tag:string):Promise<Tag> {
        let tagUrl = '/api/v2/repository/tag/' + group + '---' + name + '---' + tag;
        return this.http.get(tagUrl)
            .toPromise().then(response => response.json())
            .catch(this.handleError)
    }

    private handleError(error:any) {
        console.error('An error occurred', error);
        return Promise.reject(error.message || error);
    }
}
