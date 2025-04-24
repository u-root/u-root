// Copyright 2018-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/qnetwork"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

func tftpVM(t *testing.T, name, script string, net *qnetwork.InterVM, mods ...uimage.Modifier) *qemu.VM {
	fixedMods := []uimage.Modifier{
		uimage.WithBusyboxCommands(
			"github.com/u-root/u-root/cmds/core/cat",
			"github.com/u-root/u-root/cmds/core/echo",
			"github.com/u-root/u-root/cmds/core/grep",
			"github.com/u-root/u-root/cmds/core/ip",
			"github.com/u-root/u-root/cmds/core/mkdir",
			"github.com/u-root/u-root/cmds/core/shasum",
			"github.com/u-root/u-root/cmds/core/sleep",
		),
		uimage.WithCoveredCommands(
			"github.com/u-root/u-root/cmds/exp/tftp",
			"github.com/u-root/u-root/cmds/exp/tftpd",
		),
	}

	return scriptvm.Start(t, name, script,
		scriptvm.WithUimage(append(fixedMods, mods...)...),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			net.NewVM(),
		),
	)
}

func TestTFTPTransfer(t *testing.T) {
	net := qnetwork.NewInterVM()

	// Prepare test files for VMs
	testFilesForServer := uimage.WithFiles(
		"testdata/tftp/small.txt:/small.txt",
		"testdata/tftp/large.txt:/large.txt",
		"testdata/tftp/binary.dat:/binary.dat",
		"testdata/tftp/file_a.txt:/file_a.txt",
		"testdata/tftp/file_b.txt:/file_b.txt",
		"testdata/tftp/file_c.txt:/file_c.txt",
	)

	testFilesForClient := uimage.WithFiles(
		"testdata/tftp/small.txt:/upload_small.txt",
		"testdata/tftp/large.txt:/upload_large.txt",
		"testdata/tftp/binary.dat:/upload_binary.dat",
		"testdata/tftp/file_a.txt:/upload_file_a.txt",
		"testdata/tftp/file_b.txt:/upload_file_b.txt",
		"testdata/tftp/file_c.txt:/upload_file_c.txt",
	)

	serverScript := `#!/bin/sh
# Server script for TFTP testing
# Exit on any error
set -e

# Setup networking
ip addr add 192.168.0.2/24 dev eth0
ip -6 addr add fd51:3681:1eb4::2/126 dev eth0
ip link set eth0 up
ip route add 0.0.0.0/0 dev eth0
ip -6 route add ::/0 dev eth0
echo "192.168.0.1 tftp_client" >>/etc/hosts
echo "192.168.0.2 tftp_server" >>/etc/hosts

# Make sure test files exist
if [ ! -f /small.txt ] || [ ! -f /large.txt ] || [ ! -f /binary.dat ] || [ ! -f /file_a.txt ] || [ ! -f /file_b.txt ] || [ ! -f /file_c.txt ]; then
    echo "ERROR: Test files not found"
    exit 1
fi

# Create directory for multi file upload test
mkdir -p /multi
if [ ! -d /multi ]; then
	echo "ERROR: Failed to create upload directory"
	exit 1
fi

# Start TFTP server in root directory
echo "Starting TFTP server..."
tftpd -port 69 -root / -v &
SERVER_PID=$!
sleep 2
echo "TFTP server started with PID $SERVER_PID"

# Wait for client to complete uploads (60 seconds max)
sleep 60

# Verify uploaded files
echo "Verifying uploaded files from client..."
FAILED=0

# Check small file
if [ -f /upload_small.txt ]; then
    echo "upload_small.txt was received"
    shasum /upload_small.txt > /upload_small.txt.sha1
    if grep -q "9470c442585479bfad86d9a0a0daf01779c020b3" /upload_small.txt.sha1; then
        echo "upload_small.txt hash is correct"
    else
        echo "ERROR: upload_small.txt hash mismatch"
        echo "Expected: 9470c442585479bfad86d9a0a0daf01779c020b3"
        echo "Got: $(cat /upload_small.txt.sha1)"
        FAILED=1
    fi
else
    echo "ERROR: upload_small.txt was not received"
    FAILED=1
fi

# Check renamed small file
if [ -f /upload_renamed.txt ]; then
    echo "upload_renamed.txt was received"
    shasum /upload_renamed.txt > /upload_renamed.txt.sha1
    if grep -q "9470c442585479bfad86d9a0a0daf01779c020b3" /upload_renamed.txt.sha1; then
        echo "upload_renamed.txt hash is correct"
    else
        echo "ERROR: upload_renamed.txt hash mismatch"
        echo "Expected: 9470c442585479bfad86d9a0a0daf01779c020b3"
        echo "Got: $(cat /upupload_renamed.txt.sha1)"
        FAILED=1
    fi
else
    echo "ERROR: upload_renamed.txt was not received"
    FAILED=1
fi

# Check multi-uploaded small file
if [ -f /multi/upload_file_a.txt ] || [ -f multi/upload_file_b.txt ] || [ -f /multi/upload_file_c.txt ]; then
    echo "multi/upload_file_a.txt, multi/upload_file_b.txt and multi/upload_file_c.txt were received"
    shasum multi/upload_file_a.txt > /upload_file_a.txt.sha1
	shasum multi/upload_file_b.txt > /upload_file_b.txt.sha1
	shasum multi/upload_file_c.txt > /upload_file_c.txt.sha1

    if grep -q "0f8735a9dfe23fcd2195db61ca8f96f3150c5f9b" /upload_file_a.txt.sha1 && grep -q "a7ac17275300cb68ca88742cc87904c290128557" /upload_file_b.txt.sha1 && grep -q "f99c3d0a167bb79fb9c6f7c38f85dd67b28c0558" /upload_file_c.txt.sha1; then
        echo "multi/upload_file_a.txt, multi/upload_file_b.txt and multi/upload_file_c.txt  hashes are correct"
    else
        echo "ERROR: hash mismatch"
        echo "Expected: 0f8735a9dfe23fcd2195db61ca8f96f3150c5f9b"
		echo "Expected: a7ac17275300cb68ca88742cc87904c290128557"
		echo "Expected: f99c3d0a167bb79fb9c6f7c38f85dd67b28c0558"
        echo "Got: $(cat /upload_file_a.txt.sha1)"
		echo "Got: $(cat /upload_file_b.txt.sha1)"
		echo "Got: $(cat /upload_file_c.txt.sha1)"
        FAILED=1
    fi
else
    echo "ERROR: multi/upload_file_a.txt, multi/upload_file_b.txt or multi/upload_file_c.txt was not received"
    FAILED=1
fi

# Check large file
if [ -f /upload_large.txt ]; then
    echo "upload_large.txt was received"
    shasum /upload_large.txt > /upload_large.txt.sha1
    if grep -q "379278b59fbb7ce439a17344c95be5afd335d540" /upload_large.txt.sha1; then
        echo "upload_large.txt hash is correct"
    else
        echo "ERROR: upload_large.txt hash mismatch"
        echo "Expected: 379278b59fbb7ce439a17344c95be5afd335d540"
        echo "Got: $(cat /upload_large.txt.sha1)"
        FAILED=1
    fi
else
    echo "ERROR: upload_large.txt was not received"
    FAILED=1
fi

# Check binary file
if [ -f /upload_binary.dat ]; then
    echo "upload_binary.dat was received"
    shasum /upload_binary.dat > /upload_binary.dat.sha1
    if grep -q "b359ebb365d1abc5a927518dcdd42fc00f012523" /upload_binary.dat.sha1; then
        echo "upload_binary.dat hash is correct"
    else
        echo "ERROR: upload_binary.dat hash mismatch"
        echo "Expected: b359ebb365d1abc5a927518dcdd42fc00f012523" 
        echo "Got: $(cat /upload_binary.dat.sha1)"
        FAILED=1
    fi
else
    echo "ERROR: upload_binary.dat was not received"
    FAILED=1
fi

if [ $FAILED -eq 1 ]; then
    echo "ERROR: Server test failed"
    exit 1
fi

# Clean up
#kill $SERVER_PID #not needed to kill tftpd background process. VM is killed by scriptvm.
#wait

echo "TESTS PASSED MARKER"
`

	clientScript := `#!/bin/sh
# Client script for TFTP testing
# Exit on any error
set -e

# Setup networking
ip addr add 192.168.0.1/24 dev eth0
ip -6 addr add fd51:3681:1eb4::1/126 dev eth0
ip link set eth0 up
ip route add 0.0.0.0/0 dev eth0
ip -6 route add ::/0 dev eth0
echo "192.168.0.1 tftp_client" >>/etc/hosts
echo "192.168.0.2 tftp_server" >>/etc/hosts

# Make sure test files exist
if [ ! -f /upload_small.txt ] || [ ! -f /upload_large.txt ] || [ ! -f /upload_binary.dat ] || [ ! -f /upload_file_a.txt ] || [ ! -f /upload_file_b.txt ] || [ ! -f /upload_file_c.txt ]; then
    echo "ERROR: Upload test files not found"
    exit 1
fi

# Wait for server to start
sleep 5
echo "Starting TFTP client tests..."

# Track failures
FAILED=0

# Test 1: Download small text file (IPv4)
echo "Test 1: Downloading small.txt (IPv4)..."
if ! tftp 192.168.0.2 -c get small.txt small_download.txt; then
    echo "ERROR: Test 1 - Download failed"
    FAILED=1
else
    if [ ! -f small_download.txt ]; then
    	echo "ERROR: Test 1 - File not created"
    	FAILED=1
	else
    	shasum small_download.txt > small_download.txt.sha1
    
    	if grep -q "9470c442585479bfad86d9a0a0daf01779c020b3" small_download.txt.sha1; then
        	echo "Test 1: PASS - Hash matches"
    	else
        	echo "ERROR: Test 1 - Hash mismatch"
        	echo "Expected: 9470c442585479bfad86d9a0a0daf01779c020b3"
        	echo "Got: $(cat small_download.txt.sha1)"
        	FAILED=1
    	fi
	fi
fi

# Test 1a: Download small text file without renaming (IPv4)
echo "Test 1a: Download small text file without renaming (IPv4)..."
if ! tftp 192.168.0.2 -c get small.txt; then
    echo "ERROR: Test 1 - Download failed"
    FAILED=1
else
    if [ ! -f small.txt ]; then
    	echo "ERROR: Test 1a - File not created"
    	FAILED=1
	else
    	shasum small.txt > small.txt.sha1
    
    	if grep -q "9470c442585479bfad86d9a0a0daf01779c020b3" small.txt.sha1; then
        	echo "Test 1a: PASS - Hash matches"
    	else
        	echo "ERROR: Test 1a - Hash mismatch"
        	echo "Expected: 9470c442585479bfad86d9a0a0daf01779c020b3"
        	echo "Got: $(cat small.txt.sha1)"
        	FAILED=1
    	fi
	fi
fi

# Test 1b: Download multiple small text files (IPv4)
echo "Test 1b: Download multiple small text files (IPv4)..."
if ! tftp 192.168.0.2 -c get file_a.txt file_b.txt file_c.txt; then
    echo "ERROR: Test 1b - Download failed"
    FAILED=1
else
    if [ ! -f /upload_file_a.txt ] || [ ! -f /upload_file_b.txt ] || [ ! -f /upload_file_c.txt ]; then
    	echo "ERROR: Test 1b - Files not created"
    	FAILED=1
	else
    	shasum file_a.txt > file_a.txt.sha1
		shasum file_b.txt > file_b.txt.sha1
		shasum file_c.txt > file_c.txt.sha1
    
    	if grep -q "0f8735a9dfe23fcd2195db61ca8f96f3150c5f9b" file_a.txt.sha1 && grep -q "a7ac17275300cb68ca88742cc87904c290128557" file_b.txt.sha1 grep -q "f99c3d0a167bb79fb9c6f7c38f85dd67b28c0558" file_c.txt.sha1; then
        	echo "Test 1b: PASS - Hash matches"
    	else
        	echo "ERROR: Test 1b - Hash mismatch"
        	echo "Expected: 0f8735a9dfe23fcd2195db61ca8f96f3150c5f9b"
			echo "Expected: a7ac17275300cb68ca88742cc87904c290128557"
			echo "Expected: f99c3d0a167bb79fb9c6f7c38f85dd67b28c0558"
        	echo "Got: $(cat file_a.txt.sha1)"
			echo "Got: $(cat file_b.txt.sha1)"
			echo "Got: $(cat file_c.txt.sha1)"
        	FAILED=1
    	fi
	fi
fi

# Test 2: Download large text file (IPv4)
echo "Test 2: Downloading large.txt (IPv4)..."
if ! tftp 192.168.0.2 -c get large.txt large_download.txt; then
    echo "ERROR: Test 2 - Download failed"
    FAILED=1
else
    if [ ! -f large_download.txt ]; then
        echo "ERROR: Test 2 - File not created"
        FAILED=1
    else
        shasum large_download.txt > large_download.txt.sha1
    
    	if grep -q "379278b59fbb7ce439a17344c95be5afd335d540" large_download.txt.sha1; then
        	echo "Test 2: PASS - Hash matches"
    	else
        	echo "ERROR: Test 2 - Hash mismatch"
        	echo "Expected: 379278b59fbb7ce439a17344c95be5afd335d540"
        	echo "Got: $(cat large_download.txt.sha1)"
        	FAILED=1
    	fi
    fi
fi

# Test 3: Download binary file (IPv4)
echo "Test 3: Downloading binary.dat (IPv4)..."
if ! tftp 192.168.0.2 -m binary -c get binary.dat binary_download.dat; then
    echo "ERROR: Test 3 - Download failed"
    FAILED=1
else
    if [ ! -f binary_download.dat ]; then
        echo "ERROR: Test 3 - File not created"
        FAILED=1
    else
        shasum binary_download.dat > binary_download.dat.sha1
    
    	if grep -q "b359ebb365d1abc5a927518dcdd42fc00f012523" binary_download.dat.sha1; then
        	echo "Test 3: PASS - Hash matches"
    	else
        	echo "ERROR: Test 3 - Hash mismatch"
        	echo "Expected: b359ebb365d1abc5a927518dcdd42fc00f012523"
        	echo "Got: $(cat binary_download.dat.sha1)"
        	FAILED=1
    	fi
    fi
fi

# Test 4: Upload small text file (IPv4)
echo "Test 4: Uploading small text file (IPv4)..."
if ! tftp 192.168.0.2 -c put upload_small.txt; then
    echo "ERROR: Test 4 - Upload failed"
    FAILED=1
else
    echo "Test 4: PASS - File uploaded"
fi

# Test 4a: Upload small text file with renaming (IPv4)
echo "Test 4a: Uploading small text file with renaming (IPv4)..."
if ! tftp 192.168.0.2 -c put upload_small.txt upload_renamed.txt; then
    echo "ERROR: Test 4a - Upload failed"
    FAILED=1
else
    echo "Test 4a: PASS - File uploaded"
fi

# Test 4b: Upload multiple small text files to a named directory (IPv4)
echo "Test 4b: Uploading sma small text files to a named directory (IPv4)..."
if ! tftp 192.168.0.2 -c put upload_file_a.txt upload_file_b.txt upload_file_c.txt multi; then
    echo "ERROR: Test 4b - Upload failed"
    FAILED=1
else
    echo "Test 4b: PASS - File uploaded"
fi

# Test 5: Upload large text file (IPv4)
echo "Test 5: Uploading large text file (IPv4)..."
if ! tftp 192.168.0.2 -c put upload_large.txt; then
    echo "ERROR: Test 5 - Upload failed"
    FAILED=1
else
    echo "Test 5: PASS - File uploaded"
fi

# Test 6: Upload binary file (IPv4)
echo "Test 6: Uploading binary file (IPv4)..."
if ! tftp 192.168.0.2 -m binary -c put upload_binary.dat; then
    echo "ERROR: Test 6 - Upload failed"
    FAILED=1
else
    echo "Test 6: PASS - File uploaded"
fi

# Test 7: Download small file using IPv6
echo "Test 7: Downloading small.txt (IPv6)..."
if ! tftp fd51:3681:1eb4::2 -c get small.txt small_ipv6.txt; then
    echo "ERROR: Test 7 - IPv6 download failed"
    FAILED=1
else
    if [ ! -f small_ipv6.txt ]; then
        echo "ERROR: Test 7 - IPv6 file not created"
        FAILED=1
    else
        shasum small_ipv6.txt > small_ipv6.txt.sha1
    
    	if grep -q "9470c442585479bfad86d9a0a0daf01779c020b3" small_ipv6.txt.sha1; then
        	echo "Test 7: PASS - IPv6 download matches expected hash"
    	else
        	echo "ERROR: Test 7 - IPv6 download hash mismatch"
        	echo "Expected: 9470c442585479bfad86d9a0a0daf01779c020b3"
        	echo "Got: $(cat small_ipv6.txt.sha1)"
        	FAILED=1
    	fi
    fi
fi

# Test 8: Error handling - nonexistent file
echo "Test 8: Testing error handling (nonexistent file)..."
if tftp 192.168.0.2 -c get nonexistent.txt nonexistent.txt 2>&1 | grep -q "File not found"; then
    echo "Test 8: PASS - Proper error handling"
else
    echo "ERROR: Test 8 - Failed to handle nonexistent file error"
    FAILED=1
fi

if [ $FAILED -eq 1 ]; then
    echo "ERROR: Client test failed"
    exit 1
fi

echo "All TFTP client tests completed successfully"
echo "TESTS PASSED MARKER"
`

	serverVM := tftpVM(t, "tftp_server", serverScript, net, testFilesForServer)
	clientVM := tftpVM(t, "tftp_client", clientScript, net, testFilesForClient)

	if _, err := serverVM.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("serverVM: %v", err)
	}
	if _, err := clientVM.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("clientVM: %v", err)
	}

	clientVM.Wait()
	serverVM.Wait()
}
