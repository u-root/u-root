# Roadmap (alpha)

## Finish commands and tests

| Command        | Flags         | Flags TODO      | Comments               |
| -------------- | ------------- | --------------- | ---------------------- |
| ansi           |               |                 | u-root specific        |
| archive        |               |                 | u-root specific        |
| builtin        | -d            |                 | u-root specific        |
| cat            | -u            |                 |                        |
| chmod          |               | -R, --reference | More mode forms        |
| :x: chroot     |               |                 | Not implemented yet!   |
| cmp            | -lLs          |                 |                        |
| comm           | -123h         |                 |                        |
| cp             | -fiPRrvw      |                 |                        |
| cpio           | -oitv         |                 |                        |
| date           | -ur           | -ds             |                        |
| dd             |               |                 |                        |
| dhcp           |               |                 | u-root specific        |
| dmesg          | -c            | -Clr            |                        |
| echo           | -ne           |                 |                        |
| ectool         |               |                 | u-root specific        |
| exit           |               |                 | Rush builtin           |
| false          |               |                 |                        |
| fmap           | -s            | -crudV          | u-root specific        |
| free           | -bkmgth       |                 |                        |
| freq           | -cdorx        |                 | From plan 9            |
| :x: gitclone   |               |                 | Not implemented yet!   |
| gopxe          |               |                 | u-root specific        |
| gpgv           | -v            |                 |                        |
| grep           | -glrv         | -cno            | RE2-compatible only    |
| gzip           |               |                 | Not implemented yet!   |
| hexdump        |               |                 |                        |
| hostname       |               |                 |                        |
| init           |               |                 |                        |
| insmod         |               |                 |                        |
| installcommand |               |                 | u-root specific        |
| ip             |               |                 |                        |
| kexec          |               |                 |                        |
| kill           | -ls           |                 |                        |
| ldd            |               |                 |                        |
| :x: less       |               |                 | Not implemented yet!   |
| ln             | -fiLPrsTtv    |                 |                        |
| losetup        | -Ad           |                 |                        |
| ls             | -lRr          | -hFfS           | -r is raw not reverse  |
| lsmod          |               |                 |                        |
| :x: man        |               | -k              | Not implemented yet!   |
| mkdir          | -mpv          |                 |                        |
| mkfifo         |               |                 |                        |
| mknod          |               |                 |                        |
| modprobe       | -n            |                 | Further options?       |
| mount          | -rt           |                 |                        |
| mv             |               | -nu             |                        |
| netcat         |               |                 |                        |
| pflask         |               |                 | u-root specific        |
| ping           | -6chisVw      |                 |                        |
| printenv       |               |                 |                        |
| :x: printf     |               |                 | Not implemented yet!   |
| ps             | -Aaex         |                 |                        |
| pwd            | -LP           |                 |                        |
| readlink       | -fv           | -emnqsz         |                        |
| rm             | -iRrv         | -I              |                        |
| rmmod          |               | -fsv            |                        |
| run            |               |                 | u-root specific        |
| rush           |               | -c              |                        |
| seq            | -s            |                 |                        |
| shutdown       | halt reboot suspend |           |
| sleep          |               |                 |                        |
| sort           | -or           | -bcfmnRu        |                        |
| srvfiles       | -dhp          |                 | u-root specific        |
| sync           |               |                 |                        |
| tail           | -n            | -f              | u-root specific        |
| tcz            | -ahpv         |                 | u-root specific        |
| tee            | -ai           |                 |                        |
| time           |               | -p              | Rush builtin           |
| :x: tr         |               |                 | Not implemented yet!   |
| true           |               |                 |                        |
| truncate       | -cs           | -or             |                        |
| umount         | -fl           |                 |                        |
| uname          | -admnrsv      |                 |                        |
| uniq           | -cdfu, --cn   | -i              |                        |
| unshare        | -muin         |                 | Different flag names   |
| validate       |               |                 | u-root specific        |
| wc             | -cblrw        |                 |                        |
| wget           |               |                 | No args yet...         |
| which          | -a            |                 |                        |

(Commands marked with an :x: are not yet implemented.)

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
