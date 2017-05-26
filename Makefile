.DEFAULT: plurality

plurality:
	go build -i -o plurality .

all: plurailty

deps:
	vndr

clean:
	rm -f plurality
