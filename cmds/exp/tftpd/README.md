# TFTP Server

A simple TFTP server implementation that can handle both read and write requests.
Mainly added for integrations tests for tftp (client).

## Overview

This TFTP server uses the same underlying package as the u-root TFTP client (`pack.ag/tftp`).
It implements the basic TFTP protocol as defined in [RFC 1350](https://tools.ietf.org/html/rfc1350) with support for:

- READ requests (GET)
- WRITE requests (PUT)
- Security against directory traversal attacks
- Both binary and netascii transfer modes

## Usage

```
tftpd [-port PORT] [-root DIRECTORY] [-v]
```

### Options

- `-port`: Port to listen on (default 69)
- `-root`: Root directory to serve files from (default current directory)
- `-v`: Enable verbose logging

### Example

Start a TFTP server that listens on port 6969 and serves files from `/tmp/tftp`:

```
tftpd -port 6969 -root /tmp/tftp
```

## Testing with the u-root TFTP Client

To test file transfer between the client and server:

1. Start the server in one terminal:
   ```
   tftpd -root /tmp/tftp
   ```

2. In another terminal, use the u-root TFTP client to transfer files:
   ```
   # Get a file
   tftp localhost -c get test.txt

   # Put a file
   tftp localhost -c put test.txt
   ```
