package main

import (
	"fmt"

	"github.com/mnaufalhilmym/bookshelf/internal/config"
)

func main() {
	conf := config.NewViper()

	config.ConfigureTracing(
		conf.GetString("log.print_level"),
		conf.GetString("log.stacktrace_level"),
		conf.GetUint("log.max_pc"),
	)

	db := config.NewDatabase(
		conf.GetInt("db.pool.idle"),
		conf.GetInt("db.pool.max"),
		conf.GetInt("db.pool.lifetime"),
	)

	router := config.NewGin(conf.GetString("app.mode"))

	config.Bootstrap(
		router,
		db,
		conf.GetString("jwt.key"),
		conf.GetDuration("jwt.duration"),
	)

	if err := router.Run(conf.GetString("web.address")); err != nil {
		panic(fmt.Errorf("failed to start server: %w", err))
	}
}
