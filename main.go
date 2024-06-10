package main

import (
	"fmt"
	"log"

	"github.com/rivo/tview"
)

type post struct {
	username string
	message  string
	time     string
}

var posts = []post{
	{username: "john", message: "hi, what's up?", time: "14:41"},
	{username: "jane", message: "not much", time: "14:43"},
}

func main() {
	app := tview.NewApplication()

	sidebar := tview.NewList().
		AddItem("> node1", "11ms", 0, nil).
		AddItem("> node2", "22ms", 0, nil)


	sidebar.SetBorder(true).SetTitle("active nodes")

	history := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	history.SetBorder(true).SetTitle("global chat")

	for _, m := range posts {
		fmt.Fprintf(history, "[yellow]%s [white]<%s> %s\n", m.time, m.username, m.message)
	}

	input := tview.NewInputField().
		SetLabel("@node1: ").
		SetFieldWidth(0).SetFieldBackgroundColor(0)

	input.SetBorder(true)

	chatLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(history, 0, 1, false).
		AddItem(input, 3, 1, true)

	root := tview.NewFlex().AddItem(chatLayout, 0, 3, true).AddItem(sidebar, 25, 1, false)

	app.SetRoot(root, true).SetFocus(input)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
