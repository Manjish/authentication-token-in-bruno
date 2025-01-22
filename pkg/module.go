package pkg

import (
	"bruno_authentication/pkg/framework"
	"bruno_authentication/pkg/infrastructure"
	"bruno_authentication/pkg/middlewares"
	"bruno_authentication/pkg/services"

	"go.uber.org/fx"
)

var Module = fx.Module("pkg",
	framework.Module,
	services.Module,
	middlewares.Module,
	infrastructure.Module,
)
