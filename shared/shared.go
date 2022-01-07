package shared

var Config = configuration{
	MONGOURL:    "",
	POSTGRESURL: "",
}

type configuration struct {
	MONGOURL    string
	POSTGRESURL string
}
