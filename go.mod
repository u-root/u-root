module github.com/u-root/u-root

go 1.19

require (
	github.com/ProtonMail/go-crypto v0.0.0-20221026131551-cf6655e29de4
	github.com/beevik/ntp v0.3.0
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/creack/pty v1.1.18
	github.com/davecgh/go-spew v1.1.1
	github.com/dustin/go-humanize v1.0.1
	github.com/gliderlabs/ssh v0.1.2-0.20181113160402-cbabf5414432
	github.com/gojuno/minimock/v3 v3.0.8
	github.com/google/go-cmp v0.5.9
	github.com/google/go-tpm v0.9.1-0.20230914180155-ee6cbcd136f8
	github.com/google/goexpect v0.0.0-20191001010744-5b6988669ffa
	github.com/google/uuid v1.3.0
	github.com/insomniacslk/dhcp v0.0.0-20230220063916-5369909a5de7
	github.com/intel-go/cpuid v0.0.0-20200819041909-2aa72927c3e2
	github.com/kevinburke/ssh_config v1.1.0
	github.com/klauspost/compress v1.10.6
	github.com/klauspost/pgzip v1.2.4
	github.com/knz/bubbline v0.0.0-20230717192058-486954f9953f
	github.com/kr/pty v1.1.8
	github.com/nanmu42/limitio v1.0.0
	github.com/orangecms/go-framebuffer v0.0.0-20200613202404-a0700d90c330
	github.com/pborman/getopt/v2 v2.1.0
	github.com/pierrec/lz4/v4 v4.1.18
	github.com/rck/unit v0.0.3
	github.com/rekby/gpt v0.0.0-20200219180433-a930afbc6edc
	github.com/safchain/ethtool v0.0.0-20200218184317-f459e2d13664
	github.com/spf13/pflag v1.0.5
	github.com/u-root/gobusybox/src v0.0.0-20230817123913-21a9727f4316
	github.com/u-root/iscsinl v0.1.1-0.20210528121423-84c32645822a
	github.com/u-root/uio v0.0.0-20230305220412-3e8cd9d6bf63
	github.com/ulikunitz/xz v0.5.8
	github.com/vishvananda/netlink v1.1.1-0.20211118161826-650dca95af54
	github.com/vtolstov/go-ioctl v0.0.0-20151206205506-6be9cced4810
	golang.org/x/crypto v0.14.0
	golang.org/x/sys v0.13.0
	golang.org/x/term v0.13.0
	golang.org/x/text v0.13.0
	golang.org/x/tools v0.14.0
	gopkg.in/yaml.v2 v2.2.8
	mvdan.cc/sh/v3 v3.7.0
	pack.ag/tftp v1.0.1-0.20181129014014-07909dfbde3c
	src.elv.sh v0.16.0-rc1.0.20220116211855-fda62502ad7f
)

require (
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/bubbles v0.15.1-0.20230123181021-a6a12c4a31eb // indirect
	github.com/charmbracelet/bubbletea v0.24.1 // indirect
	github.com/charmbracelet/lipgloss v0.7.1 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/containerd/console v1.0.4-0.20230706203907-8f6c4e4faef5 // indirect
	github.com/google/goterm v0.0.0-20200907032337-555d40f16ae2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/jsimonetti/rtnetlink v0.0.0-20201110080708-d2c240429e6c // indirect
	github.com/kaey/framebuffer v0.0.0-20140402104929-7b385489a1ff // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mdlayher/ethernet v0.0.0-20190606142754-0394541c37b7 // indirect
	github.com/mdlayher/netlink v1.1.1 // indirect
	github.com/mdlayher/raw v0.0.0-20191009151244-50f2db8cc065 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/sahilm/fuzzy v0.1.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	golang.org/x/mod v0.13.0 // indirect
	golang.org/x/net v0.16.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	google.golang.org/grpc v1.27.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

retract (
	// Published v7 too early (before migrating to go modules)
	v7.0.0+incompatible
	// Published v6 too early (before migrating to go modules)
	v6.0.0+incompatible
	// Published v5 too early (before migrating to go modules)
	v5.0.0+incompatible
	// Published v4 too early (before migrating to go modules)
	v4.0.0+incompatible
	// Published v3 too early (before migrating to go modules)
	v3.0.0+incompatible
	// Published v2 too early (before migrating to go modules)
	v2.0.0+incompatible
	// Published v1 too early (before migrating to go modules)
	[v1.0.0, v1.0.1]
)
