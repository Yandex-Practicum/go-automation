package grader

type MemoryAmount int

const (
	MemoryAmountByte  MemoryAmount = (iota + 1) * 1024
	MemoryAmountKByte
	MemoryAmountMByte
)

