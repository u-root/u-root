module github.com/u-root/u-root

go 1.17

require (
	github.com/beevik/ntp v0.3.0
	github.com/cenkalti/backoff/v4 v4.0.2
	github.com/creack/pty v1.1.11
	github.com/davecgh/go-spew v1.1.1
	github.com/dustin/go-humanize v1.0.0
	github.com/gliderlabs/ssh v0.1.2-0.20181113160402-cbabf5414432
	github.com/gojuno/minimock/v3 v3.0.8
	github.com/google/go-cmp v0.5.2
	github.com/google/go-tpm v0.2.1-0.20200615092505-5d8a91de9ae3
	github.com/google/goexpect v0.0.0-20191001010744-5b6988669ffa
	github.com/google/goterm v0.0.0-20200907032337-555d40f16ae2
	github.com/insomniacslk/dhcp v0.0.0-20210817203519-d82598001386
	github.com/intel-go/cpuid v0.0.0-20200819041909-2aa72927c3e2
	github.com/klauspost/pgzip v1.2.4
	github.com/kr/pty v1.1.8
	github.com/orangecms/go-framebuffer v0.0.0-20200613202404-a0700d90c330
	github.com/pborman/getopt/v2 v2.1.0
	github.com/pierrec/lz4/v4 v4.1.11
	github.com/rck/unit v0.0.3
	github.com/rekby/gpt v0.0.0-20200219180433-a930afbc6edc
	github.com/safchain/ethtool v0.0.0-20200218184317-f459e2d13664
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/u-root/gobusybox/src v0.0.0-20210529132627-adc854ea4425
	github.com/u-root/iscsinl v0.1.1-0.20210528121423-84c32645822a
	github.com/u-root/uio v0.0.0-20210528151154-e40b768296a7
	github.com/ulikunitz/xz v0.5.8
	github.com/vishvananda/netlink v1.1.1-0.20211118161826-650dca95af54
	github.com/vtolstov/go-ioctl v0.0.0-20151206205506-6be9cced4810
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/sys v0.0.0-20210820121016-41cdb8703e55
	golang.org/x/term v0.0.0-20210317153231-de623e64d2a6
	golang.org/x/text v0.3.3
	golang.org/x/tools v0.1.1
	gopkg.in/yaml.v2 v2.2.2
	pack.ag/tftp v1.0.1-0.20181129014014-07909dfbde3c
	src.elv.sh v0.16.3
)

require (
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be
	github.com/jsimonetti/rtnetlink v0.0.0-20201110080708-d2c240429e6c
	github.com/klauspost/compress v1.10.6
	github.com/mdlayher/ethernet v0.0.0-20190606142754-0394541c37b7
	github.com/mdlayher/netlink v1.1.1
	github.com/mdlayher/raw v0.0.0-20191009151244-50f2db8cc065
	github.com/pmezard/go-difflib v1.0.0
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f
	golang.org/x/mod v0.4.2
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

require (
	github.com/kaey/framebuffer v0.0.0-20140402104929-7b385489a1ff // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
)

retract (
	v7.0.0+incompatible
	v6.0.0+incompatible
	v5.0.0+incompatible
	v4.0.0+incompatible
	v3.0.0+incompatible
	v2.0.0+incompatible
	v1.0.0
)
