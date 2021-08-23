package adapter

import (
	"crud/dto/customerrors"
	"crud/dto/responses"
	"crud/env"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func GetStateFromLocation(lat, long string) (string, error) {
	url := env.GetValue("REVERSE_GEO_LOC_API")
	url = fmt.Sprintf(url, lat, long)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "", customerrors.GetBaseErrorWithDefaultMessage("REVLOC-API")
	}

	log.Println("response:", resp)
	var revGeoLocResponse responses.RevGeoLocResponse
	json.NewDecoder(resp.Body).Decode(&revGeoLocResponse)
	log.Println("location response:", revGeoLocResponse)
	if len(revGeoLocResponse.Results) == 0 {
		log.Println("no details found for location")
		return "", customerrors.GetBaseError("REVLOC-404", "No Details found for the given co-ordinates")
	}
	country := revGeoLocResponse.Results[0].Area
	if country != "India" {
		log.Println("Location is not of India")
		return "", customerrors.GetBaseError("REVLOC-404", "Location is not of India.")
	}
	return strings.ToUpper(revGeoLocResponse.Results[0].State), nil
}
