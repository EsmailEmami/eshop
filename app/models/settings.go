package models

type SettingsReqModel struct {
	Column string `json:"column"`
	Value  any    `json:"value,omitempty"`
}
