package models

type Traffic struct {
	VehicleCount int    `json:"vehicle_count"`
	LaneNo       int    `json:"lane_no"`
	SignalId     int    `json:"signal_id"`
	CurrentColor string `json:"current_color"`
}
