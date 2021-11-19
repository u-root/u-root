//go:generate embedvar -file=./bbmain/cmd/main.go -varname=bbMainSource -p=bb -o=bbmain_src.go
//go:generate embedvar -file=./bbmain/register.go -varname=bbRegisterSource -p=bb -o=bbregister_src.go

package bb
