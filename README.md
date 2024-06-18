# go2pee: decentralized irs using pubsub

## Abstract

Decentralized networks is a hot topic in the current era of digitalisation, especially where privacy and security meets its highest requirements.

The program I made named "go2pee", is a terminal-based peer to peer messaging network inspired by the infamous IRC Weechat (internet relay chat), which was designed as a client-server model. This project was made so that I could blend the features of IRC and Slack to make a p2p workspace.

<br>

<img src='assets/weechat_img.png' width="400"></img>

## Overview

The program was implemented using Pub-Sub pattern, which is a messaging system where a node in a network is able to share messages asynchronously between senders and receivers via a subscribe system. Essentially, when you send a message to a topic, every that is subscribed to the topic will receive your message.

The codebase is written in Golang.

## Features:

- [x] Global chat
- [x] User status
- [ ] Store old messages locally
- [ ] Switch rooms
- [ ] Local bot
- [ ] Reminders and todo list

## Installation

Click on the releases tab to download the latest version of the application. It's already compiled as an executable.

## Usage

Upon running the application, you will be brought into the global chat. This is where you meet people in your network. Make sure you are on the same room to text each other.

```
go2pee.exe --name="john" --room="doe"
```
