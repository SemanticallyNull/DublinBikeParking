package apiv0

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type APIUser struct {
	gorm.Model
	Email  string
	APIKey string
}

func (a *api) AuthenticateUser(apiKey string) (APIUser, error) {
	user := APIUser{}
	count := 0

	a.DB.Where("api_key = ?", apiKey).First(&user).Count(&count)
	if count == 0 {
		return APIUser{}, fmt.Errorf("could not find user for api key")
	}

	return user, nil
}
