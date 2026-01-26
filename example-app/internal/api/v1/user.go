package v1

import (
	"net/http"

	"example-app/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/project/go-framework/pkg/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserHandler struct {
	Repo   *database.Repository[models.User]
	Logger *zap.Logger
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	h := &UserHandler{
		Repo:   database.NewRepository[models.User](db),
		Logger: logger,
	}

	g := r.Group("/api/v1/users")
	{
		g.POST("/", h.Create)
		g.GET("/:id", h.Get)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Create(&user); err != nil {
		h.Logger.Error("failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
