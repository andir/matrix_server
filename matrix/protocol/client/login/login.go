package login

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	p "github.com/andir/matrix_server/matrix/protocol"
	m "github.com/andir/matrix_server/matrix"
	"github.com/andir/matrix_server/matrix/database"
)

type ClientLoginRequest struct {
	Password string `json:"password"`
	Medium string `json:"medium,omitempty"`
	Type string `json:"type"`
	User string `json:"type"`
	Address string `json:"address"`
}

type ClientLoginResponse struct {

}

func ClientLoginRequestFromHTTPRequest(r *http.Request) (ClientLoginRequest, error) {
	clientLoginRequest := ClientLoginRequest{}
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return clientLoginRequest, err
	}

	err = json.Unmarshal(buffer, &clientLoginRequest)
	return clientLoginRequest, err
}

func (loginRequest *ClientLoginRequest) Handle() (status_code int, response interface{}) {

	if loginRequest.Type != "m.login.password" {
		return 400, p.ErrorResponse{
			Errcode: m.M_UNKNOWN,
			Error: "Bad login type.",
		}
	}
	if loginRequest.User == "" {
		return 400, p.ErrorResponse{
			Errcode: m.M_INVALID_USERNAME,
			Error: "You must supply a user.",
		}
	}

	if loginRequest.Password == "" {
		return 403, p.ErrorResponse{
			Errcode: m.M_FORBIDDEN,
			Error: "Please supply a password",
		}
	}

	if loginRequest.Address != "" {
		return 400, p.ErrorResponse{
			Errcode: m.M_UNKNOWN,
			Error: "Address not yet supported", // FIXME
		}
	}
	var account database.Account
	database.DB.Where(&database.Account{Username: loginRequest.User}).First(&account)


	return 200, ClientLoginResponse{

	}
}

