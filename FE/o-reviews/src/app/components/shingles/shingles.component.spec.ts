import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShinglesComponent } from './shingles.component';

describe('ShinglesComponent', () => {
  let component: ShinglesComponent;
  let fixture: ComponentFixture<ShinglesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShinglesComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShinglesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
