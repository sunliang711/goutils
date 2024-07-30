package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sunliang711/goutils/http/types"
	"github.com/sunliang711/goutils/http/utils"
)

const (
	jwtHeaderName = "Authorization"
)

func JwtChecker(secret string) func(c *gin.Context) {
	return func(c *gin.Context) {

		// requestId, _ := c.Get(consts.ContextKeyRequestId)

		// Get header
		token := c.Request.Header.Get(jwtHeaderName)
		if token == "" {
			// logger.Error().Str("header", jwtHeaderName).Msg("cannot get jwt header ")
			c.AbortWithStatusJSON(http.StatusBadRequest, types.Response{
				// RequestId: requestId.(string),
				// Success: false,
				Code: types.CodeGeneralError,
				Msg:  "need jwt token",
				Data: nil,
			})
			return
		}

		// Parse token
		parsedToken, err := utils.ParseJwtToken(token, secret)
		if err != nil {
			// logger.Error().Err(err).Msg("parse jwt token")
			c.AbortWithStatusJSON(http.StatusBadRequest, types.Response{
				// RequestId: requestId.(string),
				// Success:   false,
				Code: types.CodeGeneralError,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		// Type assertion
		_, OK := parsedToken.Claims.(jwt.MapClaims)
		if !OK {
			// logger.Error().Msg("parsed token type assertion failed")
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.Response{
				// RequestId: requestId.(string),
				// Success:   false,
				Code: types.CodeGeneralError,
				Msg:  "unable to parse claims",
				Data: nil,
			})
			return
		}
		c.Next()
	}

}
