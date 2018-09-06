load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "70d0204f1e834d14fa9eef1e9b97160917a48957cd1e3a39b5ef9acdbdde6972",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.15.2/rules_go-0.15.2.tar.gz"],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "c0a5739d12c6d05b6c1ad56f2200cb0b57c5a70e03ebd2f7b87ce88cabf09c7b",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.14.0/bazel-gazelle-0.14.0.tar.gz"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_beevik_ntp",
    commit = "62c80a04de2086884d8296004b6d74ee1846c582",
    importpath = "github.com/beevik/ntp",
)

go_repository(
    name = "ag_pack_tftp",
    commit = "d0fca786a8ea95ad4515bb83b2c34a11ec6a3fc0",
    importpath = "pack.ag/tftp",
)

go_repository(
    name = "com_github_dustin_go_humanize",
    commit = "9f541cc9db5d55bce703bd99987c9d5cb8eea45e",
    importpath = "github.com/dustin/go-humanize",
)

go_repository(
    name = "com_github_go_test_deep",
    commit = "6592d9cc0a499ad2d5f574fde80a2b5c5cc3b4f5",
    importpath = "github.com/go-test/deep",
)

go_repository(
    name = "com_github_google_go_tpm",
    commit = "f9bda7c79425630507590804d9ecb42603d6291e",
    importpath = "github.com/google/go-tpm",
)

go_repository(
    name = "com_github_google_goexpect",
    commit = "9db06cbbaed691242d13045e8e5f5037f72db9e7",
    importpath = "github.com/google/goexpect",
)

go_repository(
    name = "com_github_google_goterm",
    commit = "70e1c263818522aa1ecbdd34d7e143639d447ced",
    importpath = "github.com/google/goterm",
)

go_repository(
    name = "com_github_google_netstack",
    commit = "8b4eef0f44a9713d81e56056edf3b96f28f7a8e3",
    importpath = "github.com/google/netstack",
)

go_repository(
    name = "com_github_gorilla_context",
    commit = "08b5f424b9271eedf6f9f0ce86cb9396ed337a42",
    importpath = "github.com/gorilla/context",
)

go_repository(
    name = "com_github_gorilla_mux",
    commit = "e3702bed27f0d39777b0b37b664b6280e8ef8fbf",
    importpath = "github.com/gorilla/mux",
)

go_repository(
    name = "com_github_klauspost_compress",
    commit = "b939724e787a27c0005cabe3f78e7ed7987ac74f",
    importpath = "github.com/klauspost/compress",
)

go_repository(
    name = "com_github_klauspost_cpuid",
    commit = "ae7887de9fa5d2db4eaa8174a7eff2c1ac00f2da",
    importpath = "github.com/klauspost/cpuid",
)

go_repository(
    name = "com_github_klauspost_crc32",
    commit = "cb6bfca970f6908083f26f39a79009d608efd5cd",
    importpath = "github.com/klauspost/crc32",
)

go_repository(
    name = "com_github_klauspost_pgzip",
    commit = "0bf5dcad4ada2814c3c00f996a982270bb81a506",
    importpath = "github.com/klauspost/pgzip",
)

go_repository(
    name = "com_github_mdlayher_dhcp6",
    commit = "e26af0688e455a82b14ebdbecf43f87ead3c4624",
    importpath = "github.com/mdlayher/dhcp6",
)

go_repository(
    name = "com_github_mdlayher_ethernet",
    commit = "0a1564b57aeaf8387d273d4bdffbb34b1bb8f626",
    importpath = "github.com/mdlayher/ethernet",
)

go_repository(
    name = "com_github_mdlayher_eui64",
    commit = "eee6532bb9adf30c2a17e9963ed90aa14161dae9",
    importpath = "github.com/mdlayher/eui64",
)

go_repository(
    name = "com_github_mdlayher_raw",
    commit = "1d2cec5bb8cce1f029119448aa6760cf9421bea4",
    importpath = "github.com/mdlayher/raw",
)

go_repository(
    name = "com_github_rck_unit",
    commit = "16d9ed3d60d943bbb0ed704264795a2541457601",
    importpath = "github.com/rck/unit",
)

go_repository(
    name = "com_github_spf13_pflag",
    commit = "9a97c102cda95a86cec2345a6f09f55a939babf5",
    importpath = "github.com/spf13/pflag",
)

go_repository(
    name = "com_github_u_root_dhcp4",
    commit = "f78158b7c380d2ae515105471716468098373c72",
    importpath = "github.com/u-root/dhcp4",
)

go_repository(
    name = "com_github_vishvananda_netlink",
    commit = "8aa85bfa77a45236ae842cab3a91853e2b74e07a",
    importpath = "github.com/vishvananda/netlink",
)

go_repository(
    name = "com_github_vishvananda_netns",
    commit = "13995c7128ccc8e51e9a6bd2b551020a27180abd",
    importpath = "github.com/vishvananda/netns",
)

go_repository(
    name = "org_golang_google_grpc",
    commit = "32fb0ac620c32ba40a4626ddf94d90d12cce3455",
    importpath = "google.golang.org/grpc",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "aabede6cba87e37f413b3e60ebfc214f8eeca1b0",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "org_golang_x_net",
    commit = "aaf60122140d3fcf75376d319f0554393160eb50",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "1c9583448a9c3aa0f9a6a5241bf73c0bd8aafded",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "org_golang_x_tools",
    commit = "7d1dc997617fb662918b6ea95efc19faa87e1cf8",
    importpath = "golang.org/x/tools",
)
