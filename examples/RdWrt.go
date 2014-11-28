//
// Read and Write to an Instrument Example
//
// This code demonstrates synchronous read and write commands to a
// GPIB, serial or message-based VXI instrument using VISA.
//
// The general flow of the code is
//      Open Resource Manager
//      Open VISA Session to an Instrument
//      Write the Identification Query Using viWrite
//      Try to Read a Response With viRead
//      Close the VISA Session
//

package main

import (
	"fmt"

	vi "github.com/jpoirier/visa"
)

// In every source code or header file that you use it is necessary to prototype
// your VISA variables at the beginning of the file. You need to declare the VISA
// session, VISA integers, VISA strings, VISA pointers, and VISA floating variables.
// Remember that if you are prototyping variables that are to be used as part of the
// VISA session that need this prototyping. As an example, above retCount has been
// prototyped as a static variable to this particular module.   It is an integer of
// bit length 32. If you are uncertain how to declare your VISA prototypes refer
// to the VISA help under the Section titled Type Assignments Table. The VISA
// help is located in your NI-VISA directory or folder.

func main() {
	// First we must call viOpenDefaultRM to get the resource manager
	// handle.  We will store this handle in defaultRM.
	rm, status := vi.OpenDefaultRM()
	if status < vi.SUCCESS {
		fmt.Println("Could not open a session to the VISA Resource Manager!")
		return
	}
	defer rm.Close()

	// Now we will open a VISA session to a device at Primary Address 2.
	// You can use any address for your instrument. In this example we are
	// using GPIB Primary Address 2.
	//
	// We must use the handle from viOpenDefaultRM and we must
	// also use a string that indicates which instrument to open.  This
	// is called the instrument descriptor.  The format for this string
	// can be found in the NI-VISA User Manual.
	// After opening a session to the device, we will get a handle to
	// the instrument which we will use in later VISA functions.
	// The two parameters in this function which are left blank are
	// reserved for future functionality.  These two parameters are
	// given the value VI_NULL.
	//
	// This example will also work for serial or VXI instruments by changing
	// the instrument descriptor from GPIB0::2::INSTR to ASRL1::INSTR or
	// VXI0::2::INSTR depending on the necessary descriptor for your
	// instrument.
	instr, status := rm.Open("GPIB0::2::INSTR", vi.NULL, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Println("An error occurred opening the session to GPIB0::2::INSTR")
		return
	}
	defer instr.Close()

	// Set timeout value to 5000 milliseconds (5 seconds).
	instr.SetAttribute(vi.ATTR_TMO_VALUE, 5000)

	// At this point we now have a session open to the instrument at
	// Primary Address 2.  We can use this session handle to write
	// an ASCII command to the instrument.  We will use the viWrite function
	// to send the string "*IDN?", asking for the device's identification.
	b := []byte("*IDN?")
	_, status = instr.Write(b, uint32(len(b)))
	if status < vi.SUCCESS {
		fmt.Println("Error writing to the device, exiting...")
		return
	}

	instr.Write(b, uint32(len(b)))

	// Now we will attempt to read back a response from the device to
	// the identification query that was sent.  We will use the viRead
	// function to acquire the data.  We will try to read back 100 bytes.
	// After the data has been read the response is displayed.

	buffer, retCount, status := instr.Read(100)
	if status < vi.SUCCESS {
		fmt.Println("Error reading a response from the device")
	} else {
		fmt.Printf("Count: %d Data: %*s\n", retCount, buffer)
	}

	// Now we will close the session to the instrument using
	// viClose. This operation frees all system resources.

	fmt.Println("Closing Sessions...")
}
