package configs

var (
	GinMode      string
	IsGinInDebug bool
)

func ConfigGin(mode string) {
	GinMode = mode
	IsGinInDebug = GinMode == "debug"
}
