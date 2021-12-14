package main

import (
	ssh "github.com/nanobox-io/golang-ssh"
	"github.com/rivo/tview"
	"strconv"
)

func connectFunction(session *Session, app *tview.Application, user string) {
	app.Stop()
	dpass := DecryptDES([]byte(*userPassword), session.Password)

	nanPass := ssh.Auth{Passwords: []string{dpass}}
	port, _ := strconv.Atoi(session.Port)
	client, _ := ssh.NewNativeClient(session.User, session.IPAddr, "", port, &nanPass, nil)
	_ = client.Shell()

	app = tview.NewApplication()
	HomeView(app, user)
	app.Run()
}
