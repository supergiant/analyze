import { Component, OnInit } from '@angular/core';
import { HttpClient }        from "@angular/common/http";
import { map }               from "rxjs/operators";
import { Observable }        from "rxjs";

@Component({
  selector: 'app-checks',
  templateUrl: './checks.component.html',
  styleUrls: ['./checks.component.scss']
})
export class ChecksComponent implements OnInit {

  checks$: Observable<any[]>;

  constructor(
    private http: HttpClient
  ) {
  }

  ngOnInit() {
    const apiV1Check = 'http://ec2-52-53-217-176.us-west-1.compute.amazonaws.com:31146/api/v1/check';
    this.checks$ = this.http.get(apiV1Check).pipe(
      map((checks: any[]) => checks.map(check => {
        try {
          const description = JSON.parse(check.description);
          return {
            ...check, description
          }
        } catch (e) {
        //  is not json
          return check;
        }
      }))
    );
  }

  isObject(val) {
    return typeof val === 'object';
  }


}
