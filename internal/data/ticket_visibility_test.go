package data

import (
	"context"
	"io"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/npanel-dev/NPanel-backend/ent"
	"github.com/npanel-dev/NPanel-backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
)

const closedTicketStatus int8 = 4

func TestAdminTicketListIncludesClosedByDefault(t *testing.T) {
	ctx := context.Background()
	client := enttest.Open(t, "sqlite3", "file:admin_ticket_closed_visible?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	createTicketRow(t, client, 1, 10, 1)
	createTicketRow(t, client, 2, 10, closedTicketStatus)

	repo := NewTicketRepo(&Data{db: client}, log.NewStdLogger(io.Discard))
	total, list, err := repo.GetTicketList(ctx, 1, 10, 0, nil, "")
	if err != nil {
		t.Fatalf("get admin ticket list: %v", err)
	}

	if total != 2 || len(list) != 2 {
		t.Fatalf("expected default admin list to include closed tickets, total=%d len=%d", total, len(list))
	}
}

func TestPublicTicketListIncludesClosedByDefault(t *testing.T) {
	ctx := context.Background()
	client := enttest.Open(t, "sqlite3", "file:public_ticket_closed_visible?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	createTicketRow(t, client, 1, 10, 1)
	createTicketRow(t, client, 2, 10, closedTicketStatus)
	createTicketRow(t, client, 3, 11, closedTicketStatus)

	repo := NewPublicTicketRepo(&Data{db: client}, log.NewStdLogger(io.Discard))
	total, list, err := repo.GetTicketList(ctx, 10, 1, 10, nil, nil)
	if err != nil {
		t.Fatalf("get public ticket list: %v", err)
	}

	if total != 2 || len(list) != 2 {
		t.Fatalf("expected default public list to include current user's closed tickets, total=%d len=%d", total, len(list))
	}
}

func createTicketRow(t *testing.T, client *ent.Client, id int64, userID int64, status int8) {
	t.Helper()
	if err := client.ProxyTicket.Create().
		SetID(id).
		SetUserID(userID).
		SetTitle("ticket").
		SetDescription("description").
		SetStatus(status).
		Exec(context.Background()); err != nil {
		t.Fatalf("create ticket row: %v", err)
	}
}
