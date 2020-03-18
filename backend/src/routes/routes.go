package routes


import (
  "articles"
  "auth"
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/mongo"
  // "fmt"
)

func AuthMiddleware(c *gin.Context) {
  // Check authorization
  c.Next()
}

func SetupRouter(client *mongo.Client) *gin.Engine{
  router := gin.Default()
  router.GET("/api/article/:article_id", articles.GetArticle(client))
  // router.GET("/article", views.GetAllArticles(collection)) get all articles
  router.POST("/api/article", articles.CreateArticle(client))
  router.PUT("/api/article/:article_id", articles.UpdateArticle(client))
  router.DELETE("/api/article/:article_id", articles.DeleteArticle(client))

  router.GET("/api/user/:user_id", auth.GetUser(client))
  router.POST("/api/user/login", auth.Login(client))
  router.POST("/api/user/registration", auth.CreateUser(client))

  authorized := router.Group("/")
  authorized.Use(AuthMiddleware)
  {
    authorized.POST("/api/comment/create", articles.CreateComment(client))
    authorized.GET("api/comment/:article_id", articles.GetComments(client))
    // authorized.PUT("api/comment/:comment_id", articles.UpdateComment(client))
    authorized.DELETE("api/comment", articles.DeleteComment(client))
  }
  return router
}
