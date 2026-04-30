package db

import (
	"fmt"
	"testing"

	"github.com/huweiup/go-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Id       int64  `gorm:"primaryKey;index"`
	Name     string `gorm:"column:name;"`
	Email    string `gorm:"column:email;"`
	Password string `gorm:"column:password;"`
}

func TestNew(t *testing.T) {
	logCfg := logger.Config{
		Level:    "info",
		Encoding: "console",
	}
	err := logger.New(logCfg)
	if err != nil {
		t.Fatalf("failed to initialize logger: %v", err)
	}

	// Test SQLite in-memory
	cfg := Config{
		Driver: "mysql",
		Source: Source{
			Master: "root:123456@tcp(localhost:3306)/laravel?charset=utf8mb4&parseTime=True&loc=Local",
			Slave:  []string{"root:123456@tcp(localhost:3306)/laravel2?charset=utf8mb4&parseTime=True&loc=Local"},
		},
	}

	err = New(cfg)
	if err != nil {
		t.Fatalf("failed to initialize db: %v", err)
	}

	users := make([]User, 0)
	users = append(users, User{Id: 1, Name: "test", Email: "test@example.com", Password: "123456"})
	GetDB().Create(&users)
	fmt.Println(users)
}

func TestRepository(t *testing.T) {
	// Setup DB
	cfg := Config{
		Driver: "sqlite",
		Source: Source{
			Master: "root:123456@tcp(localhost:3306)/laravel?charset=utf8mb4&parseTime=True&loc=Local",
			Slave:  []string{"root:123456@tcp(localhost:3306)/laravel2?charset=utf8mb4&parseTime=True&loc=Local"},
		},
	}
	_ = New(cfg)
	GetDB().AutoMigrate(&User{})

	repo := NewRepository[User](GetDB())

	// Test Create
	entity := &User{Name: "test"}
	err := repo.Create(entity)
	assert.NoError(t, err)
	assert.NotZero(t, entity.Id)

	// Test FindByID
	found, err := repo.FindByID(entity.Id)
	assert.NoError(t, err)
	assert.Equal(t, entity.Name, found.Name)

	// Test Update
	entity.Name = "updated"
	err = repo.Update(entity)
	assert.NoError(t, err)

	found, _ = repo.FindByID(entity.Id)
	assert.Equal(t, "updated", found.Name)

	// Test List
	list, count, err := repo.List(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Len(t, list, 1)

	// Test Delete
	err = repo.Delete(entity.Id)
	assert.NoError(t, err)

	_, err = repo.FindByID(entity.Id)
	assert.Error(t, err)
}
