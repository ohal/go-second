import { CommonModule } from "@angular/common";
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from "@angular/forms";
import { HttpClientModule } from '@angular/common/http';
import { AppComponent } from './app.component';
import { MyDateRangePickerModule } from 'mydaterangepicker';
import { ReportService } from './report.service';


@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
  	CommonModule,
  	FormsModule,
    BrowserModule,
    MyDateRangePickerModule,
    HttpClientModule
  ],
  providers: [ReportService],
  bootstrap: [AppComponent]
})
export class AppModule { }
