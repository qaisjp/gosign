# gosign: CoSign library for Golang [![Go Report Card](https://goreportcard.com/badge/github.com/qaisjp/gosign)](https://goreportcard.com/report/github.com/qaisjp/gosign)

gosign is an **experimental** library that provides an interface to the CoSign daemon. It works well, but the API might change in the future.

[CoSign](http://weblogin.org) is a "secure single sign-on web authentication system".

This only maintains a living connection and can handle `CHECK`. There are no plans to support further protocol
commands. This library is only built to support CoSign protocol version 2 (in use as of Cosign v2.x). Contributions are welcome.
