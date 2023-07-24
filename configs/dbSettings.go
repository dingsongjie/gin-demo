package configs

import "github.com/namsral/flag"

var (
	UserInfoConnectionString string
	NewPathConnectionString  string
)

func ConfigDb(commandSet *flag.FlagSet) {
	if command := commandSet.Lookup("user-info-db-connection"); command != nil {
		UserInfoConnectionString = command.Value.String()
	}

	if command := commandSet.Lookup("new-path-db-connection"); command != nil {
		NewPathConnectionString = command.Value.String()
	}
}
