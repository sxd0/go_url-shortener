package configs

import "fmt"

type Db struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (d *Db) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		d.Host, d.User, d.Password, d.Name, d.Port,
	)
}
