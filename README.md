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

Building the visa package on Windows
------------------------------------

Additional dependencies to build the visa wrapper

* [mingw-w64] (http://sourceforge.net/projects/mingw-w64/?source=recommended)


An example Windows build: 

    - install Go
    - install the latest NI Visa drivers
    - install a 64 bit GCC via mingw-w64 and put mingw-w64's bin folder in your PATH
    - open a Windows shell and run> go get -d github.com/jpoirier/visa
    - copy C:\Program Files\IVI Foundation\VISA\Win64\Lib_x64\msc\visa64.lib to the same folder as the visa package source files
    - open a Windows shell in the visa package folder
    - run> go build -o %GOPATH%/pkg/windows_amd64/github.com/jpoirier/visa.a -ldflags "-linkmode external -extldflags -static" visa.go exports.go defs.go
    - cd to examples and run> go run FindRsrc.go
