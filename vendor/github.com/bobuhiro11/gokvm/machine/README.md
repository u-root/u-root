# machine

## Memory layout
```
InitialRegState GuestPhysAddr                      Binary files [+ offsets in the file]

                 0x00000000    +------------------+
                               |                  |
 RSI -->         0x00010000    +------------------+ bzImage [+ 0]
                               |                  |
                               |  boot param      |
                               |                  |
                               +------------------+
                               |                  |
                 0x00020000    +------------------+
                               |                  |
                               |   cmdline        |
                               |                  |
                               +------------------+
                               |                  |
 RIP -->         0x00100000    +------------------+ bzImage [+ 512 x (setup_sects in boot param header + 1)]
                               |                  |
                               |   64bit kernel   |
                               |                  |
                               +------------------+
                               |                  |
                 0x0f000000    +------------------+ initrd [+ 0]
                               |                  |
                               |   initrd         |
                               |                  |
                               +------------------+
                               |                  |
                 0x40000000    +------------------+
```