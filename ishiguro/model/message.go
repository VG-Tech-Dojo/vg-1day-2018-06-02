package model

import (
	"database/sql"
)

// Message はメッセージの構造体です
type Message struct {
	ID   int64  `json:"id"`
	Body string `json:"body"`
	Username string `json:"username"`
	// Tutorial 1-2. ユーザー名を表示しよう
}

// MessagesAll は全てのメッセージを返します
func MessagesAll(db *sql.DB) ([]*Message, error) {

	rows, err := db.Query(`select id, body, username from message`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []*Message
	for rows.Next() {
		m := &Message{}
		if err := rows.Scan(&m.ID, &m.Body, &m.Username); err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ms, nil
}

// MessageByID は指定されたIDのメッセージを1つ返します
func MessageByID(db *sql.DB, id string) (*Message, error) {
	m := &Message{}

	if err := db.QueryRow(`select id, body from message where id = ?`, id).Scan(&m.ID, &m.Body); err != nil {
		return nil, err
	}

	return m, nil
}

// Insert はmessageテーブルに新規データを1件追加します
func (m *Message) Insert(db *sql.DB) (*Message, error) {
	
	if m.Username == "" {
		m.Username = "anonymous"
	}

	res, err := db.Exec(`insert into message (body, username) values (?, ?)`, m.Body, m.Username)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:   id,
		Body: m.Body,
		Username: m.Username,
	}, nil
}


func (m *Message) UpdateByID(db *sql.DB, id string) (*Message, error) {

	_, err := db.Exec(`update message set body = ? where id = ?`, m.Body, id)
	if err != nil {
		return nil, err
	}

	m, err = MessageByID(db, id)
	if err != nil {
		return nil, err
	}

	return m, nil
}


func (m *Message) DeleteByID(db *sql.DB, id string) (error) {

	_, err := db.Exec(`delete from message where id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}