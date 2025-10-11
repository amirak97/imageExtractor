package main

import (
	"fmt"
	"imageExtractor/internal/web"
	"testing"
)

func TestUser(t *testing.T) {
	ud := web.NewUserData()
	ud.New("1", "akdjfljsaf")

	id, err := ud.GetId("akdjfljsaf")
	if err != nil {
		t.Errorf("GetId failed: %v", err)
	} else {
		t.Logf("Got ID: %s", id)
	}
}
