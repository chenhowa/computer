package assembly

/*TokenType is an enum for the types of tokens that RISC-V assembly supports*/
type TokenType uint

/*These constants are a TokenType enum for the types of tokens that RISC-V assembly supports*/
const (
	Label TokenType = iota + 1 // start at 1 to avoid default initialization match problems
	Error
	EndOfInput
	Newline
	NumericConstant
	Register

	RDCYCLE
	RDCYCLEH
	RDTIME
	RDTIMEH
	RDINSTRET
	RDINSTRETH
	ADDI
	SLTI
	SLTIU
	ANDI
	ORI
	XORI
	SLLI
	SRLI
	SRAI
	LUI
	AUIPC
	ADD
	SLT
	SLTU
	AND
	OR
	XOR
	SLL
	SRL
	SUB
	SRA
	NOP
	JAL
	JALR
	BEQ
	BNE
	BLT
	BLTU
	BGE
	BGEU
	LW
	LH
	LHU
	LB
	LBU
	SW
	SH
	SB
	CSRRW
	CSRRS
	CSRRC
	CSRRWI
	CSRRSI
	CSRRCI
	CSRR
	CSRW
	CSRWI
	CSRS
	CSRSI
	CSRC
	CSRCI
	ECALL
	EBREAK
	MV
	SEQZ
	NOT
	SNEZ
	J
	BGT
	BGTU
	BLE
	BLEU
)
