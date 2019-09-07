package main

import (
	// The confg lib may be changed to some golang existing library such as viper etc.

	"strings"

	"./app"
	"./command"
	"./consumer"
	"./goconfig"
	"./log"
	"./service"
	"github.com/joho/godotenv"
)

func main() {
	// Loading .env
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	// Load configurations
	if err := goconfig.Load(); err != nil {
		panic(err.Error())
	}

	// Initialze logging
	log.Init()
	log.Info(goconfig.Get("name"))

	// Consume kafka
	consumerObject := struct {
		Group struct {
			ID string
		}
		Topic string
	}{}
	goconfig.GetObject("consumer", &consumerObject)
	consumer.Init()
	kafkaBrokers := strings.Split(goconfig.Get("kafka.brokers").(string), ",")
	consumer.Consume(kafkaBrokers, consumerObject.Group.ID, consumerObject.Topic, 0)

	// Setup service
	var notificationURL, userURL string
	var requestPoolSize int
	goconfig.GetObject("api.url.notification", &notificationURL)
	goconfig.GetObject("api.url.user", &userURL)
	goconfig.GetObject("requestPoolSize", &requestPoolSize)
	service.Init(requestPoolSize, notificationURL, userURL)

	// Start App
	app.Init()
	command.Setup()
	var port int
	goconfig.GetObject("port", &port)
	app.Serve(port)
}
