package stand

import (
	"github.com/jinzhu/gorm"
)

type Stand struct {
	gorm.Model
	StandID        string
	Lat            float64 `validate:"min=-90,max=90"`
	Lng            float64 `validate:"min=-180,max=180"`
	Source         string
	SourceID       string
	Name           string
	Type           string `validate:"nonzero"`
	NumberOfStands int
	ImageID        string
	PublicImageURL string
	Notes          string
	Checked        string
	Verified       bool
	LastUpdateBy   string
	Token          string
}
