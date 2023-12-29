package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(username, password string) error
	Login(username, password string) error
	doesExist(username string) bool
	GetCount(username string) int
	Increment(username string) error
}

type service struct {
	db *sql.DB
}

var (
	dburl = os.Getenv("DB_URL")
)

func New() Service {
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL,
		count INTEGER DEFAULT 0,
		created_at TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatal(err)
	}

	s := &service{db: db}
	return s
}

func (s *service) CreateUser(username, password string) error {
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	doesExist := s.doesExist(username)
	if doesExist {
		return errors.New("user already exists")
	}

	_, err = db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?);`,
		username, hashedPassword, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *service) doesExist(username string) bool {
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		return false
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE username=?;`, username).Scan(&count)
	if err != nil {
		return false
	}

	err = db.Close()
	if err != nil {
		return false
	}

	return count > 0
}

func (s *service) Login(username, password string) error {
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		return err
	}

	var hashedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetCount(username string) int {
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		return 0
	}

	var count int
	err = db.QueryRow("SELECT count FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return 0
	}

	err = db.Close()
	if err != nil {
		return 0
	}

	return count
}

func (s *service) Increment(username string) error {
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE users SET count = count + 1 WHERE username = ?;`, username)
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}
