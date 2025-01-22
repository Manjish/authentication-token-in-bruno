package middlewares

import (
	"bruno_authentication/pkg/api_errors"
	"bruno_authentication/pkg/framework"
	"bruno_authentication/pkg/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PermissionMiddleware struct {
	logger framework.Logger
	env    *framework.Env
}

func NewPermissionMiddleware(logger framework.Logger, env *framework.Env) PermissionMiddleware {
	return PermissionMiddleware{
		logger: logger,
		env:    env,
	}
}

func (pm PermissionMiddleware) BasicAuthPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, password, hasAuth := ctx.Request.BasicAuth()
		if !hasAuth {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}
		if username != pm.env.BasicAuthUsername || password != pm.env.BasicAuthPassword {
			responses.ErrorJSON(ctx, http.StatusUnauthorized, api_errors.ErrUnauthorizedAccess)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
