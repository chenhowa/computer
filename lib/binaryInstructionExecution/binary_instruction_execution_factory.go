package binaryInstructionExecution

import "fmt"
import Producer "github.com/chenhowa/os/computer/binaryInstructionExecution/executionFactoryProducers"
import Parser "github.com/chenhowa/os/computer/binaryInstructionExecution/instructionParsing"

/*
This file is tightly coupled with the `instructionParsing` and `executionFactoryProducers` package.
*/

/*RiscVBinaryInstructionExecutionFactory is responsible for producing executors of instructions that
can be executed at later times.
*/
type RiscVBinaryInstructionExecutionFactory struct {
	executor Producer.RiscVExecutor
}

/*MakeRiscVInstructionExecutionFactory is a constructor of an instance of RiscVBinaryInstructionExecutionFactory*/
func MakeRiscVInstructionExecutionFactory(executor Producer.RiscVExecutor) RiscVBinaryInstructionExecutionFactory {
	factory := RiscVBinaryInstructionExecutionFactory{
		executor: executor,
	}

	return factory
}

type binaryExecutor interface {
	Execute()
}

/*Produce generates an executor that will execute the `instruction`
 */
func (factory *RiscVBinaryInstructionExecutionFactory) Produce(instruction uint32) binaryExecutor {
	parser := Parser.RiscVBinaryInstructionParser{}
	result := parser.Parse(instruction)

	var executor binaryExecutor

	switch result.InstructionType {
	case Parser.I:
		executor = factory.produceI(result, factory.executor)
	case Parser.U:
		executor = factory.produceU(result, factory.executor)
	case Parser.R:
		executor = factory.produceR(result, factory.executor)
	case Parser.J:
		executor = factory.produceJ(result, factory.executor)
	case Parser.B:
		executor = factory.produceB(result, factory.executor)
	case Parser.S:
		executor = factory.produceS(result, factory.executor)
	default:
		panic(fmt.Sprintf("unrecognized instruction type: %d", result.InstructionType))
	}

	return executor
}

func (factory *RiscVBinaryInstructionExecutionFactory) produceI(result Parser.RiscVBinaryParseResult, ex Producer.RiscVExecutor) binaryExecutor {

	var executor = Producer.ExecutorI{
		Executor: ex,
		Result:   result,
	}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

func (factory *RiscVBinaryInstructionExecutionFactory) produceU(result Parser.RiscVBinaryParseResult, ex Producer.RiscVExecutor) binaryExecutor {
	var executor = Producer.ExecutorU{
		Executor: ex,
		Result:   result,
	}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

func (factory *RiscVBinaryInstructionExecutionFactory) produceR(result Parser.RiscVBinaryParseResult, ex Producer.RiscVExecutor) binaryExecutor {

	var executor = Producer.ExecutorR{
		Executor: ex,
		Result:   result,
	}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

func (factory *RiscVBinaryInstructionExecutionFactory) produceJ(result Parser.RiscVBinaryParseResult, ex Producer.RiscVExecutor) binaryExecutor {

	var executor = Producer.ExecutorJ{
		Executor: ex,
		Result:   result,
	}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

func (factory *RiscVBinaryInstructionExecutionFactory) produceB(result Parser.RiscVBinaryParseResult, ex Producer.RiscVExecutor) binaryExecutor {

	var executor = Producer.ExecutorB{
		Executor: factory.executor,
		Result:   result,
	}
	// This works because Go does Pointer Escape analysis.
	return &executor
}

func (factory *RiscVBinaryInstructionExecutionFactory) produceS(result Parser.RiscVBinaryParseResult, ex Producer.RiscVExecutor) binaryExecutor {

	var executor = Producer.ExecutorS{
		Executor: factory.executor,
		Result:   result,
	}
	// This works because Go does Pointer Escape analysis.
	return &executor
}
