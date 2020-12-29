import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PeerCallingComponent } from './peer-calling.component';

describe('PeerCallingComponent', () => {
  let component: PeerCallingComponent;
  let fixture: ComponentFixture<PeerCallingComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PeerCallingComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PeerCallingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
