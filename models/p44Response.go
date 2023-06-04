package models

import "time"

type Response struct {
	Data []struct {
		Bols []struct {
			ID     string `json:"id"`
			Number string `json:"number"`
			URL    string `json:"url"`
		} `json:"bols"`
		Bookings []struct {
			ID     string `json:"id"`
			Number string `json:"number"`
			URL    string `json:"url"`
		} `json:"bookings"`
		Carrier struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Scac string `json:"scac"`
		} `json:"carrier"`
		Equipment struct {
			Category string `json:"category"`
			Number   string `json:"number"`
			Type     string `json:"type"`
		} `json:"equipment"`
		Events []struct {
			ActionType string    `json:"action_type"`
			ActualTime time.Time `json:"actual_time"`
			Location   struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Unlocode string `json:"unlocode"`
			} `json:"location"`
			LocationType  string      `json:"location_type"`
			PredictedTime interface{} `json:"predicted_time"`
			StageType     string      `json:"stage_type"`
			Timezone      string      `json:"timezone"`
			Vehicle       struct {
			} `json:"vehicle"`
			VehicleType string `json:"vehicle_type"`
		} `json:"events"`
		ID        string        `json:"id"`
		Orders    []interface{} `json:"orders"`
		Shipments []interface{} `json:"shipments"`
		Statuses  []struct {
			ActualTime interface{} `json:"actual_time"`
			Location   struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Unlocode string `json:"unlocode"`
			} `json:"location"`
			LocationType  string      `json:"location_type"`
			PredictedTime interface{} `json:"predicted_time"`
			StageType     string      `json:"stage_type"`
			StatusType    string      `json:"status_type"`
			Timezone      string      `json:"timezone"`
		} `json:"statuses"`
		URL string `json:"url"`
	} `json:"data"`
	NextURL string `json:"next_url"`
	URL     string `json:"url"`
}
