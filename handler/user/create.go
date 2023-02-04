package user

import (
	. "apiserver/handler"
	"apiserver/model"
	"apiserver/pkg/errno"

	"github.com/gin-gonic/gin"
)

// Create creates a new user account.
func Create(c *gin.Context) {
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	if err := r.checkParam(); err != nil {
		SendResponse(c, err, nil)
		return
	}

	u := model.UserModel{
		Username: r.Username,
		Password: r.Password,
	}

	// Validate the data
	if err := u.Validate(); err != nil {
		SendResponse(c, errno.ErrValidation, nil)
	}

	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	// Insert the user to the database.
	if err := u.Create(); err != nil {
		SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	rsp := CreateResponse{
		Username: r.Username,
	}

	// Show the user information.
	SendResponse(c, nil, rsp)
}

func (r *CreateRequest) checkParam() error {
	if r.Username == "" {
		return errno.New(errno.ErrValidation, nil).Add("username is empty.")
	}

	if r.Password == "" {
		return errno.New(errno.ErrValidation, nil).Add("password is empty.")
	}
	return nil
}
