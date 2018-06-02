package model

import (
  "database/sql"
)

type Room struct {
  ID       int64  `json:"id"`
  Name     string `json:"name"`
  Image_url string `json:"image_url"`
  Birth string `json:"birth"`
}

func RoomsAll(db *sql.DB) ([]*Room, error) {
  rows, err := db.Query(`select id, name, image_url, birth from room`)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var rm []*Room
  for rows.Next() {
    r := &Room{}
    if err := rows.Scan(&r.ID, &r.Name, &r.Image_url, &r.Birth); err != nil {
      return nil, err
    }
    rm = append(rm, r)
  }
  if err := rows.Err(); err != nil {
    return nil, err
  }
  return rm, nil
}

func RoomByID(db *sql.DB, id string) (*Room, error) {
  r := &Room{}

  if err := db.QueryRow(`select id, name, image_url, birth from room where id = ?`, id).Scan(&r.ID, &r.Name, &r.Image_url, &r.Birth); err != nil {
    return nil, err
  }

  return r, nil
}