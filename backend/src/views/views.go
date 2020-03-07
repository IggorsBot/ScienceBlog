package views

import (
  "github.com/gin-gonic/gin"
  "log"
  "context"
  "models"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
)


func GetArticle(collection *mongo.Collection) gin.HandlerFunc {
    fn := func(c *gin.Context) {
        // Your handler code goes in here - e.g.
        var result models.Article
        filter := bson.D{{"title", "Test Title"}}
        err := collection.FindOne(context.TODO(), filter).Decode(&result)

        if err != nil {
          log.Println(err)
          c.JSON(500, gin.H{"error": err.Error(),})
        } else {
          c.JSON(200, gin.H{"result": result,})
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
