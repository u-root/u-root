all: fork hello sleep

fork: fork.c
	gcc -static fork.c -o fork

hello: hello.c
	gcc -static hello.c -o hello

sleep: sleep.c
	gcc -static sleep.c -o sleep

clean:
	rm fork hello sleep
