# Roadmap (alpha)

## Finish commands and tests

| Command        | Flags         | Flags TODO      | Comments               |
| -------------- | ------------- | --------------- | ---------------------- |
| ansi           |               |                 | u-root specific        |
| archive        |               |                 | u-root specific        |
| builtin        | -d            |                 | u-root specific        |
| cat            | -u            |                 | Still need - for stdin |
| chmod          |               | -R, --reference | More mode forms        |
| cmp            | -lLs          |                 |                        |
| comm           | -123hi        |                 |                        |
| cp             | -fiPRrvw      |                 |                        |
| date           | -u            | -drs            |                        |
| dd             |               |                 |                        |
| dhcp           |               |                 | u-root specific        |
| dmesg          | -c            | -Clr            |                        |
| echo           | -n            | -e              |                        |
| ectool         |               |                 | u-root specific        |
| :x: free       |               | -bkmght         | Not implemented yet!   |
| freq           | -cdorx        |                 | From plan 9            |
| :x: gitclone   |               |                 | Not implemented yet!   |
| gopxe          |               |                 | u-root specific        |
| gpgv           | -v            |                 |                        |
| grep           | -glrv         | -cno            | RE2-compatible only    |
| hostname       |               |                 |                        |
| init           |               |                 |                        |
| installcommand |               |                 | u-root specific        |
| ip             |               |                 |                        |
| kexec          |               |                 |                        |
| kill           | -ls           |                 |                        |
| ldd            |               |                 |                        |
| :x: less       |               |                 | Not implemented yet!   |
| ln             | -fiLPrsTtv    |                 |                        |
| losetup        | -Ad           |                 |                        |
| ls             | -lRr          | -hFfS           | -r is raw not reverse  |
| :x: man        |               | -k              | Not implemented yet!   |
| mkdir          | -mpv          |                 |                        |
| mount          | -rt           |                 |                        |
| mv             |               | -nu             |                        |
| netcat         |               |                 |                        |
| pflask         |               |                 | u-root specific        |
| ping           | -6chisVw      |                 |                        |
| printenv       |               |                 |                        |
| :x: printf     |               |                 | Not implemented yet!   |
| ps             | -Aaex         |                 |                        |
| pwd            | -LP           |                 |                        |
| rm             | -iRrv         | -I              |                        |
| rush           |               | -c              |                        |
| script         |               |                 | u-root specific        |
| seq            | -s            |                 |                        |
| sort           | -or           | -bcfmnRu        |                        |
| srvfiles       | -dhp          |                 | u-root specific        |
| sync           |               |                 |                        |
| tcz            | -ahpv         |                 | u-root specific        |
| tee            | -ai           |                 |                        |
| :x: tr         |               |                 | Not implemented yet!   |
| uname          | -admnrsv      |                 |                        |
| uniq           | -cdfu, --cn   | -i              |                        |
| unshare        | -muin         |                 | Different flag names   |
| validate       |               |                 | u-root specific        |
| wc             | -cblrw        |                 |                        |
| wget           |               |                 | No args yet...         |
| which          | -a            |                 |                        |
| :x: xxd        |               |                 |                        |

## New Goal
- [ ] Dealing with filenames containing newlines, spaces and dashes
- [ ] Get enough basic commands working to support a container mechanism.
- [ ] Determine what commands we might need for "New ChromeOS"
- [ ] Bring in Go readline package for the u-root shell
- [ ] Finish implementation of the ip command

## Figure out a container solution
Options:

* Docker
* Rocket
* wget + unpack (cpio? tar?) + u-root pflask
* implement a gitclone command and use u-root pflask
