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
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"

	cognitosrp "github.com/alexrudd/cognito-srp/v4"
)

var jwkURL = ""
var issuer = ""
var keySet jwk.Set = jwk.NewSet()

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
	issuer = "https://cognito-idp." + env.AWSRegion + ".amazonaws.com/" + env.UserPoolID
	jwkURL = issuer + "/.well-known/jwks.json"

	keySet, _ = jwk.Fetch(context.Background(), jwkURL)
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

func (cg *CognitoAuthService) VerifyToken(tokenString string) (jwt.Token, error) {
	parsedToken, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
		jwt.WithIssuer(issuer),
	)

	if err != nil {
		return nil, err
	}

	return parsedToken, nil
}
