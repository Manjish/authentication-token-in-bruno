package bootstrap

import (
	"bruno_authentication/domain"
	"bruno_authentication/migrations"
	"bruno_authentication/pkg"
	"bruno_authentication/seeds"

	"go.uber.org/fx"
)

var CommonModules = fx.Module("common",
	fx.Options(
		pkg.Module,
		seeds.Module,
		migrations.Module,
		domain.Module,
	),
)
