package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ServerAddress    string `mapstructure:"SERVER_ADDRESS"`
	PostgresDB       string `mapstucture:"POSTGRES_DB"`
	PostgresUser     string `mapstucture:"POSTGRES_USER"`
	PostgresPassword string `mapstucture:"POSTGRES_PASSWORD"`
}

func LoadEnvVariables() *Config {
	var c Config
	//tell viper the path of env file
	viper.AddConfigPath("./")
	//tell viper the name of file
	viper.SetConfigName("app")
	//tell viper type of file
	viper.SetConfigType("env")

	//reads all the variables from env file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Error reading env file", err)
	}

	//unmarshal the loaded env variables
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(&c)
	return &c
}
