package business

import "github.com/jinzhu/gorm"

//noinspection SpellCheckingInspection
type Card struct {
	gorm.Model
	IMSI     string `gorm:"unique"`
	ICCID    string `gorm:"column:iccid"`
	Phone    string
	Name     string
	Note     string
	Metadata string `gorm:"type:text"`
}
