all: clean coverage

coverage:
	gocov convert coverage.out | gocov-xml > coverage.xml
	gocov convert coverage.out | gocov-html > coverage.html
	echo "total test coverage:"
	go tool cover -func=coverage.out | tail -n 1


clean:
	@rm -f coverage.*

