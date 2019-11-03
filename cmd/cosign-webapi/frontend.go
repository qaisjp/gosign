package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Valid redirects the user to a thing to do a thing - frontend
func (i *API) Valid(c *gin.Context) {
	key := "cosign-" + i.Config.CoSign.Service
	loginCookie, exists := c.GetQuery(key)
	if !exists {
		fmt.Println("Redirecting because of a missing cookie")
		c.Redirect(http.StatusPermanentRedirect, "https://"+i.Config.CoSign.CGIAddress+"/cosign/validation_error.html")
		return
	}

	// loginCookie should be exactly what was passed to the address bar
	loginCookie = url.QueryEscape(loginCookie)

	// localhost:8080/cosign/valid?cosign-betterinformatics.com=sIDrTIml5hWfK5uu9TAQ8-mdwmif6An81-8vgs2qjzupJ5w4rldrWzyxgsUKNLEY3Ovsjd-doqO9xXdQ421h6dA+k5tiQkhbek79PczciT590awVKvFviT9gQUIY&https://RETURN_ADDRESS_HERE
	redirect := strings.TrimPrefix(c.Request.URL.RawQuery, key+"="+loginCookie+"&")
	_, err := url.ParseRequestURI(redirect)
	if err != nil {
		fmt.Println("Redirecting because of a bad URI parse; ", err.Error())
		c.Redirect(http.StatusPermanentRedirect, "https://"+i.Config.CoSign.CGIAddress+"/cosign/validation_error.html")
		return
	}

	host := c.Request.Host
	if strings.Contains(host, ":") {
		host, _, err = net.SplitHostPort(c.Request.Host)
		if err != nil {
			fmt.Println("Redirected because of a SplitHostPort error")
			c.Redirect(http.StatusPermanentRedirect, "https://"+i.Config.CoSign.CGIAddress+"/cosign/validation_error.html")
			return
		}
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     key,
		Value:    strings.Replace(url.QueryEscape(loginCookie), "%2B", "+", -1),
		MaxAge:   43200, // 12 hour hard timeout period for CoSign
		Path:     "/",
		Domain:   host,
		Secure:   !i.Config.Insecure, // limit to secure webpages?
		HttpOnly: true,               // meaning that client js cannot access the cookie
	})

	c.Redirect(http.StatusPermanentRedirect, redirect)
}
