package entity

import "time"

type CovidCases struct {
	Region      string     `json:"region" bson:"_id"`
	ActiveCases int        `json:"activeCases,omitempty" bson:"activeCases,omitempty"`
	ModifiedAt  *time.Time `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
	Msg         string     `json:"errorMessage,omitempty" bson:"msg,omitempty"`
}
