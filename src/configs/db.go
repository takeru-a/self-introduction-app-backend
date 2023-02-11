package configs

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
	"github.com/takeru-a/self-introduction-app-backend/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    "github.com/golang-jwt/jwt"
    "github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"

)

type DB struct {
    client *mongo.Client
}

func ConnectDB() *DB  {
    client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
    if err != nil {
        log.Fatal(err)
    }

    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    //ping the database 応答確認
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB")
    return &DB{client: client}
}


func GetCollection(db *DB, collectionName string) *mongo.Collection {
    collection := db.client.Database("game").Collection(collectionName)
    return collection
}

func GenerateToken(n int) string {
    var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

func (db *DB) CreateRoom(input *model.NewRoom, eC echo.Context) (*model.Room, error){
    RoomCollection := GetCollection(db, "room")
    UserCollection := GetCollection(db, "user")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    host := &model.User{
        ID: primitive.NewObjectID().Hex(),
		Name: input.HostName,
        Score: 0,
        Answer: "",
	}
	room := &model.Room{
        ID: primitive.NewObjectID().Hex(),
		Host: host,
		Token: GenerateToken(32),
		Players: []*model.User{host},
	}
    _, err := RoomCollection.InsertOne(ctx, room)
    if err != nil {
        return nil, err
    }
    _, err = UserCollection.InsertOne(ctx, host)
    if err != nil {
        return nil, err
    }
    session := getSession(eC)
    session.Values["username"] = host.Name
    session.Values["auth"] = true                                       // ログイン有無の確認用
    if err := sessions.Save(eC.Request(), eC.Response()); err != nil {  // Session情報の保存
        log.Fatal("Failed save session", err)
    }
    claims := JWTCustomClaims{
		host.ID,
		host.Name,
		jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString(signingKey)
    if err != nil{
        return nil, err
    }
    session.Values["token"] = t
    return room, err
}

func (db *DB) CreateUser(input *model.NewUser, eC echo.Context) (*model.User, error) {
    collection := GetCollection(db, "user")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    user := &model.User{
        ID: primitive.NewObjectID().Hex(),
        Name: input.Name,
        Score: 0,
        Answer: "",
    }
    _, err := collection.InsertOne(ctx, user)

    if err != nil {
        return nil, err
    }
    session := getSession(eC)
    session.Values["username"] = user.Name
    session.Values["auth"] = true                                       // ログイン有無の確認用
    if err := sessions.Save(eC.Request(), eC.Response()); err != nil {  // Session情報の保存
        log.Fatal("Failed save session", err)
    }
    claims := &JWTCustomClaims{
        user.ID,
        user.Name,
        jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
    t, err := token.SignedString(signingKey)
    if err != nil{
        return nil, err
    }
    session.Values["token"] = t
    return user, err
}

func (db *DB) GetUses() ([]*model.User, error) {
    collection := GetCollection(db, "user")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var users []*model.User
    defer cancel()

    res, err := collection.Find(ctx, bson.M{})

    if err != nil {
        return nil, err
    }

    defer res.Close(ctx)
    for res.Next(ctx) {
        var user *model.User
        if err = res.Decode(&user); err != nil {
            log.Fatal(err)
        }
        users = append(users, user)
    }

    return users, err
}

func (db *DB) GetRooms() ([]*model.Room, error) {
    collection := GetCollection(db, "room")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var rooms []*model.Room
    defer cancel()

    res, err := collection.Find(ctx, bson.M{})

    if err != nil {
        return nil, err
    }

    defer res.Close(ctx)
    for res.Next(ctx) {
        var room *model.Room
        if err = res.Decode(&room); err != nil {
            log.Fatal(err)
        }
        rooms = append(rooms, room)
    }

    return rooms, err
}

func (db *DB) SingleUser(ID string) (*model.User, error) {
    collection := GetCollection(db, "user")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var user *model.User
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(ID)

    err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)

    return user, err
}

func (db *DB) SingleRoom(ID string) (*model.Room, error) {
    collection := GetCollection(db, "room")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var room *model.Room
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(ID)

    err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&room)

    return room, err
}