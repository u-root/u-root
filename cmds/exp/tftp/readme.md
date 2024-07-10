# TFTP

## Usage
```
tftp [ options... ] [host [port]] [-c command]
```

### Options
- [ ] `-4` Connect with IPv4 only, even if IPv6 support was compiled in.
- [ ] `-6` Connect with IPv4 only, even if IPv6 support was compiled in.
- [ ] `-l` Default to literal mode. Used to avoid special processing of ':' in a file name.
- [ ] `-R <PORT:PORT>` Force the originating port number to be in the specified range of port numbers.
- [x] `-m <ascii/binary>` Set the default transfer mode to mode.  This is usually used with -c.
- [ ] `-v` Default to verbose mode.
- [ ] `-V` Print the version number and configuration to standard output, then exit gracefully.

The flag `-c` is a positional argument and if set must be placed at the end.
Is must have one argument, which is the actual command to execute. Some of these commands have arguments as well.

### Commands
- [x] `?,h,help` - Print help information
- [x] `ascii` - Set mode to netascii
- [x] `binary` - Set mode to binary
- [x] `connect <host> [<port>]` - Set Host and port to respective values for future connection
- [x] `get <file>` - Get file from a set host
- [x] `get <remotefile> <localfile>` - Get remotefile and save it in remotefile
- [x] `get <file1> <file2> <file3>` - Get all the file
- [ ] `literal` - Set literal mode: Treads `:` in filenames differently. (Windows path support)
- [x] `mode <ascii/binary>` - Set mode to netascii or binary
- [x] `put <file>` - Put file on set host
- [x] `put <localfile> <remotefile>` - Put localfile to host in remotefile
- [x] `put <file1> <file2> <file3> .... <remote-directory>` - Put files into remote-directory of host
- [x] `quit` - Quit immediatly
- [x] `rexmt <int>` - Set per-packet retransmission
- [x] `status` - Prints hostname, port and status of Mode, Literal, verbose, rexmt, timeout
- [x] `timeout <int>` - Set timeout value
- [ ] `trace` - Switch trace mode
- [ ] `verbose` - Switch verbose mode

Even though some of the commands are available in the program, they have no functionality implemented.

### Missing options/commands
- Hostname will be resolved by the `pack.ag/tftp` library dynamically and a connection is managed by it. (Missing -4 and -6)
- The library is incapable of setting the originating port (Missing -R port:port)
- Running tftp server on windows is for mad men. Not supporting such crimes against humanity (Missing -l support)
- There is very little to report on this program, if something goes wrong, the user will be notified by error (Missing -v)


### Interactive mode
All commands available via commandline flag `-c` are also available in interactive mode.
