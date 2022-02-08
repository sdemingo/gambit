package netcon

import (
	"bufio"
	"fmt"
	"net"

	dt "github.com/dylhunn/dragontoothmg"
)

const (
	SERVER = "localhost"
	PORT   = "22022"
)

var conn net.Conn
var user string

func ConnectToServer(username string) {
	conn, _ = net.Dial("tcp", SERVER+":"+PORT)
	user = username
}

// Join to a created match or create a new match in the server
// if match has no text
func JoinMatch(match string) string {
	cmd := fmt.Sprintf("/join:%s\n", match)
	conn.Write([]byte(cmd))

	resp := bufio.NewReader(conn)
	respBody, _ := resp.ReadString('\n')
	return respBody
}

// Send a move to the server
func SendMove(move dt.Move) {
	conn.Write([]byte(move.String() + "\n"))
}

/*
func ReceiveRoutine(func() tea.Msg) {
	for {

	}
}
*/
