package user

import (
	"bruno_authentication/pkg/api_errors"
	"bruno_authentication/pkg/framework"
	"bruno_authentication/pkg/responses"
	"bruno_authentication/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
	logger  framework.Logger
}

func NewController(service *Service, logger framework.Logger) *Controller {
	return &Controller{service: service, logger: logger}
}

func (c *Controller) Login(ctx *gin.Context) {
	c.logger.Info("[Common Controller...Login]")
	var loginSerializer LoginSerializer
	if err := ctx.ShouldBindJSON(&loginSerializer); err != nil {
		utils.HandleValidationError(c.logger, ctx, err)
		return
	}

	user, err := c.service.Login(loginSerializer)
	if err != nil {
		if errors.Is(err, api_errors.ErrIncorrectCredentials) {
			responses.ErrorJSON(ctx, http.StatusUnauthorized, api_errors.ErrIncorrectCredentials)
			return
		}
		utils.HandleError(c.logger, ctx, err)
		return
	}

	responses.JSON(ctx, http.StatusOK, user)
}

func (c *Controller) TestAdminRoute(ctx *gin.Context) {
	c.logger.Info("[Common Controller...TestAdminRoute]")
	responses.JSON(ctx, http.StatusOK, "Admin Route OK")
}

func (c *Controller) TestStudentRoute(ctx *gin.Context) {
	c.logger.Info("[Common Controller...TestStudentRoute]")
	responses.JSON(ctx, http.StatusOK, "Student Route OK")
}
