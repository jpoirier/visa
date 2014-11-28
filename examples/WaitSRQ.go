// Synchronous SRQ Event Handling Example
//
//  This example shows how to enable VISA to detect SRQ events.
//  The program writes a command to a device and then waits to receive
//  an SRQ event before trying to read the response.
//
//  Open A Session To The VISA Resource Manager
//  Open A Session To A GPIB Device
//  Enable SRQ Events
//  Write A Command To The Instrument
//  Wait to receive an SRQ event
//  Read the Data
//  Print Out The Data
//  Close The Instrument Session
//  Close The Resource Manager Session

package main

import (
	"fmt"

	vi "github.com/jpoirier/visa"
)

func main() {
	// First we open a session to the VISA resource manager.  We are
	// returned a handle to the resource manager session that we must
	// use to open sessions to specific instruments.
	rm, status := vi.OpenDefaultRM()
	if status < vi.SUCCESS {
		fmt.Println("Could not open a session to the VISA Resource Manager!")
		return
	}
	defer rm.Close()

	// Next we use the resource manager handle to open a session to a
	// GPIB instrument at address 2.  A handle to this session is
	// returned in the handle inst.
	instr, status := rm.Open("GPIB0::2::INSTR", vi.NULL, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Println("1. An error occurred opening the session to GPIB0::2::INSTR...")
		return
	}
	defer instr.Close()

	// Now we must enable the service request event so that VISA
	// will receive the events.  Note: one of the parameters is
	// VI_QUEUE indicating that we want the events to be handled by
	// a synchronous event queue.  The alternate mechanism for handling
	// events is to set up an asynchronous event handling function using
	// the VI_HNDLR option.  The events go into a queue which by default
	// can hold 50 events.  This maximum queue size can be changed with
	// an attribute but it must be called before the events are enabled.
	status = instr.EnableEvent(vi.EVENT_SERVICE_REQ, vi.QUEUE, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Println("The SRQ event could not be enabled")
		return
	}

	// Now the VISA write command is used to send a request to the
	// instrument to generate a sine wave and assert the SRQ line
	// when it is finished.  Notice that this is specific to one
	// particular instrument.
	b := []byte("SRE 0x10; SOUR:FUNC SIN\n")
	_, status = instr.Write(b, uint32(len(b)))
	if status < vi.SUCCESS {
		fmt.Println("Error writing to the instrument")
		return
	}

	// Now we wait for an SRQ event to be received by the event queue.
	// The timeout is in milliseconds and is set to 30000 or 30 seconds.
	// Notice that a handle to the event is returned by the viWaitOnEvent
	// call.  This event handle can be used to obtain various
	// attributes of the event.  Finally, the event handle should be closed
	// to free system resources.
	fmt.Println("\nWaiting for an SRQ Event")
	_, ehandle, status := instr.WaitOnEvent(vi.EVENT_SERVICE_REQ, 30000)

	// If an SRQ event was received we first read the status byte with
	// the viReadSTB function.  This should always be called after
	// receiving a GPIB SRQ event, or subsequent events will not be
	// received properly.  Then the data is read and the event is closed
	// and the data is displayed.  Otherwise sessions are closed and the
	// program terminates.
	if status >= vi.SUCCESS {
		_, status := instr.ReadSTB()
		if status < vi.SUCCESS {
			fmt.Println("There was an error reading the status byte")
			return
		}

		b := []byte("SENS: DATA?\n")
		rcount, status := instr.Write(b, uint32(len(b)))
		if status < vi.SUCCESS {
			fmt.Println("There was an error writing the command to get the data")
			return
		}

		data, rcount, status := instr.Read(3000)
		if status < vi.SUCCESS {
			fmt.Println("There was an error reading the data")
			return
		}
		fmt.Println("Count: %d, Data: %s\n", rcount, data)

		vi.Close(ehandle)
	}
}
