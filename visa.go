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
//
package visa

// TODO:
// -

/*
#cgo linux LDFLAGS: -lgpibapi
#cgo darwin CFLAGS: -I.
#cgo darwin LDFLAGS: -framework VISA
#cgo windows CFLAGS: -I.
#cgo windows LDFLAGS: -lgpib-32 -LC:/WINDOWS/system32
#include <stdlib.h>
#if defined(__amd64) || defined(__amd64__) || defined(__x86_64) || defined(__x86_64__) && !defined(__APPLE__)
#define size_g size_t
#include <ni4882.h>
#else
#define size_g long
#include <ni488.h>
#endif
*/
import "C"
import "unsafe"

var PackageVersion string = "v0.1"

/*- VISA Types --------------------------------------------------------------*/

typedef ViObject             ViEvent;
typedef ViEvent      _VI_PTR ViPEvent;
typedef ViObject             ViFindList;
typedef ViFindList   _VI_PTR ViPFindList;

#if defined(_VI_INT64_UINT64_DEFINED) && defined(_VISA_ENV_IS_64_BIT)
typedef ViUInt64             ViBusAddress;
typedef ViUInt64             ViBusSize;
typedef ViUInt64             ViAttrState;
#else
typedef ViUInt32             ViBusAddress;
typedef ViUInt32             ViBusSize;
typedef ViUInt32             ViAttrState;
#endif

#if defined(_VI_INT64_UINT64_DEFINED)
typedef ViUInt64             ViBusAddress64;
typedef ViBusAddress64 _VI_PTR ViPBusAddress64;
#endif

typedef ViUInt32             ViEventType;
typedef ViEventType  _VI_PTR ViPEventType;
typedef ViEventType  _VI_PTR ViAEventType;
typedef void         _VI_PTR ViPAttrState;
typedef ViAttr       _VI_PTR ViPAttr;
typedef ViAttr       _VI_PTR ViAAttr;

typedef ViString             ViKeyId;
typedef ViPString            ViPKeyId;
typedef ViUInt32             ViJobId;
typedef ViJobId      _VI_PTR ViPJobId;
typedef ViUInt32             ViAccessMode;
typedef ViAccessMode _VI_PTR ViPAccessMode;
typedef ViBusAddress _VI_PTR ViPBusAddress;
typedef ViUInt32             ViEventFilter;

typedef va_list              ViVAList;

typedef ViStatus (_VI_FUNCH _VI_PTR ViHndlr)
   (ViSession vi, ViEventType eventType, ViEvent event, ViAddr userHandle);

/*- Resource Manager Functions and Operations -------------------------------*/

// ViStatus _VI_FUNC viOpenDefaultRM(ViPSession vi);
func ViOpenDefaultRM(vi ViPSession) (status ViStatus) {
	status = C.viOpenDefaultRM(vi)
}

// ViStatus _VI_FUNC  viFindRsrc(ViSession sesn, ViString expr, ViPFindList vi,
                                    // ViPUInt32 retCnt, ViChar _VI_FAR desc[]);
func ViFindRsrc(sesn ViSession, expr ViString, vi ViPFindList,
	retCnt ViPUInt32, desc []ViChar) (status ViStatus) {
	status = C.viFindRsrc(sesn, expr, vi, retCnt, desc)
}

// ViStatus _VI_FUNC  viFindNext      (ViFindList vi, ViChar _VI_FAR desc[]);
func ViFindNext(vi ViFindList, desc []ViChar) (status ViStatus) {
	status = C.viFindNext(vi, desc)
}

// ViStatus _VI_FUNC  viParseRsrc     (ViSession rmSesn, ViRsrc rsrcName,
                                    // ViPUInt16 intfType, ViPUInt16 intfNum);
func ViParseRsrc(rmSesn ViSession, rsrcName ViRsrc, intfType ViPUInt16,
	intfNum ViPUInt16) (status ViStatus) {
	status = C.viParseRsrc(rmSesn, rsrcName, intfType, intfNum)
}

// ViStatus _VI_FUNC  viParseRsrcEx   (ViSession rmSesn, ViRsrc rsrcName, ViPUInt16 intfType,
                                    // ViPUInt16 intfNum, ViChar _VI_FAR rsrcClass[],
                                    // ViChar _VI_FAR expandedUnaliasedName[],
                                    // ViChar _VI_FAR aliasIfExists[]);
func ViParseRsrcEx(rmSesn ViSession, rsrcName ViRsrc, intfType ViPUInt16,
	intfNum ViPUInt16, rsrcClass []ViChar, expandedUnaliasedName []ViChar,
	aliasIfExists []ViChar) (status ViStatus) {
	status = C.viParseRsrcEx(rmSesn, rsrcName, intfType, intfNum, rsrcClass,
		expandedUnaliasedName, aliasIfExists)
}

// ViStatus _VI_FUNC  viOpen          (ViSession sesn, ViRsrc name, ViAccessMode mode,
                                    // ViUInt32 timeout, ViPSession vi);
func ViOpen(timeout ViUInt32, vi ViPSession) (status ViStatus) {
	status = C.viOpen(timeout, vi)
}

/*- Resource Template Operations --------------------------------------------*/

// ViStatus _VI_FUNC  viClose         (ViObject vi);
// ViStatus _VI_FUNC  viSetAttribute  (ViObject vi, ViAttr attrName, ViAttrState attrValue);
// ViStatus _VI_FUNC  viGetAttribute  (ViObject vi, ViAttr attrName, void _VI_PTR attrValue);
// ViStatus _VI_FUNC  viStatusDesc    (ViObject vi, ViStatus status, ViChar _VI_FAR desc[]);
// ViStatus _VI_FUNC  viTerminate     (ViObject vi, ViUInt16 degree, ViJobId jobId);

// ViStatus _VI_FUNC  viLock          (ViSession vi, ViAccessMode lockType, ViUInt32 timeout,
//                                     ViKeyId requestedKey, ViChar _VI_FAR accessKey[]);
// ViStatus _VI_FUNC  viUnlock        (ViSession vi);
// ViStatus _VI_FUNC  viEnableEvent   (ViSession vi, ViEventType eventType, ViUInt16 mechanism,
//                                     ViEventFilter context);
// ViStatus _VI_FUNC  viDisableEvent  (ViSession vi, ViEventType eventType, ViUInt16 mechanism);
// ViStatus _VI_FUNC  viDiscardEvents (ViSession vi, ViEventType eventType, ViUInt16 mechanism);
// ViStatus _VI_FUNC  viWaitOnEvent   (ViSession vi, ViEventType inEventType, ViUInt32 timeout,
//                                     ViPEventType outEventType, ViPEvent outContext);
// ViStatus _VI_FUNC  viInstallHandler(ViSession vi, ViEventType eventType, ViHndlr handler,
//                                     ViAddr userHandle);
// ViStatus _VI_FUNC  viUninstallHandler(ViSession vi, ViEventType eventType, ViHndlr handler,
//                                       ViAddr userHandle);

/*- Basic I/O Operations ----------------------------------------------------*/

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

/*- Formatted and Buffered I/O Operations -----------------------------------*/

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

/*- Memory I/O Operations ---------------------------------------------------*/

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

/*- Shared Memory Operations ------------------------------------------------*/

// ViStatus _VI_FUNC  viMemAlloc      (ViSession vi, ViBusSize size, ViPBusAddress offset);
// ViStatus _VI_FUNC  viMemFree       (ViSession vi, ViBusAddress offset);

// #if defined(_VI_INT64_UINT64_DEFINED)
// ViStatus _VI_FUNC  viMemAllocEx    (ViSession vi, ViBusSize size, ViPBusAddress64 offset);
// ViStatus _VI_FUNC  viMemFreeEx     (ViSession vi, ViBusAddress64 offset);
// #endif

/*- Interface Specific Operations -------------------------------------------*/

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
