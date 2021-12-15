package main

import (
	"github.com/rivo/tview"
	"os/exec"
	"os/user"
)

var (
	currentUser     *user.User = nil
	directoryFolder            = "/.sshman"
	appHome         *string    = nil
	userPassword    *string    = nil
)

func checkAppHome() {
	if currentUser != nil {
		homePath := currentUser.HomeDir
		_, err := exec.Command("mkdir", "-p", homePath+directoryFolder).Output()
		if err != nil {
			panic(err)
		}
		path := homePath + directoryFolder
		appHome = &path
	}
}

func main() {
	currentUser, _ = user.Current()
	checkAppHome()
	initDb()
	acct := getUser(currentUser.Username)
	app := tview.NewApplication()
	app.EnableMouse(true)
	if acct != nil {
		loginView(app, currentUser.Username)
	} else {
		signUpView(app, currentUser.Username)
	}
	app.Run()
}
