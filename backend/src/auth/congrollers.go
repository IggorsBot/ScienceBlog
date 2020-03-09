package auth

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

  "golang.org/x/crypto/bcrypt"  // for crypt password
  jwt "github.com/dgrijalva/jwt-go" // for token
)

func CreateUser(client *mongo.Client) gin.HandlerFunc {
  fn := func(c *gin.Context) {

    if (len(c.PostForm("email")) == 0 || len(c.PostForm("password")) == 0){
      c.JSON(400, gin.H{"message": "not valid data",})
      return
    }
    collection := client.Database("appdb").Collection("users")


    pass, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
  	if err != nil {
      c.JSON(400, gin.H{"message": "Password Encryption  failed",})
      return
  	}
    password := string(pass)

      user := User {
      ID: primitive.NewObjectID(),
      Name: c.PostForm("name"),
      Surname: c.PostForm("surname"),
      Email: c.PostForm("email"),
      Gender: c.PostForm("gender"),
      Password: password,
      TimeCreated: time.Now(),
    }

    result, err := collection.InsertOne(context.TODO(), user)
    if err != nil {
      c.JSON(400, gin.H{"message": "can't post data to database", "body": nil})
    } else {
      c.JSON(200, gin.H{"message": "create user sucess", "body": result})
    }

  }
  return gin.HandlerFunc(fn)
}


func GetUser(client *mongo.Client) gin.HandlerFunc {
    fn := func(c *gin.Context) {
      collection := client.Database("appdb").Collection("users")

      id := c.Param("user_id")
      objID, err := primitive.ObjectIDFromHex(id)
      if err != nil {
        c.JSON(400, gin.H{"message": "error id", "body": nil})
        return
      }
      var result User
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


func GetToken(id, name, email string) (string, error) {
  expiresAt := time.Now().Add(time.Minute * 100000).Unix()
  tk := &Token{
    UserID: id,
    Name:   name,
    Email:  email,
    StandardClaims: &jwt.StandardClaims{
      ExpiresAt: expiresAt,
    },
  }
  token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

  tokenString, error := token.SignedString([]byte("secret"))
  if error != nil {
    return "", error
  }
  return tokenString, nil
}

func Login(client *mongo.Client) gin.HandlerFunc {
  fn := func(c *gin.Context) {
    if (len(c.PostForm("email")) == 0 || len(c.PostForm("password")) == 0){
      c.JSON(400, gin.H{"message": "not valid data",})
      return
    }
    collection := client.Database("appdb").Collection("users")

    var user User
    filter := bson.M{"email": c.PostForm("email")}

    err := collection.FindOne(context.TODO(), filter).Decode(&user)
    if err != nil {
      c.JSON(400, gin.H{"message": "user with this email doesn't exist", "body": nil})
      return
    }

    errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.PostForm("password")))
    if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
      c.JSON(400, gin.H{"message": "wrong password", "body": nil})
  	} else {
      token, err := GetToken(user.ID.Hex(), user.Name, user.Email)
      if err != nil {
        c.JSON(400, gin.H{"message": "error"})
      }
      c.JSON(200, gin.H{"message": "login sucess", "body": user, "token": token})
    }
  }
  return gin.HandlerFunc(fn)
}
