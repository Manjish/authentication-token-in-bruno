package user

import (
	"bruno_authentication/pkg/api_errors"
	"bruno_authentication/pkg/framework"
	"bruno_authentication/pkg/services"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// Service handles the business logic of the module
type Service struct {
	repo           *Repository
	logger         framework.Logger
	cognitoService services.CognitoAuthService
}

// NewService creates a new instance of TestService
func NewService(repo *Repository, logger framework.Logger, cognitoService services.CognitoAuthService) *Service {
	return &Service{
		repo:           repo,
		logger:         logger,
		cognitoService: cognitoService,
	}
}

func (s *Service) Login(loginDetails LoginSerializer) (*types.AuthenticationResultType, error) {
	s.logger.Info("[User...Service...Login]")

	loginResponse, err := s.cognitoService.Login(loginDetails.Email, loginDetails.Password)

	if err != nil {
		return nil, err
	}

	if loginResponse == nil {
		return nil, api_errors.ErrUnauthorizedAccess
	}

	return loginResponse, nil
}
