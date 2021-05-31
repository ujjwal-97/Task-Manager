package db

import (
	"app/db"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestDBInstance(t *testing.T) {
	godotenv.Load("../.env")

	client := db.DBInstance()
	assert.NotEqual(t, client, nil)

	collection := db.OpenCollection(client, "task")
	assert.NotEqual(t, collection, nil)

}

func TestEstablishConnection(t *testing.T) {
	godotenv.Load("../.env")

	assert.Equal(t, db.CollectionName, "task")

	db.EstablishConnection()
	assert.NotEqual(t, db.Collection, nil)

}
