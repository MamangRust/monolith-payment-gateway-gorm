package outbox

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeExecutor struct {
	calls int
	query string
	args  []any
}

func (f *fakeExecutor) Exec(_ context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	f.calls++
	f.query = query
	f.args = args
	return pgconn.NewCommandTag("INSERT 0 1"), nil
}

func (*fakeExecutor) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, errors.New("query not used by enqueue tests")
}

func (*fakeExecutor) QueryRow(context.Context, string, ...any) pgx.Row {
	return nil
}

func TestEnqueueValidatesAndPersistsEvent(t *testing.T) {
	exec := &fakeExecutor{}
	event := Event{EventKey: "topup:42", Topic: "topup.created", Key: "42", Payload: []byte(`{"id":42}`)}
	if err := Enqueue(context.Background(), exec, event); err != nil {
		t.Fatalf("Enqueue returned error: %v", err)
	}
	if exec.calls != 1 || len(exec.args) != 4 {
		t.Fatalf("expected one four-argument insert, calls=%d args=%d", exec.calls, len(exec.args))
	}
	if !contains(exec.query, "ON CONFLICT (event_key) DO NOTHING") {
		t.Fatalf("enqueue is not idempotent: %s", exec.query)
	}
}

func TestEnqueueRejectsIncompleteEvent(t *testing.T) {
	for name, event := range map[string]Event{
		"missing key":     {Topic: "topic", Payload: []byte("x")},
		"missing topic":   {EventKey: "key", Payload: []byte("x")},
		"missing payload": {EventKey: "key", Topic: "topic"},
		"invalid json":    {EventKey: "key", Topic: "topic", Payload: []byte("not-json")},
	} {
		t.Run(name, func(t *testing.T) {
			if err := Enqueue(context.Background(), &fakeExecutor{}, event); !errors.Is(err, ErrInvalidEvent) {
				t.Fatalf("expected ErrInvalidEvent, got %v", err)
			}
		})
	}
}

func contains(value, needle string) bool {
	for i := 0; i+len(needle) <= len(value); i++ {
		if value[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
