package responses

type AddressDetail struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Pincode string `json:"pincode"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
	Area    string `json:"area"`
}
