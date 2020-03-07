package main
import (
  "github.com/gin-gonic/gin"
  "log"
  "fmt"
  "context"
  "routes"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  // "go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitConnectionDB() (*mongo.Client, error) {
  // Create client
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
  if err != nil {
      log.Fatal(err)
  }

  // Create connect
  err = client.Connect(context.TODO())
  if err != nil {
      log.Fatal(err)
  }

  // Check the connection
  err = client.Ping(context.TODO(), nil)
  if err != nil {
      log.Fatal(err)
  }

  fmt.Println("Connected to MongoDB!")

  return client, nil
}

func main(){

  client, err := InitConnectionDB()
  if err != nil {
    log.Fatal(err)
  }

  collection := client.Database("appdb").Collection("articles")

  r := gin.Default()
  routes.RoutesInit(collection, r)
  r.Run(":8090")
}
