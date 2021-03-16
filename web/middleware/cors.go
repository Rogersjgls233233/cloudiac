package middleware

import (
	"net/http"
	"cloudiac/libs/ctx"
)

var (
	allowHeaders  = "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token"
	exposeHeaders = "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type"
)

func Cors(c *ctx.GinRequestCtx) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", allowHeaders)
	c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
	c.Header("Access-Control-Expose-Headers", exposeHeaders)
	c.Header("Access-Control-Allow-Credentials", "true")

	//放行所有OPTIONS方法
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}

	c.Next()
}
