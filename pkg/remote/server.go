package remote

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/google/uuid"
)

type Tracker struct {
	*sync.Mutex

	// All ongoing Sessions.
	sessions map[string]Session

	// User connections to the Tracker.
	connections map[string]*net.Conn
}

func (t *Tracker) StartTracker(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		bytes, _ := uuid.New().MarshalBinary()
		t.Lock()
		t.connections[string(bytes)] = &conn
		t.Unlock()

		// Spawn a handler for messages.
		decoder, encoder := json.NewDecoder(conn), json.NewEncoder(conn)

		// First message must be an announce.
		msg := Message{}
		if err := decoder.Decode(&msg); err != nil || msg.Type != Announce {
			encoder.Encode(Message{
				Type: Error,
				Info: "First message to the server must be an announce.",
			})
			return err
		}

		// Following this, handle messages accordingly.
		switch msg.Type {
		case Announce:
			// Assign a UID.
			uid := uuid.New()
			bytes, _ := uid.MarshalBinary()
			conn.Write(bytes)

		case JobUpdate:
			// Read the update and respond accordingly.

		default:
			return nil
		}
	}
}
