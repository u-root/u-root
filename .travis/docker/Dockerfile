# {docker build -t systemboottest -f Dockerfile ../..}
FROM yarikk/systemboot-test-image

COPY . /go/src/github.com/systemboot/systemboot

RUN set -x; \
	sudo chmod -R a+w /go/src && \
	cd /go/src/github.com/systemboot/systemboot && \
	go get -v ./...  && \
	u-root -build=bb core uinit localboot # netboot

CMD ./qemu-system-x86_64 \
	-M q35 \
	-L pc-bios/ `# for vga option rom` \
	-kernel bzImage -initrd /tmp/initramfs.linux_amd64.cpio \
	-m 1024 \
	-nographic \
	-append 'console=ttyS0 earlyprintk=ttyS0' \
	-object 'rng-random,filename=/dev/urandom,id=rng0' \
	-device 'virtio-rng-pci,rng=rng0' \
	-hda disk.img
