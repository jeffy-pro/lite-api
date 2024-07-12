package model

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrCurrencyNotFound    = errors.New("currency not allowed")
	ErrCountryNotAllowed   = errors.New("country not allowed")
	ErrMinOneRoomRequired  = errors.New("at least one room is required")
	ErrMinOneAdultRequired = errors.New("at least one adult is required")
	ErrNegativeValue       = errors.New("negative value not allowed")
	ErrEmptyOccupancies    = errors.New("empty occupancies")
)

var (
	// USD represents US Dollars.
	USD = Currency("USD")
	// EUR represents Euro.
	EUR = Currency("EUR")

	allowedCurrencies = []Currency{
		USD,
		EUR,
	}
)

var (
	// US represents United States.
	US = Country("US")

	// UK represents United Kingdom.
	UK = Country("UK")

	// ES represents the Kingdom of Spain
	ES = Country("ES")

	allowedCountries = []Country{
		US,
		UK,
		ES,
	}
)

// Currency is a concrete type to represent currency codes.
type Currency string

// String returns Currency as a string.
func (c Currency) String() string {
	return string(c)
}

// Validate checks if search in the specified currency is allowed.
func (c Currency) Validate() error {
	for _, currency := range allowedCurrencies {
		if c == currency {
			return nil
		}
	}

	return ErrCurrencyNotFound
}

// Country is a concrete type to represent currency codes.
type Country string

// String returns Country as a string.
func (c Country) String() string {
	return string(c)
}

// Validate checks if search in the specified country is allowed.
func (c Country) Validate() error {
	for _, country := range allowedCountries {
		if c == country {
			return nil
		}
	}

	return ErrCountryNotAllowed
}

// DateString represents a date string DateOnly format.
type DateString string

// String returns DateString as string.
func (d DateString) String() string {
	return string(d)
}

// Parse parses a date string.
func (d DateString) Parse() (time.Time, error) {
	return time.Parse(time.DateOnly, string(d))
}

// IntegerList represents a list of integers that can be initialized from a comma-separated string
type IntegerList string

// String returns IntegerList as string.
func (i IntegerList) String() string {
	return string(i)
}

// Parse converts the IntegerList to an array of integers
func (i IntegerList) Parse() ([]int, error) {
	nums := make([]int, 0)
	if err := json.Unmarshal([]byte(i), &nums); err != nil {
		return nil, err
	}

	return nums, nil
}

// Occupancy represents the requested occupancy details in Hotelbeds request.
type Occupancy struct {
	Adults   int `json:"adults"`
	Children int `json:"children"`
	Rooms    int `json:"rooms"`
}

// Occupancies is a collection of Occupancy.
type Occupancies []Occupancy

// OccupancyList is Occupancies serialized to string.
type OccupancyList string

// String returns OccupancyList as string.
func (o OccupancyList) String() string {
	return string(o)
}

// Parse parses OccupancyList into Occupancies.
func (o OccupancyList) Parse() (Occupancies, error) {
	occupancies := make(Occupancies, 0)
	err := json.Unmarshal([]byte(o), &occupancies)
	if err != nil {
		return nil, err
	}

	return occupancies, nil
}

// Validate validates each occupancy entry in the list.
func (o Occupancies) Validate() error {
	if len(o) == 0 {
		return ErrEmptyOccupancies
	}

	for _, occupancy := range o {
		if occupancy.Rooms < 1 {
			return ErrMinOneRoomRequired
		}

		if occupancy.Adults < 1 {
			return ErrMinOneAdultRequired
		}

		if occupancy.Children < 0 {
			return ErrNegativeValue
		}

	}

	return nil
}
