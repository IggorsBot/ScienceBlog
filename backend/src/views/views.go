package views

import (
  "github.com/gin-gonic/gin"
  "log"
  "context"
  "models"
  // "fmt"
  // "reflect" // Get an object type
  "time"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID

)

func CreateArticle(collection *mongo.Collection) gin.HandlerFunc {
    fn := func(c *gin.Context) {
      article := models.Article {
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

      insertResult, err := collection.InsertOne(context.TODO(), article)
      if err != nil {
          log.Fatal(err)
      } else {
        c.JSON(200, gin.H{"result": insertResult,})
      }
    }
    return gin.HandlerFunc(fn)
}

func GetArticle(collection *mongo.Collection) gin.HandlerFunc {
    fn := func(c *gin.Context) {
      var result models.Article
      id := c.Param("article_id")
      objID, err := primitive.ObjectIDFromHex(id)

      if err != nil {
        log.Fatal(err)
      }

      filter := bson.M{"_id": objID}
      err = collection.FindOne(context.TODO(), filter).Decode(&result)

      if err != nil {
        log.Println(err)
        c.JSON(500, gin.H{"error": err.Error(),})
      } else {
        c.JSON(200, gin.H{"result": result,})
      }
    }
    return gin.HandlerFunc(fn)
  }

func UpdateArticle(collection *mongo.Collection) gin.HandlerFunc {
    fn := func(c *gin.Context) {
      id := c.Param("article_id")
      objID, err := primitive.ObjectIDFromHex(id)
      filter := bson.M{"_id": bson.M{"$eq": objID}}
      update := bson.M{"$set": bson.M{"body": "UpdateArticle"}}
      result, err := collection.UpdateOne(
              context.Background(),
              filter,
              update,
          )
    if err != nil {
      c.JSON(500, gin.H{"error": err.Error(),})
    } else {
      c.JSON(200, gin.H{"result": result,})
    }
  }
    return gin.HandlerFunc(fn)
}



func DeleteArticle(collection *mongo.Collection) gin.HandlerFunc {
    fn := func(c *gin.Context) {
      id := c.Param("article_id")
      oid, err := primitive.ObjectIDFromHex(id)

      if err != nil {
        log.Fatal("Primitive error\n", err)
      }

      filter := bson.M{"_id": oid}
      res, err := collection.DeleteOne(context.TODO(), filter)
      // fmt.Println("DeleteOne Result TYPE:", reflect.TypeOf(res))

      if err != nil {
        c.JSON(500, gin.H{"error": err.Error(),})
      } else {
        if res.DeletedCount == 0 {
          c.JSON(404, gin.H{"result": res,})
          // fmt.Println("DeleteOne() document not found:", res)
        } else {
          // fmt.Println("DeleteOne Result:", res)
          c.JSON(200, gin.H{"result": res,})
        }
      }
    }
    return gin.HandlerFunc(fn)
  }


  // article := models.Article {
  //   Title: "Test Title",
  //   Desctiption: "Test Desctiption",
  //   Images: []string{"https://naked-science.ru/article/media/v-rossii-vyyavili-eshhe-shest-chelovek-bolnyh-novym-koronavirusom", "https://naked-science.ru/wp-content/uploads/2020/03/vaccine0.jpg"},
  //   Body: "Test Body of article",
  //   Tags: []string{"вакцинация", "укол"},
  //   Category: "Медицина",
  //   TimeCreated: time.Now(),
  // }

  // insertResult, err := collection.InsertOne(context.TODO(), article)
  // if err != nil {
  //     log.Fatal(err)
  // }
