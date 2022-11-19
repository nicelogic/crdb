package crdb

import (
	"context"
	"testing"
	"time"
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

//bulk insert
//https://github.com/jackc/pgx/issues/764#issuecomment-685249471
func TestExec(t *testing.T) {
	client := Client{}
	ctx := context.Background()
	err := client.Init(ctx, "/Users/bryan.wu/code/secret/config-crdb.yml", "contacts", 4)
	if err != nil {
		t.Errorf("cient init error: %v\n", err)
	}
	const UpsertAddContactsApply = `
	UPSERT INTO add_contacts_apply (
		user_id,
		contacts_id,
		message,
		update_time
	)
	VALUES ($1, $2, $3, $4)
	`
	updateTime := time.Now().Format(time.RFC3339)
	err = client.Exec(ctx, UpsertAddContactsApply, "4", "2", "please  me", updateTime)
	if err != nil{
		t.Errorf("exec error: %v\n", err)
	}
}
