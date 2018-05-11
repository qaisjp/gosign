# gosign: CoSign library for Go [![Go Report Card](https://goreportcard.com/badge/github.com/qaisjp/gosign)](https://goreportcard.com/report/github.com/qaisjp/gosign) [![GoDoc](https://godoc.org/github.com/qaisjp/gosign?status.svg)](https://godoc.org/github.com/qaisjp/gosign)

gosign is an **experimental** library that provides an interface to a CoSign daemon. It works well, but the API might change in the future.

[CoSign](http://weblogin.org) is a "secure single sign-on web authentication system".

This only maintains a living connection and can handle the `CHECK` command (this project was created for a "CoSign filter"). There are no plans to support further protocol
commands. This library is only built to support CoSign protocol version 2 (in use as of Cosign v2.x). Contributions are welcome.

## Example

**Creating a CoSign client**

```go
client, err := gosign.Dial(&gosign.Config{
  Address: "www.ease.ed.ac.uk:6663",
  Service: "betterinformatics.com",
  TLSConfig: &tls.Config{
    ServerName:         "www.ease.ed.ac.uk",
    Certificates:       []tls.Certificate{cert},
    RootCAs:            pool,
  },
})
```

- `Address` is the address of your CoSign daemon. It is usually the same address of your university's login portal.
- `Service` is the name of your service, assigned to you by the daemon operators (this is the domain name of your service).
- `TLSConfig` uses the stdlib [`tls.Config`](https://golang.org/pkg/crypto/tls/#Config):
  - `ServerName` is the name of the domain, required if you want the client to verify the server's certificate chain and host name (default)
  - `Certificates` should contain the service certificate given to you by the daemon operators
  - `RootCA` is required as CoSign certificates don't use regular website root CAs
  - (see the _Certificates_ section below for more info)

**Certificates**

You can get `cert` for `Certificates` by doing the following:

```go
cert, err := tls.LoadX509KeyPair("service.crt", "service.key")
if err != nil {
  panic("could not read certfile+keyfile")
}
```

You can get `pool` for `RootCAs` by doing the following:

```go
// Read CAFile containing multiple certs
certs, err := ioutil.ReadFile("cosign.CA.crt")
if err != nil {
  panic("could not read CAFile")
}

// Build a cert pool based from the CAFile
pool := x509.NewCertPool()
pool.AppendCertsFromPEM(certs)
```

**Checking CoSign cookies**

Once you have retrieved a `cosign-service.com` (e.g `cosign-betterinformatics.com`) cookie from a (web) client,
you can then verify the logged in state of the cookie and retrieve information about that user.

```go
response, err := client.Check(cookie, false)

// The only gosign related error is ErrLoggedOut.
if err == gosign.ErrLoggedOut {
  panic("not logged in due to various reasons")
}

// There could be some other error, like a network issue.
if err != nil {
  panic(err.Error())
}

// Success! Print out the response.
fmt.Println(response)
```

The response printed out in this example code is just [this CheckResponse struct](https://godoc.org/github.com/qaisjp/gosign#CheckResponse).

## Projects using gosign

- [cosign-webapi](https://github.com/qaisjp/cosign-webapi/) is a web service that exposes the CHECK command over a REST API to save you from reimplementing CoSign in other languages. It is designed for firewalled access and also authenticates based on defined API keys.
