package main

import (
	ssh "github.com/nanobox-io/golang-ssh"
	"github.com/rivo/tview"
	"strconv"
)

func connectFunction(session *Session, app *tview.Application, user string) {
	app.Stop()
	var nanPass *ssh.Auth

	if session.AuthMethod == 0 {
		decryptedPassword := decryptDES([]byte(*userPassword), session.Password)
		nanPass = &ssh.Auth{Passwords: []string{decryptedPassword}}
	} else {
		nanPass = &ssh.Auth{Keys: []string{session.KeyFile}}
	}

	port, _ := strconv.Atoi(session.Port)
	client, _ := ssh.NewNativeClient(session.User, session.IPAddr, "", port, nanPass, nil)
	_ = client.Shell()

	app = tview.NewApplication()
	homeView(app, user)
	app.Run()
}
