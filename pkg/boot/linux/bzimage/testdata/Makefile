init.cpio: init
	echo init | cpio -H newc -o > init.cpio
init: init.S
	gcc -o init -static -nostdlib init.S

