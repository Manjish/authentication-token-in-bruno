package commands

import (
	"bruno_authentication/migrations"
	"bruno_authentication/pkg/framework"

	"github.com/spf13/cobra"
)

type MigrateCommand struct {
}

func (s *MigrateCommand) Short() string {
	return "run migrate command"
}

func NewMigrateCommand() *MigrateCommand {
	return &MigrateCommand{}
}

func (s *MigrateCommand) Setup(cmd *cobra.Command) {
	// left empty intentionally
}

func (s *MigrateCommand) Run() framework.CommandRunner {
	return func(
		migrator *migrations.Migrator,
	) {
		if err := migrator.Exec(); err != nil {
			framework.GetLogger().Fatal(err)
		}
	}
}
