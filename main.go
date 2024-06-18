package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
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

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const DiscoveryServiceTag = "global-chat"


func main() {
	nickFlag := flag.String("nick", "", "nickname to use in chat. will be generated if empty")
	roomFlag := flag.String("room", "skibidi", "name of chat room to join")
	flag.Parse()

	ctx := context.Background()
	// create a new libp2p Host that listens on a random TCP port
	h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}

	if err := setupDiscovery(h); err != nil {
		panic(err)
	}

	nick := *nickFlag
	if len(nick) == 0 {
		nick = defaultNick(h.ID())
	}

	room := *roomFlag

	cr, err := JoinChatRoom(ctx, ps, h.ID(), nick, room)
	if err != nil {
		panic(err)
	}

	// draw the UI
	ui := NewChatUI(cr)
	if err = ui.Run(); err != nil {
		printErr("error running text UI: %s", err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}


func shortID(p peer.ID) string {
	pretty := p.String()
	return pretty[len(pretty)-8:]
}

func printErr(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
}


func defaultNick(p peer.ID) string {
	return fmt.Sprintf("%s-%s", os.Getenv("USER"), shortID(p))
}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID)
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID, err)
	}
}


func setupDiscovery(h host.Host) error {
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}