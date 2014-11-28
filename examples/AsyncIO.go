//
//                 Asynchronous I/O Completion Example
//
//  This example shows how to use an asynchronous event handling function
//  that is called when an asynchronous input/output operation completes.
//  Compare this to viRead and viWrite which block the application until
//  either the call returns successfully or a timeout occurs.  Read and
//  write operations can be quite slow sometimes, so these asynchronous
//  operations will allow you processor to perform other tasks.
//  The code uses VISA functions and sets a flag in the callback upon
//  completion of an asynchronous read from a GPIB device to break out of
//  an otherwise infinite loop.  The flow of the code is as follows:
//
//  Open A Session To The VISA Resource Manager
//  Open A Session To A GPIB Device
//  Install A Handler For Asynchronous IO Completion Events
//  Enable Asynchronous IO Completion Events
//  Write A Command To The Instrument
//  Call The Asynchronous Read Command
//  Start A Loop That Can Only Be Broken By A Handler Flag Or Timeout
//  Print Out The Returned Data
//  Close The Instrument Session
//  Close The Resource Manager Session
//

package main

import (
	"fmt"
	"unsafe"

	vi "github.com/jpoirier/visa"
)

var rdCount uint32
var rcount uint32
var stopflag vi.Bool
var statusSession vi.Status

// The handler function. The instrument session, the type of event, and a
// handle to the event are passed to the function along with a user handle
// which is basically a label that could be used to reference the handler.
// The only thing done in the handler is to set a flag that allows the
// program to finish.  Always return VI_SUCCESS from your handler.
func userCB(instr vi.Object, etype, eventContext uint32) {
	fmt.Printf("instr: %d, etype: %d, eventContext: %d\n", instr, etype, eventContext)
	instr.GetAttribute(vi.ATTR_STATUS, unsafe.Pointer(&statusSession))
	instr.GetAttribute(vi.ATTR_RET_COUNT, unsafe.Pointer(&rdCount))
	stopflag = vi.TRUE
}

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
	// GPIB instrument at device 2.  A handle to this session is
	// returned in the handle inst.  Please consult the NI-VISA User Manual
	// for the syntax for using other instruments.
	instr, status := rm.Open("GPIB::2::INSTR", vi.NULL, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Println("An error occurred opening the session to GPIB::2::INSTR")
		return
	}
	defer instr.Close()

	// Now we install the handler for asynchronous i/o completion events.
	// To install the handler, we must pass our instrument session, the type of
	// event to handle, the handler function name and a user handle
	// which acts as a handle to the handler function.
	instr.InstallHandler(vi.EVENT_IO_COMPLETION, userCB)

	// Now we must actually enable the I/O completion event so that our
	// handler will see the events.  Note, one of the parameters is
	// VI_HNDLR indicating that we want the events to be handled by
	// an asynchronous event handler.  The alternate mechanism for handling
	// events is to queue them and read events off of the queue when
	// you want to check them in your program.
	instr.EnableEvent(vi.EVENT_IO_COMPLETION, vi.HNDLR, vi.NULL)

	// Now the VISA write command is used to send a request to the
	// instrument to generate a sine wave.  This demonstrates the
	//  synchronous read operation that blocks the application until viRead()
	//  returns.  Note that the command syntax is instrument specific.

	// Here you specify which string you wish to send to your instrument.
	// The command listed below is device specific. You may have to change
	// command to accommodate your instrument.
	b := []byte("SOUR:FUNC SIN; SENS: DATA?\n")
	instr.Write(b, uint32(len(b)))

	// Next the asynchronous read command is called to read back the
	// date from the instrument.  Immediately after this is called
	// the program goes into a loop which will terminate
	// on an i/o completion event triggering the asynchronous callback.
	// Note that the asynchronous read command returns a job id that is
	// a handle to the asynchronous command.  We can use this handle
	// to terminate the read if too much time has passed.

	buf, job, status := instr.ReadAsync(4096)

	fmt.Printf("\n\nHit enter to continue...")
	var resp int
	fmt.Scanf("%c", &resp)

	// If the asynchronous callback was called and the flag was set
	// we print out the returned data otherwise we terminate the
	// asynchronous job.
	if stopflag == vi.TRUE {
		// rdCount was set in the callback
		fmt.Printf("Count: %d data:  %s", rdCount, string(buf))
	} else {
		instr.Terminate(vi.NULL, uint16(job))
		fmt.Println("The asynchronous read did not complete.")
	}

	fmt.Printf("Exiting...")
}
