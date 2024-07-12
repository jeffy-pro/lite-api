package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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
		// Add more test cases here if there are other specific scenarios you want to test
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, tt.currency.Validate())
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
			assert.Equal(t, tt.wantErr, tt.country.Validate())
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
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIntegerList_Parse(t *testing.T) {
	tests := []struct {
		name    string
		il      IntegerList
		want    []int
		wantErr bool
	}{
		{
			name:    "Valid integer list",
			il:      IntegerList("1,2,3"),
			want:    []int{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "Valid integer list with spaces",
			il:      IntegerList(" 1, 2, 3 "),
			want:    []int{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "Empty string",
			il:      IntegerList(""),
			want:    []int{},
			wantErr: false,
		},
		{
			name:    "Single number",
			il:      IntegerList("42"),
			want:    []int{42},
			wantErr: false,
		},
		{
			name:    "Invalid number in list",
			il:      IntegerList("1,2,a"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid format",
			il:      IntegerList("1,2,3,"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Negative numbers",
			il:      IntegerList("-1,0,1"),
			want:    []int{-1, 0, 1},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.il.Parse()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, got)
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
				assert.Error(t, err)
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
