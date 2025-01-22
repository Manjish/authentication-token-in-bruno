package infrastructure

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Module(
	"infrastructure",
	fx.Options(
		fx.Provide(NewRouter), 
		fx.Provide(NewAWSConfig), 
		fx.Provide(NewCognitoClient), 
		fx.Provide(NewS3Client), 
		fx.Provide(NewS3PresignClient), 
		fx.Provide(NewS3Uploader),
	),
)
