package articles

import (
  "github.com/gin-gonic/gin"
  // "log"
  "context"
  // "fmt"
  // "reflect" // reflect.TypeOf(data)
  "time"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
)

func CreateArticle(client *mongo.Client) gin.HandlerFunc {
    fn := func(c *gin.Context) {
      collection := client.Database("appdb").Collection("articles")

      if (len(c.PostForm("author")) == 0 || len(c.PostForm("body")) == 0 || len(c.PostForm("title")) == 0){
        c.JSON(400, gin.H{"message": "not valid data",})
        return
      }
      article := Article {
        ID: primitive.NewObjectID(),
        Author: c.PostForm("author"),
        Title: c.PostForm("title"),
        Slug: c.PostForm("slug"),
        Desctiption: c.PostForm("description"),
        Images: c.PostFormArray("images"),
        Videos: c.PostFormArray("videos"),
        Body: c.PostForm("body"),
        Tags: c.PostFormArray("tags"),
        Category: c.PostForm("category"),
        TimeCreated: time.Now(),
      }

      result, err := collection.InsertOne(context.TODO(), article)
      if err != nil {
        c.JSON(400, gin.H{"message": "can't post data to database", "body": nil})
      } else {
        c.JSON(200, gin.H{"message": "create article sucess", "body": result})
      }
    }
    return gin.HandlerFunc(fn)
}

func GetArticle(client *mongo.Client) gin.HandlerFunc {
    fn := func(c *gin.Context) {
      collection := client.Database("appdb").Collection("articles")

      id := c.Param("article_id")
      objID, err := primitive.ObjectIDFromHex(id)
      if err != nil {
        c.JSON(400, gin.H{"message": "error id", "body": nil})
        return
      }
        var result Article
        filter := bson.M{"_id": objID}
        err = collection.FindOne(context.TODO(), filter).Decode(&result)

        if err != nil {
          c.JSON(400, gin.H{"message": "can't get data from database", "body": nil})
        } else {
          c.JSON(200, gin.H{"message": "get data sucess", "body": result})
        }
    }
    return gin.HandlerFunc(fn)
}

func UpdateArticle(client *mongo.Client) gin.HandlerFunc {
    fn := func(c *gin.Context) {
      collection := client.Database("appdb").Collection("articles")
      id := c.Param("article_id")
      objID, err := primitive.ObjectIDFromHex(id)

      if err != nil {
        c.JSON(400, gin.H{"message": "error id", "body": nil})
        return
      }

      query := bson.M{}
      c.PostForm("author")
      for key, _ := range c.Request.PostForm {
            query[key] = c.PostForm(key)
        }

      filter := bson.M{"_id": bson.M{"$eq": objID}}
      update := bson.M{"$set": query}

      result, err := collection.UpdateOne(
              context.Background(),
              filter,
              update,
            )
      if err != nil{
        c.JSON(200, gin.H{"message": "ca't update data", "body": nil})
      } else {
        c.JSON(200, gin.H{"message": "update data sucess", "body": result})
      }
    }
    return gin.HandlerFunc(fn)
  }



func DeleteArticle(client *mongo.Client) gin.HandlerFunc {
    fn := func(c *gin.Context) {
      collection := client.Database("appdb").Collection("articles")

      id := c.Param("article_id")
      objID, err := primitive.ObjectIDFromHex(id)

      if err != nil {
        c.JSON(400, gin.H{"message": "error id", "body": nil})
        return
      }

      filter := bson.M{"_id": objID}
      result, err := collection.DeleteOne(context.TODO(), filter)

      if err != nil {
        c.JSON(400, gin.H{"message": "can't delete", "body": nil})
      } else {
        if result.DeletedCount == 0 {
          c.JSON(404, gin.H{"message": "data not found", "body": result})
        } else {
          c.JSON(200, gin.H{"message": "delete data sucess", "body": result})
        }
      }
    }
    return gin.HandlerFunc(fn)
  }
