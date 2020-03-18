package auth

import (
  "time"
  "go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	ID primitive.ObjectID `bson:"_id"`
	Name     string  `bson:"name"`
	Surname string `bson:"surname"`
  Image string `bson:"image"`
	Email    string `bson:"email"`
	Gender   string `bson:"gender"`
	Password string `bson:"password"`
	TimeCreated time.Time `bson:"time_created"`
}



//Token struct declaration
type Token struct {
	UserID string
	Name   string
	Email  string
	*jwt.StandardClaims
}
