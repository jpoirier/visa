//
// Keithley S46 RF Switch example
//
//
//

package main

import (
	"fmt"

	vi "github.com/jpoirier/visa"
	ke "github.com/jpoirier/visa/keithley"
)

func main() {
	// First we must call viOpenDefaultRM to get the resource manager
	// handle.  We will store this handle in defaultRM.
	rm, status := vi.OpenDefaultRM()
	if status < vi.SUCCESS {
		fmt.Println("Could not open a session to the VISA Resource Manager!")
		return
	}
	defer rm.Close()

	//
	instr, status := ke.OpenGpib(rm, 0, 2, vi.NULL, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Println("An error occurred opening the session to GPIB 0:2")
		return
	}
	defer instr.Close()

	instr.Reset()
	instr.OpenAllChans()
	instr.OpenChan(1)
	instr.CloseChan(1)
	list, _ := instr.ClosedChanList()
	fmt.Println("Closed channel list: ", list)
	fmt.Println("Closing Sessions...")
}
