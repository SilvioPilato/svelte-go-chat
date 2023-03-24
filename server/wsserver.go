package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WSServer struct {
	conns map[*websocket.Conn]WSUser
	ctx context.Context
}

type WSUser struct{
	tag int
	logged bool
	connected bool
	username string
	FullUsername string `json:"fullUsername"`
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
	tag:=len(s.conns) + 1
	s.conns[c] = WSUser{connected: true, username: "", tag: tag, FullUsername: ""}
	log.Printf("%s has entered the room and is tagged as %d", r.RemoteAddr, tag)
	s.ctx = r.Context()
	defer c.Close(websocket.StatusInternalError, "Server closing")
	channel := make(chan error)
	go s.readLoop(c, channel)
	<-channel
}

func (s WSServer) ServeUsers(w http.ResponseWriter, r *http.Request) {
	users := filterLoggedUsers(s.conns, true)
	out, err := json.Marshal(users)
	if (err != nil) {
		log.Println(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(out)
}

func filterLoggedUsers(users map[*websocket.Conn]WSUser, logged bool) []WSUser {
	loggedUsers := []WSUser{}
	for _, user := range users {
		if (user.logged == logged) {
			loggedUsers = append(loggedUsers, user)
		}
	}
	return loggedUsers
}

func (s WSServer) readLoop(conn *websocket.Conn, channel chan <- error) {
	for {
		data, err := s.read(conn)
		
		mustClose := websocket.CloseStatus(err) > -1
		if mustClose {
			log.Printf("Websocket closed with code: %v", websocket.CloseStatus(err))
			s.logout(conn)
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

	switch msg.Topic {
	case "login":
		err = s.login(c, msg.Message, user.logged)
	case "chat":
		if (user.logged) {
			msg.Username = user.FullUsername;
			err = s.broadCast(c, msg)
		} else {
			err = errors.New("this user is not logged in")
		}
	default:
		err = errors.New("unknown topic")
	}
	return err
}

func (s WSServer) broadCast(c *websocket.Conn, msg WSMessage) (error) {
	var err error = nil
	for conn, user := range s.conns {
		loggedIn := user.logged
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

	user.username = username
	user.FullUsername = username + "#" + strconv.Itoa(user.tag)
	user.logged = true

	s.conns[conn] = user
	s.sendLoginSuccess(conn, username)
	return nil
}

func (s *WSServer) logout(conn *websocket.Conn) (error) {
	user, connected := s.conns[conn]

	if (!connected) {
		return errors.New("connection not found")
	}

	user.logged = false

	s.conns[conn] = user
	return nil
}

func (s WSServer) read(conn *websocket.Conn) (*WSMessage, error) {
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
	msg := WSMessage{Topic: "login_success", Message: "logged in successfully", Username: username}
	wsjson.Write(s.ctx, conn, &msg)
}