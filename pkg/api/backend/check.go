package backend

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Valid redirects the user to a thing to do a thing
func (i *Impl) Check(c *gin.Context) {
	cookie := c.Param("login_cookie")
	msg, err := i.Filter.Check(cookie)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": msg,
	})
}
