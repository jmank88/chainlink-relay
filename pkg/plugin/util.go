package plugin

import (
	"net"

	"github.com/pkg/errors"
)

func MustCall(err error) {
	if err != nil {
		panic(errors.Wrap(err, "plugin Call failed"))
	}
}

func MustDial(conn net.Conn, err error) net.Conn {
	if err != nil {
		panic(errors.Wrap(err, "plugin Dial failed"))
	}
	return conn
}
