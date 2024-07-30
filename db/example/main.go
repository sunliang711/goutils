package main

import (
	"fmt"

	"github.com/sunliang711/goutils/db"
	"gorm.io/gorm"
)

type MyTable struct {
	Name string

	gorm.Model
}

func main() {
	var configs []db.DatabaseConfig
	configs = append(configs, db.DatabaseConfig{
		Name:   "mysql1",
		Dsn:    "root:root@tcp(10.1.9.120:3306)/test?charset=utf8&parseTime=True&loc=Local&timeout=1000ms",
		Driver: "mysql",
		Tables: []db.Table{
			{
				Name:       "test",
				Definition: &MyTable{},
			},
		},
	})

	database := db.NewDatabase(configs)

	err := database.Init()
	if err != nil {
		panic(fmt.Sprintf("init database error: %v", err))
	}

	data := &MyTable{Name: "haha"}
	database.GetDatabase("mysql1").Create(data)

}
