all: bin/mwan3-notify-fcgi ;

bin/mwan3-notify-fcgi:
	go build -mod=vendor -v -o bin/mwan3-notify-fcgi ./cmd/mwan3-notify-fcgi

vendor:
	go mod vendor

clean:
	rm -fv bin/mwan3-notify-fcgi

.PHONY: all bin/* vendor clean
