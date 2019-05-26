package main

import (
	"fmt"
	"os"
)

/*main describes an application simulates a 32-bit Risc-V CPU with 16-bit memory (only 65 KiB RAM).
Once the simulation is done, the application will write the state of the memory and the registers
to standard output, so be sure to redirect.
The simulation run in two modes: interactive and file.

- In interactive mode, the user is repeatedly prompted to choose between executing the
current instruction referenced by the Program Counter, and the instruction they can choose
to write directly to the address of the current instruction. These instructions will be RISC-V
assembly instructions for both input and output.

- In file mode, the user must enter just one command line argument that is a valid path to a
file. The application will evaluate the contents of the file for a valid Risc-V assembly program, and
if valid, it will load and run the program as binary instructions to the simulated CPU.
*/
func main() {
	fmt.Println("Hello, world!")
	argsWithoutProg := os.Args[1:]
	fmt.Println(len(argsWithoutProg))
	for i := uint(0); i < uint(len(argsWithoutProg)); i++ {
		fmt.Println(argsWithoutProg[i])
	}
}
