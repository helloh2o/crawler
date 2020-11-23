package protect

import (
	"log"
	"testing"
)

func TestDial(t *testing.T) {
	conn, err := Dial("tcp","baidu.com:443")
	if err != nil {
		panic(err)
	}
	log.Printf("Local::%s => Remote::%s",conn.LocalAddr(),conn.RemoteAddr())
}


