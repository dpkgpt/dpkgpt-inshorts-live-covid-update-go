package responses

import (
	"crud/entity"
	"time"
)

type LiveCovidDataResponse struct {
	ActiveCases     int                  `json:"activeCases"`
	LastUpdatedTime time.Time            `json:"lastUpdatedAtApify"`
	RegionData      []*entity.CovidCases `json:"regionData"`
}
