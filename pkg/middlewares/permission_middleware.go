package middlewares

import (
	"bruno_authentication/pkg/api_errors"
	"bruno_authentication/pkg/framework"
	"bruno_authentication/pkg/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PermissionMiddleware struct {
	logger  framework.Logger
	cognito CognitoAuthMiddleware
	env     *framework.Env
}

func NewPermissionMiddleware(logger framework.Logger, cognito CognitoAuthMiddleware, env *framework.Env) PermissionMiddleware {
	return PermissionMiddleware{
		logger:  logger,
		env:     env,
		cognito: cognito,
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

func (pm PermissionMiddleware) checkIsRoleSatisfied(roles []string, role interface{}) bool {
	for _, val := range roles {
		if r, ok := role.(string); ok {
			if r == val {
				return true
			}
		}
	}
	return false
}

func (pm PermissionMiddleware) allow(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, check := ctx.Get("isPassed")
		if !check {
			if err := pm.cognito.addClaimsToContext(ctx); err != nil {
				responses.ErrorJSON(ctx, http.StatusUnauthorized, err.Error())
				ctx.Abort()
				return
			}
		}

		ctx.Set("Role", ctx.MustGet("Role").(string))
		ok := pm.checkIsRoleSatisfied(roles, ctx.MustGet("Role"))
		if ok {
			ctx.Next()
		} else {
			pm.logger.Info("unauthorized access")
			responses.ErrorJSON(ctx, http.StatusForbidden, api_errors.ErrUnauthorizedAccess)
			ctx.Abort()
		}
	}
}

func (pm PermissionMiddleware) IsAdmin() gin.HandlerFunc {
	return pm.allow("admin")
}

func (pm PermissionMiddleware) IsStudent() gin.HandlerFunc {
	return pm.allow("student")
}
