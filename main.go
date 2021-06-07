package main

import (
	"strings"

	"wss/app"
	"wss/app/command"
	"wss/app/mongodb"
	"wss/app/rest"
	"wss/consumer"
	"wss/libs/goconfig"
	"wss/libs/kafka"
	"wss/libs/log"

	"github.com/joho/godotenv"
)

func main() {
	// Initialze logging
	log.Init()

	// Loading .env
	if err := godotenv.Load(); err != nil {
		log.Info(err.Error())
	}

	// Load configurations
	if err := goconfig.Load(); err != nil {
		panic(err.Error())
	}

	log.Info(goconfig.Get("name"))

	// Consume kafka topics
	consumerObject := struct {
		Topic string
		Group struct {
			ID string
		}
	}{}
	goconfig.GetObject("consumer", &consumerObject)
	kafkaBrokers := strings.Split(goconfig.Get("kafka.brokers").(string), ",")

	kafka.AddConsumerHandler("update_subscription", consumer.HandleSubscriptions)
	log.Info("Consuming kafka messages from", kafkaBrokers)
	go kafka.Consume(kafkaBrokers, "", consumerObject.Topic, 0)

	// Setup service
	var requestPoolSize int
	goconfig.GetObject("requestPoolSize", &requestPoolSize)

	rest.Init(requestPoolSize)

	var notificationURL string
	goconfig.GetObject("api.url.notification", &notificationURL)
	command.NotificationService = rest.NewNotificationService(notificationURL)

	var userURL string
	goconfig.GetObject("api.url.user", &userURL)
	command.UserService = rest.NewUserService(userURL)

	var scheduleURL string
	goconfig.GetObject("api.url.scheduleMobile", &scheduleURL)
	command.ScheduleServiceMobile = rest.NewScheduleService(scheduleURL)

	var tournamentMobileURL string
	goconfig.GetObject("api.url.tournamentMobile", &tournamentMobileURL)
	command.TournamentServiceMobile = rest.NewTournamentService(tournamentMobileURL)

	// Setup application handler for mobile endpoint
	mobileHandler := app.NewApplicationHandler()
	mobileHandler.AddCommandHandler("authentication", command.NewAuthenticationHandler())
	mobileHandler.AddCommandHandler("subscription", command.NewSubscriptionHandler())
	mobileHandler.AddCommandHandler("getMyNotifications", command.NewNotificationHandler())
	mobileHandler.AddCommandHandler("getMyNotificationsCount", command.NewNotificationCountHandler())
	mobileHandler.AddCommandHandler("updateNotifications", command.NewUpdateNotificationHandler())
	mobileHandler.AddCommandHandler("getSchedules", command.NewScheduleHandlerMobile())
	mobileHandler.AddCommandHandler("getTournament", command.NewTournamentHandlerMobile())

	// Setup application handler for web endpoint
	webHandler := app.NewApplicationHandler()
	webHandler.AddCommandHandler("authentication", command.NewAuthenticationHandler())
	webHandler.AddCommandHandler("subscription", command.NewSubscriptionHandler())
	webHandler.AddCommandHandler("getMyNotifications", command.NewNotificationHandler())
	webHandler.AddCommandHandler("getMyNotificationsCount", command.NewNotificationCountHandler())
	webHandler.AddCommandHandler("updateNotifications", command.NewUpdateNotificationHandler())

	// Add application handlers
	app.AddHandler("/mobile", mobileHandler)
	app.AddHandler("/web", webHandler)

	// Setup subscription service
	var dbURL, dbName, dbPassword string
	goconfig.GetObject("db.url", &dbURL)
	goconfig.GetObject("db.name", &dbName)
	goconfig.GetObject("db.pwd", &dbPassword)
	dbURL = strings.Replace(dbURL, "{pwd}", dbPassword, 1)
	subscriptionService := mongodb.NewSubscriptionService(dbURL, dbName)

	log.Info("Connecting to database:", dbURL)
	if err := subscriptionService.Connect(); err != nil {
		log.Fatal("Unable to connect to mongodb:", err.Error())
	}
	defer subscriptionService.Disconnect()

	app.SetSubscriptionService(subscriptionService)

	// Start App
	var port int
	goconfig.GetObject("port", &port)
	app.Serve(port)
}
