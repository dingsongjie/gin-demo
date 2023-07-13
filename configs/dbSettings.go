package configs

var (
	UserInfoConnectionString string
	NewPathConnectionString  string
)

func ConfigDb(userInfoConnectionString string, newPathConnectionString string) {
	UserInfoConnectionString = userInfoConnectionString
	NewPathConnectionString = newPathConnectionString
}
