//
// Agilent MXA Spectrum Analyzer example
//
//
//

package main

import (
	"fmt"
	"time"

	vi "github.com/jpoirier/visa"
	mxa "github.com/jpoirier/visa/mxa"
)

func main() {
	rm, status := vi.OpenDefaultRM()
	if status < vi.SUCCESS {
		fmt.Println("Could not open a session to the VISA Resource Manager!")
		return
	}
	defer rm.Close()

	//
	instr, status := mxa.OpenGpib(rm, 0, 2, vi.NULL, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Println("An error occurred opening the session to GPIB 0:2")
		return
	}
	defer instr.Close()

	instr.SetScreenTitle("MXA Example")
	instr.ShowSpectrumAnalyzer()
	time.Sleep(10 * time.Second)
	instr.ShowLTEACP()
	time.Sleep(10 * time.Second)
	fmt.Println("Closing Sessions...")
}
