package main
import (
  // "github.com/gin-gonic/gin"
  "log"
  "fmt"
  "context"
  "routes"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  // "go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitConnectionDB() (*mongo.Client) {
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
  return client
}

func main(){
  client := InitConnectionDB()
  collection := client.Database("appdb").Collection("articles")
  router := routes.SetupRouter(collection)
  router.Run(":8090")
}
