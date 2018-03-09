import { Component } from '@angular/core';
import { IMyDrpOptions, IMyDateRangeModel } from 'mydaterangepicker';
import { Observable } from 'rxjs/Observable';
import { ReportService } from './report.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  // title = 'app';

  public reviews;
  // : Review[];

  constructor(private _reportService: ReportService) {}

  myDateRangePickerOptions: IMyDrpOptions = {
      dateFormat: 'dd.mm.yyyy',
  };
  private model: any = {beginDate: {year: 2018, month: 10, day: 9},
                        endDate: {year: 2018, month: 10, day: 19}};

  onDateRangeChanged(event: IMyDateRangeModel) {
  // event properties are: event.beginDate, event.endDate, event.formatted,
  // event.beginEpoc and event.endEpoc
    console.log(event);
  }

  onSubmitNgModel(): void {
    this._reportService.postRange(this.model).subscribe(
      data => {
        this.reviews = data;
        console.log(this.reviews);
      },
      err => console.error(err),
      () => console.log('done loading ranged reviews')
    );
  }

  // ngOnInit() {
    // this.getReviews();
  // }

  getReviews() {
    this._reportService.getReviews().subscribe(
      data => {
        this.reviews = data;
        console.log(this.reviews);
      },
      err => console.error(err),
      () => console.log('done loading reviews')
    );
  }
}

interface Review {
  id: string;
  time_stamp: string;
  author: string;
  date: string;
  post: string;
  link: string;
}
