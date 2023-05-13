package main

import (
	"fmt"
	"strings"

	internel "github.com/jbliao/proj-quiz1/internal"
	"github.com/spf13/viper"
)

func main() {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("quiz1")
	v.AutomaticEnv()
	v.BindEnv("app.port", "db.user", "db.passwd", "db.host", "db.name")

	app := internel.NewApp(
		v.GetInt("app.port"),
		fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true",
			v.Get("db.user"), v.Get("db.passwd"), v.Get("db.host"), v.Get("db.name")),
	)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
