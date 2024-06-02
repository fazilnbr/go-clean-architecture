package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	service    Service
	middleware MiddlewareIf
}

func NewHttpHandler(service Service, middleware MiddlewareIf) HttpHandler {
	return HttpHandler{
		service:    service,
		middleware: middleware,
	}
}

func (h HttpHandler) InitRoutes(router *gin.RouterGroup) {
	router.POST("/signup", h.SignUp)
	router.POST("/login", h.Login)
	router.Use(h.middleware.ValidateToken()).PUT("/", h.Update)
	router.Use(h.middleware.ValidateToken()).DELETE("/", h.Delete)
}

func (h HttpHandler) SignUp(c *gin.Context) {
	var req User
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	err = h.service.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user signup failed with " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
}

func (h HttpHandler) Login(c *gin.Context) {
	var req User
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	user, err := h.service.Login(req.UserName, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user login failed with " + err.Error()})
		return
	}
	token, err := h.middleware.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create token failed with " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user login successfully", "token": token})
}

func (h HttpHandler) Update(c *gin.Context) {
	user_id, err := strconv.Atoi(c.GetHeader(USER_ID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid input"})
		return
	}
	var req User
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	req.ID = uint(user_id)
	err = h.service.UpdateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user login failed with " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user update successfully"})
}

func (h HttpHandler) Delete(c *gin.Context) {
	user_id, err := strconv.Atoi(c.GetHeader(USER_ID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid input"})
		return
	}
	err = h.service.DeleteUser(uint(user_id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user login failed with " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user delete successfully"})
}
