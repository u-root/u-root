# {docker build -t yarikk/systemboot-test-image -f Dockerfile.base .}
FROM uroottest/test-image-amd64:v3.2.4

# Install dependencies
RUN sudo apt-get update &&                          \
	sudo apt-get install -y --no-install-recommends \
		`# tools for creating bootable disk images` \
		gdisk \
		e2fsprogs \
		qemu-utils \
		&& \
	sudo rm -rf /var/lib/apt/lists/*

# Get u-root
RUN go get github.com/u-root/u-root

# Get Linux kernel
# 
# Config taken from:
#	curl -s https://raw.githubusercontent.com/linuxboot/demo/master/20190203-FOSDEM-barberio-hendricks/config/linux-config |
#	sed \
#		-e '/^# CONFIG_RELOCATABLE / s!.*!CONFIG_RELOCATABLE=y!' `# for kexec` \
#		-e '/^CONFIG_INITRAMFS_SOURCE=/ s!^!#!' \
#		> linux-config
# 	
COPY linux-config .
RUN set -x; \
	git clone -q --depth 1 -b v4.19 https://github.com/torvalds/linux.git && \
	mv linux-config linux/.config && \
	(cd linux/ && exec make -j$(nproc)) && \
	cp linux/arch/x86/boot/bzImage bzImage && \
	rm -r linux/

# Create a bootable disk image to test localboot; the init there simply shuts down.
RUN set -x; \
	mkdir rootfs && \
	cp bzImage rootfs/ && \
	u-root -build=bb -o rootfs/ramfs.cpio -initcmd shutdown  && \
	xz --check=crc32 --lzma2=dict=512KiB rootfs/ramfs.cpio && \
	{ \
		echo menuentry; \
		echo linux bzImage; \
		echo initrd ramfs.cpio.xz; \
	} > rootfs/grub2.cfg && \
	du -a rootfs/ && \
	qemu-img create -f raw disk.img 20m && \
	sgdisk --clear --new 1::-0 --typecode=1:8300 --change-name=1:'Linux root filesystem' \
		disk.img && \
	mkfs.ext2 -F -E 'offset=1048576' -d rootfs/ disk.img 18m && \
	gdisk -l disk.img && \
	qemu-img convert -f raw -O qcow2 disk.img disk.qcow2 && \
	mv disk.qcow2 disk.img && \
	rm -r rootfs/
