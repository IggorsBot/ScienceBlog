package articles

import (
  "time"
  "go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID

)

type Article struct {
  ID primitive.ObjectID `bson:"_id"`
  Author string `bson:"author"`
  Title string `bson:"title"`
  Slug string `bson:"slug"`
  Desctiption string `bson:"desctiption"`
  Images []string `bson:"images"`
  Videos []string `bson:"videos"`
  Body string `bson:"body"`
  Tags []string `bson:"tags"`
  Category string `bson:"category"`
  TimeCreated time.Time `bson:"time_created"`
}


type Comment struct {
  ID primitive.ObjectID `bson:"_id"`
  Author string `bson:"author"`
  Body string `bson:"body"`
  TimeCreated time.Time `bson:"time_created"`
}
