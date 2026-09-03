---
name: unit-test-pattern
description: Add table-driven unit tests with gomock mocks to Go microservice service layers
source: auto-skill
extracted_at: '2026-06-28T05:11:57.637Z'
---

# Unit Test Pattern for Go Microservice Service Layers

## When to use

When a Go microservice (in `service/<name>/`) has interfaces with `//go:generate mockgen` directives but no fast unit tests — only slow integration tests with `testcontainers-go`. Also use when you need to add CI test execution, race detection, and coverage.

## Goal

Create fast (< 1ms), isolated, race-checked unit tests covering both happy path and every error path, so developers can validate business logic without spinning up PostgreSQL/Redis containers (~30s per suite).

## Procedure

### Step 1 — Export deps structs (if needed for external tests)

Before writing tests, check if the service's constructor takes an **unexported** deps struct (e.g. `cardCommandServiceDeps`). If so, rename it to an **exported** name (e.g. `CardCommandServiceDeps`) so external test packages can use it. This is a backward-compatible change:

```bash
# Example: export cardCommandServiceDeps → CardCommandServiceDeps
sed -i 's/cardCommandServiceDeps/CardCommandServiceDeps/g' service/card/service/*.go
sed -i 's/cardQueryServiceDeps/CardQueryServiceDeps/g' service/card/service/*.go
```

Update all references in `service.go` / `services.go` wiring files too. Verify with `go build ./service/<name>/...`.

**Why:** This lets you place tests in `tests/<service>/unit/` (external package, project convention) instead of inside the service module. No more test files polluting the production module.

### Step 2 — Ensure mocks are generated

For every interface file (`service/<name>/{service,repository,redis}/interfaces.go`), check for a `//go:generate mockgen` directive. If missing, add one after the `package` line:

```go
package repository


```

Then add `go.uber.org/mock` to the service's `go.mod`, and run:

```bash
cd service/<name>
go get go.uber.org/mock@v0.6.0
go mod tidy
go generate ./...
```

Also add `go.uber.org/mock` to the tests module if not present:

```bash
cd tests && go get go.uber.org/mock@v0.6.0
```

### Step 3 — Create fake logger + observability

In the test file, add no-op implementations of `logger.LoggerInterface` and `observability.TraceLoggerObservability`. These prevent OpenTelemetry span boilerplate from polluting test logic. Place them at the top of the test file (or in a shared helper if the service has multiple test files).

Pattern for fakes:

```go
type fakeLogger struct{}
func (f *fakeLogger) Info(_ string, _ ...zap.Field)             {}
func (f *fakeLogger) Fatal(_ string, _ ...zap.Field)            {}
func (f *fakeLogger) Debug(_ string, _ ...zap.Field)            {}
func (f *fakeLogger) Error(_ string, _ ...zap.Field)            {}
func (f *fakeLogger) Warn(_ string, _ ...zap.Field)             {}
func (f *fakeLogger) Check(_ zapcore.Level, _ string) *zapcore.CheckedEntry { return nil }
func (f *fakeLogger) With(_ ...zap.Field) logger.LoggerInterface { return f }
func (f *fakeLogger) Sync() error                                 { return nil }

type fakeObservability struct{}
func (f *fakeObservability) StartTracingAndLogging(
    ctx context.Context, _ string, _ ...attribute.KeyValue,
) (context.Context, trace.Span, func(string), string, func(string, ...zap.Field)) {
    _, span := trace.NewNoopTracerProvider().Tracer("test").Start(ctx, "test")
    end := func(_ string) {}
    status := "success"
    logSuccess := func(_ string, _ ...zap.Field) {}
    return ctx, span, end, status, logSuccess
}
func (f *fakeObservability) RecordMetrics(_ context.Context, _, _ string, _ time.Time) {}
```

### Step 4 — Write factory functions

For each service interface under test, write a `new<Name>Service(t)` function. This goes in the test file, not the production code:

```go
type commandMocks struct {
    userRepo *mock_repo.MockUserRepository
    cardCmd  *mock_repo.MockCardCommandRepository
    cache    *mock_cache.MockCardCommandCache
}

func newCommandService(t *testing.T) (service.CardCommandService, *commandMocks) {
    t.Helper()
    ctrl := gomock.NewController(t)
    m := &commandMocks{
        userRepo: mock_repo.NewMockUserRepository(ctrl),
        cardCmd:  mock_repo.NewMockCardCommandRepository(ctrl),
        cache:    mock_cache.NewMockCardCommandCache(ctrl),
    }
    svc := service.NewCardCommandService(&service.CardCommandServiceDeps{
        Cache:                 m.cache,
        Kafka:                 nil,
        UserRepository:        m.userRepo,
        CardCommandRepository: m.cardCmd,
        Logger:                &fakeLogger{},
        Observability:         &fakeObservability{},
    })
    return svc, m
}
```

The factory returns both the service interface and a mocks struct so each test case can set up expectations inline. This keeps the test scope explicit and eliminates shared state.

### Step 5 — Write table-driven tests

Each exported function gets its own test function with a `tests []struct` table. Place test files at `tests/<service>/unit/<name>_test.go` with package `<service>_test`:

```go
func TestCardCommandService_CreateCard(t *testing.T) {
    t.Parallel()

    validReq := &requests.CreateCardRequest{...}
    createdRow := &db.CreateCardRow{
        CardID:     1,
        ExpireDate: pgtype.Date{Time: fixedExpiry(t), Valid: true},
        Cvv:        "123",
        ...
    }

    tests := []struct {
        name    string
        req     *requests.CreateCardRequest
        mock    func(m *commandMocks)
        wantErr bool
    }{
        {
            name: "success_without_kafka",
            req:  validReq,
            mock: func(m *commandMocks) {
                m.userRepo.EXPECT().FindById(gomock.Any(), 1).
                    Return(&db.GetUserByIDRow{UserID: 1}, nil)
                m.cardCmd.EXPECT().CreateCard(gomock.Any(), validReq).
                    Return(createdRow, nil)
            },
            wantErr: false,
        },
        {
            name: "error_user_not_found",
            req:  validReq,
            mock: func(m *commandMocks) {
                m.userRepo.EXPECT().FindById(gomock.Any(), 1).
                    Return(nil, errors.New("user not found"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc, m := newCommandService(t)
            tt.mock(m)
            got, err := svc.CreateCard(context.Background(), tt.req)
            if tt.wantErr {
                if err == nil {
                    t.Fatal("expected error but got nil")
                }
                return
            }
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if got == nil || got.CardID != createdRow.CardID {
                t.Errorf("unexpected result: %+v", got)
            }
        })
    }
}
```

### Step 6 — Cover every code path

For each method, test:
- **Happy path** — all mocks return success, verify cache invalidation
- **Each error path** — one test per dependency that can fail
- **Edge cases** — empty results, nil pointers, zero values, pagination defaults
- **Bulk operations** — RestoreAll, DeleteAll (they often differ from single-item ops)
- **Complex multi-step operations** — Test the entire flow end-to-end in one test (e.g. transaction Create which does: merchant lookup → card lookup → balance check → saldo deduction → txn creation → status update → merchant balance credit). For rollback scenarios, verify both the original action AND the rollback mocks are called.

For cache-based query services, test:
- **Cache hit** — verify repo is NOT called
- **Cache miss + repo success** — verify repo IS called and cache IS populated
- **Repo error** — verify error is propagated

### Step 7 — Handle pgtype types in test data

SQLC-generated models use `pgtype.Date`, `pgtype.Timestamp`, etc. Create them explicitly:

```go
expireDate := pgtype.Date{Time: time.Date(2030, 12, 31, 0, 0, 0, 0, time.UTC), Valid: true}
createdAt := pgtype.Timestamp{Time: time.Now(), Valid: true}
```

Use a helper function to avoid repetition:

```go
func fixedExpiry(t *testing.T) time.Time {
    t.Helper()
    return time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC)
}
```

### Step 8 — Run with race detection

```bash
cd tests && go test -race -count=1 -short ./<service>/unit/...
```

Race detection adds ~1s per package but catches real data races in production code.

### Step 9 — CI integration

Add a `test` job to `.github/workflows/build_and_push.yaml` that runs **before** the `docker` job:

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: testdb
          POSTGRES_USER: testuser
          POSTGRES_PASSWORD: testpass
        ports:
          - 5432:5432
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
          cache: true
      - run: go install go.uber.org/mock/mockgen@v0.6.0
      - run: for svc in auth card merchant ...; do (cd service/$svc && go generate ./...); done
      - run: cd tests && go test -race -count=1 -short -coverprofile=coverage.out ./<service>/unit/...
      - uses: actions/upload-artifact@v4

  docker:
    needs: [test]
    ...
```

Use `needs: [test]` so the Docker build only runs if tests pass. This prevents broken code from being deployed.

## Important constraints

- **Prefer external tests** (`package <service>_test` in `tests/<service>/unit/`). First export the deps struct so the constructor is accessible. Only use internal `package service` tests when the deps struct absolutely cannot be exported (rare).
- **Do NOT mock the logger or observability** — use no-op fakes. They add noise without value and the real implementations require OpenTelemetry setup.
- **`gomock.Any()` is acceptable** for `context.Context` and for pointer arguments whose exact value is irrelevant to the assertion.
- **For paginated queries**, the mock must expect the **original unmodified request**, not the normalised values. The production code normalises `page`/`pageSize` into **local variables**, leaving the struct unchanged.
- **For mock error types**, use simple `errors.New("...")` strings. The production code wraps these into typed `AppError` instances via the errorhandler, so the test just needs to verify an error was returned, not the specific error type.
- **`TotalCount` fields** on paginated SQLC rows are `int64`. Set them explicitly: `TotalCount: 1`. Zero values cause `*total == 0` assertions to fail.
- **`t.Parallel()`** is safe because each test function creates its own `gomock.NewController(t)` and mock instances. Never share mocks across test functions.
- **Coverage measurement** in a multi-module workspace is noisy. The unit tests live in the `tests/` module but cover code in other modules. Use coverage as a trend indicator, not a gate.
