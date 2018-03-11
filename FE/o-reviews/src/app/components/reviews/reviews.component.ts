import { Component, OnInit } from '@angular/core';
import { IMyDrpOptions, IMyDateRangeModel } from 'mydaterangepicker';
import { Observable } from 'rxjs/Observable';
import { ReportService } from '../../report.service';

@Component({
  selector: 'app-reviews',
  templateUrl: './reviews.component.html',
  styleUrls: ['./reviews.component.css']
})
export class ReviewsComponent implements OnInit {

  public reviews;

  constructor(private _reportService: ReportService) {}

  myDateRangePickerOptions: IMyDrpOptions = {
      dateFormat: 'dd.mm.yyyy',
  };
  private model: any = {beginDate: {year: 2017, month: 1, day: 1},
                        endDate: {year: 2017, month: 12, day: 31}};

  ngOnInit() {
    // this.getReviews();
  }

  onSubmitNgModelRange(): void {
    this._reportService.postRange(this.model).subscribe(
      data => {
        this.reviews = data;
        console.log(this.reviews);
      },
      err => console.error(err),
      () => console.log('done loading ranged reviews')
    );
  }
}
