//
// Read and Write to a USBTMC Instrument
//
// This code demonstrates sending synchronous read & write commands
// to an USB Test & Measurement Class (USBTMC) instrument using
// NI-VISA
// The example writes the "*IDN?\n" string to all the USBTMC
// devices connected to the system and attempts to read back
// results using the write and read functions.
//
// The general flow of the code is
//    Open Resource Manager
//    Open VISA Session to an Instrument
//    Write the Identification Query Using viWrite
//    Try to Read a Response With viRead
//    Close the VISA Session

package main

import (
	"fmt"

	vi "github.com/jpoirier/visa"
)

func main() {
	// First we must call viOpenDefaultRM to get the manager
	// handle.  We will store this handle in defaultRM.
	rm, status := vi.OpenDefaultRM()
	if status < vi.SUCCESS {
		fmt.Println("Could not open a session to the VISA Resource Manager!")
		return
	}
	defer rm.Close()

	// Find all the USB TMC VISA resources in our system and store the
	// number of resources in the system in numInstrs.
	findList, numInstrs, instrResourceString, status := rm.FindRsrc("USB?*INSTR")
	if status < vi.SUCCESS {
		fmt.Println("An error occurred while finding resources.\n")
		return
	}
	defer vi.Close(findList)

	// Now we will open VISA sessions to all USB TMC instruments.
	// We must use the handle from viOpenDefaultRM and we must
	// also use a string that indicates which instrument to open.  This
	// is called the instrument descriptor.  The format for this string
	// can be found in the function panel by right clicking on the
	// descriptor parameter. After opening a session to the
	// device, we will get a handle to the instrument which we
	// will use in later VISA functions.  The AccessMode and Timeout
	// parameters in this function are reserved for future
	// functionality.  These two parameters are given the value VI_NULL.

	for i := 0; i < int(numInstrs); i++ {
		if i > 0 {
			instrResourceString, _ = vi.FindNext(findList)
		}

		instr, status := rm.Open(instrResourceString, vi.NULL, vi.NULL)
		if status < vi.SUCCESS {
			fmt.Printf("Cannot open a session to the device %d.\n", i+1)
			continue
		}

		// At this point we now have a session open to the USB TMC instrument.
		// We will now use the viWrite function to send the device the string "*IDN?\n",
		// asking for the device's identification.
		b := []byte("*IDN?\n")
		_, status = instr.Write(b, uint32(len(b)))
		if status < vi.SUCCESS {
			fmt.Printf("Error writing to the device %d.\n", i+1)
			instr.Close()
			continue
		}

		// Now we will attempt to read back a response from the device to
		// the identification query that was sent.  We will use the viRead
		// function to acquire the data.  We will try to read back 100 bytes.
		// This function will stop reading if it finds the termination character
		// before it reads 100 bytes.
		// After the data has been read the response is displayed.
		buffer, retCount, status := instr.Read(100)
		if status < vi.SUCCESS {
			fmt.Printf("Error reading a response from the device %d.\n", i+1)
		} else {
			fmt.Printf("\nDevice %d: %*s\n", i+1, retCount, buffer)
		}
		instr.Close()
	}
}
