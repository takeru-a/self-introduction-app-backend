package configs

import(
	"context"
    "net/http"
    "github.com/go-redis/redis/v8"
    "github.com/golang-jwt/jwt"
    "github.com/gorilla/sessions"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/labstack/gommon/log"
    "github.com/rbcervilla/redisstore/v8"

)

type JWTCustomClaims struct {
    UID  string    `json:"uid"`
    Name string `json:"name"`
    jwt.StandardClaims
}

var (
	redisEndpoint = GetRedisEndpoint()
	sessionKey = GetSESSION_KEY()
	signingKey = []byte(GetSIGNING_KEY())
)

// JWTConfig https://echo.labstack.com/middleware/jwt/#configuration
var JWTConfig = middleware.JWTConfig{
    Claims:     &JWTCustomClaims{},
    SigningKey: signingKey,
}

type JWTResponse struct {
    UID  string    `json:"uid"`
    Name string    `json:"name"`
}

// ResponseMessageJSON レスポンスメッセージ
type ResponseMessageJSON struct {
    Message string `json:"message"`
}
type ResponseJTWTokenJSON struct {
    Name  string `json:"name"`
    Token string `json:"token"`
}

func getSession(c echo.Context) *sessions.Session {
    client := redis.NewClient(&redis.Options{
        Addr: redisEndpoint,
    })
    store, err := redisstore.NewRedisStore(context.Background(), client)
    if err != nil {
        log.Fatal("Failed cannot connect redis", err)
        //return err
    }
    store.KeyPrefix("session_")
    store.Options(sessions.Options{
        MaxAge:   600,
        HttpOnly: true,
    })
    session, err := store.Get(c.Request(), sessionKey)
    if err != nil {
        log.Fatal("Failed cannot get session", err)
    }
    return session
}

// Logout ログアウト処理
func Logout(c echo.Context) error {
    session := getSession(c)
    // ログアクト
    session.Values["auth"] = false
    // セッション削除
    session.Options.MaxAge = -1
    if err := sessions.Save(c.Request(), c.Response()); err != nil {
        log.Fatal("Failed cannot delete session", err)
    }
    return c.Redirect(http.StatusFound, "/") 
}

//  ログイン済みユーザに表示する
func ShowPage(c echo.Context) error {
    session := getSession(c)
    // ログイン確認
    if session.Values["auth"] != true {
        return c.String(http.StatusUnauthorized, "401")
    } else {
        return c.String(http.StatusOK, session.Values["username"].(string))
    }
}

// ShowData JWT認証の確認用
func ShowData(c echo.Context) error {
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*JWTCustomClaims)
    responseJSON := JWTResponse{
        UID:  claims.UID,
        Name: claims.Name,
    }
    return c.JSON(http.StatusOK, responseJSON)
}