package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run atlasTest.go <ATLAS_URI>")
		fmt.Println("  ATLAS_URI: something like mongodb+srv://<username>:<password>@cluster0-zzart.mongodb.net/test?retryWrites=true&w=majority")
		os.Exit(1)
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	
	defer client.Disconnect(ctx)
}