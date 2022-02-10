package netcon

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"

	dt "github.com/dylhunn/dragontoothmg"
)

const (
	SERVER = "localhost"
	PORT   = "22022"
)

var conn net.Conn
var user string

var logFile *os.File

func InitLog(user string) {
	var logFileName = "/tmp/gambit-" + user + ".log"
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
}

func ConnectToServer(username string) (err error) {
	conn, err = net.Dial("tcp", SERVER+":"+PORT)
	user = username
	return err
}

// Join to a created match and receive the white player name
func JoinMatch(game string) string {
	msg := NewMsg(JOIN, user)
	msg.Args = game
	b, _ := json.Marshal(msg)
	conn.Write(b)

	resp := bufio.NewReader(conn)
	message, _ := resp.ReadString('\n')
	return strings.Trim(message, "\n")
}

// Create a new match
func CreateMatch() string {
	msg := NewMsg(CREATE, user)
	b, _ := json.Marshal(msg)
	conn.Write(b)

	resp := bufio.NewReader(conn)
	message, _ := resp.ReadString('\n')
	return strings.Trim(message, "\n")
}

// Send a move to the server
func SendMove(move dt.Move, game string) {
	msg := NewMsg(MOVE, user)
	msg.Args = game + ":" + move.String()
	b, _ := json.Marshal(msg)
	conn.Write(b)
}

// Receive the player black player name
func ReceivePlayerName() string {
	resp := bufio.NewReader(conn)
	message, _ := resp.ReadString('\n')
	return strings.Trim(message, "\n")
}

// Receive a message from the server
func ReciveMsg() (*Msg, error) {
	msg, err := UnpackMsg(conn)
	if err != nil {
		log.Println("Error: receive msg")
		return nil, err
	}
	return msg, nil
}
