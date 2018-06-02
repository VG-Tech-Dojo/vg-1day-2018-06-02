package model

import (
	"database/sql"
)

// Message はメッセージの構造体です
type Message struct {
	ID   int64  `json:"id"`
	Body string `json:"body"`
	// Tutorial 1-1. ユーザー名を表示しよう
	Username string `json:"username"`
}

// MessagesAll は全てのメッセージを返します
func MessagesAll(db *sql.DB) ([]*Message, error) {

	// Tutorial 1-1. ユーザー名を表示しよう
	rows, err := db.Query(`select id, body, username from message`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []*Message
	for rows.Next() {
		m := &Message{}
		// Tutorial 1-1. ユーザー名を表示しよう
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

	// Tutorial 1-1. ユーザー名を表示しよう
	if err := db.QueryRow(`select id, body, username from message where id = ?`, id).Scan(&m.ID, &m.Body, &m.Username); err != nil {
		return nil, err
	}

	return m, nil
}

// Insert はmessageテーブルに新規データを1件追加します
func (m *Message) Insert(db *sql.DB) (*Message, error) {
	// Tutorial 1-2. ユーザー名を追加しよう
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
		// Tutorial 1-2. ユーザー名を追加しよう
		Username: m.Username,
	}, nil
}

// Mission 1-1. メッセージを編集しよう
func UpdateMessageBody(db *sql.DB, msg *Message, newBody string) (*Message, error) {
	if _, err := db.Exec(`update message set body = ? where id = ?`, newBody, msg.ID); err != nil {
		return nil, err
	}
	return &Message{
		ID:       msg.ID,
		Body:     newBody,
		Username: msg.Username,
	}, nil
}

// Mission 1-2. メッセージを削除しよう
func DeleteMessage(db *sql.DB, id int64) error {
	if _, err := db.Exec(`delete from message where id = ?`, id); err != nil {
		return err
	}
	return nil
}
