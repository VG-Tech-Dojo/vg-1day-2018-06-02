package controller

import (
  "database/sql"
  "net/http"

  "github.com/VG-Tech-Dojo/vg-1day-2018-06-02/tockn/httputil"
  "github.com/VG-Tech-Dojo/vg-1day-2018-06-02/tockn/model"
  "github.com/gin-gonic/gin"
)

type Room struct {
  DB     *sql.DB
  Stream chan *model.Room
}

func (m *Room) RoomAll(c *gin.Context) {
  rooms, err := model.RoomsAll(m.DB)
  if err != nil {
    resp := httputil.NewErrorResponse(err)
    c.JSON(http.StatusInternalServerError, resp)
    return
  }

  if len(rooms) == 0 {
    c.JSON(http.StatusOK, make([]*model.Room, 0))
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "result": rooms,
    "error":  nil,
  })
}

func (m *Room) GetByID(c *gin.Context) {
  rooms, err := model.RoomByID(m.DB, c.Param("room_id"))

  switch {
  case err == sql.ErrNoRows:
    resp := httputil.NewErrorResponse(err)
    c.JSON(http.StatusNotFound, resp)
    return
  case err != nil:
    resp := httputil.NewErrorResponse(err)
    c.JSON(http.StatusInternalServerError, resp)
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "result": rooms,
    "error":  nil,
  })
}