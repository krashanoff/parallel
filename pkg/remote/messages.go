package remote

import "net"

type MessageType int

const (
	Announce  MessageType = iota
	JobUpdate MessageType = iota
	Error     MessageType = iota
)

type Message struct {
	Type     MessageType `json:"type"`
	Uid      string      `json:"uid"`
	Program  string      `json:"program"`
	Args     []string    `json:"args"`
	NeedFile bool        `json:"needFile"`
	Info     string      `json:"info"`
}

type Job struct {
	Program string   `json:"program"`
	Args    []string `json:"args"`
}

type Session struct {
	Users map[string][]*net.Conn
	Jobs  map[string]Job
}
