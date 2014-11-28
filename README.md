goVISA
======

A Go wrapper around National Instruments VISA (Virtual Instrument Software Architecture) driver.

http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go

sudo mount -o loop NI-VISA-14.0.0.iso /mnt/disk

to cross compile
----------------
- use dave cheney's crosscompile to to build the 386 compiler
- sudo su
- export CGO_ENABLED=1
- export GOARCH=386
- go build visa,go defs.go

for compile checks I created a stubbed 64-bit version of libvisa.so

https://code.google.com/p/go-wiki/wiki/cgo

Notes to me:
- saved /usr/local/lib/libvisa.so to /usr/local/lib/libvisa.so.orig
- moved my lib stubs to /usr/local/lib/
