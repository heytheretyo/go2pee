package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type post struct {
	username string
	message  string
	time     string
}

var (
	posts   []post
	app     *tview.Application
	history *tview.TextView
	input   *tview.InputField
)


func main() {
	initializeApp()
	initializeInterface()
	initializeInputHandling()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func initializeApp() {
	app = tview.NewApplication()

	posts = []post{
		{username: "jeb", message: "hi, what's up?", time: "14:41"},
		{username: "monty", message: "not much", time: "14:43"},
	}
}

func initializeInterface() {
	sidebar := tview.NewList().
		AddItem("> jeb", "11ms", 0, nil).
		AddItem("> monty", "22ms", 0, nil)
	sidebar.SetBorder(true).SetTitle("active users")

	history = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	history.SetBorder(true).SetTitle("global chat")

	updateHistory()

	input = tview.NewInputField().
		SetLabel("@node1: ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetDoneFunc(handleInput)
	input.SetBorder(true)

	chatLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(history, 0, 1, false).
		AddItem(input, 3, 1, true)

	root := tview.NewFlex().
		AddItem(chatLayout, 0, 3, true).
		AddItem(sidebar, 25, 1, false)

	app.SetRoot(root, true).SetFocus(input)
}

func initializeInputHandling() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			scrollHistory(-1)
		case tcell.KeyDown:
			scrollHistory(1)
		case tcell.KeyPgUp:
			scrollHistory(-10)
		case tcell.KeyPgDn:
			scrollHistory(10)
		}
		return event
	})
}

func updateHistory() {
	history.Clear()
	for _, m := range posts {
		fmt.Fprintf(history, "[yellow]%s [white]<%s> %s\n", m.time, m.username, m.message)
	}
	history.ScrollToEnd()
}

func handleInput(key tcell.Key) {
	if key == tcell.KeyEnter {
		text := input.GetText()
		if text != "" {
			newMessage := post{
				username: "alice",
				message:  text,
				time:     time.Now().Format("15:04"),
			}
			posts = append(posts, newMessage)
			updateHistory()
			input.SetText("")
		}
	}
}

func scrollHistory(offset int) {
	row, _ := history.GetScrollOffset()
	history.ScrollTo(row+offset, 0)
}
