package subsonic

import (
	"net/http"

	"github.com/deluan/navidrome/model/request"
	"github.com/deluan/navidrome/server/subsonic/responses"
)

type UsersController struct{}

func NewUsersController() *UsersController {
	return &UsersController{}
}

// TODO This is a placeholder. The real one has to read this info from a config file or the database
func (c *UsersController) GetUser(w http.ResponseWriter, r *http.Request) (*responses.Subsonic, error) {
	loggedUser, ok := request.UserFrom(r.Context())
	if !ok {
		return nil, newError(responses.ErrorGeneric, "Internal error")
	}
	response := newResponse()
	response.User = &responses.User{}
	response.User.Username = loggedUser.UserName
	response.User.AdminRole = loggedUser.IsAdmin
	response.User.Email = loggedUser.Email
	response.User.StreamRole = true
	response.User.DownloadRole = true
	response.User.ScrobblingEnabled = true
	return response, nil
}

func (c *UsersController) GetUsers(w http.ResponseWriter, r *http.Request) (*responses.Subsonic, error) {
	loggedUser, ok := request.UserFrom(r.Context())
	if !ok {
		return nil, newError(responses.ErrorGeneric, "Internal error")
	}
	user := responses.User{}
	user.Username = loggedUser.Name
	user.AdminRole = loggedUser.IsAdmin
	user.Email = loggedUser.Email
	user.StreamRole = true
	user.DownloadRole = true
	user.ScrobblingEnabled = true
	response := newResponse()
	response.Users = &responses.Users{User: []responses.User{user}}
	return response, nil
}
