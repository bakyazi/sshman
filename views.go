package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Action connect 1/delete 2/ edit3/ other 0
type Action uint8

var action Action

func modal(p tview.Primitive) tview.Primitive {
	return tview.NewGrid().SetColumns(0, 0, 0).SetRows(0, 0, 0).AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func loginView(app *tview.Application, user string) {
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
							loginView(app, user)
						}
					})
				app.SetRoot(popmodal, true)
			} else {
				userPassword = &password
				homeView(app, user)
			}
		}).
		AddButton("Cancel", func() {
			app.Stop()
		})

	title := fmt.Sprintf("Login for user: %s", user)
	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

func signUpView(app *tview.Application, user string) {
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
							signUpView(app, user)
						}
					})
				app.SetRoot(popmodal, false)
			} else {
				createUser(user, password)
				loginView(app, user)
			}
		}).
		AddButton("Cancel", func() {
			app.Stop()
		})
	title := fmt.Sprintf("Sign up for user: %s", user)
	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

func homeView(app *tview.Application, user string) {

	var actionList *tview.List
	var sessionList *tview.Table
	app.EnableMouse(false)
	actionList = tview.NewList().ShowSecondaryText(false).
		AddItem("New Connection", "", 'n', func() {
			action = 0
			newConnection(app, user)
		}).
		AddItem("Connect", "", 'c', func() {
			action = 1
			app.SetFocus(sessionList)
		}).
		AddItem("Delete", "", 'd', func() {
			action = 2
			app.SetFocus(sessionList)
		}).
		AddItem("Edit", "", 'e', func() {
			action = 3
			app.SetFocus(sessionList)
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	actionList.SetBorder(true).SetTitle("Actions")

	sessionList = sessionTable(app, user)

	sessionList.SetBorder(true).SetTitle("Saved Sessions")
	flex := tview.NewFlex().
		AddItem(actionList, 0, 1, true).
		AddItem(sessionList, 0, 3, true)

	flex.SetFullScreen(true)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight {
			action = 0
			app.SetFocus(sessionList)
		} else if event.Key() == tcell.KeyLeft {
			action = 0
			app.SetFocus(actionList)
		}
		return event
	})

	app.SetRoot(flex, true)
}

func sessionTable(app *tview.Application, user string) *tview.Table {
	style := tcell.Style{}
	style.Foreground(tcell.ColorBlack)
	style.Background(tcell.ColorWhite)
	table := tview.NewTable().SetFixed(1, 1).SetSelectable(true, false).
		SetSelectedStyle(style)
	sessions := getAllSessions(getUser(user))
	headers := []string{"Name", "Host", "Username", "AuthMethod", "Port", "Public Key"}
	for i := 0; i < len(headers); i++ {
		table.SetCell(0, i, &tview.TableCell{
			Text:          headers[i],
			Color:         tcell.ColorYellow,
			Align:         tview.AlignLeft,
			NotSelectable: true,
		})
	}
	for i := 0; i < len(sessions); i++ {
		sessionCells := createSessionRow(&sessions[i])
		for j := 0; j < len(headers); j++ {
			table.SetCell(i+1, j, (&tview.TableCell{
				Text:          sessionCells[j],
				Color:         tcell.ColorWhite,
				Align:         tview.AlignLeft,
				NotSelectable: false,
			}).SetExpansion(1))
		}
	}
	table.SetSelectedFunc(func(row, column int) {
		if action == 1 {
			connectFunction(&sessions[row-1], app, user)
		} else if action == 2 {
			deleteFunction(&sessions[row-1], app, user)
		} else if action == 3 {
			editConnection(app, user, &sessions[row-1])
		}
	})
	return table

}

func createSessionRow(session *Session) []string {
	var authMethod string
	if session.AuthMethod == 0 {
		authMethod = "Password"
	} else {
		authMethod = "PubKey"
	}
	return []string{
		session.Name,
		session.IPAddr,
		session.User,
		authMethod,
		session.Port,
		session.KeyFile,
	}
}

func newConnection(app *tview.Application, user string) {

	var fieldWidth int = 15

	var name string
	var host string
	var username string
	var authMethod int
	var keyFile string
	var password string
	var port string
	form := tview.NewForm().
		AddInputField("Connection Name", "", fieldWidth, nil, func(text string) { name = text }).
		AddInputField("Host IP Address", "", fieldWidth, nil, func(text string) { host = text }).
		AddInputField("Username", "", fieldWidth, nil, func(text string) { username = text }).
		AddDropDown("Auth Method", []string{"Password", "Public Key"}, 0, func(option string, optionIndex int) {
			authMethod = optionIndex
		}).
		AddInputField("Key File", "", fieldWidth, nil, func(text string) {
			keyFile = text
		}).
		AddPasswordField("Password", "", fieldWidth, '*', func(text string) { password = text }).
		AddInputField("Port", "", fieldWidth, nil, func(text string) { port = text }).
		AddButton("Save", func() {
			if err := validConnectionParams(authMethod, name, host, username, keyFile, password, port); err != nil {
				popmodal := tview.NewModal().
					SetText(err.Error()).
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "OK" {
							newConnection(app, user)
						}
					})
				app.SetRoot(popmodal, false)
			} else {
				acct := getUser(user)
				createSession(acct, authMethod, name, host, username, keyFile, password, port, *userPassword)
				homeView(app, user)

			}
		}).
		AddButton("Cancel", func() {
			homeView(app, user)
		})
	title := fmt.Sprintf("New Connection")
	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

func editConnection(app *tview.Application, user string, session *Session) {

	var fieldWidth int = 15

	var name string = session.Name
	var host string = session.IPAddr
	var username string = session.User
	var authMethod int = session.AuthMethod
	var keyFile string = session.KeyFile
	var password string
	if authMethod == 0 {
		password = decryptDES([]byte(*userPassword), session.Password)
	}
	var port string = session.Port
	form := tview.NewForm().
		AddInputField("Connection Name", session.Name, fieldWidth, nil, func(text string) { name = text }).
		AddInputField("Host IP Address", session.IPAddr, fieldWidth, nil, func(text string) { host = text }).
		AddInputField("Username", session.User, fieldWidth, nil, func(text string) { username = text }).
		AddDropDown("Auth Method", []string{"Password", "Public Key"}, session.AuthMethod, func(option string, optionIndex int) {
			authMethod = optionIndex
		}).
		AddInputField("Key File", session.KeyFile, fieldWidth, nil, func(text string) {
			keyFile = text
		}).
		AddPasswordField("Password", password, fieldWidth, '*', func(text string) { password = text }).
		AddInputField("Port", session.Port, fieldWidth, nil, func(text string) { port = text }).
		AddButton("Save", func() {
			if err := validConnectionParams(authMethod, name, host, username, keyFile, password, port); err != nil {
				popmodal := tview.NewModal().
					SetText(err.Error()).
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "OK" {
							newConnection(app, user)
						}
					})
				app.SetRoot(popmodal, false)
			} else {
				editSession(session, authMethod, name, host, username, keyFile, password, port)
				homeView(app, user)
			}
		}).
		AddButton("Cancel", func() {
			homeView(app, user)
		})
	title := fmt.Sprintf("Edit Connection")
	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}
