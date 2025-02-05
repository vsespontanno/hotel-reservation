package main

import (
	"context"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	db := db.NewMongoHotelStore(client, db.DBNAME)
	hotel := types.Hotel{
		Name:     "Bellucia",
		Location: "France",
	}

	room := types.Room{
		Type:      types.SingleRoomType,
		BasePrice: 99.9,
	}
}
