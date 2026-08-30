package shipping

import (
	"errors"
	"testing"
)

func TestCalculate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		request Request
		want Quote
		wantErr error
	}{
		{
			name: "local package under one kilogram",
			request: Request{
				WeightGrams: 500,
				Zone:        ZoneLocal,
				Service:     ServiceStandard,
			},
			want: Quote{
				BasePrice:        10000,
				WeightSurcharge:  0,
				ServiceSurcharge: 0,
				Total:            10000,
			},
		},
		{
			name: "partial extra kilogram rounds up",
			request: Request{
				WeightGrams: 1001,
				Zone:        ZoneLocal,
				Service:     ServiceStandard,
			},
			want: Quote{
				BasePrice:        10000,
				WeightSurcharge:  2000,
				ServiceSurcharge: 0,
				Total:            12000,
			},
		},
		{
			name: "regional express package",
			request: Request{
				WeightGrams: 2500,
				Zone:        ZoneRegional,
				Service:     ServiceExpress,
			},
			want: Quote{
				BasePrice:        20000,
				WeightSurcharge:  10000,
				ServiceSurcharge: 15000,
				Total:            45000,
			},
		},
		{
			name: "international package at maximum weight",
			request: Request{
				WeightGrams: 20000,
				Zone:        ZoneInternational,
				Service:     ServiceStandard,
			},
			want: Quote{
				BasePrice:        50000,
				WeightSurcharge:  285000,
				ServiceSurcharge: 0,
				Total:            335000,
			},
		},
		{
			name: "zero weight",
			request: Request{
				WeightGrams: 0,
				Zone:        ZoneLocal,
				Service:     ServiceStandard,
			},
			wantErr: ErrInvalidWeight,
		},
		{
			name: "weight exceeds maximum",
			request: Request{
				WeightGrams: 20001,
				Zone:        ZoneLocal,
				Service:     ServiceStandard,
			},
			wantErr: ErrInvalidWeight,
		},
		{
			name: "invalid zone",
			request: Request{
				WeightGrams: 1000,
				Zone:        Zone("UNKNOWN"),
				Service:     ServiceStandard,
			},
			wantErr: ErrInvalidZone,
		},
		{
			name: "invalid service level",
			request: Request{
				WeightGrams: 1000,
				Zone:        ZoneLocal,
				Service:     ServiceLevel("UNKNOWN"),
			},
			wantErr: ErrInvalidService,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()
				
				got, err := Calculate(
					test.request,
				)

				if test.wantErr != nil {
					if !errors.Is(
						err,
						test.wantErr,
					) {
						t.Fatalf(
							"expected error %v, got %v",
							test.wantErr,
							err,
						)
					}

					return
				}

				if err != nil {
					t.Fatalf(
						"unexpected error: %v",
						err,
					)
				}

				if got != test.want {
					t.Errorf(
						"expected %+v, got %+v",
						test.want,
						got,
					)
				}
			},
		)
	}
}