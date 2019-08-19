# Running Windows over u-root

## Prerequisites
**WARNING**: the scripts will require sudo priviliges. It is highly advised that
you read them carefully before running them!

1. Golang, see `https://golang.org/`.
   Don't forget to have go/bin in your path via `export PATH=${HOME}/go/bin:${PATH}`
1. Packages required to build Linux kernel
1. A functional, **already installed**, Windows Server 2019 or Windows 10 **raw** image, asumed to exist at
   `"${WORKSPACE}"/windows.img`. A raw image can be created via 
   
   `qemu-img create -f raw "${WORKSPACE}"/windows.img 20G`. 
   
   `setup.sh` will create a masking image over it,
   so the original image will not be modified. Windows boot manager,
   bootmgfw.efi is assumed to exist in the 2nd partition of the image. See
   `install_windows.sh` for an example.
1. An environment variable `EFI_WORKSPACE`, where files will be downloaded to or
   otherwise created.
1. kpartx . Install via `sudo apt-get install kpartx`
1. alien. Install via `sudo apt-get install alien`
## Installing the Modified u-root
1.  Install u-root:

    ```shell
    go get github.com/u-root/u-root
    ```

1.  Change the uroot github remote to our modified one:

    ```shell
    pushd ~/go/src/github.com/u-root/u-root
    git remote add oweisse https://github.com/oweisse/u-root  # our revised uroot repo
    git fetch oweisse
    git checkout -b kexec_test oweisse/kexec_test
    go install
    popd
    ```

## Setting up the Kernel.
Setup Linux kernel source tree with our modifications. We modified
kexec_load syscall to launch EFI applications. **Read the script before running!**

The script will:
1. Download an EFI loader image.
1. Extract windows boot-manager from the windows image (see prerequisites above).
1. Clone our forked linux kernel from `https://github.com/oweisse/linux/`
   into $EFI_WORKSPACE.
1. Install prerequisites (sudo required)
1. Build Linux kernel

```
./setup.sh
```

## Running u-Root and booting Windows
The script command line arguments:
1. rebuild_uroot: will also add bootmgfw.efi, extracted from the Windows image
   to the filesystem.
1. rebuild_kernel: Only necessary if you modified the kernel at
   $EFI_WORKSPACE/linux

```
./run_vm.sh rebuild_uroot rebuild_kernel
```

After u-root has loaded, launch Windows bootmanager"
```
pekexec bootmgfw.efi
```

## Attaching gdb
In the following command, `vmlinux` is the Linux kernel we built.
`launch_efi_app` is the function jumping into `bootmgfw.efi` entry point.
Port 1234 is QEMU default port when using the `-s` flag (see `run_vm.sh`).
```
gdb vmlinux -ex "target remote :1234"     \
            -ex "hbreak launch_efi_app"   \
            -ex "layout regs"             \
            -ex "focus next"              \
            -ex "focus next" -ex "c"
```
## Debugging Crashes

Windows loading has 3 stages:

1. Boot loader, e.g., `bootmgfw.efi`
1. Windows loader, .e.g., `Winload.efi`, located  in `C:\Windows\System32`
1. Windows kernel, always `ntoskrnl.exe`, located in  `C:\Windows\System32`

When windows launch crashes, we want to pin-point the exact location of the crash
(or hang) to be able to fix the problem. This section is a **very** short
tutorial on how to find the crash location and navigate around the loaded
binaries using *gdb*. We are using *objdump* and *IDA* for static analysis
(https://www.hex-rays.com/products/ida/). *IDA* is extremely helpfull since it
automatically pulls symbols from Microsoft symbol server, which makes much
easier debugging. It is also beneficial to look at ReactOS
(https://github.com/reactos), since their code looks very similar to what you'll
find in Windows binaries.

### Initial Binary Analysis with *objdump*
1. First we need to extract `Winload.efi` & `ntoskrnl.exe`. This can be done in
   a similar way to what `setup.sh` is doing:

    ```shell
    # Location to store the analyzed binaries:
    WINDOWS_BINARIES=~/windows_binaries/

    # Loop device. You may need to choose another loop device, see `losetup --list`
    LOOP_DEVICE=loop1

    sudo losetup "${LOOP_DEVICE}" "${WINDOWS_DISK}"  # Attach raw disk to loop1
    sudo kpartx -a /dev/"${LOOP_DEVICE}" # Create /dev/mapper/loop1* partitions

    sudo mkdir -p /mnt/win_disk2
    sudo mkdir -p /mnt/win_disk3

    # The boot-manager is typically on the 2nd partition
    # The "C drive"  is typically on the 3rd partition
    sudo mount /dev/mapper/"${LOOP_DEVICE}"p2 /mnt/win_disk2
    sudo mount /dev/mapper/"${LOOP_DEVICE}"p3 /mnt/win_disk3

    cp /mnt/win_disk2/EFI/Microsoft/Boot/bootmgfw.efi "${WINDOWS_BINARIES}"
    cp /mnt/win_disk3/Windows/System32/Winload.efi "${WINDOWS_BINARIES}"
    cp /mnt/win_disk3/Windows/System32/ntoskrnl "${WINDOWS_BINARIES}"

    sudo umount /mnt/win_disk2
    sudo umount /mnt/win_disk3

    sudo kpartx -d /dev/"${LOOP_DEVICE}"         # Remove /dev/mapper paritions
    sudo losetup -d /dev/"${LOOP_DEVICE}"        # Dettach WINDOWS_DISK
    ```

1. Let's analyze the binaries:
	```shell
	cd "${WINDOWS_BINARIES}"
	objdump -xd bootmgfw.efi > bootmgfw.efi.disas
	objdump -xd Winload.efi > Winload.efi.disas
	objdump -xd ntoskrnl.exe > ntoskrnl.exe.disas
	```
1. In the beginning of every `*.disas` file you'll find three important numbers:
  - *AddressOfEntryPoint* - offset of first instruction to be executed
    from the **very beginning** of the code segment.
  - *ImageBase* - "The preferred address of the first byte of image when
    loaded into memory". Read more at
    https://docs.microsoft.com/en-us/windows/win32/debug/pe-format.
  - *start address* - which is the *ImageBase* + *AddressOfEntryPoint*
  
    ```
    # From bootmgfw.efi.disas
    start address       0x000000001001edd0
    AddressOfEntryPoint 000000000001edd0
    ImageBase           0000000010000000
    ```
  We will need the *AddressOfEntryPoint* and *ImageBase* for the following steps.

### Dynamic Analysis with *gdb* - Analyzing *bootmgfw.efi*
1. We will have a total of 3 terminal windows:
	  1. Running the VM with u-root, via `run_vm.sh`
	  1. Running *gdb*
	  1. Running `debugging_tools.py` which opens an IPython dynamic shell.
1. In one terminal, launch u-root via `run_vm.sh`. Wait for *u-root* shell.
1. In another terminal, launch *gdb*
   ```shell
   # replace vmlinux with the full path to the linux kernel.
   gdb "${EFI_WORKSPACE}"/linux/vmlinux  \
       -ex "target remote :1234"         \
       -ex "hbreak launch_efi_app"       \
       -ex "layout regs"                 \
       -ex "focus next"                  \
       -ex "focus next" -ex "c"
   ```
   Hit ENTER to resume the VM
1. In the *u-root* shell run `pekexec bootmgfw.efi`
1. gdb will break on `launch_efi_app`. Scroll through the output in the *u-root*
   shell. Look for a line that looks like
   `Entry point: (64 bytes @   0x000000001001edd0)`.
   This line means that the **dynamic** entry point of
   *bootmgfw.efi*  is `0x000000001001edd0`. In our current implementation this
   is identical to "start address" in *bootmgfw.efi.disas* but this may change
   in the future.
1. Open yet another terminal, and run `debugging_tools.py`. This will launch an
   IPython shell. Let's intialize our *Reverser* to help "reverse engineer"
   windows loading process:
   ```python
   In [1]: bootmg = Reverser( live_entry_point  = 0x1001edd0,
                              image_base        = 0x10000000,
                              image_entry_point = 0x1edd0 )
   ```
1. Open IDA (ida64, to be exact) and load *bootmgfw.efi*. IDA will prompt to ask
   if you want it to load symbols from Microsoft symbol server. The answer is
   definitely YES. A quick look around will show us that the very first steps of
   execution are:
   1. Our entry point:
   `.text:000000001001EDD0 EfiEntry proc near`
   1. Which calls
   `.text:000000001001F0EC BmMain proc near`
   1. Which then calls
   `.text:0000000010082620 BlInitializeLibrary proc near`
1. Say we crash and we don't know why.  Let's see if we even get to
  `BlInitializeLibrary`. In our IPython shell try this:
	```python
	In [2]: bootmg.breakpointFromImageAddr( 0x010082620 )
    b *0x10082620
	```
	We call the address we see in *IDA* the "image address*. These are the preferred addresses for loading, but the actual **dynamic** addresses at run-time may be different.

	Copy the output line `b *0x10082620` and paste it in gdb. Continue execution in
  gdb via the `c`. gdb will break at the entrance to `BlInitializeLibrary`.

	NOTE: The *Reverser* class may now seem very unimpressive, since for bootmgfw.efi the dynamic and
  static addresses are the same. This will be much more interesting when we'll
  get to dig  into Winload.efi and ntoskrnl.exe which have drastically
  different dynamic addresses.

1. We now want to understand where in `BlInitializeLibrary` are we crashing. In
   a normal debugging setting, we whould just use the *step* and *next*
   instructions in *gdb*. However, since there is no source code available, this
   is not possible. We therefore resort to simply place a breakpoint on every
   *call* instruction and it's following instruction to simulate this *step*
   functionality.  From the IPython shell execute the following:
      ```python
      In [3]: bootmg.generateBreakpoinCmdsOnCalls( imageDisasPath = "bootmg.efi.disas",
                                                   liveStartAddress = 0x10082620)
      Will search for code starting at 0x10082620
      Starting analysis at: 10082620:	48 89 5c 24 08       	mov    %rbx,0x8(%rsp)
      Found call at: 10082685:	e8 3e 03 00 00       	callq  0x100829c8
      Found call at: 100826a4:	e8 db de fa ff       	callq  0x10030584
      Found call at: 100826ab:	e8 b8 3d 00 00       	callq  0x10086468
      Found call at: 100826b5:	e8 8a 88 02 00       	callq  0x100aaf44
      Found call at: 100826c1:	e8 7a ae fe ff       	callq  0x1006d540
      Found call at: 100826cd:	e8 8e 0f 00 00       	callq  0x10083660
      Reached end of function at: 100826e5:	c3                   	retq
      b *0x10082685
      b *0x1008268a
      b *0x100826a4
      b *0x100826a9
      b *0x100826ab
      b *0x100826b0
      b *0x100826b5
      b *0x100826ba
      b *0x100826c1
      b *0x100826c6
      b *0x100826cd
      b *0x100826d2
      ```
    The command translate the **dynamic** (live) address 0x10082620  into a **static** one. It then inspects the dissasembly in `bootmg.efi.disas`, look for the instruction in that **static** address (a.k.a. the "image address") and searches for *call*s until it hit a *ret* instruction.  
    
      Copy-paste all the breakpoint commands into *gdb*, to set all the breakpoints.
      
      IMPORTANT: This is obviously
    just a heuristic. What it does is linearly searching for calls until it hit a
    ret instruction. It may fail to find all the "calls" in a function. Inspect
    the binary in IDA to see if you are missing any breakpoints.

1. You may now realize that Windows loader is crahing when calling (for example)
   `100826C1 call    BlpDisplayInitialize`. What you want to do now is to go to
   the implemantation of BlpDisplayInitialize (double click it in IDA) and
   continue the process from there.  For instance, if `BlpDisplayInitialize` is
   in **dynamic** address `0x1006D540`, try running:
    ```python
    In [3]: bootmg.generateBreakpoinCmdsOnCalls( imageDisasPath = "bootmg.efi.disas",
                                                 liveStartAddress = 0x1006D540)
    Starting analysis at: 1006d540:	48 83 ec 28          	sub    $0x28,%rsp
    Found call at: 1006d549:	e8 16 07 00 00       	callq  0x1006dc64
    Found call at: 1006d567:	ff d0                	callq  *%rax
    Found call at: 1006d583:	ff d0                	callq  *%rax
    Reached end of function at: 1006d58c:	c3                   	retq
    b *0x1006d549
    b *0x1006d54e
    b *0x1006d567
    b *0x1006d569
    b *0x1006d583
    b *0x1006d585
    ```

1. This process is effectively a "binary" search for the offending instruction.
    We start from the top (BmMain) and we find the child call that crashes. We
    repeat this process recursively until we get to a function with no "calls".
    Then we can try understanding what went wrong.

### Finding the Dynamic Addresses of Winload.efi

So you think *bootmgfw.efi* if running to completion. But maybe Windows boot
process crashes during *Winload.efi* execution.

1. In IDA (for `bootmgfw.efi`), locate a function called `Archpx64TransferTo64BitApplicationAsm`. In
   this function, you'll see an instruction like this:
   
   `.text:0000000010137EE4   call    rax ; ArchpChildAppEntryRoutine`
   
   (potentially in a different address).
1. Use the IPython shell to generate a breakpoing command:
	```python
	In [4]: bootmg.breakpointFromImageAddr( 0x10137EE4 )
    b *0x10137EE4
	```
1. Copy-paste the breakpoint command and continue `c`. When you hit the
   breakpoint, execute a single instruction via `si`.
1. Look at that! We just jumpt into *Winload.efi*. Which means we now know the
   **dynamic** address of the entry point. For example, let's say that *gdb* is
   now at address `0x100920090`. If we inspect `Winload.efi.disas` we see that:
	```
	start address 0x0000000180001090

	Characteristics 0x2022
			executable
			large address aware
			DLL

	Time/Date               Sat May 11 15:47:20 2019
	Magic                   020b    (PE32+)
	MajorLinkerVersion      14
	MinorLinkerVersion      13
	SizeOfCode              00158c00
	SizeOfInitializedData   00040200
	SizeOfUninitializedData 00000000
	AddressOfEntryPoint     0000000000001090
	BaseOfCode              0000000000001000
	ImageBase               0000000180000000
	```

  	Based on that information, we can initialize a new Reverser object in our IPython shell:
    ```
    In [5]: winload = Reverser( live_entry_point  = 0x100920090,
                                image_base        = 0x180000000,
                                image_entry_point = 0x1090 )
    ```

1. Now that the **dynamic** entry point is vastly different the **static**
  entry point, we can really see the benefit of the *Reverser* class.  In IDA, we
  see, for example, the following function of Winload.efi:
  
  	`.text:0000000180002174 OslpMain proc near`

	We can now create a dynamic breakpoint for it based only on the address we see
  in IDA:
	```python
	In [6]: winload.breakpointFromImageAddr( 0x180002174 )
  	b *0x100921174
	```
	Copy-paste this breakpoint command and hit `c`.
1. Let's try to generate breakpoints on all *calls* inside OslpMain. Notice that
   we use the **dynamic** address that we see in gdb.
    ```python
    In [7]: winload.generateBreakpoinCmdsOnCalls( "winload.efi.disas",
                                                  liveStartAddress = 0x100921174)
    Will search for code starting at 0x180002174
    Starting analysis at: 180002174:	48 89 5c 24 08       	mov    %rbx,0x8(%rsp)
    Found call at: 1800021a5:	e8 1a 19 0f 00       	callq  0x1800f3ac4
    Found call at: 1800021e7:	e8 20 1c 0f 00       	callq  0x1800f3e0c
    Found call at: 1800021fe:	e8 2d 41 14 00       	callq  0x180146330
    Found call at: 18000224d:	e8 6e 1c 0f 00       	callq  0x1800f3ec0
    Found call at: 18000227a:	e8 b9 03 00 00       	callq  0x180002638
    Found call at: 18000228e:	e8 39 c1 00 00       	callq  0x18000e3cc
    Found call at: 180002295:	e8 c6 39 01 00       	callq  0x180015c60
    Found call at: 1800022a4:	e8 cb 7a 02 00       	callq  0x180029d74
    Reached end of function at: 1800022bb:	c3                   	retq
    b *0x1009211a5
    b *0x1009211aa
    b *0x1009211e7
    b *0x1009211ec
    b *0x1009211fe
    b *0x100921203
    b *0x10092124d
    b *0x100921252
    b *0x10092127a
    b *0x10092127f
    b *0x10092128e
    b *0x100921293
    b *0x100921295
    b *0x10092129a
    b *0x1009212a4
    b *0x1009212a9
    ```
    
1. IMPORTANT: *bootmgfw.efi* dynamically loads *Winload.efi*. In theory, this
can be in a different location on every invocation (though currently it is
always loaded to the same address). Regardless, you can only ask gdb to place
breakpoints in address belonging to Winload.efi **after** Winload.efi was loaded
into memory. This is because software breakpoints are implemented by replacing
an instruction with an `int 3` instruction  gdb cannot replace an instruction
with `int 3` if the instruction was not yet loaded into memory.

### Finding the Dynamic Address of ntoskrnl.exe
After Winload.efi was loaded into memory, we can now identify where
*ntoskrnl.exe* is loaded.

1. Load Winload.efi into IDA. When prompted if IDA should fetch symbols from
Microsoft server - answer "yes please".
1. In IDA, locate a function called `OslArchTransferToKernel`:
	```
	.text:000000018014EBF0 OslArchTransferToKernel proc near
	.....
	.text:000000018014EC60 retfq
	.text:000000018014EC60 OslArchTransferToKernel endp
	```
	The *ret* instruction will actually jump to the entry point of *ntoskrnl.exe*.
1. We need a breakpoint at a location corresponding to image-address of
0x18014EC60. In the IPython shell execute the following to get the breakpoint
command:
    ```
    In [8]: winload.breakpointFromImageAddr( 0x18014EC60 )
    b *0x100a6dc60
    ```
	Continue execution via `c` and then do a step-instruction (`si`).
1. Congrats! We're now inside Windows kernel! And we now know the **dynamic**
   entry point. In our example we see that *gdb* jumped to `0xfffff80746e05010`.
   If we inspect `ntoskrnl.exe.disas` we see the following:

	```
	start address 0x0000000140566010

	Characteristics 0x22
			executable
			large address aware

	Time/Date               Sat Jul 31 08:45:35 2021
	Magic                   020b    (PE32+)
	MajorLinkerVersion      14
	MinorLinkerVersion      13
	SizeOfCode              007ca200
	SizeOfInitializedData   00185000
	SizeOfUninitializedData 00000000
	AddressOfEntryPoint     0000000000566010
	BaseOfCode              0000000000001000
	ImageBase               0000000140000000
	```
  	Let's initialize a Reverser object:

	```
	In [9]: ntoskrnl = Reverser( live_entry_point  = 0xfffff80746e05010,
                                  image_base        = 0x140000000,
                                  image_entry_point = 0x566010 )
	```

1. Say you're suspecting something is going south in `KdInitSystem`. Inspecting
it in IDA shows the **static** address:

	`0000000140917140 KdInitSystem proc near`

  	First, lste's create a breakpoint. In the IPython shell:

    ```
    In [10]: ntoskrnl.breakpointFromImageAddr( 0x140917140 )
    b *0xfffff807471b6140
    ```

	Now let's place breakpoints on all the calls:
    ```
    In [11]: ntoskrnl.generateBreakpoinCmdsOnCalls( "ntoskrnl.exe.disas",
                                                    liveStartAddress = 0xfffff807471b6140
    Will search for code starting at 0x140917140
    Starting analysis at: 140917140:	48 89 5c 24 18       	mov    %rbx,0x18(%rsp)
    Found call at: 14091718d:	48 ff 15 e4 7f c2 ff 	rex.W callq *-0x3d801c(%rip)        # 0x14053f178
    Found call at: 1409171d2:	e8 89 71 87 ff       	callq  0x14018e360
    Reached end of function at: 1409171f1:	c3                   	retq
    b *0xfffff807471b618d
    b *0xfffff807471b6194
    b *0xfffff807471b61d2
    b *0xfffff807471b61d7
      ```

    **That's weird**. We only identified 2 calls, but inspecting the function in
    IDA clearly shows there are more. We need to be creative. We only found 2
    calls because of the *ret* intstuction at `0x1409171f1`. We can see that the
    following instruction starts at `0x1409171f2`. Let's translate this **static**
    address (image address) into a **dynamic** one and continue searching for calls from there:

    ```python
      In [12]: ntoskrnl.imageAddr2Addr( 0x1409171f2 ) 
      Out[12]: (18446735308874277362, '0xfffff807471b61f2')

      In [13]: ntoskrnl.generateBreakpoinCmdsOnCalls( "ntoskrnl.exe.disas",
                                                    liveStartAddress = 0xfffff807471b61f2)

      Will search for code starting at 0x1409171f2
      Starting analysis at: 1409171f2:	44 38 2d 3a c0 b0 ff 	cmp    %r13b,-0x4f3fc6(%rip)        # 0x140423233
      Found call at: 140917245:	e8 46 ff 72 ff       	callq  0x140047190
      Found call at: 140917289:	e8 92 01 00 00       	callq  0x140917420
      Found call at: 140917316:	e8 05 7c 87 ff       	callq  0x14018ef20
      Found call at: 14091732f:	e8 dc 7e 87 ff       	callq  0x14018f210
      Found call at: 140917347:	e8 c4 7e 87 ff       	callq  0x14018f210
      Found call at: 14091735f:	e8 ac 7e 87 ff       	callq  0x14018f210
      Found call at: 14091738c:	e8 7f 7e 87 ff       	callq  0x14018f210
      Found call at: 1409173a4:	e8 67 7e 87 ff       	callq  0x14018f210
      Found call at: 140917434:	e8 47 3a 80 ff       	callq  0x14011ae80
      Found call at: 14091748c:	e8 3f f7 6f ff       	callq  0x140016bd0
      Reached end of function at: 1409174b3:	c3                   	retq
      b *0xfffff807471b6245
      b *0xfffff807471b624a
      b *0xfffff807471b6289
      b *0xfffff807471b628e
      b *0xfffff807471b6316
      b *0xfffff807471b631b
      b *0xfffff807471b632f
      b *0xfffff807471b6334
      b *0xfffff807471b6347
      b *0xfffff807471b634c
      b *0xfffff807471b635f
      b *0xfffff807471b6364
      b *0xfffff807471b638c
      b *0xfffff807471b6391
      b *0xfffff807471b63a4
      b *0xfffff807471b63a9
      b *0xfffff807471b6434
      b *0xfffff807471b6439
      b *0xfffff807471b648c
      b *0xfffff807471b6491
          ```

### Dynamic Debugging Summary

Finding the function that crashes the entire machine is clearly just the
beginning of the rabbit hole. If there is a NULL-pointer access (or access to
0x42) one must dig deeper and query how did a register received a value of 0x42.
I woukd highly recommend searching the function names in ReactOS source code.
While the code there is different from what you'll find in Windows, it is
sufficiently similar to be helpful to understand the boot process. If you're
lucky enough to have a licence to IDA decompiler, you can press F5 in IDA to
"decompile" the binary code.
