package protect

import (
	"log"
	"net"
	"os"
	"syscall"
)

func Dial(network, addr string) (net.Conn, error) {
	var xfd uintptr
	d := &net.Dialer{}
	d.Control = func(network, address string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			// Access socket fd
			xfd = fd
		})
	}
	_, err := d.Dial(network, addr)
	if err != nil {
		panic(err)
	}
	log.Printf("fd:: %v", xfd)
	file := os.NewFile(xfd, "Socket")
	defer func() {_ = file.Close()}()
	return net.FileConn(file)
}
