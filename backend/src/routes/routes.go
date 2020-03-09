package routes


import (
  "articles"
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(collection *mongo.Collection) *gin.Engine{
  router := gin.Default()
  router.GET("/api/article/:article_id", articles.GetArticle(collection))
  // router.GET("/article", views.GetAllArticles(collection)) get all articles
  router.POST("/api/article", articles.CreateArticle(collection))
  router.PUT("/api/article/:article_id", articles.UpdateArticle(collection))
  router.DELETE("/api/article/:article_id", articles.DeleteArticle(collection))

  // router.POST("/api/login", views.)
  return router
}
