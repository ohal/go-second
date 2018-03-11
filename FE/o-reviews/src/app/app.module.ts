import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule } from '@angular/common/http';

import { AppComponent } from './app.component';
import { ReviewsComponent } from './components/reviews/reviews.component';
import { ShinglesComponent } from './components/shingles/shingles.component';
import { ListComponent } from './components/list/list.component';

import { MyDateRangePickerModule } from 'mydaterangepicker';
import { ReportService } from './report.service';


@NgModule({
  declarations: [
    AppComponent,
    ReviewsComponent,
    ShinglesComponent,
    ListComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    CommonModule,
    FormsModule,
    MyDateRangePickerModule,
    HttpClientModule
  ],
  providers: [ReportService],
  bootstrap: [AppComponent]
})
export class AppModule { }
