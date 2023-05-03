package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Result struct {
	Posts []Post `json:"posts"`
}

var Posts = []Post{
	{ID: 1, Title: "Judul Postingan Pertamax", Content: "Ini adalah postingan pertama di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Judul Postingan Kedua", Content: "Ini adalah postingan kedua di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

var jwtKey = []byte("secret-key")

type Claims struct {
	Username  string `json:"username"`
	Role      string `json:"role"`
	CompanyID string `json:"company_id"`
	jwt.StandardClaims
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/posts", func(ctx *gin.Context) {
		result := Result{
			Posts: Posts,
		}
		ctx.JSON(200, result)
	})

	router.GET("/posts/:id", func(ctx *gin.Context) {

		id := ctx.Param("id")
		for _, post := range Posts {
			if fmt.Sprintf("%d", post.ID) == id {
				ctx.JSON(200, post)
				return
			}
		}
		ctx.JSON(404, gin.H{"message": "post not found"})
	})

	router.POST("/posts", func(ctx *gin.Context) {
		var post Post
		if err := ctx.ShouldBindJSON(&post); err != nil {
			ctx.JSON(400, gin.H{"message": err.Error()})
			return
		}
		Posts = append(Posts, post)
		ctx.JSON(201, post)
	})

	router.PUT("/posts/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var postPayload Post
		if err := ctx.ShouldBindJSON(&postPayload); err != nil {
			ctx.JSON(400, gin.H{"message": err.Error()})
			return
		}
		for i, post := range Posts {
			if fmt.Sprintf("%d", post.ID) == id {
				Posts[i].Title = postPayload.Title
				Posts[i].Content = postPayload.Content
				Posts[i].UpdatedAt = time.Now()
				ctx.JSON(200, Posts[i])
				return
			}
		}
		ctx.JSON(404, gin.H{"message": "post not found"})
	})

	router.DELETE("/posts/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		for i, post := range Posts {
			if fmt.Sprintf("%d", post.ID) == id {
				Posts = append(Posts[:i], Posts[i+1:]...)
				ctx.JSON(200, gin.H{"message": "post deleted"})
				return
			}
		}
		ctx.JSON(404, gin.H{"message": "post not found"})
	})

	router.GET("/token", func(ctx *gin.Context) {
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Username:  "john",
			Role:      "admin123",
			CompanyID: "1",
			StandardClaims: jwt.StandardClaims{
				// expiry time menggunakan time millisecond
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(jwtKey)
		ctx.JSON(200, gin.H{
			"token":   tokenString,
			"expired": expirationTime,
		})
	})

	router.POST("/token", func(ctx *gin.Context) {
		token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			ctx.JSON(401, gin.H{"error": err.Error()})
		} else {
			fmt.Println(tkn.Valid, tkn.Claims)
			ctx.JSON(200, claims)
		}

	})

	return router
}

func main() {
	r := SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
