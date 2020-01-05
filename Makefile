test:
	go test ./...

cover-web:
	@go test -coverprofile=cover.prof ./...
	@go tool cover -html=cover.prof
	@rm cover.prof

clear:
	@rm cover.prof

bench-old:
	make BENCH_FILE=bench.old bench-delta

bench-new:
	make BENCH_FILE=bench.new bench-delta

benchcmp:
	benchcmp bench.old bench.new

benchstat:
	benchstat bench.old bench.new

BENCH_FILE=bench.old
BENCH_COUNT=5
bench-delta:
	@rm -f ${BENCH_FILE}
	for i in {1..${BENCH_COUNT}}; do /bin/echo -n "$$i.."; go test -bench=. >> ${BENCH_FILE}; done
	@/bin/echo "Done"
