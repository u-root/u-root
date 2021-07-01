# Roadmap

## Missing commands and flags

Before starting work on a command, please open a GitHub issue and assign to
yourself.

| Command        | Flags TODO      | Comments               |
| -------------- | --------------- | ---------------------- |
| :x: base64     | -d              | Not implemented yet!   |
| :x: flashrom   | -p internal     |                        |
| :x: gitclone   |                 | Not implemented yet!   |
| grep           | -nF             | RE2-compatible only    |
| gzip           | -d              |                        |
| :x: printf     |                 | Not implemented yet!   |
| ps             |                 | Fix race conditions    |
| readlink       | -em             |                        |
| :x: sed        | -ie             | Not implemented yet!   |
| sort           | -bcfmnRu        |                        |
| srvfiles       |                 | Serve files with TLS   |
| :x: time       | -p              |                        |
| truncate       | -o              |                        |
| uniq           | -i              |                        |
| unshare        |                 | Different flag names   |

(Commands marked with an :x: are not yet implemented.)
