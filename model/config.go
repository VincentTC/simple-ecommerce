package model

import (
	"github.com/VincentTC/simple-ecommerce/util/database/mysql"
)

type Config struct {
	Server       ServerConfig       `envconfig:"SERVER"`
	Database     mysql.Config       `envconfig:"DATABASE"`
	Apps         AppsConfig         `envconfig:"APPS"`
	CronSchedule CronScheduleConfig `envconfig:"CRON_SCHEDULE"`
}

type ServerConfig struct {
	HTTPPort string `envconfig:"HTTP_PORT"`
}

type AppsConfig struct {
	AuthSecret string `envconfig:"AUTH_SECRET"`
}

type CronScheduleConfig struct {
	OrderReminder string `envconfig:"ORDER_REMINDER"`
}
