package gosign

import (
	"crypto/tls"
	"net"
	"net/textproto"
	"strings"

	"github.com/pkg/errors"
)

// A daemon represents an internal connection to a single CoSign daemon
type daemon struct {
	// text is the textproto.Conn used by clients
	text *textproto.Conn

	closed bool
}

// Internal dial to connect
func dialConn(config *Config) (d *daemon, err error) {
	d = &daemon{closed: true}

	conn, err := net.Dial("tcp", config.Address)
	if err != nil {
		return
	}

	d.text = textproto.NewConn(conn)
	_, message, err := d.text.ReadResponse(220)
	if err != nil {
		d.text.Close()
		return
	}

	// Make sure this is PROTOCOL v2
	if strings.HasPrefix(message, "Collaborative Web Single Sign-On ") {
		return nil, errors.New("daemon has protocol version 1, expected protocol version 2")
	} else if !strings.HasPrefix(message, "2 Collaborative Web Single Sign-On ") {
		return nil, errors.Errorf("daemon supplied unknown welcome message: %s", message)
	}

	_, _, err = d.cmd(220, "STARTTLS 2")
	if err != nil {
		return nil, err
	}

	tlsConn := tls.Client(conn, config.TLSConfig)
	d.text = textproto.NewConn(tlsConn)

	code, message, err := d.text.ReadResponse(220)
	if err != nil {
		d.text.Close()
		return nil, errors.Wrapf(err, "expected code 200, got %d %s", code, message)
	}

	defer func() {
		rerr := recover()
		if rerr != nil {
			err = errors.Wrapf(rerr.(error), "noop was unsuccessful... did you provide the right key/crt?")
			d.text.Close()
		}
	}()

	// Make sure the NOOP works
	_, _, err = d.cmd(250, "NOOP")
	if err != nil {
		d.text.Close()
		return nil, errors.Wrap(err, "noop was unsuccessful")
	}

	d.closed = false

	return d, nil
}

// quit sends the QUIT command and closes the connection to the daemon.
// If the connection is already closed, this returns nil.
func (d *daemon) quit() error {
	if d.closed {
		return nil
	}

	_, msg, err := d.cmd(221, "QUIT")
	if err != nil {
		return errors.Wrap(err, "QUIT failed")
	}

	if msg != "Service closing transmission channel" {
		return errors.Errorf("unexpected response: %s", msg)
	}
	return d.text.Close()
}

// close closes the connection to the daemon.
func (d *daemon) close() error {
	if d.closed {
		return errors.New("connection to daemon already closed")
	}
	d.closed = true
	return d.text.Close()
}

// cmd is a convenience function that sends a command and returns the response
func (d *daemon) cmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	id, err := d.text.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}
	d.text.StartResponse(id)
	defer d.text.EndResponse(id)
	code, msg, err := d.text.ReadResponse(expectCode)
	return code, msg, err
}
