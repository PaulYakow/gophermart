test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

run_accrual:
	./cmd/accrual/accrual.exe