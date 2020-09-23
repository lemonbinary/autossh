default:
	go build -v -o assh ./*.go

clean:
	rm -f assh *.o

install:
	make
	cp -f assh /usr/local/bin/