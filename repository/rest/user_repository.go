package rest

import (
	"github.com/sampila/go-utils/rest_errors"
	"github.com/sampila/oauth2-server-pg/domain/user"
	"github.com/mercadolibre/golang-restclient/rest"
	"time"
	"encoding/json"
	"errors"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:9001",
		Timeout: 3000 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (map[string]interface{}, rest_errors.RestErr)
}

type usersRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(username string, password string) (map[string]interface{}, rest_errors.RestErr) {
	request := user.UserLoginRequest{
		Email: username,
		Password: password,
	}

	response := usersRestClient.Post("/login", request)

	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalServerError("invalid restclient response when trying to login user", errors.New("restclient error"))
	}

	if response.StatusCode > 299 {
		apiErr, err := rest_errors.NewRestErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to login user", err)
		}
		return nil, apiErr
	}

	var respondJson map[string]interface{}
	if err := json.Unmarshal(response.Bytes(), &respondJson); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal users login response", errors.New("json parsing error"))
	}
	return respondJson, nil
}
