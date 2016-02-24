package app

import (
	"time"
)

type ErrorTo struct {
	Code int `json:"code"`
}

type SchemaTO struct {
	Id        int          `json:"id"`
	Name      string       `json:"name"`
	State     int          `json:"state"`
	Timeslots []TimeslotTO `json:"timeslots"`
}

type TimeslotTO struct {
	Id        int    `json:"id"`
	Schema    int    `json:"schema"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Locations []int  `json:"locations"`
}

type LocationTO struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PeriodTO struct {
	Id        int       `json:"id"`
	Schema    int       `json:"schema"`
	Status    int       `json:"status"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type ReservationTO struct {
	Id       int    `json:"id"`
	Timeslot int    `json:"timeslot"`
	Location int    `json:"location"`
	Period   int    `json:"period"`
	User     string `json:"user"`
}
