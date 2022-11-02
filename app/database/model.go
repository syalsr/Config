package database

type Config struct{
	key1 string	`json: key1`
	key2 string	`json: key2`
}

type Request struct{
	version int	`json: version`
	config Config
}