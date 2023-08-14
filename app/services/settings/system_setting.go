package settings

var systemSetting = new(SystemSetting)

type SystemSetting struct {
	FileExpireTimeStampts *int `column:"file_expire_time_stampts"`
}

func (SystemSetting) TableName() string {
	return "system_setting"
}

func (SystemSetting) SchemaName() string {
	return "public"
}

func (s *SystemSetting) LoadDefaultValues() {
	s.FileExpireTimeStampts = func() *int {
		value := 4
		return &value
	}()
}

func GetSystemSettings() (*SystemSetting, error) {
	if systemSetting == nil {
		bindSystemSettings()
	}

	return systemSetting, nil
}
