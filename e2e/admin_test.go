package e2e

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"yip/e2e/samples"
	"yip/src/api/services/dto"
)

func TestSignIn(t *testing.T) {
	statusCode, response, err := client.SignIn(samples.AdminUser)
	assert.NoError(t, err, "didn't expect error")
	assert.Equal(t, 200, statusCode, "status code should be ok")
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.IdToken, "token should not be empty")
	client.SetToken(response.IdToken)
}

func TestRegisterUser(t *testing.T) {
	statusCode, response, err := client.RegisterUser(samples.AnotherUser)
	fmt.Println(response)
	assert.NoError(t, err, "didn't expect error")
	assert.Equal(t, 201, statusCode, "status code should be ok")
	assert.NotNil(t, response)
	assert.Equal(t, samples.AnotherUser.Email, response.Email, "email should be correct")
	assert.True(t, response.IsAdmin(), "should not be admin")
	newStatusCode, account, err := client.GetAccount(response.ID.String())
	assert.Equal(t, 200, newStatusCode, "status code should be ok")
	fmt.Println(account)
}

func TestRegisterUserAgain(t *testing.T) {
	statusCode, _, err := client.RegisterUser(samples.AnotherUser)
	assert.Error(t, err, "didn't expect error")
	assert.Equal(t, 409, statusCode, "status code should be ok")
}

func TestSignInRegisteredUser(t *testing.T) {
	// try login in to YIP -> should fail
	statusCode, _, err := client.SignIn(dto.SignInRequest{
		Email:     samples.AnotherUser.Email,
		Password:  samples.AnotherUser.Password,
		Audiences: samples.AdminUser.Audiences,
	})

	assert.Equal(t, 403, statusCode, "status code should be 403")
	assert.Error(t, err, "didn't expect error")

	// login in to resource service
	statusCode, resourceToken, err := client.SignIn(dto.SignInRequest{
		Email:     samples.AnotherUser.Email,
		Password:  samples.AnotherUser.Password,
		Audiences: []string{"http://localhost:8081"},
	})

	assert.NoError(t, err, "didn't expect error")
	assert.Equal(t, 200, statusCode, "status code should be ok")
	assert.NotNil(t, resourceToken)
}

func TestGetUsers(t *testing.T) {
	statusCode, response, err := client.GetUsers()
	assert.NoError(t, err, "didn't expect error")
	assert.Equal(t, 200, statusCode, "status code should be ok")
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Users), "email should be correct")
}

func TestSetEmail(t *testing.T) {
	statusCode, response, err := client.GetUsers()
	status, user, err := client.SetEmail(dto.SetEmailRequest{
		UserID: response.Users[0].ID,
		Email:  "some@email.com",
	})
	assert.Equal(t, 200, status, "status code should be ok")
	assert.NoError(t, err, "didn't expect error")
	assert.Equal(t, 200, statusCode, "status code should be ok")
	assert.NotNil(t, response)
	assert.Equal(t, "some@email.com", user.Email, "email should be correct")
	assert.True(t, user.IsEmailVerified, "email should be verified")
}
