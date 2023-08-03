package models

type Address struct {
	Model

	FirstName    string `gorm:"column:first_name" json:"firstName"`
	LastName     string `gorm:"column:last_name"         json:"lastName"`
	Plaque       int    `gorm:"column:plaque"            json:"plaque"`
	PhoneNumber  string `gorm:"column:phone_number"      json:"phoneNumber"`
	NationalCode string `gorm:"column:national_code"     json:"nationalCode"`
	PostalCode   string `gorm:"column:postal_code"     json:"postalCode"`
	Address      string `gorm:"column:address"              json:"address"`
}

func (Address) TableName() string {
	return "address"
}
