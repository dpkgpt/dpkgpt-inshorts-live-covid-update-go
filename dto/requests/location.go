package requests

import _ "github.com/pdrum/swagger-automation/api"

// swagger:parameters requestParams
type Location struct {
	//this will be description for request params
	//in:body
	Lat  string `json:"lat"`
	Long string `json:"long"`
}
