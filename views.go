package main

import (
	"fmt"
	"github.com/rivo/tview"
)

// Action connect 1/delete 2/ other 0
type Action uint8

var action Action

func modal(p tview.Primitive) tview.Primitive {
	return tview.NewGrid().SetColumns(0, 0, 0).SetRows(0, 0, 0).AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func LoginView(app *tview.Application, user string) {
	var password string
	form := tview.NewForm().
		AddPasswordField("Password", "", 15, '*', func(text string) { password = text }).
		AddButton("Log In", func() {
			if !loginVerify(user, password) {
				popmodal := tview.NewModal().
					SetText("Password is wrong or User has been deleted in DB").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "OK" {
							LoginView(app, user)
						}
					})
				app.SetRoot(popmodal, true)
			} else {
				userPassword = &password
				HomeView(app, user)
			}
		}).
		AddButton("Cancel", func() {
			app.Stop()
		})

	title := fmt.Sprintf("Login for user: %s", user)
	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

func SignUpView(app *tview.Application, user string) {
	var password string
	var confirm string
	form := tview.NewForm().
		AddPasswordField("Password", "", 15, '*', func(text string) { password = text }).
		AddPasswordField("Confirm", "", 15, '*', func(text string) { confirm = text }).
		AddButton("Sign Up", func() {
			if password != confirm {
				popmodal := tview.NewModal().
					SetText("Password and confirmation is not match").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "OK" {
							SignUpView(app, user)
						}
					})
				app.SetRoot(popmodal, false)
			} else {
				createUser(user, password)
				LoginView(app, user)
			}
		}).
		AddButton("Cancel", func() {
			app.Stop()
		})
	title := fmt.Sprintf("Sign up for user: %s", user)
	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

func HomeView(app *tview.Application, user string) {

	var actionList *tview.List
	var sessionList *tview.List
	app.EnableMouse(false)
	actionList = tview.NewList().ShowSecondaryText(false).
		AddItem("New Connection", "", 'n', func() {
			action = 0
			NewConnection(app, user)
		}).
		AddItem("Connect", "", 'c', func() {
			action = 1
			app.SetFocus(sessionList)
		}).
		AddItem("Delete", "", 'd', func() {
			action = 2
			app.SetFocus(sessionList)
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	actionList.SetBorder(true).SetTitle("Actions")
	sessionList = tview.NewList().ShowSecondaryText(false)

	sessions := getAllSessions(getUser(user))
	for i := 0; i < len(sessions); i++ {
		sessionList.AddItem(fmt.Sprintf("[ %s ] |  %s@%s:%s",
			sessions[i].Name,
			sessions[i].User,
			sessions[i].IPAddr,
			sessions[i].Port), "", 'x',
			onSessionSelect(app, user, &sessions[i]),
		)
	}

	sessionList.SetBorder(true).SetTitle("Saved Sessions")
	flex := tview.NewFlex().
		AddItem(actionList, 0, 1, true).
		AddItem(sessionList, 0, 3, true)

	app.SetRoot(flex, true)
}

func onSessionSelect(app *tview.Application, user string, session *Session) func() {
	return func() {
		if action == 2 {
			deleteFunction(session, app, user)
		} else if action == 1 {
			connectFunction(session, app, user)
		} else {
			app.Stop()
		}
	}

}

func NewConnection(app *tview.Application, user string) {
	var name string
	var host string
	var username string
	var password string
	var port string
	form := tview.NewForm().
		AddInputField("Connection Name", "", 15, nil, func(text string) { name = text }).
		AddInputField("Host IP Address", "", 15, nil, func(text string) { host = text }).
		AddInputField("Username", "", 15, nil, func(text string) { username = text }).
		AddPasswordField("Password", "", 15, '*', func(text string) { password = text }).
		AddInputField("Port", "", 15, nil, func(text string) { port = text }).
		AddButton("Save", func() {
			if !validConnectionParams(name, host, username, password, port) {
				popmodal := tview.NewModal().
					SetText("Password and confirmation is not match").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "OK" {
							NewConnection(app, user)
						}
					})
				app.SetRoot(popmodal, false)
			}
			acct := getUser(user)
			createSession(acct, name, host, username, password, port, *userPassword)
			HomeView(app, user)
		}).
		AddButton("Cancel", func() {
			HomeView(app, user)
		})
	title := fmt.Sprintf("New Connection")
	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

// TODO implement
func validConnectionParams(name string, host string, username string, password string, port string) bool {
	return true
}
