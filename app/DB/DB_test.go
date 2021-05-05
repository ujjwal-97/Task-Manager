package DB

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestDBInstance(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}

	client := DBInstance()

	if client == nil {
		t.Error()
	}
	collection := OpenCollection(client, "task")
	if collection == nil {
		t.Error()
	}
}

func TestEstablishConnection(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}
	if collectionName != "task" {
		t.Error()
	}

	EstablishConnection()
	if Collection == nil {
		t.Error()
	}
}
