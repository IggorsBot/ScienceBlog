package routes


import (
  "views"
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/mongo"

)

func RoutesInit(collection *mongo.Collection, r *gin.Engine){
  r.GET("/test", views.GetArticle(collection))
  // r.POST("/test", views.GetArticle(collection))
  // r.PUT("/test", views.GetArticle(collection))
  // r.DELETE("/test", views.GetArticle(collection))
}
