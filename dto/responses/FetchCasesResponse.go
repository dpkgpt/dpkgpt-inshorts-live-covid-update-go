package responses

import (
	"crud/entity"

	_ "github.com/pdrum/swagger-automation/api"
)

// This text will appear as description of your response body.
// swagger:response fetchCasesResponse
type FetchCasesResponse struct {
	//in:body
	Result [2]entity.CovidCases `json:"result"`
}
