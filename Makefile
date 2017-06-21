.DEFAULT: plurality

plurality:
	go build -i -o plurality .

all: plurality

deps:
	vndr

clean:
	sudo rm -rf /home/$(USER)/.runc/container/echo_test
	rm -f plurality example/echo-server
	docker rm -f echo_test 2>/dev/null || true
	docker rmi echo-server 2>/dev/null || true

example/echo-server:
	cd example && CGO_ENABLED=0 go build  -a -ldflags '-s' -o echo-server server.go

build: example/echo-server
	cd example && docker build -t echo-server .

create:
	./plurality --debug --nosudo create --nopull echo-server:latest echo_test

run:
	./plurality run echo_test /echo-server

test: plurality build create run
