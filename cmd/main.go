package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"log"
	"ridenow/matcher"
	"ridenow/matcher/models"
	"ridenow/matcher/queue"
)

type Env struct {
	db   models.Datastore
	cons *queue.QueueConsumer
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	// get config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
	dbUser := viper.GetString("DB_USER")
	dbPass := viper.GetString("DB_PASSWORD")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/users?sslmode=disable", dbUser, dbPass, dbHost, dbPort)

	qUser := viper.GetString("AMQP_USER")
	qPass := viper.GetString("AMQP_PASSWORD")
	qHost := viper.GetString("AMQP_HOST")
	qPort := viper.GetString("AMQP_PORT")
	qUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", qUser, qPass, qHost, qPort)

	db, err := models.NewDB(dbUrl)
	if err != nil {
		log.Panic(err)
	}
	cons, err := queue.NewQueueConsumer(qUrl)
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db, cons}

	forecasts, err := env.cons.Subscribe("ridenow.forecasts.create", "ridenow.forecasts.update")

	wait := make(chan bool)

	go func() {
		for f := range forecasts {
			forecast := &matcher.Forecast{}
			err := proto.Unmarshal(f.Body, forecast)
			if err != nil {
				log.Panic(err)
			}
			log.Printf("%s", forecast)
			matches, err := env.db.MatchUsers(forecast)
			if err != nil {
				log.Panic(err)
			}
			for _, match := range matches {
				// TODO: send match to the notifier component
				fmt.Printf(" * Match: %v @ %v [wave height: %v m, time: %v]\n", match.User.Username, match.Location.Id, match.WaveHeightM, match.Time)
			}
		}
	}()
	log.Printf(" [*] Running `matcher` service . To exit press CTRL+C")
	<-wait
}
