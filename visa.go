// Copyright (c) 2014 Joseph D Poirier
// Distributable under the terms of The simplified BSD License
// that can be found in the LICENSE file.

// Package visa wraps National Instruments VISA (Virtual Instrument Software
// Architecture) driver. The driver allows a client application to communicate
// with most instrumentation buses including GPIB, USB, Serial, and Ethernet.
// VISA is an industry standard for instrument communications.
//
// The package is low level and, for the most part, is one-to-one with the
// exported C functions it wraps. Clients would typically build instrument
// drivers around the package but it can also be used directly.
//
// NI-VISA Drivers:
//     http://www.ni.com/downloads/ni-drivers/
//
// NI-VISA Overview:
//     http://www.ni.com/white-paper/3702/en/

// export CGO_ENABLED=1
// export GOARCH=386

package visa

/*
#cgo linux CFLAGS: -I.
#cgo linux LDFLAGS: -L. -lvisa
#cgo darwin CFLAGS: -I.
#cgo darwin LDFLAGS: -framework VISA
#cgo windows CFLAGS: -I.
#cgo windows LDFLAGS: -LC:/WINDOWS/system32 -lvisa
#include <stdlib.h>
#include <visa.h>

extern void go_cb(ViSession, ViEventType, ViEvent, ViAddr);
ViHndlr get_go_cb() {
	return (ViHndlr)go_cb;
}
*/
import "C"
import "unsafe"

var PackageVersion string = "v0.1"

type ViStatus int32
type Session uint32
type Object uint32
type ViBusAddress C.ViBusAddress
type ViBusSize C.ViBusSize
type ViAttrState C.ViAttrState

type UserCallback func(instr Object, etype, eventContext uint32)
type PUserCallback *UserCallback

// Resource Manager Functions and Operations

// ViOpenDefaultRM returns a session to the Default Resource Manager resource.
func ViOpenDefaultRM() (rm Session, status ViStatus) {
	status = ViStatus(C.viOpenDefaultRM((*C.ViSession)(unsafe.Pointer(&rm))))
	return rm, status
}

var ViGetDefaultRM = ViOpenDefaultRM

// ViFindRsrc queries a VISA system to locate the resources associated with a specified interface.
func (rm Session) ViFindRsrc(expr string) (findList, retCnt uint32, desc string, status ViStatus) {
	cexpr := (*C.ViChar)(C.CString(expr))
	defer C.free(unsafe.Pointer(cexpr))
	d := make([]byte, 257)
	status = ViStatus(C.viFindRsrc(C.ViSession(rm),
		cexpr,
		(*C.ViFindList)(unsafe.Pointer(&findList)),
		(*C.ViUInt32)(unsafe.Pointer(&retCnt)),
		(*C.ViChar)(unsafe.Pointer(&d[0]))))
	return findList, retCnt, string(d), status
}

// ViFindNext gets the next resource from the list of resources found during a
// previous call to viFindRsrc.
func ViFindNext(findList uint32) (string, ViStatus) {
	d := make([]byte, 257)
	status := ViStatus(C.viFindNext((C.ViFindList)(findList),
		(*C.ViChar)(unsafe.Pointer(&d[0]))))
	return string(d), status
}

// ViParseRsrc parses a resource string to get the interface information.
func (rm Session) ViParseRsrc(rsrcName string) (intfType, intfNum uint16, status ViStatus) {
	crsrcName := (*C.ViChar)(C.CString(rsrcName))
	defer C.free(unsafe.Pointer(crsrcName))
	status = ViStatus(C.viParseRsrc(C.ViSession(rm),
		crsrcName,
		(*C.ViUInt16)(unsafe.Pointer(&intfType)),
		(*C.ViUInt16)(unsafe.Pointer(&intfNum))))
	return intfType, intfNum, status
}

// ViParseRsrcEx parses a resource string to get extended interface information.
func (rm Session) ViParseRsrcEx(rsrcName string) (intfType, intfNum uint16,
	rsrcClass, expandedUnaliasedName, aliasIfExists string, status ViStatus) {

	crsrcName := (*C.ViChar)(C.CString(rsrcName))
	defer C.free(unsafe.Pointer(crsrcName))
	r := make([]byte, 257)
	e := make([]byte, 257)
	a := make([]byte, 257)
	status = ViStatus(C.viParseRsrcEx(C.ViSession(rm),
		crsrcName,
		(*C.ViUInt16)(unsafe.Pointer(&intfType)),
		(*C.ViUInt16)(unsafe.Pointer(&intfNum)),
		(*C.ViChar)(unsafe.Pointer(&r[0])),
		(*C.ViChar)(unsafe.Pointer(&e[0])),
		(*C.ViChar)(unsafe.Pointer(&a[0]))))
	return intfType, intfNum, string(r), string(e), string(a), status
}

// ViOpen opens a session to the specified resource.
func (rm Session) ViOpen(name string, mode, timeout uint32) (instr Object, status ViStatus) {
	cname := (*C.ViChar)(C.CString(name))
	defer C.free(unsafe.Pointer(cname))
	status = ViStatus(C.viOpen(C.ViSession(rm),
		cname,
		(C.ViAccessMode)(mode),
		(C.ViUInt32)(timeout),
		(*C.ViSession)(unsafe.Pointer(&instr))))
	return instr, status
}

// Resource Template Operations

// ViClose Closes the specified session, event, or find list.
func (rm Session) ViClose() ViStatus {
	return ViStatus(C.viClose((C.ViObject)(rm)))
}

func (instr Object) ViClose() ViStatus {
	return ViStatus(C.viClose((C.ViObject)(instr)))
}

// ViSetAttribute Sets the state of an attribute.
func (instr Object) ViSetAttribute(attribute, attrState uint32) ViStatus {
	return ViStatus(C.viSetAttribute((C.ViObject)(instr),
		(C.ViAttr)(attribute),
		(C.ViAttrState)(attrState)))
}

// ViGetAttribute Retrieves the state of an attribute.
//
func (instr Object) ViGetAttribute(attrName uint32, attrValue unsafe.Pointer) ViStatus {
	return ViStatus(C.viGetAttribute((C.ViObject)(instr),
		(C.ViAttr)(attrName),
		attrValue))
}

// ViStatusDesc Returns a user-readable description of the status code
// passed to the operation.
func (instr Object) ViStatusDesc(status_in ViStatus) (string, ViStatus) {
	d := make([]byte, 257)
	status := ViStatus(C.viStatusDesc((C.ViObject)(instr),
		(C.ViStatus)(status_in),
		(*C.ViChar)(unsafe.Pointer(&d[0]))))
	return string(d), status
}

// ViTerminate Requests a VISA session to terminate normal execution of an operation.
func (instr Object) ViTerminate(degree, jobId uint16) ViStatus {
	return ViStatus(C.viTerminate((C.ViObject)(instr),
		(C.ViUInt16)(degree),
		(C.ViJobId)(jobId)))
}

// ViLock Establishes an access mode to the specified resource.
func (instr Object) ViLock(lockType, timeout uint32, requestedKey string) (string, ViStatus) {
	crequestedKey := (*C.ViChar)(C.CString(requestedKey))
	defer C.free(unsafe.Pointer(crequestedKey))
	a := make([]byte, 257)
	status := ViStatus(C.viLock((C.ViSession)(instr),
		(C.ViAccessMode)(lockType),
		(C.ViUInt32)(timeout),
		crequestedKey,
		(*C.ViChar)(unsafe.Pointer(&a[0]))))
	return string(a), status
}

// ViUnlock Relinquishes a lock for the specified resource.
func (instr Object) ViUnlock() ViStatus {
	return ViStatus(C.viUnlock((C.ViSession)(instr)))
}

// ViEnableEvent Enables notification of a specified event.
func (instr Object) ViEnableEvent(eventType uint32, mechanism uint16, context uint32) ViStatus {

	return ViStatus(C.viEnableEvent((C.ViSession)(instr),
		(C.ViEventType)(eventType),
		(C.ViUInt16)(mechanism),
		(C.ViEventFilter)(context)))
}

// ViDisableEvent Disables notification of the specified event type(s)
// via the specified mechanism(s).
func (instr Object) ViDisableEvent(eventType uint32, mechanism uint16) ViStatus {
	return ViStatus(C.viDisableEvent((C.ViSession)(instr),
		(C.ViEventType)(eventType),
		(C.ViUInt16)(mechanism)))
}

// ViDiscardEvents Discards event occurrences for specified event types
// and mechanisms in a session.
func (instr Object) ViDiscardEvents(eventType uint32, mechanism uint16) ViStatus {
	return ViStatus(C.viDiscardEvents((C.ViSession)(instr),
		(C.ViEventType)(eventType),
		(C.ViUInt16)(mechanism)))
}

// ViWaitOnEvent Waits for an occurrence of the specified event for a given session.
func (instr Object) ViWaitOnEvent(inEventType, timeout uint32) (outEventType,
	outContext uint32, status ViStatus) {

	status = ViStatus(C.viWaitOnEvent((C.ViSession)(instr),
		(C.ViEventType)(inEventType),
		(C.ViUInt32)(timeout),
		(*C.ViEventType)(unsafe.Pointer(&outEventType)),
		(*C.ViEvent)(unsafe.Pointer(&outContext))))
	return outEventType, outContext, status
}

// ViInstallHandler Installs handlers for event callbacks.
func (instr Object) ViInstallHandler(eventType uint32, userHandle UserCallback) ViStatus {
	return ViStatus(C.viInstallHandler((C.ViSession)(instr),
		(C.ViEventType)(eventType),
		(C.ViHndlr)(C.get_go_cb()),
		(C.ViAddr)(unsafe.Pointer(&userHandle))))
}

// ViUninstallHandler Uninstalls handlers for events.
// Note that VISA identifies handlers uniquely using the userHandle reference.
func (instr Object) ViUninstallHandler(eventType uint32, userHandle UserCallback) ViStatus {
	return ViStatus(C.viUninstallHandler((C.ViSession)(instr),
		(C.ViEventType)(eventType),
		(C.ViHndlr)(C.get_go_cb()),
		(C.ViAddr)(unsafe.Pointer(&userHandle))))
}

// Basic I/O Operations

// ViRead Reads data from device or interface synchronously.
func (instr Object) ViRead(cnt uint32) (buf []byte, retCnt uint32, status ViStatus) {
	buf = make([]byte, cnt)
	status = ViStatus(C.viRead((C.ViSession)(instr),
		(*C.ViByte)(unsafe.Pointer(&buf[0])),
		(C.ViUInt32)(cnt),
		(*C.ViUInt32)(unsafe.Pointer(&retCnt))))
	return buf, retCnt, status
}

// ViReadAsync Reads data from device or interface asynchronously.
func (instr Object) ViReadAsync(cnt uint32) (buf []byte, jobId uint32, status ViStatus) {
	buf = make([]byte, cnt)
	status = ViStatus(C.viReadAsync((C.ViSession)(instr),
		(*C.ViByte)(unsafe.Pointer(&buf[0])),
		(C.ViUInt32)(cnt),
		(*C.ViJobId)(unsafe.Pointer(&jobId))))
	return buf, jobId, status
}

// ViReadToFile Reads data synchronously and stores the transferred data in a file.
func (instr Object) ViReadToFile(filename string, cnt uint32) (retCnt uint32, status ViStatus) {
	cfilename := (*C.ViChar)(C.CString(filename))
	defer C.free(unsafe.Pointer(cfilename))
	status = ViStatus(C.viReadToFile((C.ViSession)(instr),
		cfilename,
		(C.ViUInt32)(cnt),
		(*C.ViUInt32)(unsafe.Pointer(&retCnt))))
	return retCnt, status
}

// ViWrite Writes data to a device or interface synchronously.
func (instr Object) ViWrite(buf []byte, cnt uint32) (retCnt uint32, status ViStatus) {
	status = ViStatus(C.viWrite((C.ViSession)(instr),
		(C.ViBuf)(unsafe.Pointer(&buf[0])),
		(C.ViUInt32)(cnt),
		(*C.ViUInt32)(unsafe.Pointer(&retCnt))))
	return retCnt, status
}

// ViWriteAsync Writes data to a device or interface asynchronously.
func (instr Object) ViWriteAsync(buf []byte, cnt uint32) (jobId uint32, status ViStatus) {
	status = ViStatus(C.viWriteAsync((C.ViSession)(instr),
		(C.ViBuf)(unsafe.Pointer(&buf[0])),
		(C.ViUInt32)(cnt),
		(*C.ViJobId)(unsafe.Pointer(&jobId))))
	return jobId, status
}

// ViWriteFromFile Take data from a file and write it out synchronously.
func (instr Object) ViWriteFromFile(filename string, cnt uint32) (retCnt uint32, status ViStatus) {
	cfilename := (*C.ViChar)(C.CString(filename))
	defer C.free(unsafe.Pointer(cfilename))
	status = ViStatus(C.viWriteFromFile((C.ViSession)(instr),
		cfilename,
		(C.ViUInt32)(cnt),
		(*C.ViUInt32)(unsafe.Pointer(&retCnt))))
	return retCnt, status
}

// ViAssertTrigger Asserts software or hardware trigger.
func (instr Object) ViAssertTrigger(protocol uint16) ViStatus {
	return ViStatus(C.viAssertTrigger((C.ViSession)(instr),
		(C.ViUInt16)(protocol)))
}

// ViReadSTB Reads a status byte of the service request.
func (instr Object) ViReadSTB() (stb_stat uint16, status ViStatus) {
	status = ViStatus(C.viReadSTB((C.ViSession)(instr),
		(*C.ViUInt16)(unsafe.Pointer(&stb_stat))))
	return stb_stat, status
}

// ViClear Clears a device.
func (instr Object) ViClear() ViStatus {
	return ViStatus(C.viClear((C.ViSession)(instr)))
}

// Formatted and Buffered I/O Operations

// ViSetBuf Sets the size for the formatted I/O and/or low-level
// I/O communication buffer(s).
func (instr Object) ViSetBuf(mask uint16, size uint32) ViStatus {
	return ViStatus(C.viSetBuf((C.ViSession)(instr),
		(C.ViUInt16)(mask),
		(C.ViUInt32)(size)))
}

// ViFlush Manually flushes the specified buffers associated with
// formatted I/O operations and/or serial communication.
func (instr Object) ViFlush(mask uint16) ViStatus {
	return ViStatus(C.viFlush((C.ViSession)(instr),
		(C.ViUInt16)(mask)))
}

// ViBufWrite Writes data to a formatted I/O write buffer synchronously.
func (instr Object) ViBufWrite(buf []byte, cnt uint32) (retCnt uint32, status ViStatus) {
	status = ViStatus(C.viBufWrite((C.ViSession)(instr),
		(C.ViBuf)(unsafe.Pointer(&buf[0])),
		(C.ViUInt32)(cnt),
		(*C.ViUInt32)(unsafe.Pointer(&retCnt))))
	return retCnt, status
}

// ViBufRead Reads data from a device or interface through the use of a formatted I/O read buffer.
func (instr Object) ViBufRead(cnt uint32) (buf []byte, retCnt uint32, status ViStatus) {
	buf = make([]byte, cnt)
	status = ViStatus(C.viBufRead((C.ViSession)(instr),
		(C.ViBuf)(unsafe.Pointer(&buf[0])),
		(C.ViUInt32)(cnt),
		(*C.ViUInt32)(unsafe.Pointer(&retCnt))))
	return buf, retCnt, status
}

// ViPrintf Converts, formats, and sends the parameters (designated by ...)
// to the device as specified by the format string.
// ViStatus _VI_FUNCC viPrintf        (ViSession vi, ViString writeFmt, ...);

// ViVPrintf Converts, formats, and sends the parameters designated by params
// to the device or interface as specified by the format string.
// ViStatus _VI_FUNC  viVPrintf       (ViSession vi, ViString writeFmt, ViVAList params);

// ViSPrintf Converts, formats, and sends the parameters (designated by ...)
// to a user-specified buffer as specified by the format string.
// ViStatus _VI_FUNCC viSPrintf       (ViSession vi, ViPBuf buf, ViString writeFmt, ...);

// ViVSPrintf Converts, formats, and sends the parameters designated by params
// to a user-specified buffer as specified by the format string.
// ViStatus _VI_FUNC  viVSPrintf      (ViSession vi, ViPBuf buf, ViString writeFmt,
//                                     ViVAList parms);

// ViScanf Reads, converts, and formats data using the format specifier.
// Stores the formatted data in the parameters (designated by ...).
// ViStatus _VI_FUNCC viScanf         (ViSession vi, ViString readFmt, ...);

// ViVScanf Reads, converts, and formats data using the format specifier.
// Stores the formatted data in the parameters designated by params.
// ViStatus _VI_FUNC  viVScanf        (ViSession vi, ViString readFmt, ViVAList params);

// ViSScanf Reads, converts, and formats data from a user-specified buffer
// using the format specifier. Stores the formatted data in the parameters
// (designated by ...).
// ViStatus _VI_FUNCC viSScanf        (ViSession vi, ViBuf buf, ViString readFmt, ...);

// ViVSScanf Reads, converts, and formats data from a user-specified buffer
// using the format specifier. Stores the formatted data in the parameters
// designated by params.
// ViStatus _VI_FUNC  viVSScanf       (ViSession vi, ViBuf buf, ViString readFmt,
//                                     ViVAList parms);

// ViQueryf Performs a formatted write and read through a single call
// to an operation.
// ViStatus _VI_FUNCC viQueryf        (ViSession vi, ViString writeFmt, ViString readFmt, ...);

// ViVQueryf Performs a formatted write and read through a single call
// to an operation.
// ViStatus _VI_FUNC  viVQueryf       (ViSession vi, ViString writeFmt, ViString readFmt,
//                                     ViVAList params);

// Memory I/O Operations

// ViIn8 Reads in an 8-bit value from the specified memory space and offset.
func (instr Object) ViIn8(space uint16, offset ViBusAddress) (val uint8, status ViStatus) {
	status = ViStatus(C.viIn8((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(*C.ViUInt8)(&val)))
	return val, status
}

// ViOut8 Writes an 8-bit value to the specified memory space and offset.
func (instr Object) viOut8(space uint16, offset ViBusAddress, val uint8) ViStatus {
	return ViStatus(C.viOut8((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(C.ViUInt8)(val)))
}

// ViIn16 Reads in an 16-bit value from the specified memory space and offset.
func (instr Object) ViIn16(space uint16, offset ViBusAddress) (val uint16, status ViStatus) {
	status = ViStatus(C.viIn16((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(*C.ViUInt16)(&val)))
	return val, status
}

// ViOut16 Writes an 16-bit value to the specified memory space and offset.
func (instr Object) viOut16(space uint16, offset ViBusAddress, val uint16) ViStatus {
	return ViStatus(C.viOut16((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(C.ViUInt16)(val)))
}

// ViIn32 Reads in an 32-bit value from the specified memory space and offset.
func (instr Object) ViIn32(space uint16, offset ViBusAddress) (val uint32, status ViStatus) {
	status = ViStatus(C.viIn32((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(*C.ViUInt32)(&val)))
	return val, status
}

// ViOut32 Writes an 32-bit value to the specified memory space and offset.
func (instr Object) viOut32(space uint16, offset ViBusAddress, val uint32) ViStatus {
	return ViStatus(C.viOut32((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(C.ViUInt32)(val)))
}

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViIn64 Reads in an 64-bit value from the specified memory space and offset.
// ViStatus _VI_FUNC  viIn64          (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViPUInt64 val64);
// ViOut64 Writes an 64-bit value to the specified memory space and offset.
// ViStatus _VI_FUNC  viOut64         (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViUInt64  val64);
// ViStatus _VI_FUNC  viIn8Ex         (ViSession vi, ViUInt16 space,
//                                     ViBusAddress64 offset, ViPUInt8  val8);
// ViStatus _VI_FUNC  viOut8Ex        (ViSession vi, ViUInt16 space,
//                                     ViBusAddress64 offset, ViUInt8   val8);
// ViStatus _VI_FUNC  viIn16Ex        (ViSession vi, ViUInt16 space,
//                                     ViBusAddress64 offset, ViPUInt16 val16);
// ViStatus _VI_FUNC  viOut16Ex       (ViSession vi, ViUInt16 space,
//                                     ViBusAddress64 offset, ViUInt16  val16);
// ViStatus _VI_FUNC  viIn32Ex        (ViSession vi, ViUInt16 space,
//                                     ViBusAddress64 offset, ViPUInt32 val32);
// ViStatus _VI_FUNC  viOut32Ex       (ViSession vi, ViUInt16 space,
//                                     ViBusAddress64 offset, ViUInt32  val32);
// ViStatus _VI_FUNC  viIn64Ex        (ViSession vi, ViUInt16 space,
//                                     ViBusAddress64 offset, ViPUInt64 val64);
// ViStatus _VI_FUNC  viOut64Ex       (ViSession vi, ViUInt16 space,
//                                     ViBusAddress64 offset, ViUInt64  val64);
// #endif

// ViMoveIn8 Moves a block of data from the specified address space and offset to local memory.
func (instr Object) ViMoveIn8(space uint16, offset ViBusAddress, length ViBusSize) ([]uint8, ViStatus) {
	buf := make([]uint8, length)
	status := ViStatus(C.viMoveIn8((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(C.ViBusSize)(length),
		(C.ViAUInt8)(unsafe.Pointer(&buf[0]))))
	return buf, status
}

// ViMoveOut8 Moves a block of data from local memory to the specified
func (instr Object) ViMoveOut8(space uint16, offset ViBusAddress, length ViBusSize, buf []uint8) ViStatus {
	return ViStatus(C.viMoveOut8((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(C.ViBusSize)(length),
		(C.ViAUInt8)(unsafe.Pointer(&buf[0]))))
}

// ViMoveIn16 Moves a block of data from the specified address space and offset to local memory.
func (instr Object) ViMoveIn16(space uint16, offset ViBusAddress, length ViBusSize) ([]uint16, ViStatus) {
	buf := make([]uint16, length)
	status := ViStatus(C.viMoveIn16((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(C.ViBusSize)(length),
		(C.ViAUInt16)(unsafe.Pointer(&buf[0]))))
	return buf, status
}

// ViMoveOut16 Moves a block of data from local memory to the specified address space and offset.
func (instr Object) ViMoveOut16(space uint16, offset ViBusAddress, length ViBusSize, buf []uint16) ViStatus {
	return ViStatus(C.viMoveOut16((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(C.ViBusSize)(length),
		(C.ViAUInt16)(unsafe.Pointer(&buf[0]))))
}

// ViMoveIn32 Moves a block of data from the specified address space and offset to local memory.
func (instr Object) ViMoveIn32(space uint16, offset ViBusAddress, length ViBusSize) ([]uint32, ViStatus) {
	buf := make([]uint32, length)
	status := ViStatus(C.viMoveIn32((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(C.ViBusSize)(length),
		(C.ViAUInt32)(unsafe.Pointer(&buf[0]))))
	return buf, status
}

// ViMoveOut32 Moves a block of data from local memory to the specified address space and offset.
func (instr Object) ViMoveOut32(space uint16, offset ViBusAddress, length ViBusSize, buf []uint32) ViStatus {
	return ViStatus(C.viMoveOut32((C.ViSession)(instr),
		(C.ViUInt16)(space),
		(C.ViBusAddress)(offset),
		(C.ViBusSize)(length),
		(C.ViAUInt32)(unsafe.Pointer(&buf[0]))))
}

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViStatus _VI_FUNC  viMoveIn64      (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt64 buf64);
// ViStatus _VI_FUNC  viMoveOut64     (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt64 buf64);
// ViStatus _VI_FUNC  viMoveIn8Ex     (ViSession vi, ViUInt16 space, ViBusAddress64 offset,
//                                     ViBusSize length, ViAUInt8  buf8);
// ViStatus _VI_FUNC  viMoveOut8Ex    (ViSession vi, ViUInt16 space, ViBusAddress64 offset,
//                                     ViBusSize length, ViAUInt8  buf8);
// ViStatus _VI_FUNC  viMoveIn16Ex    (ViSession vi, ViUInt16 space, ViBusAddress64 offset,
//                                     ViBusSize length, ViAUInt16 buf16);
// ViStatus _VI_FUNC  viMoveOut16Ex   (ViSession vi, ViUInt16 space, ViBusAddress64 offset,
//                                     ViBusSize length, ViAUInt16 buf16);
// ViStatus _VI_FUNC  viMoveIn32Ex    (ViSession vi, ViUInt16 space, ViBusAddress64 offset,
//                                     ViBusSize length, ViAUInt32 buf32);
// ViStatus _VI_FUNC  viMoveOut32Ex   (ViSession vi, ViUInt16 space, ViBusAddress64 offset,
//                                     ViBusSize length, ViAUInt32 buf32);
// ViStatus _VI_FUNC  viMoveIn64Ex    (ViSession vi, ViUInt16 space, ViBusAddress64 offset,
//                                     ViBusSize length, ViAUInt64 buf64);
// ViStatus _VI_FUNC  viMoveOut64Ex   (ViSession vi, ViUInt16 space, ViBusAddress64 offset,
//                                     ViBusSize length, ViAUInt64 buf64);
// #endif

// ViMove Moves a block of data.
func (instr Object) ViMove(srcSpace uint16, srcOffset ViBusAddress, srcWidth uint16,
	destSpace uint16, destOffset ViBusAddress, destWidth uint16, srcLength ViBusSize) ViStatus {

	return ViStatus(C.viMove((C.ViSession)(instr),
		(C.ViUInt16)(srcSpace),
		(C.ViBusAddress)(srcOffset),
		(C.ViUInt16)(srcWidth),
		(C.ViUInt16)(destSpace),
		(C.ViBusAddress)(destOffset),
		(C.ViUInt16)(destWidth),
		(C.ViBusSize)(srcLength)))
}

// ViMoveAsync Moves a block of data asynchronously.
func (instr Object) ViMoveAsync(srcSpace uint16, srcOffset ViBusAddress, srcWidth, destSpace uint16,
	destOffset ViBusAddress, destWidth uint16, srcLength ViBusSize) (jobId uint32, status ViStatus) {

	status = ViStatus(C.viMoveAsync((C.ViSession)(instr),
		(C.ViUInt16)(srcSpace),
		(C.ViBusAddress)(srcOffset),
		(C.ViUInt16)(srcWidth),
		(C.ViUInt16)(destSpace),
		(C.ViBusAddress)(destOffset),
		(C.ViUInt16)(destWidth),
		(C.ViBusSize)(srcLength),
		(*C.ViJobId)(unsafe.Pointer(&jobId))))
	return jobId, status
}

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViStatus _VI_FUNC  viMoveEx        (ViSession vi, ViUInt16 srcSpace, ViBusAddress64 srcOffset,
//                                     ViUInt16 srcWidth, ViUInt16 destSpace,
//                                     ViBusAddress64 destOffset, ViUInt16 destWidth,
//                                     ViBusSize srcLength);
// ViStatus _VI_FUNC  viMoveAsyncEx   (ViSession vi, ViUInt16 srcSpace, ViBusAddress64 srcOffset,
//                                     ViUInt16 srcWidth, ViUInt16 destSpace,
//                                     ViBusAddress64 destOffset, ViUInt16 destWidth,
//                                     ViBusSize srcLength, ViPJobId jobId);
// #endif

// ViMapAddress Maps the specified memory space into the process’s address space.
func (instr Object) ViMapAddress(mapSpace uint16, mapOffset ViBusAddress, mapSize ViBusSize,
	access uint16, suggested *byte) (address *byte, status ViStatus) {

	status = ViStatus(C.viMapAddress((C.ViSession)(instr),
		(C.ViUInt16)(mapSpace),
		(C.ViBusAddress)(mapOffset),
		(C.ViBusSize)(mapSize),
		(C.ViBoolean)(access),
		(C.ViAddr)(unsafe.Pointer(suggested)),
		(*C.ViAddr)(unsafe.Pointer(&address))))
	return address, status
}

// ViUnmapAddress Unmaps memory space previously mapped by ViMapAddress().
func (instr Object) ViUnmapAddress() ViStatus {
	return ViStatus(C.viUnmapAddress((C.ViSession)(instr)))
}

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViStatus _VI_FUNC  viMapAddressEx  (ViSession vi, ViUInt16 mapSpace, ViBusAddress64 mapOffset,
//                                     ViBusSize mapSize, ViBoolean access,
//                                     ViAddr suggested, ViPAddr address);
// #endif

// ViPeek8 Reads an 8-bit value from the specified address.
func (instr Object) ViPeek8(address unsafe.Pointer) (val uint8) {
	C.viPeek8((C.ViSession)(instr), (C.ViAddr)(address), (*C.ViUInt8)(&val))
	return val
}

// ViPoke8 Writes an 8-bit value to the specified address.
func (instr Object) ViPoke8(address unsafe.Pointer, val uint8) {
	C.viPoke8((C.ViSession)(instr), (C.ViAddr)(address), (C.ViUInt8)(val))
}

// ViPeek16 Reads an 16-bit value from the specified address.
func (instr Object) ViPeek16(address unsafe.Pointer) (val uint16) {
	C.viPeek16((C.ViSession)(instr), (C.ViAddr)(address), (*C.ViUInt16)(&val))
	return val
}

// ViPoke16 Writes an 16-bit value to the specified address.
func (instr Object) ViPoke16(address unsafe.Pointer, val uint16) {
	C.viPoke16((C.ViSession)(instr), (C.ViAddr)(address), (C.ViUInt16)(val))
}

// ViPeek32 Reads an 32-bit value from the specified address.
func (instr Object) ViPeek32(address unsafe.Pointer) (val uint32) {
	C.viPeek32((C.ViSession)(instr), (C.ViAddr)(address), (*C.ViUInt32)(&val))
	return val
}

// ViPoke32 Writes an 32-bit value to the specified address.
// void     _VI_FUNC  viPoke32        (ViSession vi, ViAddr address, ViUInt32  val32);
func (instr Object) ViPoke32(address unsafe.Pointer, val uint32) {
	C.viPoke32((C.ViSession)(instr), (C.ViAddr)(address), (C.ViUInt32)(val))
}

// #if defined(_VI_INT64_UINT64_DEFINED)
// Reads an 64-bit value from the specified address.
// void     _VI_FUNC  viPeek64        (ViSession vi, ViAddr address, ViPUInt64 val64);

// Writes an 64-bit value to the specified address.
// void     _VI_FUNC  viPoke64        (ViSession vi, ViAddr address, ViUInt64  val64);
// #endif

// Shared Memory Operations

// ViMemAlloc Allocates memory from a device’s memory region.
// ViStatus _VI_FUNC  viMemAlloc      (ViSession vi, ViBusSize size, ViPBusAddress offset);

// ViMemFree Frees memory previously allocated using the viMemAlloc() operation.
// ViStatus _VI_FUNC  viMemFree       (ViSession vi, ViBusAddress offset);

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViStatus _VI_FUNC  viMemAllocEx    (ViSession vi, ViBusSize size, ViPBusAddress64 offset);
// ViStatus _VI_FUNC  viMemFreeEx     (ViSession vi, ViBusAddress64 offset);
// #endif

// Interface Specific Operations

// ViGpibControlREN Controls the state of the GPIB Remote Enable (REN)
// interface line, and optionally the remote/local state of the device.
// ViStatus _VI_FUNC  viGpibControlREN(ViSession vi, ViUInt16 mode);

// ViGpibControlATN Specifies the state of the ATN line and the local
// active controller state.
// ViStatus _VI_FUNC  viGpibControlATN(ViSession vi, ViUInt16 mode);

// ViGpibSendIFC Pulse the interface clear line (IFC) for at least
// 100 microseconds.
// ViStatus _VI_FUNC  viGpibSendIFC   (ViSession vi);

// ViGpibCommand Write GPIB command bytes on the bus.
// ViStatus _VI_FUNC  viGpibCommand   (ViSession vi, ViBuf cmd, ViUInt32 cnt, ViPUInt32 retCnt);

// ViGpibPassControl Tell the GPIB device at the specified address to
// become controller in charge (CIC).
// ViStatus _VI_FUNC  viGpibPassControl(ViSession vi, ViUInt16 primAddr, ViUInt16 secAddr);

// ViVxiCommandQuery Sends the device a miscellaneous command or query and/or
// retrieves the response to a previous query.
// ViStatus _VI_FUNC  viVxiCommandQuery(ViSession vi, ViUInt16 mode, ViUInt32 cmd,
//                                      ViPUInt32 response);

// ViAssertUtilSignal Asserts or deasserts the specified utility bus signal.
// ViStatus _VI_FUNC  viAssertUtilSignal(ViSession vi, ViUInt16 line);

// ViAssertIntrSignal Asserts the specified interrupt or signal.
// ViStatus _VI_FUNC  viAssertIntrSignal(ViSession vi, ViInt16 mode, ViUInt32 statusID);

// ViMapTrigger Map the specified trigger source line to the specified
// destination line.
// ViStatus _VI_FUNC  viMapTrigger    (ViSession vi, ViInt16 trigSrc, ViInt16 trigDest,
//                                     ViUInt16 mode);

// ViUnmapTrigger Undo a previous map from the specified trigger source
// line to the specified destination line.
// ViStatus _VI_FUNC  viUnmapTrigger  (ViSession vi, ViInt16 trigSrc, ViInt16 trigDest);

// ViUsbControlOut Performs a USB control pipe transfer to the device.
// ViStatus _VI_FUNC  viUsbControlOut (ViSession vi, ViInt16 bmRequestType, ViInt16 bRequest,
//                                     ViUInt16 wValue, ViUInt16 wIndex, ViUInt16 wLength,
//                                     ViBuf buf);

// ViUsbControlIn Performs a USB control pipe transfer from the device.
// ViStatus _VI_FUNC  viUsbControlIn  (ViSession vi, ViInt16 bmRequestType, ViInt16 bRequest,
//                                     ViUInt16 wValue, ViUInt16 wIndex, ViUInt16 wLength,
//                                     ViPBuf buf, ViPUInt16 retCnt);

// ViVersion Returns the unformatted resource version number.
func ViVersion() uint32 {
	return uint32(C.VI_SPEC_VERSION)
}

// ViVersMajor Returns the major resource version number.
func ViVersMajor() uint32 {
	return uint32((ViVersion() & 0xFFF00000) >> 20)
}

// ViVersMinor Returns the minor resource version number.
func ViVersMinor() uint32 {
	return uint32((ViVersion() & 0x000FFF00) >> 8)
}

// ViVersSubMinor Returns the sub-minor resource version number.
func ViVersSubMinor() uint32 {
	return uint32((ViVersion() & 0x000000FF))
}

// ViPxiReserveTriggers
// ViStatus _VI_FUNC  viPxiReserveTriggers(ViSession vi, ViInt16 cnt, ViAInt16 trigBuses,
//                                     ViAInt16 trigLines, ViPInt16 failureIndex);

// viVxiServantResponse
// ViStatus _VI_FUNC viVxiServantResponse(ViSession vi, ViInt16 mode, ViUInt32 resp);
