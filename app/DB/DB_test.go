package DB

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestDBInstance(t *testing.T) {
	godotenv.Load("../.env")

	client := DBInstance()
	assert.NotEqual(t, client, nil)

	collection := OpenCollection(client, "task")
	assert.NotEqual(t, collection, nil)

}

func TestEstablishConnection(t *testing.T) {
	godotenv.Load("../.env")

	assert.Equal(t, collectionName, "task")

	EstablishConnection()
	assert.NotEqual(t, Collection, nil)

}
