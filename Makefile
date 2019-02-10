test:
	go test ./...

cover-web:
	@go test -coverprofile=cover.prof ./...
	@go tool cover -html=cover.prof
	@rm cover.prof

clear:
	@rm cover.prof
