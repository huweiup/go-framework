package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	Run:   runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("module", "m", "", "Module name (defaults to project name)")
}

type ProjectData struct {
	Name       string
	ModuleName string
}

func runNew(cmd *cobra.Command, args []string) {
	projectName := args[0]
	moduleName, _ := cmd.Flags().GetString("module")
	if moduleName == "" {
		moduleName = projectName
	}

	data := ProjectData{
		Name:       projectName,
		ModuleName: moduleName,
	}

	fmt.Printf("Creating project %s (%s)...\n", projectName, moduleName)

	if err := createProjectStructure(projectName, data); err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project created successfully!")
	fmt.Printf("cd %s\n", projectName)
	fmt.Println("go mod tidy")
	fmt.Println("go run cmd/server/main.go")
}

func createProjectStructure(root string, data ProjectData) error {
	dirs := []string{
		"cmd/server",
		"configs",
		"internal/api/v1",
		"internal/models",
		"internal/service",
		"pkg",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(root, dir), 0755); err != nil {
			return err
		}
	}

	files := map[string]string{
		"go.mod":                  goModTemplate,
		"configs/config.yaml":     configTemplate,
		"cmd/server/main.go":      mainTemplate,
		"internal/api/v1/user.go": userHandlerTemplate,
		"internal/models/user.go": userModelTemplate,
	}

	for path, tmplContent := range files {
		fullPath := filepath.Join(root, path)
		tmpl, err := template.New(path).Parse(tmplContent)
		if err != nil {
			return err
		}

		f, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := tmpl.Execute(f, data); err != nil {
			return err
		}
	}

	return nil
}

const goModTemplate = `module {{.ModuleName}}

go 1.21

require (
	github.com/project/go-framework v0.0.0
)
`

const configTemplate = `app:
  name: {{.Name}}
  version: 1.0.0

server:
  port: 8080
  mode: debug

log:
  level: info
  encoding: console

db:
  driver: mysql
  source: ***:***@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
`

const mainTemplate = `package main

import (
	"log"

	"github.com/project/go-framework/pkg/app"
	"{{.ModuleName}}/internal/api/v1"
	"{{.ModuleName}}/internal/models"
)

func main() {
	a, err := app.New("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Migrate db
	if err := a.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}

	// Register routes
	v1.RegisterRoutes(a.Server.Engine, a.DB, a.Logger)

	if err := a.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
`

const userModelTemplate = `package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string ` + "`json:\"username\"`" + `
	Email    string ` + "`json:\"email\"`" + `
}
`

const userHandlerTemplate = `package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/project/go-framework/pkg/db"
	"{{.ModuleName}}/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserHandler struct {
	Repo   *db.Repository[models.User]
	Logger *zap.Logger
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB, logger *zap.Logger) {
	h := &UserHandler{
		Repo:   db.NewRepository[models.User](db),
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
`
