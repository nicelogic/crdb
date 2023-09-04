package crdb

import (
	"context"
	"testing"
)

func TestInit(t *testing.T) {
	client := Client{}
	ctx := context.Background()
	err := client.Init(ctx, "/Users/bryan.wu/code/secret/config-crdb.yml", "contacts", 4)
	if err != nil {
		t.Errorf("cient init error: %v\n", err)
	}
}
