package model

import (
	"database/sql"
	"math/rand"
	"time"
)

type User struct {
	ID   int64  `json:"id"`
  Name string `json:"name"`
  Point int64 `json:"point"`
}

func (u *User) Insert(db *sql.DB) (*User, error) {
  names := []string{
    "Winnie Lee", "Shad Halliday", "Harris Shock", "Mabelle Wunsch", "Debbie Dales",
    "Reginald Vega", "Pamela Hultgren", "Adena Franqui", "Rae Odriscoll", "Williemae Baney",
    "Callie Weingarten", "Deshawn Nilles", "Jenna Pietrowski", "Sherrell Stigall",
    "Sonia Bickel", "Julio Apel", "Ruth Nesler", "Daron Crew", "Taren Prato", "Elaina Randol",
  }
  name := names[randIntn(len(names))]
	res, err := db.Exec(`insert into users (name, point) values (?, ?)`, name, 0)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:   id,
    Name: u.Name,
	}, nil
}
func randIntn(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}
