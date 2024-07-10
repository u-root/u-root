all: fork hello addr

fork: fork.c
	gcc -static fork.c -o fork

hello: hello.c
	gcc -static hello.c -o hello

addr: addr.c
	gcc -static addr.c -o addr

clean:
	rm fork hello addr
