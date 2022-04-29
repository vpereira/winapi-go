package main

import (
	"log"
	"os"

	"golang.org/x/sys/windows"
)

func main() {

	var showCMD int32 = 1 // SW_NORMAL

	program, _ := windows.UTF16PtrFromString("notepad.exe")
	programParam, _ := windows.UTF16PtrFromString("")
	cwdParam, _ := os.Getwd()
	runas, _ := windows.UTF16PtrFromString("runas")
	err := windows.ShellExecute(0, runas, program, programParam, windows.StringToUTF16Ptr(cwdParam), showCMD)

	if err != nil {
		log.Fatal(err.Error())
	}
}
