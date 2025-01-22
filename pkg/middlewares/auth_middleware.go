package middlewares

import (
	"bruno_authentication/pkg/api_errors"
	"bruno_authentication/pkg/responses"
	"bruno_authentication/pkg/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
)

type CognitoAuthMiddleware struct {
	service services.CognitoAuthService
}

func NewCognitoAuthMiddleware(service services.CognitoAuthService) CognitoAuthMiddleware {
	return CognitoAuthMiddleware{
		service: service,
	}
}

func (am CognitoAuthMiddleware) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := am.addClaimsToContext(ctx); err != nil {
			responses.ErrorJSON(ctx, http.StatusUnauthorized, err.Error())
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (am CognitoAuthMiddleware) addClaimsToContext(ctx *gin.Context) error {
	token, err := am.getTokenFromHeader(ctx)
	if err != nil {
		return api_errors.ErrUnauthorizedAccess
	}
	claims := token.PrivateClaims()

	ctx.Set("claims", claims)

	// set role in context
	role, ok := claims["custom:role"]
	if ok {
		ctx.Set("Role", role)
	}

	ctx.Set("isPassed", true)
	return nil
}

func (am CognitoAuthMiddleware) getTokenFromHeader(gc *gin.Context) (jwt.Token, error) {
	header := gc.GetHeader("Authorization")
	idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
	token, err := am.service.VerifyToken(idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}
