package main

import (
	"errors"
	"github.com/rivo/tview"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name     string
	Password string
}

type Session struct {
	gorm.Model
	AccountID int
	Account   Account
	Name      string
	User      string
	Password  string
	IPAddr    string
	Port      string
}

var (
	db *gorm.DB
)

func initDb() {
	var err error
	dbPath := *appHome + "/sshman.db"
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Account{}, &Session{})
	if err != nil {
		panic(err)
	}
}

func getUser(user string) *Account {
	var account Account
	result := db.First(&account, "name = ?", user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &account
	}
}

func createUser(user, password string) *Account {
	db.Create(&Account{
		Name:     user,
		Password: hashPassw(password),
	})
	return getUser(user)
}

func createSession(acct *Account, name, host, username, password, port, accountPass string) {
	db.Create(&Session{
		Account:  *acct,
		Name:     name,
		IPAddr:   host,
		User:     username,
		Password: EncryptDES([]byte(accountPass), password),
		Port:     port,
	})
}

func deleteFunction(session *Session, app *tview.Application, user string) {
	db.Delete(session)
	HomeView(app, user)
}

func getAllSessions(acct *Account) []Session {
	var sessions []Session
	db.Where("account_id = ?", acct.ID).Find(&sessions)
	return sessions
}

func loginVerify(user, password string) bool {
	var account Account
	result := db.First(&account, "name = ?", user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	} else {
		bs := hashPassw(password)
		if bs == account.Password {
			return true
		} else {
			return false
		}
	}
}
