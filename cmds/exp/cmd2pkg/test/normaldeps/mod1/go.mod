module github.com/u-root/gobusybox/test/normaldeps/mod1

go 1.13

replace github.com/u-root/gobusybox/test/normaldeps/mod2/v2 => ../mod2

require (
	github.com/u-root/gobusybox/test/normaldeps/mod2/v2 v2.0.0-00010101000000-000000000000
	golang.org/x/sys v0.0.0-20200905004654-be1d3432aa8f
)
