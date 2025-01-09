package models

type SignalGroup struct {
	SignalGroup []GroupSignal `json:"signal_group"`
}

type GroupSignal struct {
	GroupSignalId string         `json:"group_id,omitempty" bson:"_id,omitempty"`
	GroupName     string         `json:"group_name"`
	SignalCount   int            `json:"signal_count"`
	Signals       []SingleSignal `json:"signals"`
}

type SingleSignal struct {
	SingleSignalId string `json:"signal_id,omitempty" bson:"_id,omitempty"`
	LaneNo         int    `json:"lane_no"`
	CurrentColor   string `json:"current_color"`
	VehicleCount   int    `json:"vehicle_count" bson:"vehicle_count"`
	GreenDuration  int    `json:"green_duration"`
	YellowDuration int    `json:"yellow_duration"`
	RedDuration    int    `json:"red_duration"`
}

type UpdateSignalCountGroup struct {
	Signals []UpdateVehicleCountRequest `json:"signals"`
}

type UpdateVehicleCountRequest struct {
	SignalSingleId string `json:"signal_id"`
	VehicleCount   int    `json:"vehicle_count"`
	GreenDuration  int    `json:"green_duration"`
	YellowDuration int    `json:"yellow_duration"`
	RedDuration    int    `json:"red_duration"`
}
