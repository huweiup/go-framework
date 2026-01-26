package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type TestModel struct {
	gorm.Model
	Name string
}

func TestNew(t *testing.T) {
	// Test SQLite in-memory
	cfg := Config{
		Driver: "sqlite",
		Source: ":memory:",
	}

	db, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test unsupported driver
	cfg.Driver = "postgres"
	db, err = New(cfg)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestRepository(t *testing.T) {
	// Setup DB
	cfg := Config{
		Driver: "sqlite",
		Source: ":memory:",
	}
	db, _ := New(cfg)
	db.AutoMigrate(&TestModel{})

	repo := NewRepository[TestModel](db)

	// Test Create
	entity := &TestModel{Name: "test"}
	err := repo.Create(entity)
	assert.NoError(t, err)
	assert.NotZero(t, entity.ID)

	// Test FindByID
	found, err := repo.FindByID(entity.ID)
	assert.NoError(t, err)
	assert.Equal(t, entity.Name, found.Name)

	// Test Update
	entity.Name = "updated"
	err = repo.Update(entity)
	assert.NoError(t, err)

	found, _ = repo.FindByID(entity.ID)
	assert.Equal(t, "updated", found.Name)

	// Test List
	list, count, err := repo.List(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Len(t, list, 1)

	// Test Delete
	err = repo.Delete(entity.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(entity.ID)
	assert.Error(t, err)
}
