package requests

import (
	"net"
	"time"
)

func DialTcp(target string) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", target, 3*time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
