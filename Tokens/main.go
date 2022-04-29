package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {

	advapi32 := windows.NewLazySystemDLL("advapi32")
	kernel32 := windows.NewLazySystemDLL("kernel32")

	GetCurrentProcess := kernel32.NewProc("GetCurrentProcess")
	GetLastError := kernel32.NewProc("GetLastError")
	OpenProdcessToken := advapi32.NewProc("OpenProcessToken")
	LookupPrivilegeValue := advapi32.NewProc("LookupPrivilegeValueW")
	AdjustTokenPrivileges := advapi32.NewProc("AdjustTokenPrivileges")

	currentProcess, _, _ := GetCurrentProcess.Call()

	const tokenAdjustPrivileges = 0x0020
	const tokenQuery = 0x0008
	const SePrivilegeEnabled uint32 = 0x00000002

	var hToken uintptr

	result, _, err := OpenProdcessToken.Call(currentProcess, tokenAdjustPrivileges|tokenQuery, uintptr(unsafe.Pointer(&hToken)))

	if result != 1 {
		fmt.Println("OpenProcessToken(): ", result, " err: ", err)
	}

	const SeDebugPrivilege = "SeDebugPrivilege"

	fmt.Println("hToken: ", hToken)

	type Luid struct {
		lowPart  uint32 // DWORD
		highPart int32  // long
	}

	type LuidAndAttributes struct {
		luid       Luid   // LUID
		attributes uint32 // DWORD
	}

	type TokenPrivileges struct {
		privilegeCount uint32 // DWORD
		privileges     [1]LuidAndAttributes
	}

	var tkp TokenPrivileges

	result, _, err = LookupPrivilegeValue.Call(uintptr(0),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(SeDebugPrivilege))), uintptr(unsafe.Pointer(&(tkp.privileges[0].luid))))
	if result != 1 {
		fmt.Println("LookupPrivilegeValue(): ", result, " err: ", err)
	}
	fmt.Println("LookupPrivilegeValue luid: ", tkp.privileges[0].luid)

	tkp.privilegeCount = 1
	tkp.privileges[0].attributes = SePrivilegeEnabled

	result, _, err = AdjustTokenPrivileges.Call(hToken, 0, uintptr(unsafe.Pointer(&tkp)), 0, uintptr(0), 0)
	if result != 1 {
		fmt.Println("AdjustTokenPrivileges() ", result, " err: ", err)
	}

	result, _, _ = GetLastError.Call()
	if result != 0 {
		fmt.Println("GetLastError() ", result)
	}
}
