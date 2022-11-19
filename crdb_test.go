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

func TestQuery(t *testing.T) {
	client := Client{}
	ctx := context.Background()
	err := client.Init(ctx, "/Users/bryan.wu/code/secret/config-crdb.yml", "contacts", 4)
	if err != nil {
		t.Errorf("cient init error: %v\n", err)
	}

	const AddContactsApply = `
	SELECT user_id,
		message,
		update_time
	from add_contacts_apply
	where contacts_id = $1
	ORDER BY update_time DESC
	`
	_, err = client.Query(ctx, AddContactsApply, "2")
	if err != nil{
		t.Errorf("query error: %v\n", err)
	}
}

