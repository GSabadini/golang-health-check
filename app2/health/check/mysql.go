package check

import (
	"database/sql"
	"log"
)

type MySQL struct {
	DNS string
}

func (m MySQL) Check() error {
	db, err := sql.Open("mysql", m.DNS)
	if err != nil {
		log.Printf("MySQL health check failed on connect: %w", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("MySQL health check failed on connect: %w", err)
		return err
	}

	_, err = db.Query(`SELECT VERSION()`)
	if err != nil {
		log.Printf("MySQL health check failed on connect: %w", err)
		return err
	}

	return nil
}
