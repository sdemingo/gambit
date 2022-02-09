package netcon

import (
	"bufio"
	"encoding/json"
	"net"
	"strings"

	dt "github.com/dylhunn/dragontoothmg"
)

const (
	SERVER = "localhost"
	PORT   = "22022"
)

var conn net.Conn
var user string

func ConnectToServer(username string) (err error) {
	conn, err = net.Dial("tcp", SERVER+":"+PORT)
	user = username
	return err
}

// Join to a created match or create a new match in the server
// if match has no text
func JoinMatch(match string) string {
	msg := NewMsg(JOIN, user)
	msg.Args = match
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
func SendMove(move dt.Move) {
	conn.Write([]byte(move.String() + "\n"))
}

func ReceivePlayerName() string {
	resp := bufio.NewReader(conn)
	message, _ := resp.ReadString('\n')
	return strings.Trim(message, "\n")
}

/*
func ReceiveRoutine(func() tea.Msg) {
	for {

	}
}
*/
