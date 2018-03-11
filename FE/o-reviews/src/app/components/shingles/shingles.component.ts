import { Component, OnInit } from '@angular/core';
import { IMyDrpOptions, IMyDateRangeModel } from 'mydaterangepicker';
import { Observable } from 'rxjs/Observable';
import { ReportService } from '../../report.service';

@Component({
  selector: 'app-shingles',
  templateUrl: './shingles.component.html',
  styleUrls: ['./shingles.component.css']
})
export class ShinglesComponent implements OnInit {

  public reviews;

  constructor(private _reportService: ReportService) {}

  myDateRangePickerOptions: IMyDrpOptions = {
      dateFormat: 'dd.mm.yyyy',
  };
  private model: any = {beginDate: {year: 2018, month: 1, day: 1},
                        endDate: {year: 2018, month: 12, day: 31}};

  onDateRangeChanged(event: IMyDateRangeModel) {
  // event properties are: event.beginDate, event.endDate, event.formatted,
  // event.beginEpoc and event.endEpoc
    console.log(event);
  }

  ngOnInit() {
    // this.getReviews();
  }

  onSubmitNgModelShingle(): void {
    this._reportService.postShingle(this.model).subscribe(
      data => {
        this.reviews = data;
        console.log(this.reviews);
      },
      err => console.error(err),
      () => console.log('done loading top shingles')
    );
  }
}
