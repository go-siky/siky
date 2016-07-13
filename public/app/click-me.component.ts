import { Component } from '@angular/core';

@Component({
  selector: 'click-me',
  template: `
    <button class="btn btn-default" (click)="onClickMe()">Click me!</button>
    {{clickMessage}}`
})
export class ClickMeComponent {
  clickMessage = '';

  onClickMe() {
    this.clickMessage = 'You are my hero!';
  }
}
