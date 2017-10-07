package gosign

import (
	"crypto/tls"
	"net"
	"net/textproto"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// A Client represents a client connection to a CoSign daemon.
type Client struct {
	// text is the textproto.Conn used by clients
	text *textproto.Conn

	// keep a reference to the connection so it can be used to create a TLS
	// connection later
	conn *tls.Conn

	config *Config // configuration passed to constructor
}

// A Config structure is used to configure a CoSign client.
// After one has been passed to a gosign function it must not be
// modified. A Config may be reused; the gosign package will also not
// modify it.
type Config struct {
	Address   string
	Service   string
	TLSConfig *tls.Config
}

// Dial returns a new Client connected to a daemon at addr.
// The addr must include a port, as in "weblogin.inf.ed.ac.uk:6663"
func Dial(conf *Config) (*Client, error) {
	f := &Client{config: conf}

	err := f.dial()
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Internal dial to connect
func (f *Client) dial() (err error) {
	conn, err := net.Dial("tcp", f.config.Address)
	if err != nil {
		return
	}

	f.text = textproto.NewConn(conn)
	_, message, err := f.text.ReadResponse(220)
	if err != nil {
		f.text.Close()
	}

	// Make sure this is PROTOCOL v2
	if strings.HasPrefix(message, "Collaborative Web Single Sign-On ") {
		return errors.New("daemon has protocol version 1, expected protocol version 2")
	} else if !strings.HasPrefix(message, "2 Collaborative Web Single Sign-On ") {
		return errors.Errorf("daemon supplied unknown welcome message: %s", message)
	}

	_, _, err = f.cmd(220, "STARTTLS 2")
	if err != nil {
		return err
	}

	f.conn = tls.Client(conn, f.config.TLSConfig)
	f.text = textproto.NewConn(f.conn)

	code, message, err := f.text.ReadResponse(220)
	if err != nil {
		f.text.Close()
		return errors.Wrapf(err, "expected code 200, got %d %s", code, message)
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
		return errors.Wrap(err, "noop was unsuccessful")
	}

	return nil
}

// Quit sends the QUIT command and closes the connection to the server.
func (f *Client) Quit() error {
	_, msg, err := f.cmd(221, "QUIT")
	if err != nil {
		return errors.Wrap(err, "QUIT failed")
	}

	if msg != "Service closing transmission channel" {
		return errors.Errorf("unexpected response: %s", msg)
	}
	return f.text.Close()
}

// Close closes the connection to the CoSign daemon.
func (f *Client) Close() error {
	return f.text.Close()
}

// Check checks the given login token
func (f *Client) Check(loginToken string) (resp Response, err error) {
	// Make sure login token is clean
	if containsWhitespace(loginToken) {
		err = errors.New("Malformed login token")
		return
	}

	code, msg, err := f.cmd(-1, "CHECK cosign-%s=%s", f.config.Service, loginToken)
	if err != nil {
		return
	}

	// fmt.Printf("%s %s", code, msg)

	// if code == 231 {
	// 	fmt.Println("Success: ", msg)
	// 	return msg, nil
	// } else if code == 430 {
	// 	fmt.Println("Failed: logged out; ", msg)
	// } else if code == 533 {
	// 	fmt.Println("Failed: unknown cookie; ", msg)
	// } else {
	// 	return "", &textproto.Error{code, msg}
	// }

	return Response{code, msg}, nil
}

// cmd is a convenience function that sends a command and returns the response
func (f *Client) cmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	id, err := f.text.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}
	f.text.StartResponse(id)
	defer f.text.EndResponse(id)
	code, msg, err := f.text.ReadResponse(expectCode)
	return code, msg, err
}

// containsWhitespace takes a string, and checks whether it contains whitespace
func containsWhitespace(token string) bool {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, token) != token
}
