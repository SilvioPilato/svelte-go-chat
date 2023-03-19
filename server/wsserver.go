package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WSServer struct {
	conns map[*websocket.Conn]WSUser
	logged map[string]bool
	ctx context.Context
}

type WSUser struct{
	connected bool
	username string
}

type WSMessage struct{
	Topic string `json:"topic"`
	Message string `json:"message"`
	Username string `json:"username,omitempty"`
}

func (s WSServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		fmt.Println("Error while accepting connection")
	}
	log.Printf("%s has entered the room", r.RemoteAddr)
	s.conns[c] = WSUser{connected: true, username: ""}
	s.ctx = r.Context()
	defer c.Close(websocket.StatusInternalError, "Server closing")
	channel := make(chan error)
	go s.readLoop(c, channel)
	<-channel
}

func (s WSServer) readLoop(conn *websocket.Conn, channel chan <- error) {
	for {
		data, err := s.read(conn)
		
		mustClose := websocket.CloseStatus(err) > -1
		if mustClose {
			log.Printf("Websocket closed with code: %v", websocket.CloseStatus(err))
			channel <- err
			close(channel)
			break
		}
		if (err != nil ) {
			s.sendError(conn, err)
			continue
		}

		err = s.handleMsg(conn, *data)
		if (err != nil) {
			s.sendError(conn, err)
			continue
		}
	}
}

func (s WSServer) handleMsg(c *websocket.Conn, msg WSMessage) (error) {
	var err error = nil
	user, found := s.conns[c]
	if (!found) {
		return errors.New("user not found")
	}
	_, loggedIn := s.logged[user.username]

	switch msg.Topic {
	case "login":
		err = s.login(c, msg.Message, loggedIn)
	case "chat":
		if (loggedIn) {
			msg.Username = user.username
			err = s.broadCast(c, msg)
		} else {
			err =  errors.New("this user is not logged in")
		}
	default:
		err = errors.New("unknown topic")
	}
	return err
}

func (s WSServer) broadCast(c *websocket.Conn, msg WSMessage) (error) {
	var err error = nil
	for conn, user := range s.conns {
		_, loggedIn := s.logged[user.username]
		if !loggedIn || !user.connected {
			continue
		}
		err = wsjson.Write(s.ctx, conn, msg)
	}
	return err
}

func (s *WSServer) login(conn *websocket.Conn, username string, alreadyLogged bool) (error) {
	user, connected := s.conns[conn]
	
	if (!connected) {
		return errors.New("connection not found")
	}

	_, found := s.logged[username]
	if (found) {
		return errors.New("username has already been taken")
	}
	oldUsername := user.username
	user.username = username
	s.logged[username] = true
	s.conns[conn] = user

	if (alreadyLogged) {
		delete(s.logged, oldUsername)
	}
	s.sendLoginSuccess(conn, username)
	return nil
}

func (s WSServer) read(conn *websocket.Conn) (*WSMessage, error) {
	// v := make(map[string]interface{})
	v := WSMessage{}
	err := wsjson.Read(s.ctx, conn, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (s WSServer) sendError(conn *websocket.Conn, err error) {
	msg := WSMessage{Topic: "error", Message: err.Error()}
	wsjson.Write(s.ctx, conn, &msg)
}
func (s WSServer) sendLoginSuccess(conn *websocket.Conn, username string) {
	msg := WSMessage{Topic: "login", Message: "logged in successfully", Username: username}
	wsjson.Write(s.ctx, conn, &msg)
}