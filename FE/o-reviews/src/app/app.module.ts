import { CommonModule } from "@angular/common";
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from "@angular/forms";

import { AppComponent } from './app.component';

import { MyDateRangePickerModule } from 'mydaterangepicker';


@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
  	CommonModule,
  	FormsModule,
    BrowserModule,
    MyDateRangePickerModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
