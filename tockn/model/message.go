package model

import (
	"database/sql"
)

// Message はメッセージの構造体です
type Message struct {
	ID       int64  `json:"id"`
	Body     string `json:"body"`
	Username string `json:"username"`
	RoomID   int64  `json:"room_id"`
	// Tutorial 1-2. ユーザー名を表示しよう
}

// MessagesAll は全てのメッセージを返します
func MessagesAll(db *sql.DB) ([]*Message, error) {

	// Tutorial 1-2. ユーザー名を表示しよう
	rows, err := db.Query(`select id, body, username, room_id from message`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []*Message
	for rows.Next() {
		m := &Message{}
		// Tutorial 1-2. ユーザー名を表示しよう
		if err := rows.Scan(&m.ID, &m.Body, &m.Username, &m.RoomID); err != nil {
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
func MessageByID(db *sql.DB, roomID, id string) (*Message, error) {
	m := &Message{}

	// Tutorial 1-2. ユーザー名を表示しよう
	if err := db.QueryRow(`select id, body, username, room_id from message where id = ? and room_id = ?`, id, roomID).Scan(&m.ID, &m.Body, &m.Username, &m.RoomID); err != nil {
		return nil, err
	}

	return m, nil
}

// Insert はmessageテーブルに新規データを1件追加します
func (m *Message) Insert(db *sql.DB) (*Message, error) {
	// Tutorial 1-2. ユーザー名を追加しよう
	res, err := db.Exec(`insert into message (body, username, room_id) values (?, ?, ?)`, m.Body, m.Username)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:       id,
		Body:     m.Body,
		Username: m.Username,
		RoomID:   m.RoomID,
		// Tutorial 1-2. ユーザー名を追加しよう
	}, nil
}

// Mission 1-1. メッセージを編集しよう
func (m *Message) Update(db *sql.DB) (*Message, error) {
	_, err := db.Exec(`update message set body = ? where id = ? and room_id = ?`, m.Body, m.ID, m.RoomID)
	if err != nil {
		return nil, err
	}
	msg := &Message{}
	if err := db.QueryRow(`select id, body, username, room_id from message where id = ? and room_id = ?`, m.ID).Scan(&msg.ID, &msg.Body, &msg.Username); err != nil {
		return nil, err
	}
	return msg, nil
}

// ...

// Mission 1-2. メッセージを削除しよう
func Delete(db *sql.DB, roomID, id string) error {
	_, err := db.Exec(`delete from message where id = ? and room_id = ?`, id, roomID)
	if err != nil {
		return err
	}
	return nil
}

// ...
