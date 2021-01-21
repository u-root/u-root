# Roadmap

## Missing commands and flags

Before starting work on a command, please open a GitHub issue and assign to
yourself.

| Command        | Flags TODO      | Comments               |
| -------------- | --------------- | ---------------------- |
| :x: flashrom   | -p internal     |                        |
| :x: gitclone   |                 | Not implemented yet!   |
| grep           | -cnF            | RE2-compatible only    |
| ls             | -hFfS           | -r is raw not reverse  |
| :x: man        | -k              | Not implemented yet!   |
| mv             | -n              |                        |
| ping           | -a              |                        |
| :x: printf     |                 | Not implemented yet!   |
| ps             |                 | Fix race conditions    |
| readlink       | -em             |                        |
| sort           | -bcfmnRu        |                        |
| srvfiles       |                 | Serve files with TLS   |
| :x: time       | -p              |                        |
| truncate       | -or             |                        |
| uniq           | -i              |                        |
| unshare        |                 | Different flag names   |
| wget           |                 | No args yet...         |

(Commands marked with an :x: are not yet implemented.)
