package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	main "a21hc3NpZ25tZW50"
	"a21hc3NpZ25tZW50/app/middleware"
	"a21hc3NpZ25tZW50/app/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ = os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/kampusmerdeka")

func SetCookie(mux *gin.Engine) *http.Cookie {
	login := model.UserLogin{
		Email:    "test@mail.com",
		Password: "testing123",
	}

	body, _ := json.Marshal(login)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(w, r)

	var cookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "session_token" {
			cookie = c
		}
	}

	return cookie
}

var _ = Describe("TestAPIHandler", Ordered, func() {
	var apiServer *gin.Engine
	var db *gorm.DB
	var userTest int

	BeforeAll(func() {
		conn, err := gorm.Open(postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        os.Getenv("DATABASE_URL"),
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		db = conn
		db.Exec("DROP TABLE IF EXISTS tweets CASCADE")
		db.Exec("DROP TABLE IF EXISTS users CASCADE")

		db.AutoMigrate(model.User{}, model.Tweet{})

		apiServer = gin.New()
		apiServer = main.RunServer(db, apiServer)

		reqRegister := model.UserRegister{
			Fullname: "test",
			Email:    "test@mail.com",
			Password: "testing123",
		}

		reqBody, _ := json.Marshal(reqRegister)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/users/register", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		apiServer.ServeHTTP(w, r)

		var resp = map[string]interface{}{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		Expect(err).To(BeNil())
		Expect(w.Result().StatusCode).To(Equal(http.StatusCreated))

		userTest = int(resp["user_id"].(float64))
	})

	AfterAll(func() {
		ctx := context.Background()

		err := db.WithContext(ctx).Exec("DELETE FROM users WHERE id = ?", userTest).Error
		if err != nil {
			panic(err)
		}

		err = db.WithContext(ctx).Exec("DELETE FROM tweets WHERE user_id = ?", userTest).Error
		if err != nil {
			panic(err)
		}
	})

	//  ==============================================
	//  ===============     USERS     ================
	//  ==============================================

	Describe("Auth Middleware", func() {
		var (
			router *gin.Engine
			w      *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			router = gin.Default()
			w = httptest.NewRecorder()
		})

		When("valid token is provided", func() {
			It("should set user ID in context and call next middleware", func() {
				// Prepare request with valid session token
				claims := &model.Claims{UserID: 123}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString(model.JwtKey)
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{Name: "session_token", Value: signedToken})

				// Attach Auth middleware to the request
				router.Use(middleware.Auth())
				router.GET("/", func(ctx *gin.Context) {
					userID := ctx.MustGet("id").(int)
					Expect(userID).To(Equal(123))
				})

				// Perform the request and assert response
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})

		When("session token is missing", func() {
			It("should return unauthorized error response", func() {
				// Prepare request with no session token
				req, _ := http.NewRequest(http.MethodGet, "/", nil)

				// Attach Auth middleware to the request
				router.Use(middleware.Auth())

				// Perform the request and assert response
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusSeeOther))
			})
		})

		When("session token is invalid", func() {
			It("should return unauthorized error response", func() {
				// Prepare request with invalid session token
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{Name: "session_token", Value: "invalid_token"})

				// Attach Auth middleware to the request
				router.Use(middleware.Auth())

				// Perform the request and assert response
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("/users/login", func() {
		When("send empty email and password with POST method", func() {
			It("should return a bad request", func() {
				loginData := model.UserLogin{
					Email:    "",
					Password: "",
				}

				body, _ := json.Marshal(loginData)
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewReader(body))
				r.Header.Set("Content-Type", "application/json")

				apiServer.ServeHTTP(w, r)

				errResp := model.ErrorResponse{}
				err := json.Unmarshal(w.Body.Bytes(), &errResp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusBadRequest))
				Expect(errResp.Error).To(Equal("invalid decode json"))
			})
		})

		When("send email and password with POST method", func() {
			It("should return a success", func() {
				loginData := model.UserLogin{
					Email:    "test@mail.com",
					Password: "testing123",
				}

				body, _ := json.Marshal(loginData)
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewReader(body))
				r.Header.Set("Content-Type", "application/json")

				apiServer.ServeHTTP(w, r)

				var resp = map[string]interface{}{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
				Expect(resp["message"]).To(Equal("login success"))
			})
		})
	})

	Describe("/users/logout", func() {
		When("hit endpoint without user login", func() {
			It("should return an error unauthorized", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/api/v1/users/logout", nil)
				r.Header.Set("Content-Type", "application/json")

				apiServer.ServeHTTP(w, r)

				var errResp = model.ErrorResponse{}
				err := json.NewDecoder(w.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(errResp.Error).To(Equal("error unauthorized user id"))
			})
		})

		When("hit endpoint with GET method", func() {
			It("should return a success", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/api/v1/users/logout", nil)

				r.AddCookie(SetCookie(apiServer))
				apiServer.ServeHTTP(w, r)

				var resp = model.SuccessResponse{}
				err := json.NewDecoder(w.Body).Decode(&resp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Message).To(Equal("logout success"))
			})
		})
	})

	//  ==============================================
	//  ==============     TWEETS     ================
	//  ==============================================
	Describe("/tweets/create", func() {
		When("hit endpoint without user login", func() {
			It("should return an error unauthorized", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/tweets/create", nil)
				r.Header.Set("Content-Type", "application/json")
				apiServer.ServeHTTP(w, r)

				var errResp = model.ErrorResponse{}
				err := json.NewDecoder(w.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(errResp.Error).To(Equal("error unauthorized user id"))
			})
		})

		When("hit endpoint with POST method and invalid JSON body", func() {
			It("should return an error invalid decode json", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/tweets/create", strings.NewReader("invalid json body"))
				r.Header.Set("Content-Type", "application/json")
				r.AddCookie(SetCookie(apiServer))

				apiServer.ServeHTTP(w, r)

				var errResp = model.ErrorResponse{}
				err := json.NewDecoder(w.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusBadRequest))
				Expect(errResp.Error).To(Equal("invalid decode json"))
			})
		})

		When("hit endpoint with POST method and valid JSON body", func() {
			It("should return a success", func() {
				w := httptest.NewRecorder()
				payload := model.UserTweet{
					Tweet: "Hallo Dunia",
					Image: "world image",
				}
				payloadBytes, _ := json.Marshal(payload)

				r := httptest.NewRequest("POST", "/api/v1/tweets/create", bytes.NewReader(payloadBytes))
				r.Header.Set("Content-Type", "application/json")
				r.AddCookie(SetCookie(apiServer))

				apiServer.ServeHTTP(w, r)

				var resp = model.SuccessResponse{}
				err := json.NewDecoder(w.Body).Decode(&resp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Message).To(Equal("tweet created"))
			})
		})
	})

	Describe("/tweets/get", func() {
		When("hit endpoint without user login", func() {
			It("should return an error unauthorized", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/api/v1/tweets/get", nil)
				r.Header.Set("Content-Type", "application/json")
				apiServer.ServeHTTP(w, r)

				var errResp = model.ErrorResponse{}
				err := json.NewDecoder(w.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(errResp.Error).To(Equal("error unauthorized user id"))
			})
		})

		When("hit endpoint with GET method", func() {
			It("should return a success", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/api/v1/tweets/get", nil)

				r.AddCookie(SetCookie(apiServer))
				apiServer.ServeHTTP(w, r)

				var resp = []model.Tweet{}
				err := json.NewDecoder(w.Body).Decode(&resp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
				Expect(resp[0].Tweet).To(Equal("Hallo Dunia"))
				Expect(resp[0].Image).To(Equal("world image"))
			})
		})
	})

	Describe("/tweets/update", func() {
		When("hit endpoint without user login", func() {
			It("should return an error unauthorized", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("PUT", "/api/v1/tweets/update", nil)
				r.Header.Set("Content-Type", "application/json")
				apiServer.ServeHTTP(w, r)

				var errResp = model.ErrorResponse{}
				err := json.NewDecoder(w.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(errResp.Error).To(Equal("error unauthorized user id"))
			})
		})

		When("hit endpoint with PUT method", func() {
			It("should return a success", func() {
				tweet := model.UserTweet{
					Tweet: "Hallo Dunia",
					Image: "world image",
				}
				payload, err := json.Marshal(tweet)
				Expect(err).To(BeNil())

				w := httptest.NewRecorder()
				r := httptest.NewRequest("PUT", "/api/v1/tweets/update?tweet_id=1", bytes.NewBuffer(payload))
				r.AddCookie(SetCookie(apiServer))
				r.Header.Set("Content-Type", "application/json")
				apiServer.ServeHTTP(w, r)

				var resp = model.SuccessResponse{}
				err = json.NewDecoder(w.Body).Decode(&resp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Message).To(Equal("tweet updated"))
				Expect(resp.UserID).To(Equal(1))
			})
		})
	})

	Describe("/tweets/delete", func() {
		When("hit endpoint without user login", func() {
			It("should return an error unauthorized", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/api/v1/tweets/delete", nil)
				r.Header.Set("Content-Type", "application/json")
				apiServer.ServeHTTP(w, r)

				var errResp = model.ErrorResponse{}
				err := json.NewDecoder(w.Body).Decode(&errResp)
				Expect(err).To(BeNil())
				Expect(w.Result().StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(errResp.Error).To(Equal("error unauthorized user id"))
			})
			When("hit endpoint with DELETE method", func() {
				It("should return a success", func() {
					w := httptest.NewRecorder()
					r := httptest.NewRequest("DELETE", "/api/v1/tweets/delete?tweet_id=1", nil)

					r.AddCookie(SetCookie(apiServer))
					apiServer.ServeHTTP(w, r)

					var resp = model.SuccessResponse{}
					err := json.NewDecoder(w.Body).Decode(&resp)
					Expect(err).To(BeNil())
					Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
					Expect(resp.UserID).To(Equal(1))
					Expect(resp.Message).To(Equal("tweet deleted"))
				})
			})
		})

	})
})
