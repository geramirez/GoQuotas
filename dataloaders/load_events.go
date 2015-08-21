package main

import (
	"log"
	"time"

	helpers "github.com/ramirezg/GoQuotas/helpers"
)

func main() {
	// Get Mongo session
	session := helpers.GetMongoSession()
	defer session.Close()
	collection := session.DB("cloudfoundry").C("events")
	// Make new token
	token := helpers.NewToken()
	// events Generator
	eventsGen := token.EventGen()
	// get event indefinitely
	endLoop := false
	for _ = range time.Tick(1 * time.Second) {
		apiResponse := eventsGen()
		if apiResponse.NextUrl == "" {
			// Break loop if there are no more urls
			break
		}
		for _, event := range apiResponse.Resources {
			info, err := collection.Upsert(event, event)
			if err != nil {
				log.Fatal(err)
			}
			if info.Updated == 1 {
				// Break loop if there are document exist in the database
				endLoop = true
				break
			}
		}
		if endLoop == true {
			break
		}
	}
}
