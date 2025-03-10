package samples

import (
	"yip/src/api/services/dto"
)

var AdminUser = dto.SignInRequest{
	Email:     "leonard.schellenberg@gmail.com",
	Password:  "sechszig",
	Audiences: []string{"http://localhost:8080"},
}

var AnotherUser = dto.RegisterRequest{
	Email:    "another.user@gmail.com",
	Password: "mybloodypassword",
}
