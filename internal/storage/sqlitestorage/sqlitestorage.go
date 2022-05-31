package sqlitestorage

import (
	"database/sql"

	"github.com/kiling91/telegram-email-assistant/internal/factory"
	"github.com/kiling91/telegram-email-assistant/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type service struct {
	db   *sql.DB
	fact factory.Factory
}

func NewSqliteStorage(fact factory.Factory) (storage.Storage, error) {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS emails(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		email TEXT,
		msgId INTEGER,
		UNIQUE (email,msgId)
	  );`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_email_msgid
		ON emails (email,msgId);`)
	if err != nil {
		return nil, err
	}

	return &service{
		fact: fact,
		db:   db,
	}, nil
}

func (s *service) SaveMsgId(email string, uid int64) error {
	if contains, err := s.MsgIdContains(email, uid); err != nil {
		return err
	} else if contains {
		return nil
	}
	_, err := s.db.Exec(`INSERT INTO emails (email, msgId) values ($1, $2)`,
		email, uid)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) MsgIdContains(email string, uid int64) (bool, error) {
	rows, err := s.db.Query("SELECT COUNT(*) FROM emails WHERE email=$1 AND msgId=$2",
		email, uid)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
	}
	return count > 0, nil
}

func (s *service) ShutDown() error {
	return s.db.Close()
}
