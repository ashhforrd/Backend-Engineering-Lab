# Table-Driven Tests

## Problem

Writing one test function per input creates repetitive setup and assertion code.
It also makes boundary cases easy to overlook.

## Design

Test cases are represented as data:

```go
tests := []struct {
    name    string
    request Request
    want    Quote
    wantErr error
}{
    // cases
}
```

A shared test loop runs every case as a named subtest.

## Implementation

The project tests a shipping-cost calculator with:

- valid zones and service levels;
- weight boundaries;
- rounded weight surcharges;
- express service charges;
- expected validation errors;
- parallel subtests;
- race detection and coverage reporting.

Each case contains a descriptive name, input, expected quote, and expected
sentinel error.

## Failure Cases

- Zero weight exposed a `< 0` versus `<= 0` boundary bug.
- Missing boundary cases can leave incorrect comparisons undetected.
- Parallel subtests must not mutate shared state.
- Coverage shows executed statements, not correctness.
- Comparing error strings is less reliable than `errors.Is`.
- A failing subtest should identify its input and expected result clearly.

## What I Learned

- Table-driven tests separate test data from test execution.
- `t.Run` creates independently named subtests.
- `t.Parallel` runs independent cases concurrently.
- `t.Fatalf` stops the current test immediately.
- `t.Errorf` records a failure and allows execution to continue.
- `errors.Is` supports sentinel and wrapped errors.
- Boundary values are often more valuable than random happy paths.
- Coverage and the race detector answer different questions.

## Running

Run all tests:

```bash
go test -v ./...
```

Run with the race detector:

```bash
go test -v -race ./...
```

Run one case:

```bash
go test -v \
  -run 'TestCalculate/partial_extra_kilogram_rounds_up' \
  ./...
```

Generate coverage:

```bash
go test \
  -coverprofile=/tmp/table-driven-tests-coverage.out \
  ./...

go tool cover \
  -func=/tmp/table-driven-tests-coverage.out
```