// Copyright (c) 2014 Joseph D Poirier
// Distributable under the terms of The simplified BSD License
// that can be found in the LICENSE file.

// Package visa wraps National Instruments
//  (VISA) driver. The driver allows a client application to communicate
// with a VISA enabled piece of test equipment remotely and/or programmatically.
// VISA is an industry standard for GPIB communications.
//
// The package is low level and, for the most part, is one-to-one with the
// exported C functions it wraps. Clients would typically build instrument
// drivers around the package but it can also be used directly.
//
// Lots of miscellaneous NI-488.2 information:
//     http://sine.ni.com/psp/app/doc/p/id/psp-356
//
// GPIB Driver Versions for Microsoft Windows and DOS:
//     http://zone.ni.com/devzone/cda/tut/p/id/5326#toc0
//
// GPIB Driver Versions for non-Microsoft Operating Systems:
//     http://zone.ni.com/devzone/cda/tut/p/id/5458
//
// Direct download: http://download.ni.com/support/softlib/gpib/

// export CGO_ENABLED=1
// export GOARCH=386

package visa

import "unsafe"

/*
#cgo linux CFLAGS: -I.
#cgo linux LDFLAGS: -L. -lvisa
#cgo darwin CFLAGS: -I.
#cgo darwin LDFLAGS: -framework VISA
#cgo windows CFLAGS: -I.
#cgo windows LDFLAGS: -lvisa -LC:/WINDOWS/system32
#include <stdlib.h>
#include <visa.h>
*/
import "C"

var PackageVersion string = "v0.1"

// VISA Types

type ViEvent uint32

// type ViPEvent C.ViPEvent

type ViFindList uint32

// type ViPFindList C.ViPFindList

type ViBusAddress uint32

// type ViBusAddress uint64

type ViBusSize uint32

// type ViBusSize uint64

type ViAttrState uint32

// type ViAttrState uint64

// #if defined(_VI_INT64_UINT64_DEFINED)
// typedef ViUInt64             ViBusAddress64;
// typedef ViBusAddress64 _VI_PTR ViPBusAddress64;
// #endif

type ViEventType uint32

// type ViPEventType C.ViPEventType

type ViAEventType C.ViAEventType
type ViPAttrState C.ViPAttrState
type ViPAttr C.ViPAttr
type ViAAttr C.ViAAttr
type ViKeyId C.ViKeyId
type ViPKeyId C.ViPKeyId
type ViJobId C.ViJobId
type ViPJobId C.ViPJobId

type ViAccessMode uint32
type ViPAccessMode C.ViPAccessMode
type ViPBusAddress C.ViPBusAddress
type ViEventFilter uint32

type ViVAList C.ViVAList
type ViStatus C.ViStatus

// type ViPSession C.ViPSession
type ViSession uint32
type ViString string

// type ViPUInt32 C.ViPUInt32
type ViChar int8
type ViRsrc string
type ViPUInt16 C.ViPUInt16
type ViUInt32 uint32

// Resource Manager Functions and Operations

// ViOpenDefaultRM returns a session to the Default Resource Manager resource.
func ViOpenDefaultRM() (vi uint32, status ViStatus) {
	status = ViStatus(C.viOpenDefaultRM((C.ViPSession)(unsafe.Pointer(&vi))))
	return vi, status
}

// ViFindRsrc queries a VISA system to locate the resources associated with a
// specified interface.
func ViFindRsrc(sesn uint32, expr string) (vi uint32, retCnt uint32,
	desc string, status ViStatus) {

	cexpr := (*C.ViChar)(C.CString(expr))
	defer C.free(unsafe.Pointer(cexpr))
	d := make([]byte, 257)
	status = ViStatus(C.viFindRsrc(C.ViSession(sesn),
		cexpr,
		(*C.ViFindList)(unsafe.Pointer(&vi)),
		(*C.ViUInt32)(unsafe.Pointer(&retCnt)),
		(*C.ViChar)(unsafe.Pointer(&d[0]))))

	return vi, retCnt, string(d), status
}

// ViFindNext gets the next resource from the list of resources found during a
// previous call to viFindRsrc.
func ViFindNext(vi uint32) (desc string, status ViStatus) {
	cdesc := (*C.ViChar)(C.CString(desc))
	defer C.free(unsafe.Pointer(cdesc))
	status = ViStatus(C.viFindNext((C.ViFindList)(vi), cdesc))

	return desc, status
}

// ViParseRsrc parses a resource string to get the interface information.
func ViParseRsrc(rmSesn uint32, rsrcName string) (intfType uint16,
	intfNum uint16, status ViStatus) {

	crsrcName := (*C.ViChar)(C.CString(rsrcName))
	defer C.free(unsafe.Pointer(crsrcName))
	status = ViStatus(C.viParseRsrc(C.ViSession(rmSesn),
		crsrcName,
		(*C.ViUInt16)(unsafe.Pointer(&intfType)),
		(*C.ViUInt16)(unsafe.Pointer(&intfNum))))

	return intfType, intfNum, status
}

// ViParseRsrcEx parses a resource string to get extended interface information.
func ViParseRsrcEx(rmSesn uint32, rsrcName string) (intfType uint16,
	intfNum uint16, rsrcClass string, expandedUnaliasedName string,
	aliasIfExists string, status ViStatus) {

	crsrcName := (*C.ViChar)(C.CString(rsrcName))
	defer C.free(unsafe.Pointer(crsrcName))

	crsrcClass := (*C.ViChar)(C.CString(rsrcClass))
	defer C.free(unsafe.Pointer(crsrcClass))

	cexpandedUnaliasedName := (*C.ViChar)(C.CString(expandedUnaliasedName))
	defer C.free(unsafe.Pointer(cexpandedUnaliasedName))

	caliasIfExists := (*C.ViChar)(C.CString(aliasIfExists))
	defer C.free(unsafe.Pointer(caliasIfExists))

	status = ViStatus(C.viParseRsrcEx(C.ViSession(rmSesn),
		crsrcName,
		(*C.ViUInt16)(unsafe.Pointer(&intfType)),
		(*C.ViUInt16)(unsafe.Pointer(&intfNum)),
		crsrcClass,
		cexpandedUnaliasedName,
		caliasIfExists))

	return intfType, intfNum, rsrcClass, expandedUnaliasedName, aliasIfExists, status
}

// ViOpen opens a session to the specified resource.
func ViOpen(sesn uint32, name string, mode, timeout uint32) (vi uint32,
	status ViStatus) {

	cname := (*C.ViChar)(C.CString(name))
	defer C.free(unsafe.Pointer(cname))
	status = ViStatus(C.viOpen(C.ViSession(sesn),
		cname,
		(C.ViAccessMode)(mode),
		(C.ViUInt32)(timeout),
		(*C.ViSession)(unsafe.Pointer(&vi))))

	return vi, status
}

// Resource Template Operations

// ViClose Closes the specified session, event, or find list.
func ViClose(vi uint32) (status ViStatus) {

	status = ViStatus(C.viClose((C.ViObject)(vi)))
	return status
}

// ViSetAttribute Sets the state of an attribute.
func ViSetAttribute(vi, attribute, attrState uint32) (status ViStatus) {
	status = ViStatus(C.viSetAttribute((C.ViObject)(vi),
		(C.ViAttr)(attribute),
		(C.ViAttrState)(attrState)))

	return status
}

// ViGetAttribute Retrieves the state of an attribute.
// vi in
// attrName in
// attrValue out
// ViStatus _VI_FUNC  viGetAttribute  (ViObject vi, ViAttr attrName, void _VI_PTR attrValue);
// func ViGetAttribute(vi, attribute uint32) {
// 	cvi := (C.ViObject)(vi)
// 	cattribute := (C.ViAttr)(attribute)

// 	return
// }

// ViStatusDesc Returns a user-readable description of the status code passed to the operation.
// vi in
// status in
// desc out
// ViStatus _VI_FUNC  viStatusDesc    (ViObject vi, ViStatus status, ViChar _VI_FAR desc[]);
// func ViStatusDesc(vi uint32) {
// 	cvi := (C.ViObject)(vi)
// 	return
// }

// ViTerminate Requests a VISA session to terminate normal execution of an operation.
func ViTerminate(vi uint32, degree, jobId uint16) (status ViStatus) {
	status = ViStatus(C.viTerminate((C.ViObject)(vi),
		(C.ViUInt16)(degree),
		(C.ViJobId)(jobId)))

	return status
}

// ViLock Establishes an access mode to the specified resource.
// vi in
// lockType in
// timeout in
// requestedKey in
// accessKey out
// ViStatus _VI_FUNC  viLock(ViSession vi, ViAccessMode lockType, ViUInt32 timeout,
//                                     ViKeyId requestedKey, ViChar _VI_FAR accessKey[]);
func ViLock(vi, lockType, timeout uint32, requestedKey string) (accessKey string,
	status ViStatus) {

	crequestedKey := (*C.ViChar)(C.CString(requestedKey))
	defer C.free(unsafe.Pointer(crequestedKey))

	a := make([]byte, 257)
	status = ViStatus(C.viLock((C.ViSession)(vi),
		(C.ViAccessMode)(lockType),
		(C.ViUInt32)(timeout),
		crequestedKey,
		(*C.ViChar)(unsafe.Pointer(&a[0]))))

	return string(a), status
}

// ViUnlock Relinquishes a lock for the specified resource.
func ViUnlock(vi uint32) (status ViStatus) {
	status = ViStatus(C.viUnlock((C.ViSession)(vi)))

	return status
}

// ViEnableEvent Enables notification of a specified event.
func ViEnableEvent(vi, eventType uint32, mechanism uint16, context uint32) (status ViStatus) {
	status = ViStatus(C.viEnableEvent((C.ViSession)(vi),
		(C.ViEventType)(eventType),
		(C.ViUInt16)(mechanism),
		(C.ViEventFilter)(context)))

	return status
}

// ViDisableEvent Disables notification of the specified event type(s) via the specified mechanism(s).
// vi in
// inEventType in
// mechanism in
// ViStatus _VI_FUNC  viDisableEvent  (ViSession vi, ViEventType eventType, ViUInt16 mechanism);
func ViDisableEvent(vi uint32) {
	cvi := (C.ViSession)(vi)
	return
}

// ViDiscardEvents Discards event occurrences for specified event types and mechanisms in a session.
// vi in
// inEventType in
// mechanism in
// ViStatus _VI_FUNC  viDiscardEvents (ViSession vi, ViEventType eventType, ViUInt16 mechanism);
func ViDiscardEvents(vi uint32) {
	cvi := (C.ViSession)(vi)
	return
}

// ViWaitOnEvent Waits for an occurrence of the specified event for a given session.
// vin in
// inEventType in
// timeout in
// outEventType out
// outContext out
// ViStatus _VI_FUNC  viWaitOnEvent   (ViSession vi, ViEventType inEventType, ViUInt32 timeout,
//                                     ViPEventType outEventType, ViPEvent outContext);
func ViWaitOnEvent(vi uint32) {
	cvi := (C.ViSession)(vi)
	return
}

// ViInstallHandler nstalls handlers for event callbacks.
// vi in
// eventType in
// handler in
// userHandle in
// ViStatus _VI_FUNC  viInstallHandler(ViSession vi, ViEventType eventType, ViHndlr handler,
//                                     ViAddr userHandle);
func ViInstallHandler(vi uint32) {
	cvi := (C.ViSession)(vi)
	return
}

// ViUninstallHandler Uninstalls handlers for events.
// vi in
// eventType in
// handler in
// userHandle in
// ViStatus _VI_FUNC  viUninstallHandler(ViSession vi, ViEventType eventType, ViHndlr handler,
//                                       ViAddr userHandle);
func ViUninstallHandler(vi uint32) {
	cvi := (C.ViSession)(vi)
	return
}

// Basic I/O Operations

// ViStatus _VI_FUNC  viRead          (ViSession vi, ViPBuf buf, ViUInt32 cnt, ViPUInt32 retCnt);
// ViStatus _VI_FUNC  viReadAsync     (ViSession vi, ViPBuf buf, ViUInt32 cnt, ViPJobId  jobId);
// ViStatus _VI_FUNC  viReadToFile    (ViSession vi, ViConstString filename, ViUInt32 cnt,
//                                     ViPUInt32 retCnt);
// ViStatus _VI_FUNC  viWrite         (ViSession vi, ViBuf  buf, ViUInt32 cnt, ViPUInt32 retCnt);
// ViStatus _VI_FUNC  viWriteAsync    (ViSession vi, ViBuf  buf, ViUInt32 cnt, ViPJobId  jobId);
// ViStatus _VI_FUNC  viWriteFromFile (ViSession vi, ViConstString filename, ViUInt32 cnt,
//                                     ViPUInt32 retCnt);
// ViStatus _VI_FUNC  viAssertTrigger (ViSession vi, ViUInt16 protocol);
// ViStatus _VI_FUNC  viReadSTB       (ViSession vi, ViPUInt16 status);
// ViStatus _VI_FUNC  viClear         (ViSession vi);

// Formatted and Buffered I/O Operations

// ViStatus _VI_FUNC  viSetBuf        (ViSession vi, ViUInt16 mask, ViUInt32 size);
// ViStatus _VI_FUNC  viFlush         (ViSession vi, ViUInt16 mask);

// ViStatus _VI_FUNC  viBufWrite      (ViSession vi, ViBuf  buf, ViUInt32 cnt, ViPUInt32 retCnt);
// ViStatus _VI_FUNC  viBufRead       (ViSession vi, ViPBuf buf, ViUInt32 cnt, ViPUInt32 retCnt);

// ViStatus _VI_FUNCC viPrintf        (ViSession vi, ViString writeFmt, ...);
// ViStatus _VI_FUNC  viVPrintf       (ViSession vi, ViString writeFmt, ViVAList params);
// ViStatus _VI_FUNCC viSPrintf       (ViSession vi, ViPBuf buf, ViString writeFmt, ...);
// ViStatus _VI_FUNC  viVSPrintf      (ViSession vi, ViPBuf buf, ViString writeFmt,
//                                     ViVAList parms);

// ViStatus _VI_FUNCC viScanf         (ViSession vi, ViString readFmt, ...);
// ViStatus _VI_FUNC  viVScanf        (ViSession vi, ViString readFmt, ViVAList params);
// ViStatus _VI_FUNCC viSScanf        (ViSession vi, ViBuf buf, ViString readFmt, ...);
// ViStatus _VI_FUNC  viVSScanf       (ViSession vi, ViBuf buf, ViString readFmt,
//                                     ViVAList parms);

// ViStatus _VI_FUNCC viQueryf        (ViSession vi, ViString writeFmt, ViString readFmt, ...);
// ViStatus _VI_FUNC  viVQueryf       (ViSession vi, ViString writeFmt, ViString readFmt,
//                                     ViVAList params);

// Memory I/O Operations

// ViStatus _VI_FUNC  viIn8           (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViPUInt8  val8);
// ViStatus _VI_FUNC  viOut8          (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViUInt8   val8);
// ViStatus _VI_FUNC  viIn16          (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViPUInt16 val16);
// ViStatus _VI_FUNC  viOut16         (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViUInt16  val16);
// ViStatus _VI_FUNC  viIn32          (ViSession vi, ViUInt16 space,
//                                     ViBusAddress offset, ViPUInt32 val32);
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

// ViStatus _VI_FUNC  viMoveIn8       (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt8  buf8);
// ViStatus _VI_FUNC  viMoveOut8      (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt8  buf8);
// ViStatus _VI_FUNC  viMoveIn16      (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt16 buf16);
// ViStatus _VI_FUNC  viMoveOut16     (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt16 buf16);
// ViStatus _VI_FUNC  viMoveIn32      (ViSession vi, ViUInt16 space, ViBusAddress offset,
//                                     ViBusSize length, ViAUInt32 buf32);
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

// ViStatus _VI_FUNC  viMove          (ViSession vi, ViUInt16 srcSpace, ViBusAddress srcOffset,
//                                     ViUInt16 srcWidth, ViUInt16 destSpace,
//                                     ViBusAddress destOffset, ViUInt16 destWidth,
//                                     ViBusSize srcLength);
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

// ViStatus _VI_FUNC  viMapAddress    (ViSession vi, ViUInt16 mapSpace, ViBusAddress mapOffset,
//                                     ViBusSize mapSize, ViBoolean access,
//                                     ViAddr suggested, ViPAddr address);
// ViStatus _VI_FUNC  viUnmapAddress  (ViSession vi);

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViStatus _VI_FUNC  viMapAddressEx  (ViSession vi, ViUInt16 mapSpace, ViBusAddress64 mapOffset,
//                                     ViBusSize mapSize, ViBoolean access,
//                                     ViAddr suggested, ViPAddr address);
// #endif

// void     _VI_FUNC  viPeek8         (ViSession vi, ViAddr address, ViPUInt8  val8);
// void     _VI_FUNC  viPoke8         (ViSession vi, ViAddr address, ViUInt8   val8);
// void     _VI_FUNC  viPeek16        (ViSession vi, ViAddr address, ViPUInt16 val16);
// void     _VI_FUNC  viPoke16        (ViSession vi, ViAddr address, ViUInt16  val16);
// void     _VI_FUNC  viPeek32        (ViSession vi, ViAddr address, ViPUInt32 val32);
// void     _VI_FUNC  viPoke32        (ViSession vi, ViAddr address, ViUInt32  val32);

// #if defined(_VI_INT64_UINT64_DEFINED)
// void     _VI_FUNC  viPeek64        (ViSession vi, ViAddr address, ViPUInt64 val64);
// void     _VI_FUNC  viPoke64        (ViSession vi, ViAddr address, ViUInt64  val64);
// #endif

// Shared Memory Operations

// ViStatus _VI_FUNC  viMemAlloc      (ViSession vi, ViBusSize size, ViPBusAddress offset);
// ViStatus _VI_FUNC  viMemFree       (ViSession vi, ViBusAddress offset);

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViStatus _VI_FUNC  viMemAllocEx    (ViSession vi, ViBusSize size, ViPBusAddress64 offset);
// ViStatus _VI_FUNC  viMemFreeEx     (ViSession vi, ViBusAddress64 offset);
// #endif

// Interface Specific Operations

// ViStatus _VI_FUNC  viGpibControlREN(ViSession vi, ViUInt16 mode);
// ViStatus _VI_FUNC  viGpibControlATN(ViSession vi, ViUInt16 mode);
// ViStatus _VI_FUNC  viGpibSendIFC   (ViSession vi);
// ViStatus _VI_FUNC  viGpibCommand   (ViSession vi, ViBuf cmd, ViUInt32 cnt, ViPUInt32 retCnt);
// ViStatus _VI_FUNC  viGpibPassControl(ViSession vi, ViUInt16 primAddr, ViUInt16 secAddr);

// ViStatus _VI_FUNC  viVxiCommandQuery(ViSession vi, ViUInt16 mode, ViUInt32 cmd,
//                                      ViPUInt32 response);
// ViStatus _VI_FUNC  viAssertUtilSignal(ViSession vi, ViUInt16 line);
// ViStatus _VI_FUNC  viAssertIntrSignal(ViSession vi, ViInt16 mode, ViUInt32 statusID);
// ViStatus _VI_FUNC  viMapTrigger    (ViSession vi, ViInt16 trigSrc, ViInt16 trigDest,
//                                     ViUInt16 mode);
// ViStatus _VI_FUNC  viUnmapTrigger  (ViSession vi, ViInt16 trigSrc, ViInt16 trigDest);
// ViStatus _VI_FUNC  viUsbControlOut (ViSession vi, ViInt16 bmRequestType, ViInt16 bRequest,
//                                     ViUInt16 wValue, ViUInt16 wIndex, ViUInt16 wLength,
//                                     ViBuf buf);
// ViStatus _VI_FUNC  viUsbControlIn  (ViSession vi, ViInt16 bmRequestType, ViInt16 bRequest,
//                                     ViUInt16 wValue, ViUInt16 wIndex, ViUInt16 wLength,
//                                     ViPBuf buf, ViPUInt16 retCnt);
// ViStatus _VI_FUNC  viPxiReserveTriggers(ViSession vi, ViInt16 cnt, ViAInt16 trigBuses,
//                                     ViAInt16 trigLines, ViPInt16 failureIndex);

// #define VI_VERSION_MAJOR(ver)       ((((ViVersion)ver) & 0xFFF00000UL) >> 20)
// #define VI_VERSION_MINOR(ver)       ((((ViVersion)ver) & 0x000FFF00UL) >>  8)
// #define VI_VERSION_SUBMINOR(ver)    ((((ViVersion)ver) & 0x000000FFUL)      )
// #define viGetDefaultRM(vi)
// #if defined(_CVI_DEBUG_)
// #pragma soft_reference (viGetAttribute);
// #endif
// ViStatus _VI_FUNC viVxiServantResponse(ViSession vi, ViInt16 mode, ViUInt32 resp);
