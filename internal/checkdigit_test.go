package internal

import (
	"testing"
)

func TestCheckDigit(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		digit     string
		algo      CheckDigit
		wantErr   bool
		mustMatch bool
	}{
		{
			name:      "ups",
			id:        "30AA33019867867",
			digit:     "8",
			algo:      NewMod10(1, 2),
			wantErr:   false,
			mustMatch: true,
		},
		{
			name:      "mod10Valid",
			id:        "123456789",
			digit:     "5",
			algo:      NewMod10(1, 2),
			wantErr:   false,
			mustMatch: true,
		},
		{
			name:      "mod10Valid0",
			id:        "1234567890",
			digit:     "0",
			algo:      NewMod10(1, 2),
			wantErr:   false,
			mustMatch: false,
		},
		{
			name:      "mod10chars",
			id:        "987654321GG",
			digit:     "1",
			algo:      NewMod10(1, 2),
			wantErr:   false,
			mustMatch: true,
		},
		{
			name:      "Amazon",
			id:        "TBA000000000000",
			digit:     "",
			algo:      NewNoop(),
			wantErr:   false,
			mustMatch: true,
		},
		{
			name:      "Canada Post Match",
			id:        "703511447713847",
			digit:     "2",
			algo:      NewMod10(3, 1),
			wantErr:   false,
			mustMatch: true,
		},
		{
			name:      "Canada Post Mismatch",
			id:        "897470275659911",
			digit:     "5",
			algo:      NewMod10(3, 1),
			wantErr:   false,
			mustMatch: false,
		},
		{
			name:      "DHL Express alph",
			id:        "JVGL099999999",
			digit:     "0",
			algo:      NewMod7(),
			wantErr:   false,
			mustMatch: true,
		},
		{
			name:      "DHL Express",
			id:        "331881002",
			digit:     "5",
			algo:      NewMod7(),
			wantErr:   false,
			mustMatch: true,
		},
		{
			name:  "Fedex Express (12)",
			id:    "98657878885",
			digit: "5",
			algo: NewSumProductWithWeightingsAndModulo(
				[]int{3, 1, 7, 3, 1, 7, 3, 1, 7, 3, 1},
				11,
				10,
			),
			wantErr:   false,
			mustMatch: true,
		},
		{
			name:  "Fedex Express (34)",
			id:    "100192133425000100030077901797269",
			digit: "7",
			algo: NewSumProductWithWeightingsAndModulo(
				[]int{1, 7, 3, 1, 7, 3, 1, 7, 3, 1, 7, 3, 1},
				11,
				10,
			),
			wantErr:   false,
			mustMatch: true,
		},
		{
			name:      "s10",
			id:        "12345678",
			digit:     "5",
			algo:      NewS10([]int{8, 6, 4, 2, 3, 5, 9, 7}, []string{"Courier"}),
			wantErr:   false,
			mustMatch: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.algo.Generate(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("check digit error = %v; expected %v", err, tt.wantErr)
				return
			}
			if tt.mustMatch {
				if got != tt.digit {
					t.Errorf("%s: check digit mismatch = %v, want %v", tt.name, got, tt.digit)
				}
			} else {
				if got == tt.digit {
					t.Errorf("check digit should not match = %v, want %v", got, tt.digit)
				}
			}
		})
	}
}
