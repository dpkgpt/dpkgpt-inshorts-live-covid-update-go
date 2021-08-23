package requests

type UpdateCovidCasesRequest struct {
	Region string `json:"region"`
	Change int    `json:"change"`
}
