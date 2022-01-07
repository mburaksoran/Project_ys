package shared

var Config = configuration{
	MONGOURL:    "mongodb+srv://Agulhas:Agulhas@cluster0.xa1nz.mongodb.net/myFirstDatabase?retryWrites=true&w=majority",
	POSTGRESURL: "postgres://cgxpmefm:67S9KutggNXm8rXP-OOs4mu4Jmlqat-2@abul.db.elephantsql.com/cgxpmefm",
}

type configuration struct {
	MONGOURL    string
	POSTGRESURL string
}
