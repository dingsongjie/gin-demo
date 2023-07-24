package configs

import "github.com/namsral/flag"

var (
	GinMode      string
	IsGinInDebug bool
)

func ConfigGin(commandSet *flag.FlagSet) {
	if command := commandSet.Lookup("gin-mode"); command != nil {
		GinMode = command.Value.String()
	}
	IsGinInDebug = GinMode == "debug"
}
