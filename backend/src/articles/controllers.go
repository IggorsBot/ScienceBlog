package articles

import (
  "github.com/gin-gonic/gin"
  "log"
  "context"
  // "fmt"
  // "reflect" // reflect.TypeOf(data)
  "time"

  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
)

// ARTICLES

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

      // Valid data

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

      // Удалить все комментарии связанные со статьей

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


// COMMENTS

func GetComments(client *mongo.Client) gin.HandlerFunc {
    fn := func(c *gin.Context) {
      collection := client.Database("appdb").Collection("comments")

      articleID, err := primitive.ObjectIDFromHex(c.Param("article_id"))
      if err != nil {
        c.JSON(400, gin.H{"message": "error id", "body": nil})
        return
      }

      findOptions := options.Find()

      var comments []*Comment

      cur, err := collection.Find(context.TODO(), bson.M{"article": articleID}, findOptions)
      if err != nil {
        log.Fatal(err)
      }

      for cur.Next(context.TODO()){
        var comment Comment
        err := cur.Decode(&comment)
        if err != nil {
          c.JSON(400, gin.H{"message": "cur error",})
        }
        comments = append(comments, &comment)
      }

      if err := cur.Err(); err != nil {
        c.JSON(400, gin.H{"message": "cur error",})
      }

      cur.Close(context.TODO())
      c.JSON(200, gin.H{"comments:": comments})

    }
    return gin.HandlerFunc(fn)
}

func CreateComment(client *mongo.Client) gin.HandlerFunc {
  fn := func(c *gin.Context) {

    if (len(c.PostForm("body")) == 0 || len(c.PostForm("author")) == 0 || len(c.PostForm("article")) == 0){
      c.JSON(400, gin.H{"message": "not valid data",})
      return
    }

    articleID, err := primitive.ObjectIDFromHex(c.PostForm("article"))
    if err != nil {
      c.JSON(400, gin.H{"message": "error id", "body": nil})
      return
    }

    // Находим статью в БД
    collection := client.Database("appdb").Collection("articles")
    var article Article
    filter := bson.M{"_id": articleID}
    err = collection.FindOne(context.TODO(), filter).Decode(&article)

    if err != nil {
      c.JSON(400, gin.H{"message": "can't get article from database", "body": nil})
      return
    }
    //

    // Добавляем комментарий в БД
    comment := Comment {
      ID: primitive.NewObjectID(),
      Author: c.PostForm("author"),
      Article: articleID,
      Body: c.PostForm("body"),
      TimeCreated: time.Now(),
    }

    collection = client.Database("appdb").Collection("comments")
    _, err = collection.InsertOne(context.TODO(), comment)
    if err != nil {
      c.JSON(400, gin.H{"message": "can't post data to database", "body": nil})
      return
    }
    //


    // В объекте статьи, в массив Comments добавляем новый коментарий
    collection = client.Database("appdb").Collection("articles")
    article.Comments = append(article.Comments, comment.ID)
    filter = bson.M{"_id": bson.M{"$eq": articleID}}
    update := bson.M{"$set": bson.M{"comments": article.Comments}}

    resultArticle, err := collection.UpdateOne(
            context.Background(),
            filter,
            update,
          )
    //
    if err != nil{
      c.JSON(400, gin.H{"message": "can't update article-comments", "body": nil})
      } else {
        c.JSON(200, gin.H{"message": "create comment sucess", "body": resultArticle})
      }
    }
    return gin.HandlerFunc(fn)
  }


func DeleteComment(client *mongo.Client) gin.HandlerFunc {
  fn := func(c *gin.Context) {

    articleID, err := primitive.ObjectIDFromHex(c.PostForm("article_id"))
    commentID, err := primitive.ObjectIDFromHex(c.PostForm("comment_id"))

    if err != nil {
      c.JSON(400, gin.H{"message": "error id", "body": nil})
      return
    }

    collection := client.Database("appdb").Collection("articles")
    var article Article
    filter := bson.M{"_id": articleID}
    err = collection.FindOne(context.TODO(), filter).Decode(&article)

    if err != nil {
      c.JSON(400, gin.H{"message": "can't get article from database", "body": nil})
      return
    }

    // Удаляем комментарий из collection articles
    var comments []primitive.ObjectID
    for i := 0; i < len(article.Comments); i++{
      if (article.Comments[i] != commentID){
        comments = append(comments, article.Comments[i])
      }
    }

    filter = bson.M{"_id": bson.M{"$eq": articleID}}
    update := bson.M{"$set": bson.M{"comments": comments}}

    result, err := collection.UpdateOne(
            context.Background(),
            filter,
            update,
          )
    if err != nil{
      c.JSON(400, gin.H{"message": "can't update article-comments", "body": nil})
      return
    }
    //

    // Удаляем комментарий из collection comments
    collection = client.Database("appdb").Collection("comments")
    filter = bson.M{"_id": commentID}
    commentResult, err := collection.DeleteOne(context.TODO(), filter)

    if err != nil {
      c.JSON(400, gin.H{"message": "can't delete comment", "body": nil})
      return
      } else {
        if commentResult.DeletedCount == 0 {
          c.JSON(404, gin.H{"message": "comment doesn't exist", "body": commentResult})
          return
          }
    }
    c.JSON(200, gin.H{"message": "delete comment sucess", "body": result})

  }
  return gin.HandlerFunc(fn)
  }
