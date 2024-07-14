package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCurrency_Validate(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		wantErr  error
	}{
		{
			name:     "Valid currency",
			currency: Currency("USD"), // Assuming USD is in allowedCurrencies
			wantErr:  nil,
		},
		{
			name:     "Invalid currency",
			currency: Currency("XYZ"),
			wantErr:  ErrCurrencyNotFound,
		},
		{
			name:     "Empty currency",
			currency: Currency(""),
			wantErr:  ErrCurrencyNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, tt.currency.Validate())
		})
	}
}

func TestCountry_Validate(t *testing.T) {
	tests := []struct {
		name    string
		country Country
		wantErr error
	}{
		{
			name:    "Valid country",
			country: Country("US"), // Assuming US is in allowedCountries
			wantErr: nil,
		},
		{
			name:    "Invalid country",
			country: Country("XYZ"),
			wantErr: ErrCountryNotAllowed,
		},
		{
			name:    "Empty country",
			country: Country(""),
			wantErr: ErrCountryNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, tt.country.Validate())
		})
	}
}

func TestDateString_Validate(t *testing.T) {
	tests := []struct {
		name    string
		date    DateString
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Valid date",
			date:    DateString("2006-01-02"),
			want:    time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Invalid date format",
			date:    DateString("2006/01/02"),
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Empty date",
			date:    DateString(""),
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Invalid month",
			date:    DateString("2006-13-01"),
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Invalid day",
			date:    DateString("2006-01-32"),
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Leap year date",
			date:    DateString("2024-02-29"),
			want:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.date.Parse()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestIntegerList_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    IntegerList
		expected []int
		wantErr  bool
	}{
		{
			name:     "Valid input",
			input:    "129410,105360,106101,1762514",
			expected: []int{129410, 105360, 106101, 1762514},
			wantErr:  false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "Single number",
			input:    "42",
			expected: []int{42},
			wantErr:  false,
		},
		{
			name:     "Whitespace",
			input:    " 1, 2 , 3 ",
			expected: []int{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "Invalid number",
			input:    "1,2,3,abc,4",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Negative numbers",
			input:    "-1,-2,-3,4",
			expected: []int{-1, -2, -3, 4},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Parse()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.expected, got)
		})
	}
}

func TestOccupancyList_Parse(t *testing.T) {
	tests := []struct {
		name    string
		o       OccupancyList
		want    Occupancies
		wantErr bool
	}{
		{
			name: "Valid occupancy list",
			o:    OccupancyList(`[{"Rooms":1,"Adults":2,"Children":1},{"Rooms":2,"Adults":1,"Children":0}]`),
			want: Occupancies{
				{Rooms: 1, Adults: 2, Children: 1},
				{Rooms: 2, Adults: 1, Children: 0},
			},
			wantErr: false,
		},
		{
			name:    "Empty occupancy list",
			o:       OccupancyList(`[]`),
			want:    Occupancies{},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			o:       OccupancyList(`[{"Rooms":1,"Adults":2,"Children":1},]`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid occupancy structure",
			o:       OccupancyList(`[{"Rooms":"1","Adults":2,"Children":1}]`),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.o.Parse()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestOccupancies_Validate(t *testing.T) {
	tests := []struct {
		name    string
		o       Occupancies
		wantErr error
	}{
		{
			name: "Valid occupancies",
			o: Occupancies{
				{Rooms: 1, Adults: 2, Children: 1},
				{Rooms: 2, Adults: 1, Children: 0},
			},
			wantErr: nil,
		},
		{
			name: "Invalid: zero rooms",
			o: Occupancies{
				{Rooms: 0, Adults: 2, Children: 1},
			},
			wantErr: ErrMinOneRoomRequired,
		},
		{
			name: "Invalid: zero adults",
			o: Occupancies{
				{Rooms: 1, Adults: 0, Children: 1},
			},
			wantErr: ErrMinOneAdultRequired,
		},
		{
			name: "Invalid: negative children",
			o: Occupancies{
				{Rooms: 1, Adults: 2, Children: -1},
			},
			wantErr: ErrNegativeValue,
		},
		{
			name:    "Empty occupancies",
			o:       Occupancies{},
			wantErr: ErrEmptyOccupancies,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.Validate()
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
			}
		})
	}
}
