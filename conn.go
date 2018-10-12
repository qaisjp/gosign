package gosign

import (
	"crypto/tls"
	"net"
	"net/textproto"
	"strings"

	"github.com/pkg/errors"
)

// A conn represents an internal connection to a single CoSign daemon
type conn struct {
	// text is the textproto.Conn used by clients
	text *textproto.Conn

	closed bool
}

// Internal dial to connect
func dialConn(config *Config) (f *conn, err error) {
	f = &conn{closed: true}

	conn, err := net.Dial("tcp", config.Address)
	if err != nil {
		return
	}

	f.text = textproto.NewConn(conn)
	_, message, err := f.text.ReadResponse(220)
	if err != nil {
		f.text.Close()
		return
	}

	// Make sure this is PROTOCOL v2
	if strings.HasPrefix(message, "Collaborative Web Single Sign-On ") {
		return nil, errors.New("daemon has protocol version 1, expected protocol version 2")
	} else if !strings.HasPrefix(message, "2 Collaborative Web Single Sign-On ") {
		return nil, errors.Errorf("daemon supplied unknown welcome message: %s", message)
	}

	_, _, err = f.cmd(220, "STARTTLS 2")
	if err != nil {
		return nil, err
	}

	tlsConn := tls.Client(conn, config.TLSConfig)
	f.text = textproto.NewConn(tlsConn)

	code, message, err := f.text.ReadResponse(220)
	if err != nil {
		f.text.Close()
		return nil, errors.Wrapf(err, "expected code 200, got %d %s", code, message)
	}

	defer func() {
		rerr := recover()
		if rerr != nil {
			err = errors.Wrapf(rerr.(error), "noop was unsuccessful... did you provide the right key/crt?")
			f.text.Close()
		}
	}()

	// Make sure the NOOP works
	_, _, err = f.cmd(250, "NOOP")
	if err != nil {
		f.text.Close()
		return nil, errors.Wrap(err, "noop was unsuccessful")
	}

	f.closed = false

	return f, nil
}

// quit sends the QUIT command and closes the connection to the daemon.
// If the connection is already closed, this returns nil.
func (f *conn) quit() error {
	if f.closed {
		return nil
	}

	_, msg, err := f.cmd(221, "QUIT")
	if err != nil {
		return errors.Wrap(err, "QUIT failed")
	}

	if msg != "Service closing transmission channel" {
		return errors.Errorf("unexpected response: %s", msg)
	}
	return f.text.Close()
}

// close closes the connection to the daemon.
func (f *conn) close() error {
	if f.closed {
		return errors.New("connection already closed")
	}
	f.closed = true
	return f.text.Close()
}

// cmd is a convenience function that sends a command and returns the response
func (f *conn) cmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	id, err := f.text.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}
	f.text.StartResponse(id)
	defer f.text.EndResponse(id)
	code, msg, err := f.text.ReadResponse(expectCode)
	return code, msg, err
}
