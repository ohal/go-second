import { Component, OnInit } from '@angular/core';
import { ReportService } from '../../report.service';

@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})

export class ListComponent implements OnInit {
  // title = 'app';

  public reviews;
  // : Review[];

  constructor(private _reportService: ReportService) {}

  ngOnInit() {
    // this.getReviews();
  }

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
