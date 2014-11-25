// Copyright (c) 2014 Joseph D Poirier
// Distributable under the terms of The simplified BSD License
// that can be found in the LICENSE file.

// Package visa wraps National Instruments VISA (Virtual Instrument Software
// Architecture) driver. The driver allows a client application to communicate
// with most instrumentation buses including GPIB, USB, Serial, and Ethernet.
// VISA is an industry standard for instrument communications.
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
#cgo linux LDFLAGS: -lvisa -L.
#cgo darwin CFLAGS: -I.
#cgo darwin LDFLAGS: -framework VISA
#cgo windows CFLAGS: -I.
#cgo windows LDFLAGS: -lvisa -LC:/WINDOWS/system32
#include <stdlib.h>
#include <visa.h>
*/
import "C"
import "unsafe"

var PackageVersion string = "v0.1"

type ViStatus int32
type Session uint32
type Object uint32

// Resource Manager Functions and Operations

// ViOpenDefaultRM returns a session to the Default Resource Manager resource.
func ViOpenDefaultRM() (sesn Session, status ViStatus) {
	status = ViStatus(C.viOpenDefaultRM((C.ViPSession)(unsafe.Pointer(&sesn))))
	return sesn, status
}

var viGetDefaultRM = ViOpenDefaultRM

// ViFindRsrc queries a VISA system to locate the resources associated with a
// specified interface.
func (sesn Session) ViFindRsrc(expr string) (findList uint32, retCnt uint32,
	desc string, status ViStatus) {

	cexpr := (*C.ViChar)(C.CString(expr))
	defer C.free(unsafe.Pointer(cexpr))
	d := make([]byte, 257)
	status = ViStatus(C.viFindRsrc(C.ViSession(sesn),
		cexpr,
		(*C.ViFindList)(unsafe.Pointer(&findList)),
		(*C.ViUInt32)(unsafe.Pointer(&retCnt)),
		(*C.ViChar)(unsafe.Pointer(&d[0]))))
	return findList, retCnt, string(d), status
}

// ViFindNext gets the next resource from the list of resources found during a
// previous call to viFindRsrc.
func ViFindNext(findList uint32) (desc string, status ViStatus) {
	d := make([]byte, 257)
	status = ViStatus(C.viFindNext((C.ViFindList)(findList),
		(*C.ViChar)(unsafe.Pointer(&d[0]))))
	return string(d), status
}

// ViParseRsrc parses a resource string to get the interface information.
func (sesn Session) ViParseRsrc(rsrcName string) (intfType uint16,
	intfNum uint16, status ViStatus) {

	crsrcName := (*C.ViChar)(C.CString(rsrcName))
	defer C.free(unsafe.Pointer(crsrcName))
	status = ViStatus(C.viParseRsrc(C.ViSession(sesn),
		crsrcName,
		(*C.ViUInt16)(unsafe.Pointer(&intfType)),
		(*C.ViUInt16)(unsafe.Pointer(&intfNum))))
	return intfType, intfNum, status
}

// ViParseRsrcEx parses a resource string to get extended interface information.
func (sesn Session) ViParseRsrcEx(rsrcName string) (intfType uint16,
	intfNum uint16, rsrcClass string, expandedUnaliasedName string,
	aliasIfExists string, status ViStatus) {

	crsrcName := (*C.ViChar)(C.CString(rsrcName))
	defer C.free(unsafe.Pointer(crsrcName))
	r := make([]byte, 257)
	e := make([]byte, 257)
	a := make([]byte, 257)
	status = ViStatus(C.viParseRsrcEx(C.ViSession(sesn),
		crsrcName,
		(*C.ViUInt16)(unsafe.Pointer(&intfType)),
		(*C.ViUInt16)(unsafe.Pointer(&intfNum)),
		(*C.ViChar)(unsafe.Pointer(&r[0])),
		(*C.ViChar)(unsafe.Pointer(&e[0])),
		(*C.ViChar)(unsafe.Pointer(&a[0]))))
	return intfType, intfNum, string(r), string(e), string(a), status
}

// ViOpen opens a session to the specified resource.
func (sesn Session) ViOpen(name string, mode, timeout uint32) (instr Object,
	status ViStatus) {

	cname := (*C.ViChar)(C.CString(name))
	defer C.free(unsafe.Pointer(cname))
	status = ViStatus(C.viOpen(C.ViSession(sesn),
		cname,
		(C.ViAccessMode)(mode),
		(C.ViUInt32)(timeout),
		(*C.ViSession)(unsafe.Pointer(&instr))))
	return instr, status
}

// Resource Template Operations

// ViClose Closes the specified session, event, or find list.
func (sesn Session) ViClose() (status ViStatus) {
	status = ViStatus(C.viClose((C.ViObject)(sesn)))
	return status
}

func (instr Object) ViClose() (status ViStatus) {
	status = ViStatus(C.viClose((C.ViObject)(instr)))
	return status
}

// ViSetAttribute Sets the state of an attribute.
func (instr Object) ViSetAttribute(attribute, attrState uint32) (status ViStatus) {
	status = ViStatus(C.viSetAttribute((C.ViObject)(instr),
		(C.ViAttr)(attribute),
		(C.ViAttrState)(attrState)))
	return status
}

// ViGetAttribute Retrieves the state of an attribute.
// vi in
// attrName in
// attrValue out
// ViStatus _VI_FUNC  viGetAttribute  (ViObject vi, ViAttr attrName, void _VI_PTR attrValue);
// func (instr Object) ViGetAttribute(attribute uint32) {
// 	cvi := (C.ViObject)(vi)
// 	cattribute := (C.ViAttr)(attribute)

// 	return
// }

// ViStatusDesc Returns a user-readable description of the status code
// passed to the operation.
func (instr Object) ViStatusDesc(status ViStatus) (desc string) {
	d := make([]byte, 257)
	status = ViStatus(C.viStatusDesc((C.ViObject)(instr),
		(C.ViStatus)(status),
		(*C.ViChar)(unsafe.Pointer(&d[0]))))
	return string(d)
}

// ViTerminate Requests a VISA session to terminate normal execution of an operation.
func (instr Object) ViTerminate(degree, jobId uint16) (status ViStatus) {
	status = ViStatus(C.viTerminate((C.ViObject)(instr),
		(C.ViUInt16)(degree),
		(C.ViJobId)(jobId)))
	return status
}

// ViLock Establishes an access mode to the specified resource.
func (instr Object) ViLock(lockType, timeout uint32, requestedKey string) (accessKey string,
	status ViStatus) {

	crequestedKey := (*C.ViChar)(C.CString(requestedKey))
	defer C.free(unsafe.Pointer(crequestedKey))

	a := make([]byte, 257)
	status = ViStatus(C.viLock((C.ViSession)(instr),
		(C.ViAccessMode)(lockType),
		(C.ViUInt32)(timeout),
		crequestedKey,
		(*C.ViChar)(unsafe.Pointer(&a[0]))))
	return string(a), status
}

// ViUnlock Relinquishes a lock for the specified resource.
func (instr Object) ViUnlock() (status ViStatus) {
	status = ViStatus(C.viUnlock((C.ViSession)(instr)))
	return status
}

// ViEnableEvent Enables notification of a specified event.
func (instr Object) ViEnableEvent(eventType uint32, mechanism uint16,
	context uint32) (status ViStatus) {

	status = ViStatus(C.viEnableEvent((C.ViSession)(instr),
		(C.ViEventType)(eventType),
		(C.ViUInt16)(mechanism),
		(C.ViEventFilter)(context)))
	return status
}

// ViDisableEvent Disables notification of the specified event type(s)
// via the specified mechanism(s).
func (instr Object) ViDisableEvent(eventType uint32, mechanism uint16) (status ViStatus) {
	status = ViStatus(C.viDisableEvent((C.ViSession)(instr),
		(C.ViEventType)(eventType),
		(C.ViUInt16)(mechanism)))
	return status
}

// ViDiscardEvents Discards event occurrences for specified event types
// and mechanisms in a session.
func (instr Object) ViDiscardEvents(eventType uint32, mechanism uint16) (status ViStatus) {
	status = ViStatus(C.viDiscardEvents((C.ViSession)(instr),
		(C.ViEventType)(eventType),
		(C.ViUInt16)(mechanism)))
	return status
}

// ViWaitOnEvent Waits for an occurrence of the specified event for a given session.
func (instr Object) ViWaitOnEvent(inEventType, timeout uint32) (outEventType,
	outContext uint32, status ViStatus) {

	status = ViStatus(C.viWaitOnEvent((C.ViSession)(instr),
		(C.ViEventType)(inEventType),
		(C.ViUInt32)(timeout),
		(C.ViPEventType)(unsafe.Pointer(&outEventType)),
		(C.ViPEvent)(unsafe.Pointer(&outContext))))
	return outEventType, outContext, status
}

// ViInstallHandler Installs handlers for event callbacks.
// vi in
// eventType in
// handler in
// userHandle in
// ViStatus _VI_FUNC  viInstallHandler(ViSession vi, ViEventType eventType, ViHndlr handler,
//                                     ViAddr userHandle);
// func (instr Object) ViInstallHandler(vi, eventType uint32) {
// 	cvi := (C.ViSession)(vi)
// 	return
// }

// ViUninstallHandler Uninstalls handlers for events.
// vi in
// eventType in
// handler in
// userHandle in
// ViStatus _VI_FUNC  viUninstallHandler(ViSession vi, ViEventType eventType, ViHndlr handler,
//                                       ViAddr userHandle);
// func (instr Object) ViUninstallHandler(vi uint32) {
// 	cvi := (C.ViSession)(vi)
// 	return
// }

// Basic I/O Operations

// ViRead Reads data from device or interface synchronously.
func (instr Object) ViRead(cnt uint32) (buf []byte, retCnt uint32, status ViStatus) {
	b := make([]byte, cnt)
	status = ViStatus(C.viRead((C.ViSession)(instr),
		(*C.ViByte)(unsafe.Pointer(&b[0])),
		(C.ViUInt32)(cnt),
		(*C.ViUInt32)(unsafe.Pointer(&retCnt))))
	return buf, retCnt, status
}

// ViReadAsync Reads data from device or interface asynchronously.
func (instr Object) ViReadAsync(cnt uint32) (buf []byte, jobId uint32, status ViStatus) {
	b := make([]byte, cnt)
	status = ViStatus(C.viReadAsync((C.ViSession)(instr),
		(*C.ViByte)(unsafe.Pointer(&b[0])),
		(C.ViUInt32)(cnt),
		(*C.ViJobId)(unsafe.Pointer(&jobId))))
	return buf, jobId, status
}

// ViReadToFile
// ViStatus _VI_FUNC  viReadToFile(ViSession vi, ViConstString filename, ViUInt32 cnt,
//                                     ViPUInt32 retCnt);

// ViWrite
// ViStatus _VI_FUNC  viWrite(ViSession vi, ViBuf  buf, ViUInt32 cnt, ViPUInt32 retCnt);

// ViWriteAsync
// ViStatus _VI_FUNC  viWriteAsync(ViSession vi, ViBuf  buf, ViUInt32 cnt, ViPJobId  jobId);

// ViWriteFromFile
// ViStatus _VI_FUNC  viWriteFromFile(ViSession vi, ViConstString filename, ViUInt32 cnt,
//                                     ViPUInt32 retCnt);

// ViAssertTrigger
// ViStatus _VI_FUNC  viAssertTrigger(ViSession vi, ViUInt16 protocol);

// ViReadSTB
// ViStatus _VI_FUNC  viReadSTB(ViSession vi, ViPUInt16 status);

// ViClear
// ViStatus _VI_FUNC  viClear(ViSession vi);

// Formatted and Buffered I/O Operations

// ViSetBuf
// ViStatus _VI_FUNC  viSetBuf        (ViSession vi, ViUInt16 mask, ViUInt32 size);

// ViFlush
// ViStatus _VI_FUNC  viFlush         (ViSession vi, ViUInt16 mask);

// ViBufWrite
// ViStatus _VI_FUNC  viBufWrite      (ViSession vi, ViBuf  buf, ViUInt32 cnt, ViPUInt32 retCnt);

// ViBufRead
// ViStatus _VI_FUNC  viBufRead       (ViSession vi, ViPBuf buf, ViUInt32 cnt, ViPUInt32 retCnt);

// ViPrintf
// ViStatus _VI_FUNCC viPrintf        (ViSession vi, ViString writeFmt, ...);

// ViVPrintf
// ViStatus _VI_FUNC  viVPrintf       (ViSession vi, ViString writeFmt, ViVAList params);

// ViSPrintf
// ViStatus _VI_FUNCC viSPrintf       (ViSession vi, ViPBuf buf, ViString writeFmt, ...);

// ViVSPrintf
// ViStatus _VI_FUNC  viVSPrintf      (ViSession vi, ViPBuf buf, ViString writeFmt,
//                                     ViVAList parms);

// ViScanf
// ViStatus _VI_FUNCC viScanf         (ViSession vi, ViString readFmt, ...);

// ViVScanf
// ViStatus _VI_FUNC  viVScanf        (ViSession vi, ViString readFmt, ViVAList params);

// ViSScanf
// ViStatus _VI_FUNCC viSScanf        (ViSession vi, ViBuf buf, ViString readFmt, ...);

// ViVSScanf
// ViStatus _VI_FUNC  viVSScanf       (ViSession vi, ViBuf buf, ViString readFmt,
//                                     ViVAList parms);

// ViQueryf
// ViStatus _VI_FUNCC viQueryf        (ViSession vi, ViString writeFmt, ViString readFmt, ...);

// ViVQueryf
// ViStatus _VI_FUNC  viVQueryf       (ViSession vi, ViString writeFmt, ViString readFmt,
//                                     ViVAList params);

// Memory I/O Operations

// ViIn8
// ViStatus _VI_FUNC  viIn8           (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViPUInt8  val8);

// ViOut8
// ViStatus _VI_FUNC  viOut8          (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViUInt8   val8);

// ViIn16
// ViStatus _VI_FUNC  viIn16          (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViPUInt16 val16);

// ViOut16
// ViStatus _VI_FUNC  viOut16         (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViUInt16  val16);

// ViIn32
// ViStatus _VI_FUNC  viIn32          (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViPUInt32 val32);

// ViOut32
// ViStatus _VI_FUNC  viOut32         (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViUInt32  val32);

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViStatus _VI_FUNC  viIn64          (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViPUInt64 val64);
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

// ViMoveIn8
// ViStatus _VI_FUNC  viMoveIn8       (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt8  buf8);

// ViMoveOut8
// ViStatus _VI_FUNC  viMoveOut8      (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt8  buf8);

// ViMoveIn16
// ViStatus _VI_FUNC  viMoveIn16      (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt16 buf16);

// ViMoveOut16
// ViStatus _VI_FUNC  viMoveOut16     (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt16 buf16);

// ViMoveIn32
// ViStatus _VI_FUNC  viMoveIn32      (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt32 buf32);

// ViMoveOut32
// ViStatus _VI_FUNC  viMoveOut32     (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt32 buf32);

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

// ViMove
// ViStatus _VI_FUNC  viMove          (ViSession vi, ViUInt16 srcSpace, ViBusAddress srcOffset,
//                                     ViUInt16 srcWidth, ViUInt16 destSpace,
//                                     ViBusAddress destOffset, ViUInt16 destWidth,
//                                     ViBusSize srcLength);

// ViMoveAsync
// ViStatus _VI_FUNC  viMoveAsync     (ViSession vi, ViUInt16 srcSpace, ViBusAddress srcOffset,
//                                     ViUInt16 srcWidth, ViUInt16 destSpace,
//                                     ViBusAddress destOffset, ViUInt16 destWidth,
//                                     ViBusSize srcLength, ViPJobId jobId);

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

// ViMapAddress
// ViStatus _VI_FUNC  viMapAddress    (ViSession vi, ViUInt16 mapSpace, ViBusAddress mapOffset,
//                                     ViBusSize mapSize, ViBoolean access,
//                                     ViAddr suggested, ViPAddr address);

// ViUnmapAddress
// ViStatus _VI_FUNC  viUnmapAddress  (ViSession vi);

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViStatus _VI_FUNC  viMapAddressEx  (ViSession vi, ViUInt16 mapSpace, ViBusAddress64 mapOffset,
//                                     ViBusSize mapSize, ViBoolean access,
//                                     ViAddr suggested, ViPAddr address);
// #endif

// ViPeek8
// void     _VI_FUNC  viPeek8         (ViSession vi, ViAddr address, ViPUInt8  val8);

// ViPoke8
// void     _VI_FUNC  viPoke8         (ViSession vi, ViAddr address, ViUInt8   val8);

// ViPeek16
// void     _VI_FUNC  viPeek16        (ViSession vi, ViAddr address, ViPUInt16 val16);

// ViPoke16
// void     _VI_FUNC  viPoke16        (ViSession vi, ViAddr address, ViUInt16  val16);

// ViPeek32
// void     _VI_FUNC  viPeek32        (ViSession vi, ViAddr address, ViPUInt32 val32);

// VSiPoke32
// void     _VI_FUNC  viPoke32        (ViSession vi, ViAddr address, ViUInt32  val32);

// #if defined(_VI_INT64_UINT64_DEFINED)
// void     _VI_FUNC  viPeek64        (ViSession vi, ViAddr address, ViPUInt64 val64);
// void     _VI_FUNC  viPoke64        (ViSession vi, ViAddr address, ViUInt64  val64);
// #endif

// Shared Memory Operations

// ViMemAlloc
// ViStatus _VI_FUNC  viMemAlloc      (ViSession vi, ViBusSize size, ViPBusAddress offset);

// ViMemFree
// ViStatus _VI_FUNC  viMemFree       (ViSession vi, ViBusAddress offset);

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViStatus _VI_FUNC  viMemAllocEx    (ViSession vi, ViBusSize size, ViPBusAddress64 offset);
// ViStatus _VI_FUNC  viMemFreeEx     (ViSession vi, ViBusAddress64 offset);
// #endif

// Interface Specific Operations

// viGpibControlREN
// ViStatus _VI_FUNC  viGpibControlREN(ViSession vi, ViUInt16 mode);

// viGpibControlATN
// ViStatus _VI_FUNC  viGpibControlATN(ViSession vi, ViUInt16 mode);

// viGpibSendIFC
// ViStatus _VI_FUNC  viGpibSendIFC   (ViSession vi);

// viGpibCommand
// ViStatus _VI_FUNC  viGpibCommand   (ViSession vi, ViBuf cmd, ViUInt32 cnt, ViPUInt32 retCnt);

// viGpibPassControl
// ViStatus _VI_FUNC  viGpibPassControl(ViSession vi, ViUInt16 primAddr, ViUInt16 secAddr);

// viVxiCommandQuery
// ViStatus _VI_FUNC  viVxiCommandQuery(ViSession vi, ViUInt16 mode, ViUInt32 cmd,
//                                      ViPUInt32 response);

// viAssertUtilSignal
// ViStatus _VI_FUNC  viAssertUtilSignal(ViSession vi, ViUInt16 line);

// viAssertIntrSignal
// ViStatus _VI_FUNC  viAssertIntrSignal(ViSession vi, ViInt16 mode, ViUInt32 statusID);

// viMapTrigger
// ViStatus _VI_FUNC  viMapTrigger    (ViSession vi, ViInt16 trigSrc, ViInt16 trigDest,
//                                     ViUInt16 mode);

// viUnmapTrigger
// ViStatus _VI_FUNC  viUnmapTrigger  (ViSession vi, ViInt16 trigSrc, ViInt16 trigDest);

// viUsbControlOut
// ViStatus _VI_FUNC  viUsbControlOut (ViSession vi, ViInt16 bmRequestType, ViInt16 bRequest,
//                                     ViUInt16 wValue, ViUInt16 wIndex, ViUInt16 wLength,
//                                     ViBuf buf);

// viUsbControlIn
// ViStatus _VI_FUNC  viUsbControlIn  (ViSession vi, ViInt16 bmRequestType, ViInt16 bRequest,
//                                     ViUInt16 wValue, ViUInt16 wIndex, ViUInt16 wLength,
//                                     ViPBuf buf, ViPUInt16 retCnt);

// viPxiReserveTriggers
// ViStatus _VI_FUNC  viPxiReserveTriggers(ViSession vi, ViInt16 cnt, ViAInt16 trigBuses,
//                                     ViAInt16 trigLines, ViPInt16 failureIndex);

// ViVersion Returns the unformatted resource version number.
func ViVersion() (vers uint32) {
	vers = uint32(C.VI_SPEC_VERSION)
	return vers
}

// ViVersMajor Returns the major resource version number.
func ViVersMajor() (versMaj uint32) {
	versMaj = (ViVersion() & 0xFFF00000) >> 20
	return versMaj
}

// ViVersMinor Returns the minor resource version number.
func ViVersMinor() (versMin uint32) {
	versMin = (ViVersion() & 0x000FFF00) >> 8
	return versMin
}

// ViVersSubMinor Returns the sub-minor resource version number.
func ViVersSubMinor() (versSubMin uint32) {
	versSubMin = (ViVersion() & 0x000000FF)
	return versSubMin
}

// viVxiServantResponse
// ViStatus _VI_FUNC viVxiServantResponse(ViSession vi, ViInt16 mode, ViUInt32 resp);
