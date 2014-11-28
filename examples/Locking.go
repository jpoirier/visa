// This example demonstrates VISA locking
//
// In VISA, applications can open multiple sessions to the same
// resource simultaneously and access that resource through these
// different sessions concurrently.
// In some cases, an application may need to restrict other
// sessions from accessing that resource.
// For example, an application may need to execute a write and a
// read operation as an atomic operation so that no other
// operations intervene between the write and read operations.
// If another application or even this same application were to
// perform another write between this first write and read, then it
// could put the instrument in an unstable state.
// This is where locking becomes convenient.  The application can
// lock the resource before invoking the write operation and unlock
// it after the read operation, to execute them as a single step.
// This prevents other applications from accessing the resource
// and prevents possible contention.  VISA defines locking
// to restrict accesses to resources for such special circumstances.
//
// The VISA locking mechanism enforces arbitration of accesses to
// resources on an individual basis. If a session locks a resource,
// operations invoked by other sessions are returned with an error.
//
// This VI opens two sessions to an instrument and locks the first
// session.  The first session then writes the String to Write to
// the instrument and reads a response of the desired length.
//
// The second session tries to do the same task but will not
// succeed unless the first session is unlocked.
//
// The general flow of the code is
//      Open Resource Manager
//      Open 2 VISA sessions to an instrument
//      Lock the first session
//      Read and write to that instrument using the first session
//      Unlock the first session
//      Now read and write to the instrument with the second session
//      Close the VISA Sessions

package main

import (
	"fmt"

	vi "github.com/jpoirier/visa"
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

	// Now we will open a session to a GPIB instrument at address 2.
	// We must use the handle from viOpenDefaultRM and we must
	// also use a string that indicates which instrument to open.  This
	// is called the instrument descriptor.  The format for this string
	// can be found in the function panel by right clicking on the
	// description parameter. After opening a session to the
	// device, we will get a handle to the instrument which we
	// will use in later Visa functions.  The two parameters in this
	// function which are reserved for future functionality.
	// These two parameters are given the value VI_NULL.
	instr1, status := rm.Open("GPIB0::2::INSTR", vi.NULL, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Println("1. An error occurred opening the session to GPIB0::2::INSTR...")
		return
	}
	defer instr1.Close()

	// Now we will open another session to a GPIB instrument at address 2.
	instr2, status := rm.Open("GPIB0::2::INSTR", vi.NULL, vi.NULL)
	if status < vi.SUCCESS {
		fmt.Println("2. An error occurred opening the session to GPIB0::2::INSTR...")
		return
	}
	defer instr2.Close()

	// Now we will lock the first session to the resource using the
	// viLock function.  Notice that the locking command takes a parameter
	// which can be set to VI_EXCLUSIVE_LOCK or VI_SHARED_LOCK.  The exclusive
	// lock will only let that session access the device until the
	// lock is released.  The shared locking option uses the last two parameters
	// which are set to VI_NULL for the exclusive lock.  The first of these
	// is a requested access key for the shared lock.  The last parameter is
	// a return value with the actual key assigned to the lock.  If the shared
	// lock is used the return actual key value can be used by another session
	// for the actual key parameter to gain access to the locked resource.
	// Please refer to the NI-VISA User Manual for more information.
	status = instr1.LockExclusive(vi.EXCLUSIVE_LOCK, vi.TMO_IMMEDIATE)
	if status < vi.SUCCESS {
		fmt.Println("Error locking the session...")
		return
	}

	// Now we will write the string "*IDN?\n" to the device and read back the
	// Response from the session that obtained a lock on the resource
	b := []byte("*IDN?\n")
	_, status = instr1.Write(b, uint32(len(b)))
	if status < vi.SUCCESS {
		fmt.Println("Error writing to the device...")
		return
	}

	// Now we will attempt to read back a response from the device to
	// the identification query that was sent.  We will use the viRead
	// function to acquire the data.  We will try to read back 100 bytes.
	// After the data has been read the response is displayed.
	buffer, retCount, status := instr1.Read(100)
	if status < vi.SUCCESS {
		fmt.Println("Error reading a response from the device...")
		return
	} else {
		fmt.Printf("\nCount: %d Data: %s\n", retCount, buffer)
	}

	// Now we will ask the user if they want to unlock the resource.
	// Then we will try the same operations with the second session.  If the
	// resource is not unlocked these operations will fail as would any attempts
	// to modify global attributes by the second session.
	fmt.Printf("Unlock the resource so the second session can also access it(y/n)?: ")
	var resp byte
	fmt.Scanf("%c", &resp)
	if resp == 'y' || resp == 'Y' {
		status = instr1.Unlock()
		if status < vi.SUCCESS {
			fmt.Println("Error unlocking the resource")
			return
		}
	}

	// Now we will attempt the same read and write sequence we attempted with the first
	// session.  If the lock was not removed this will fail
	_, status = instr2.Write(b, uint32(len(b)))
	if status == vi.ERROR_RSRC_LOCKED {
		fmt.Println("The resource is locked, you can't write to it!")
	}

	// Now we will attempt to read back a response from the device to
	// the identification query that was sent.  We will use the viRead
	// function to acquire the data.  We will try to read back 100 bytes.
	// After the data has been read the response is displayed.
	buffer, retCount, status = instr2.Read(100)
	if status == vi.ERROR_RSRC_LOCKED {
		fmt.Println("The resource is still locked you can't read from it!")
	} else {
		fmt.Printf("\nCount: %d, Data: %s\n", retCount, buffer)
	}

	fmt.Println("Finished closing sessions...")
}
