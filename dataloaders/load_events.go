package main

import (
	"log"
	"time"

	helpers "github.com/ramirezg/GoQuotas/helpers"
)

func main() {
	// Get Mongo session
	session := helpers.GetMongoSession()
	collection := session.DB("cloudfoundry").C("events")
	// Make new token
	token := helpers.NewToken()
	// events Generator
	eventsGen := token.EventGen()
	// get event indefinitely
	for _ = range time.Tick(1 * time.Second) {
		for _, event := range eventsGen().Resources {
			err := collection.Insert(event)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
