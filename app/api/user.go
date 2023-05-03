package api

import (
	"net/http"

	"a21hc3NpZ25tZW50/app/model"
	"a21hc3NpZ25tZW50/app/service"

	"github.com/gin-gonic/gin"
)

type UserAPI interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	Logout(c *gin.Context)

	GetByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	Delete(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

// Authenticate User godoc
// @Summary Authenticate User
// @Schemes https
// @Description Auth Endpoint
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} Ok
// @Param request body model.UserLogin true "user login by email and password"
// @Router /users/login [post]
func (u *userAPI) Login(c *gin.Context) {
	// TODO: answer here
}

func (u *userAPI) Logout(c *gin.Context) {
	// TODO: answer here
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	var recordUser = model.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	recordUser, err := u.userService.Register(c.Request.Context(), &recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(recordUser.ID, "register success"))
}

func (u *userAPI) Delete(c *gin.Context) {
	userId, exist := c.Get("id")

	if exist != true {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("error unauthorized user id"))
		return
	}

	err := u.userService.Delete(c.Request.Context(), userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(userId.(int), "delete success"))
}

func (u *userAPI) GetByID(c *gin.Context) {
	userId, exist := c.Get("id")
	if exist != true {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("error unauthorized user id"))
		return
	}

	res, err := u.userService.GetByID(c.Request.Context(), userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *userAPI) UpdateByID(c *gin.Context) {
	var user model.UserUpdate

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	userId, exist := c.Get("id")
	if exist != true {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("error unauthorized user id"))
		return
	}

	var recordUpdateUser = model.User{
		ID:       userId.(int),
		Fullname: user.Fullname,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Image:    user.Image,
	}

	err := u.userService.UpdateByID(c.Request.Context(), &recordUpdateUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(userId.(int), "update success"))
}
