package socket

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

type SocketIO struct {
}

func NewSocketIO() *SocketIO {
	return &SocketIO{}
}

func (sio *SocketIO) StartSocketIOServer() *socketio.Server {
	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("Connected:", s.ID())
		return nil
	})
	server.OnConnect("/csv", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("CSV to be processed:", s.ID())
		return nil
	})
	server.OnEvent("/csv", "start-bulk-update", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("Somebody just close the connection ", msg)
	})

	server.OnError("/", func(s socketio.Conn, err error) {
		fmt.Println("An error occured ", err.Error())
	})

	return server
}
