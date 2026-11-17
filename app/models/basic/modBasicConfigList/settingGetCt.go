package modBasicConfigList

type SettingGet struct {
	EventNo string `json:"eventNo" binding:"required" validate:"required" label:"事件编号"`
}
