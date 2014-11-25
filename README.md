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
-

for compile checks I created a stubbed 64-bit version of libvisa.so
