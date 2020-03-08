package routes


import (
  "views"
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/mongo"

)

func RoutesInit(collection *mongo.Collection, r *gin.Engine){
  r.GET("/article/:article_id", views.GetArticle(collection))
  r.POST("/article", views.CreateArticle(collection))
  r.PUT("/article/:article_id", views.UpdateArticle(collection))
  r.DELETE("/article/:article_id", views.DeleteArticle(collection))
}

func SetupRouter(collection *mongo.Collection) *gin.Engine{
  router := gin.Default()
  router.GET("/article/:article_id", views.GetArticle(collection))
  // router.GET("/article", views.GetAllArticles(collection)) get all articles
  router.POST("/article", views.CreateArticle(collection))
  router.PUT("/article/:article_id", views.UpdateArticle(collection))
  router.DELETE("/article/:article_id", views.DeleteArticle(collection))
  return router
}
