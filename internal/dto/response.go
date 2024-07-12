package dto

// SearchResponse is the contract to respond search response with.
type SearchResponse struct {
	Data     HotelInfos `json:"data"`
	Supplier Supplier   `json:"supplier"`
}

// HotelInfo represents the information related to hotel for the query.
type HotelInfo struct {
	HotelID  string  `json:"hotelId"`
	Currency string  `json:"currency"`
	Price    float64 `json:"price"`
}

// HotelInfos is a collection of HotelInfo.
type HotelInfos []HotelInfo

// Supplier is the information added in lite API response to have transparency.
type Supplier struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}
