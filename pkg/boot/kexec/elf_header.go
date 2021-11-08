package kexec

type MemSymbolHdr struct {
	Name      uint
	Type      uint
	Flags     uint
	Addr      uint
	Offset    uint
	Size      uint
	Link      uint
	Info      uint
	AddrAlign uint
	EntSize   uint
	Data      []byte
}

type MemSymbol struct {
	Name  uint   /* Symbol name (string tbl index) */
	Info  string /* No defined meaning, 0 */
	Other string /* Symbol type and binding */
	Shndx string /* Section index  */
	Value uint   /* Symbol value */
	Size  uint   /* Symbol size */
}
