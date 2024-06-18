package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ChatUI struct {
	cr        *ChatRoom
	app       *tview.Application
	peersList *tview.List

	msgW    *tview.TextView
	inputCh chan string
	doneCh  chan struct{}
}


func NewChatUI(cr *ChatRoom) *ChatUI {

	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := history.GetScrollOffset()

		switch event.Key() {
		case tcell.KeyUp:
			history.ScrollTo(-1+row, 0)
		case tcell.KeyDown:
			history.ScrollTo(1+row, 0)
		}
		return event
	})



	peersList := tview.NewList()
	peersList.SetBorder(true).SetTitle("active users")

	history = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	history.SetBorder(true).SetTitle("global chat")


	inputCh := make(chan string, 32)
	input = tview.NewInputField().
		SetLabel("@" + cr.nick + ":").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorDefault)

	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}
		line := input.GetText()
		if len(line) == 0 {
			return
		}

		if line == "/quit" {
			app.Stop()
			return
		}

		inputCh <- line
		input.SetText("")
		history.ScrollToEnd()
	})

	input.SetBorder(true)

	chatLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(history, 0, 1, false).
		AddItem(input, 3, 1, true)

	root := tview.NewFlex().
		AddItem(chatLayout, 0, 3, true).
		AddItem(peersList, 25, 1, false)

	app.SetRoot(root, true)

	return &ChatUI{
		cr:        cr,
		app:       app,
		peersList: peersList,
		msgW:      history,
		inputCh:   inputCh,
		doneCh:    make(chan struct{}, 1),
	}
}



func (ui *ChatUI) Run() error {
	go ui.handleEvents()
	defer ui.end()

	return ui.app.Run()
}


func (ui *ChatUI) end() {
	ui.doneCh <- struct{}{}
}

func (ui *ChatUI) refreshPeers() {
	peers := ui.cr.ListPeers()

	// clear is thread-safe
	ui.peersList.Clear()

	for _, p := range peers {
		ui.peersList.AddItem(shortID(p),"", 0, nil)
	}

	ui.app.Draw()
}

func (ui *ChatUI) displayChatMessage(cm *ChatMessage) {
	prompt := withColor("green", fmt.Sprintf("<%s>:", cm.SenderNick))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, cm.Message)
}

func (ui *ChatUI) displaySelfMessage(msg string) {
	prompt := withColor("yellow", fmt.Sprintf("<%s>:", ui.cr.nick))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
}

func (ui *ChatUI) handleEvents() {
	peerRefreshTicker := time.NewTicker(time.Second)
	defer peerRefreshTicker.Stop()

	for {
		select {
		case input := <-ui.inputCh:
			// when the user types in a line, publish it to the chat room and print to the message window
			err := ui.cr.Publish(input)
			if err != nil {
				printErr("publish error: %s", err)
			}
			ui.displaySelfMessage(input)

		case m := <-ui.cr.Messages:
			// when we receive a message from the chat room, print it to the message window
			ui.displayChatMessage(m)

		case <-peerRefreshTicker.C:
			// refresh the list of peers in the chat room periodically
			ui.refreshPeers()

		case <-ui.cr.ctx.Done():
			return

		case <-ui.doneCh:
			return
		}
	}
}

func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}

