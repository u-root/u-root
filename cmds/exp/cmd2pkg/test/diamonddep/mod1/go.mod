module github.com/u-root/gobusybox/test/diamonddep/mod1

go 1.13

replace github.com/u-root/gobusybox/test/diamonddep/mod2 => ../mod2

replace github.com/u-root/gobusybox/test/diamonddep/mod3 => ../mod3

require (
	github.com/u-root/gobusybox/test/diamonddep/mod2 v0.0.0-00010101000000-000000000000
	github.com/u-root/gobusybox/test/diamonddep/mod3 v0.0.0-00010101000000-000000000000 // indirect
)
