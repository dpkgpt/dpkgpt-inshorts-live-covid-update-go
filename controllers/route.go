package controllers

import (
	"crud/dto/requests"
	"crud/service"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/pdrum/swagger-automation/api"
)

// swagger:route GET /getcovidcases foobar-tag idOfFoobarEndpoint
// Foobar does some amazing stuff.
// responses:
//   200: fetchCovidCasesResponse
//   500: error response
func GetCovidCasesByLocation(w http.ResponseWriter, r *http.Request) {

	location := &requests.Location{Lat: r.URL.Query().Get("lat"), Long: r.URL.Query().Get("long")}
	log.Printf("user's location is lat: %s , long: %s", location.Lat, location.Long)
	response, err := service.FetchCovidDataForIndiaAndState(location)
	if response != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(err)
}

func InsertCovidCasesForStateInDB(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var request requests.UpdateCovidCasesRequest
	json.Unmarshal(requestBody, &request)
	log.Println("request cast to struct:", request)
	response, err := service.UpdateCovidData(&request)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	} else {
		println("response is: ", response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}

}
