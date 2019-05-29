package tokenizer

import (
	Assembler "github.com/chenhowa/computer/lib/assembly"
)

var mnemonicToToken = map[Mnemonic](Assembler.TokenType){
	RDCYCLE:    Assembler.RDCYCLE,
	RDCYCLEH:   Assembler.RDCYCLEH,
	RDTIME:     Assembler.RDTIME,
	RDTIMEH:    Assembler.RDTIMEH,
	RDINSTRET:  Assembler.RDINSTRET,
	RDINSTRETH: Assembler.RDINSTRETH,
	ADDI:       Assembler.ADDI,
	SLTI:       Assembler.SLTI,
	SLTIU:      Assembler.SLTIU,
	ANDI:       Assembler.ANDI,
	ORI:        Assembler.ORI,
	XORI:       Assembler.XORI,
	SLLI:       Assembler.SLLI,
	SRLI:       Assembler.SRLI,
	SRAI:       Assembler.SRAI,
	LUI:        Assembler.LUI,
	AUIPC:      Assembler.AUIPC,
	ADD:        Assembler.ADD,
	SLT:        Assembler.SLT,
	SLTU:       Assembler.SLTU,
	AND:        Assembler.AND,
	OR:         Assembler.OR,
	XOR:        Assembler.XOR,
	SLL:        Assembler.SLL,
	SRL:        Assembler.SRL,
	SUB:        Assembler.SUB,
	SRA:        Assembler.SRA,
	NOP:        Assembler.NOP,
	JAL:        Assembler.JAL,
	JALR:       Assembler.JALR,
	BEQ:        Assembler.BEQ,
	BNE:        Assembler.BNE,
	BLT:        Assembler.BLT,
	BLTU:       Assembler.BLTU,
	BGE:        Assembler.BGE,
	BGEU:       Assembler.BGEU,
	LW:         Assembler.LW,
	LH:         Assembler.LH,
	LHU:        Assembler.LHU,
	LB:         Assembler.LB,
	LBU:        Assembler.LBU,
	SW:         Assembler.SW,
	SH:         Assembler.SH,
	SB:         Assembler.SB,
	CSRRW:      Assembler.CSRRW,
	CSRRS:      Assembler.CSRRS,
	CSRRC:      Assembler.CSRRC,
	CSRRWI:     Assembler.CSRRWI,
	CSRRSI:     Assembler.CSRRSI,
	CSRRCI:     Assembler.CSRRCI,
	CSRR:       Assembler.CSRR,
	CSRW:       Assembler.CSRW,
	CSRWI:      Assembler.CSRWI,
	CSRS:       Assembler.CSRS,
	CSRSI:      Assembler.CSRSI,
	CSRC:       Assembler.CSRC,
	CSRCI:      Assembler.CSRCI,
	ECALL:      Assembler.ECALL,
	EBREAK:     Assembler.EBREAK,
	MV:         Assembler.MV,
	SEQZ:       Assembler.SEQZ,
	NOT:        Assembler.NOT,
	SNEZ:       Assembler.SNEZ,
	J:          Assembler.J,
	BGT:        Assembler.BGT,
	BGTU:       Assembler.BGTU,
	BLE:        Assembler.BLE,
	BLEU:       Assembler.BLEU,
}
