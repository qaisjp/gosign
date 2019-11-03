# cosign-webapi [![Go Report Card](https://goreportcard.com/badge/github.com/qaisjp/cosign-webapi)](https://goreportcard.com/report/github.com/qaisjp/cosign-webapi) [![GoDoc](https://godoc.org/github.com/qaisjp/cosign-webapi?status.svg)](https://godoc.org/github.com/qaisjp/cosign-webapi)

This is a fully working (albeit not unit tested) webapi to interact with CoSign. This includes access to `Check`
and also handles the frontend redirections for weblogin callbacks.

This repository uses [gosign](https://github.com/qaisjp/gosign), a library to interact with CoSign daemons.

For more information about Cosign, visit [weblogin.org](http://weblogin.org).

## What does this do?

- Exposes `/check`
- Exposes `/cosign/valid`

## What's left?

- Rate limiting on an IP by IP basis.
- Better API key handling(?)
- Better documentation.

## How to use

1. See example config. Update your server.
2. Set up nginx to proxy_pass `/cosign/valid`, like so:
   ```
   location /cosign {
       proxy_pass       http://localhost:6663;
       proxy_set_header Host $host;
       proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
   }
   ```
3. Send requests to `/check` to check login keys.

## Security

1. By default the webapi binds to `localhost:8080`, you'll want to keep it this way and make sure 8080 inbound is
   blocked from unknown addresses.
2. `/check` only works with a valid API username and API key.
3. Since the application is only exposed to the world via nginx (as above), we can guarantee
   they only have access to `/cosign/valid` (and not `/check`)

