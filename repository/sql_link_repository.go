package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"shorters/domain"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type sqlLinkRepository struct {
	client *sql.DB
}

func newSQLClient() *sql.DB {
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	db, err := sql.Open("mysql", username+":"+password+"@/shorters")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func NewSQLLinkRepository() LinkRepository {
	return &sqlLinkRepository{client: newSQLClient()}
}

func (r *sqlLinkRepository) Find(key string) (*domain.Link, error) {
	row := r.client.QueryRow("select `Key`, `URL` from `Links` where `Key` = ?", key)
	var link domain.Link
	if err := row.Scan(&link.Key, &link.URL); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &link, nil
}

func (r *sqlLinkRepository) FindByUser(email string) ([]*domain.Link, error) {
	res, err := r.client.Query("select `Key`, `URL` from `Links` where `Key` = ?", email)
	if err != nil {
		return nil, err
	}
	var links []*domain.Link
	for res.Next() {
		var link domain.Link
		if err := res.Scan(&link.Key, &link.URL); err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, nil
}

func (r *sqlLinkRepository) Store(link *domain.Link) error {
	createdTime := time.Unix(link.CreatedTime, 0).Format("2006-01-02 15:04:05")
	expiredTime := time.Unix(link.ExpiredTime, 0).Format("2006-01-02 15:04:05")
	if link.Creator == "" {
		if _, err := r.client.Query(
			"insert into `Links`(`Key`, `URL`, `Visits`, `Creator`, `CreatedTime`, `ExpiredTime`) value (?, ?, ?, ?, ?, ?)",
			link.Key, link.URL, link.Visits, nil, createdTime, expiredTime); err != nil {
			return err
		}
	} else {
		if _, err := r.client.Query(
			"insert into `Links`(`Key`, `URL`, `Visits`, `Creator`, `CreatedTime`, `ExpiredTime`) value (?, ?, ?, ?, ?, ?)",
			link.Key, link.URL, link.Visits, link.Creator, createdTime, expiredTime); err != nil {
			return err
		}
	}
	return nil
}

func (r *sqlLinkRepository) AddVisits(key string) error {
	if _, err := r.client.Query("update `Links` set `Visits` = `Visits` + 1 where `Key` = ?", key); err != nil {
		return err
	}
	return nil
}
