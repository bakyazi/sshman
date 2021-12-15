package main

import (
	"errors"
	"github.com/rivo/tview"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Account entity
type Account struct {
	gorm.Model
	Name     string
	Password string
}

// Session entity
type Session struct {
	gorm.Model
	AccountID  int
	Account    Account
	Name       string
	User       string
	AuthMethod int
	Password   string
	KeyFile    string
	IPAddr     string
	Port       string
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
	}
	return &account
}

func createUser(user, password string) *Account {
	db.Create(&Account{
		Name:     user,
		Password: hashPassw(password),
	})
	return getUser(user)
}

func createSession(acct *Account, authMethod int, name, host, username, keyFile, password, port, accountPass string) {
	if authMethod == 0 {
		db.Create(&Session{
			Account:    *acct,
			Name:       name,
			IPAddr:     host,
			User:       username,
			AuthMethod: authMethod,
			Password:   encryptDES([]byte(accountPass), password),
			Port:       port,
		})
	} else {
		db.Create(&Session{
			Account:    *acct,
			Name:       name,
			IPAddr:     host,
			User:       username,
			AuthMethod: authMethod,
			KeyFile:    keyFile,
			Port:       port,
		})
	}
}

func editSession(session *Session, method int, name string, host string, username string, file string, password string, port string) {
	session.Name = name
	session.User = username
	session.IPAddr = host
	session.AuthMethod = method
	session.KeyFile = file
	session.Password = encryptDES([]byte(*userPassword), password)
	session.Port = port
	db.Updates(&session)
}

func deleteFunction(session *Session, app *tview.Application, user string) {
	db.Delete(session)
	homeView(app, user)
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
	}
	bs := hashPassw(password)
	if bs == account.Password {
		return true
	}
	return false

}
