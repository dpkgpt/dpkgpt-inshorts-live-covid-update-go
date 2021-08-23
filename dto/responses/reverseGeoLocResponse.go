package responses

type RevGeoLocResponse struct {
	ResponseCode int             `json:"responseCode"`
	Results      []AddressDetail `json:"results"`
}
