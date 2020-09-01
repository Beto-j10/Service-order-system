package login

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "harper/database"
)

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"
var date = time.Now()

type userData struct {
	ID    int64
	Email string
	Name  string
}

// User data
type User struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

// AuthMiddleware auth middleware, used for modules that wanna be protected
var AuthMiddleware *jwt.GinJWTMiddleware

func initMiddleware() {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "harper",
		Key:         []byte("secret key"),
		Timeout:     (24 * 30) * time.Hour,
		MaxRefresh:  (24 * 30) * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*userData); ok {
				return jwt.MapClaims{
					"ID":    v.ID,
					"Email": v.Email,
					"Date":  date.String(),
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &userData{
				Email: claims["Email"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userEmail := loginVals.Email
			password := loginVals.Password

			user, err := authenticate(userEmail, password)
			if err != nil {
				// return nil, jwt.ErrFailedAuthentication
				return nil, err
			}
			c.Keys = map[string]interface{}{
				"UserID": user.ID,
				"Email":  user.Email,
				"Name":   user.Name,
			}
			return &userData{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
			}, nil

		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*userData); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			ID, ok := c.Keys["UserID"].(int64)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "error in LoginResponse, UserID not found",
				})
				return
			}
			Email, _ := c.Keys["Email"].(string)
			Name, _ := c.Keys["Name"].(string)

			c.JSON(http.StatusOK, gin.H{
				"ID":     ID,
				"Name":   Name,
				"Email":  Email,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	AuthMiddleware = authMiddleware

}

func authenticate(email string, password string) (User, error) {
	var user User

	splitEmail := strings.SplitAfter(email, "@")
	if splitEmail[1] == "harper.com" {
		result := db.Conn.Raw("SELECT * FROM technicians WHERE email = ?", email).Scan(&user)
		notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)
		if notFound {
			err := errors.New("incorrect email")
			return user, err
		}
		if password != user.Password {
			err := errors.New("incorrect password")
			return user, err
		}
	} else {
		result := db.Conn.Where("email = ?", email).First(&user)
		notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)
		if notFound {
			err := errors.New("incorrect email")
			return user, err
		}
		if password != user.Password {
			err := errors.New("incorrect password")
			return user, err
		}
	}

	return user, nil
}
