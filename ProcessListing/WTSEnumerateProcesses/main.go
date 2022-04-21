// +build windows
package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type SID struct {
	Revision            byte
	SubAuthorityCount   byte
	IdentifierAuthority SID_IDENTIFIER_AUTHORITY
}

type SID_IDENTIFIER_AUTHORITY struct {
	Value [6]byte
}

type WTS_PROCESS_INFO struct {
	SessionId    uint32
	ProcessId    uint32
	PProcessName *uint16
	PUserSid     *SID
}

func UTF16toString(p *uint16) string {
	//return syscall.UTF16ToString((*[4096]uint16)(unsafe.Pointer(p))[:])

	ptr := unsafe.Pointer(p)                   //necessary to arbitrarily cast to *[4096]uint16 (?)
	uint16ptrarr := (*[4096]uint16)(ptr)[:]    //4096 is arbitrary? could be smaller
	return syscall.UTF16ToString(uint16ptrarr) //now uint16ptrarr is in a format to pass to the builtin converter
}

func main() {
	wtsapi32DLL := windows.NewLazySystemDLL("Wtsapi32.dll")
	WTSEnumerateProcesses := wtsapi32DLL.NewProc("WTSEnumerateProcessesW")

	var processListing *WTS_PROCESS_INFO
	var processCount uint32 = 0

	defer windows.WTSFreeMemory(uintptr(unsafe.Pointer(&processListing)))

	_, _, err := WTSEnumerateProcesses.Call(uintptr(0), 0, 1, uintptr(unsafe.Pointer(&processListing)),
		uintptr(unsafe.Pointer(&processCount)))

	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Printf("ops %v\n", err.Error())
		return
	}

	size := unsafe.Sizeof(WTS_PROCESS_INFO{})
	for i := uint32(0); i < processCount; i++ {
		p := *(*WTS_PROCESS_INFO)(unsafe.Pointer(uintptr(unsafe.Pointer(processListing)) + uintptr(size)*uintptr(i)))
		fmt.Printf("%d - %v\n", p.ProcessId, UTF16toString(p.PProcessName))
	}
}
