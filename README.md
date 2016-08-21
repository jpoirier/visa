Description
===========

Package visa wraps National Instruments VISA (Virtual Instrument Software Architecture) driver. The driver
allows a client application to communicate with most instrumentation buses including GPIB, USB, Serial, and Ethernet.
The Virtual Instrument Software Architecture (VISA) is a standard for configuring, programming, and troubleshooting
instrumentation systems comprising GPIB, VXI, PXI, serial (RS232/RS485), Ethernet/LXI, and/or USB interfaces.

The package is low level and, for the most part, is one-to-one with the
exported C functions it wraps. Clients would typically build an instrument
specific driver around the package but it can also be used directly.

[NI-VISA Overview] (http://www.ni.com/white-paper/3702/en/)

Supported Platforms:
* Linux
* OS X
* Windows


Installation
============

Dependencies
------------

* [Go tools](https://golang.org)
* [NI-VISA] (http://www.ni.com/downloads/ni-drivers/)
* [git] (https://git-scm.com)

Usage
-----

All functions in libvisa are accessible from the gortlsdr package:

    go get -u github.com/jpoirier/visa
    go get -u github.com/jpoirier/visa/mxa
    go get -u github.com/jpoirier/visa/keithley

Example
-------

See the examples folder:

    go run FindRsrc.go

Windows
=======

Building visa on Windows
------------------------

Additional dependencies to build the visa wrapper

* [mingw-w64] (http://sourceforge.net/projects/mingw-w64/?source=recommended)


A example workaround for building on Windows (e.g. 64 bit Windows 7) when having NI Visa include and/or library path problems: 

    - install Go
    - install the latest NI Visa tools, e.g. version 16
    - install a 64 bit GCC via http://tdm-gcc.tdragon.net/download, add C:\TDM-GCC-64\bin to PATH
    - go get -d github.com/jpoirier/visa
    - copy visa.h, visatype.h, vpptype.h from C:/Program Files/National Instruments/shared/CompilerSupport/c/Include to %GOPATH%/src/github.com/jpoirier/visa
    - make a directory in your gopath for the visa dll, e.g. if GOPATH=C:/users/jsmith/gopath then, mkdir %GOPATH%/windows_libs
    - copy visa64.dll (or visa32.dll if on a 32 bit system) from C:/Windows/System32 to C:/users/jsmith/gopath/windows_libs
    - edit visa.go, change: #cgo windows LDFLAGS: -lvisa  to: #cgo windows LDFLAGS: -lvisaXX -LC:/users/jsmith/gopath/windows_libs where visaXX is visa32 or visa64
    - open a Window's shell in %GOPATH%/src/github.com/jpoirier/visa
    - run: go build -o %GOPATH%/pkg/windows_amd64/github.com/jpoirier/visa.a visa.go exports.go defs.go
    - from a shell window cd to visa/examples and: go run FindRsrc.go
