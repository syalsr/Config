package database

type Config struct{
	Key1 string	`json:"key1"`
	Key2 string	`json:"key2"`
}

type Request struct{
	Service string 	`json:"service"`
	Version int		`json:"version"`
	Config []byte	`json:"config"`
}