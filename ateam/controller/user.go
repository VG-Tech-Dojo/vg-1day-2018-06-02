package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/ateam/httputil"
	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/ateam/model"
	"github.com/gin-gonic/gin"
)

// User is controller for requests to users
type User struct {
	DB     *sql.DB
	Stream chan *model.User
}

func (u *User) Create(c *gin.Context) {
	var user model.User

	inserted, err := user.Insert(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// bot対応
	// u.Stream <- inserted

	c.JSON(http.StatusCreated, gin.H{
		"result": inserted,
		"error":  nil,
	})
}
