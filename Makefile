test:
	go test ./...

cover-web:
	@go test -coverprofile=cover.prof ./...
	@go tool cover -html=cover.prof
	@rm cover.prof

clear:
	@rm cover.prof

benchcmp:
	go test -benchmem -bench=. | tee  bench.new
	@git stash --quiet
	go test -benchmem -bench=. | tee  bench.old
	@git stash pop --quiet
	benchcmp bench.old bench.new
	@rm bench.new bench.old
