set -ex
CGO_ENABLED=0 GOARCH=arm go build .
mkdir -p usb
sudo mount /dev/sda1 usb
sudo mv spidev usb
sudo umount usb
rmdir usb
