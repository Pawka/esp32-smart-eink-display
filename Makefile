
.PHONY: test
test:
	go test ./...

.PHONY: mocks
mocks:
	mockgen -source=lib/clock.go -package=mocks -destination mocks/mock_clock.go
