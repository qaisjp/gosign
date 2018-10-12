package gosign

import (
	"crypto/tls"
	"fmt"
	"net/textproto"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// A Client represents a client connection to a collection of CoSign daemons.
type Client struct {
	daemon *daemon

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

// Dial returns a new Client connected to all daemons at addr.
// The addr must include a port, as in "weblogin.inf.ed.ac.uk:6663"
func Dial(conf *Config) (*Client, error) {
	f := &Client{config: conf}

	c, err := dialDaemon(conf.Address, conf)
	if err != nil {
		return nil, err
	}

	f.daemon = c

	return f, nil
}

// Quit sends the QUIT command to all servers and closes the connections.
// If all connections are already closed, this returns nil.
func (f *Client) Quit() error {
	return f.daemon.quit()
}

// Close closes the connection to the CoSign daemon.
func (f *Client) Close() error {
	return f.daemon.close()
}

// Check allows clients to retrieve information about a user based on the
// cookie presented to the daemon.
//
// This is typically used by both the CGI and the filter (service).
func (f *Client) Check(cookie string, serviceCookie bool) (resp CheckResponse, err error) {
	// Make sure login/service cookie is clean
	if containsWhitespace(cookie) {
		err = errors.New("Malformed cookie")
		return
	}

	prefix := "cosign-"
	if serviceCookie {
		prefix = "cosign="
	}

	cmd := fmt.Sprintf("CHECK %s%s=%s", prefix, f.config.Service, cookie)

	var msg string

	{
		daemon := f.daemon
		var code int

		code, msg, err = daemon.cmd(-1, cmd)
		if err != nil {
			if !daemon.closed {
				err := f.Close()
				if err != nil {
					return resp, errors.Wrap(err, "initial close failed before attempting to reconnect")
				}
			}

			daemon, err = dialDaemon(daemon.address, f.config)
			f.daemon = daemon // todo: update index
			if err != nil {
				return resp, errors.Wrap(err, "failed to reconnect")
			}

			code, msg, err = daemon.cmd(-1, cmd)
			if err != nil {
				return resp, errors.Wrap(err, "cmd failed even after reconnect")
			}
		}

		// Permitted response codes for CHECK are:
		// - 231 (for a login_cookie)
		// - 232 (for a service_cookie)
		// NOT: 233 (same method in cosignd, but only responds to REKEY)
		if code != 231 && code != 232 {
			if (code == 430) && (msg == "CHECK: Already logged out") {
				// CoSign bug: 430 is returned for another error as well
				return resp, ErrLoggedOut
			} else if code == 431 {
				return resp, ErrLoggedOut
			}

			// If code is 533 then "cookie not in db" and we need to try another daemon
			if code == 533 {
				// try another daemon
				// continue
			}

			return resp, &textproto.Error{
				Code: code,
				Msg:  msg,
			}
		}
	}

	// Lets go ahead and split the message by spaces
	segments := strings.Split(msg, " ")

	// Set the IP/Principal to the first/second segment
	resp.IP = segments[0]
	resp.Principal = segments[1]

	// Set the factors to the segment list, excluding the first two items
	factors := segments[2:]

	// Set the realm to the first factor
	resp.Realm = factors[0]

	// Sometimes `factors` will have a trailing empty string.
	// Lets deal with that.
	//
	// TODO: Investigate this. Do all messages have a trailing space?
	//		 Is this being picked up by the string split?
	if factors[len(factors)-1] == "" {
		factors = factors[:len(factors)-1]
	}

	// Actually set the response factors to our factors
	resp.Factors = factors

	// That's all folks!
	return resp, nil
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
