package backend

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Valid redirects the user to a thing to do a thing
// - success: abuse checks pass, and cosign returns a code + message
// - failure: cosign died somewhere, or the IP address is flagged for abuse
func (i *Impl) Check(c *gin.Context) {
	cookie := c.Param("login_cookie")
	response, err := i.Filter.Check(cookie)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"code":    response.Code,
		"message": response.Message,
	})
}
