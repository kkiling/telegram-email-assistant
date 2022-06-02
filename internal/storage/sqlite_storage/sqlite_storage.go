package sqlitestorage

import (
	"database/sql"
	"fmt"

	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/factory"
	"github.com/kiling91/telegram-email-assistant/internal/storage"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type service struct {
	db   *sql.DB
	fact factory.Factory
}

func NewSqliteStorage(fact factory.Factory) (storage.Storage, error) {
	logrus.Infof("Open db file: %s", fact.Config().App.SQLiteDbFile)
	db, err := sql.Open("sqlite3", fact.Config().App.SQLiteDbFile)
	if err != nil {
		return nil, err
	}

	// ***
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sent_messages(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		email TEXT,
		seqNum INTEGER,
		botUserUid INTEGER,
		UNIQUE (email,seqNum,botUserUid)
	  );`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx
		ON sent_messages (email,seqNum,botUserUid);`)
	if err != nil {
		return nil, err
	}
	// ***
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS email_info(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		email TEXT,
		seqNum INTEGER,
		fromAddress TEXT,
		UNIQUE (email,seqNum,fromAddress)
	  );`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx
		ON email_info (email,seqNum);`)
	if err != nil {
		return nil, err
	}
	// ***
	return &service{
		fact: fact,
		db:   db,
	}, nil
}

func (s *service) SaveMsgSentToBotUser(email string, seqNum int64, botUserUid int64) error {
	if contains, err := s.MsgWasSentToBotUser(email, seqNum, botUserUid); err != nil {
		return err
	} else if contains {
		return nil
	}
	_, err := s.db.Exec(`INSERT INTO sent_messages (email, seqNum, botUserUid) values ($1, $2, $3);`,
		email, seqNum, botUserUid)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) MsgWasSentToBotUser(email string, seqNum int64, botUserUid int64) (bool, error) {
	rows, err := s.db.Query("SELECT COUNT(*) FROM sent_messages WHERE email=$1 AND seqNum=$2 AND botUserUid=$3;",
		email, seqNum, botUserUid)
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

func (s *service) SaveMsgInfo(email string, msg *email.MessageEnvelope) error {
	from, _ := s.GetMsgFromAddress(email, msg.SeqNum)
	if from != "" {
		return nil
	}

	_, err := s.db.Exec(`INSERT INTO email_info (email, seqNum, fromAddress) values ($1, $2, $3);`,
		email, msg.SeqNum, msg.FromAddress)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetMsgFromAddress(email string, seqNum int64) (string, error) {
	sqlStatement := `SELECT fromAddress FROM email_info WHERE email=$1 AND seqNum=$2;`

	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	row := s.db.QueryRow(sqlStatement, email, seqNum)

	var fromAddress string
	switch err := row.Scan(&fromAddress); err {
	case sql.ErrNoRows:
		return "", fmt.Errorf("not found")
	case nil:
		return fromAddress, nil
	default:
		return "", fmt.Errorf("not found")
	}
}

func (s *service) ShutDown() error {
	return s.db.Close()
}
