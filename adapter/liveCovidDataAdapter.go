package adapter

import (
	"crud/dto/customerrors"
	"crud/dto/responses"
	"crud/env"
	"encoding/json"
	"log"
	"net/http"
)

func FetchLiveCovidData() (*responses.LiveCovidDataResponse, error) {
	url := env.GetValue("LIVE_COVID_DATA")
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, customerrors.GetBaseErrorWithDefaultMessage("COVID-F")
	}

	log.Println("response:", resp)
	var liveCovidDataResponse responses.LiveCovidDataResponse
	json.NewDecoder(resp.Body).Decode(&liveCovidDataResponse)
	log.Println("Live covid Response response:", liveCovidDataResponse)
	if len(liveCovidDataResponse.RegionData) == 0 {
		log.Println("no details fetched")
		return nil, customerrors.GetBaseError("COVID-400", "No Live Details fetched for states")
	}
	return &liveCovidDataResponse, nil
}
