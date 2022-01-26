package database

type UserCredentials struct {
	Url      string
	PSQLurl  string
	User     string
	Password string
	Host     string
	Port     string
	Database string
}