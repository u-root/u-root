# Roadmap (alpha)

## Finish commands and tests

Before starting work on a command, please open a GitHub issue and assign to
yourself.

| Command        | Flags TODO      | Comments               |
| -------------- | --------------- | ---------------------- |
| chmod          | -R, --reference | More mode forms        |
| dd             | -conv=notrunc   | Also sync, ...         |
| fmap           |                 | Move into fiano        |
| :x: gitclone   |                 | Not implemented yet!   |
| grep           | -cnF            | RE2-compatible only    |
| :x: less       |                 | Not implemented yet!   |
| ls             | -hFfS           | -r is raw not reverse  |
| :x: man        | -k              | Not implemented yet!   |
| mv             | -nu             |                        |
| ping           | -a              |                        |
| :x: printf     |                 | Not implemented yet!   |
| ps             |                 | Fix race conditions    |
| readlink       | -emn            |                        |
| rush           | -c              |                        |
| seq            |                 |                        |
| sort           | -bcfmnRu        |                        |
| srvfiles       |                 | Serve files with TLS   |
| sync           | -df             |                        |
| time           | -p              | Rush builtin           |
| :x: tr         |                 | Not implemented yet!   |
| truncate       | -or             |                        |
| uniq           | -i              |                        |
| unshare        |                 | Different flag names   |
| wget           |                 | No args yet...         |

(Commands marked with an :x: are not yet implemented.)
