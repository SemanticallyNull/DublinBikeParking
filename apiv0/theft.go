package apiv0

import "github.com/jinzhu/gorm"

type Theft struct {
	gorm.Model
	Stand Stand
	StandID *uint
}