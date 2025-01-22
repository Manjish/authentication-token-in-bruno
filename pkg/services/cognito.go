package services

import (
	"bruno_authentication/pkg/framework"
	"bruno_authentication/pkg/utils"
	"errors"
	"time"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	cognitosrp "github.com/alexrudd/cognito-srp/v4"
)

type CognitoAuthService struct {
	client *cognitoidentityprovider.Client
	env    *framework.Env
	logger framework.Logger
}

func NewCognitoAuthService(
	client *cognitoidentityprovider.Client,
	env *framework.Env,
	logger framework.Logger,
) CognitoAuthService {
	return CognitoAuthService{
		client: client,
		env:    env,
		logger: logger,
	}
}

func (cg *CognitoAuthService) Login(email, password string) (*types.AuthenticationResultType, error) {
	// generates the SRP_A required for authentication
	csrp, _ := cognitosrp.NewCognitoSRP(email, password, cg.env.UserPoolID, cg.env.ClientID, nil)

	// initiate auth
	resp, err := cg.client.InitiateAuth(context.Background(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserSrpAuth,
		ClientId:       aws.String(csrp.GetClientId()),
		AuthParameters: csrp.GetAuthParams(),
	})
	if err != nil {
		return nil, err
	}

	// respond to password verifier challenge provided by AWS Cognito
	if resp.ChallengeName == types.ChallengeNameTypePasswordVerifier {
		challengeResponses, _ := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())

		resp, err := cg.client.RespondToAuthChallenge(context.Background(), &cognitoidentityprovider.RespondToAuthChallengeInput{
			ChallengeName:      types.ChallengeNameTypePasswordVerifier,
			ChallengeResponses: challengeResponses,
			ClientId:           aws.String(csrp.GetClientId()),
		})
		if err != nil {
			if awsErr := utils.MapAWSError(cg.logger, err); awsErr != nil {
				if awsErr.ExceptionType == "NotAuthorizedException" {
					return nil, errors.New("incorrect credentials")
				}
			}
			return nil, err
		}
		return resp.AuthenticationResult, nil
	} else {
		cg.logger.Info("Failed")
	}
	return nil, errors.New("failed to login")
}
