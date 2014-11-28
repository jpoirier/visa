//
// This example demonstrates how you might query your system for
// a particular instrument.  This example queries for all
// GPIB, serial or VXI instruments.  Notice that VISA is able to
// find GPIB and VXI instruments because the instruments have a
// predefined protocol.  But serial instruments do not.  Hence,
// VISA merely indicates that a serial port is available.
//
// The general flow of the code is
//      Open Resource Manager
//      Use viFindRsrc() to query for the first available instrument
//      Open a session to this device
//      Find the next instrument using viFindNext()
//      Open a session to this device.
//      Loop on finding the next instrument until all have been found
//      Close all VISA Sessions
//

package main

import (
	"fmt"
	"os"

	vi "github.com/jpoirier/visa"
)

// static char instrDescriptor[VI_FIND_BUFLEN];
// static ViUInt32 numInstrs;
// static ViFindList findList;
// static ViSession defaultRM, instr;
// static ViStatus status;

func main() {
	rm, status := vi.OpenDefaultRM()
	if status < vi.SUCCESS {
		fmt.Println("Could not open a session to the VISA Resource Manager!")
		os.Exit(0)
	}
	defer rm.Close()

	// Find all the VISA resources in our system and store the number of resources
	// in the system in numInstrs.  Notice the different query descriptions a
	// that are available.

	//  Interface         Expression
	//------------------------------------
	//  GPIB              "GPIB[0-9]*::?*INSTR"
	//  VXI               "VXI?*INSTR"
	//  GPIB-VXI          "GPIB-VXI?*INSTR"
	//  Any VXI           "?*VXI[0-9]*::?*INSTR"
	//  Serial            "ASRL[0-9]*::?*INSTR"
	//  PXI               "PXI?*INSTR"
	//  All instruments   "?*INSTR"
	//  All resources     "?*"
	//
	findList, numInstrs, instrDescriptor, status := rm.FindRsrc("?*INSTR")
	if status < vi.SUCCESS {
		fmt.Println("An error occurred while finding resources.\n")
		return
	}
	defer vi.Close(findList)

	fmt.Printf("%d instruments, serial ports, and other resources found:\n\n", numInstrs)
	fmt.Printf("%s \n", instrDescriptor)

	// Now we will open a session to the instrument we just found.
	instr, status := rm.Open(instrDescriptor, vi.NULL, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Printf("An error occurred opening a session to %s\n", instrDescriptor)
	} else {
		// Now close the session we just opened.
		// In actuality, we would probably use an attribute to determine
		// if this is the instrument we are looking for.
		instr.Close()
	}

	for {
		// stay in this loop until we find all instruments
		instrDescriptor, status = vi.FindNext(findList) // find next desriptor
		if status < vi.SUCCESS {
			// did we find the next resource? */
			fmt.Println("An error occurred finding the next resource.")
			return
		}
		fmt.Printf("%s \n", instrDescriptor)

		// Now we will open a session to the instrument we just found
		instr, status = rm.Open(instrDescriptor, vi.NULL, vi.NULL)
		if status < vi.SUCCESS {
			fmt.Printf("An error occurred opening a session to %s\n", instrDescriptor)
		} else {
			// Now close the session we just opened.
			// In actuality, we would probably use an attribute to determine
			// if this is the instrument we are looking for.
			instr.Close()
		}
		numInstrs -= 1
		if numInstrs == 0 {
			break
		}
	}
}
