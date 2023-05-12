package middleware

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		// claims := &Claims{}
		// token, err := ctx.Cookie("session_token")
		// if err != nil {
		// 	jwt
		// }
		// TODO: answer here
	})
}
