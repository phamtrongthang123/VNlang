package builtin

import (
	"log"
	"syscall"
	"unsafe"
)

func str_to_uintptr(s string) uintptr {
	b := append([]byte(s), 0)
	return uintptr(unsafe.Pointer(&b[0]))
}

var MAX_STRING int = 200

func uintptr_to_str_res(resptr uintptr) string {
	myres := ""
	i := 0
	for true {
		ptr := unsafe.Pointer(resptr)
		stringPtr := (*byte)(ptr)
		stringRes := *stringPtr
		myres += string(stringRes)
		if stringRes == 0 {
			break
		}
		resptr++
		i++
		if i > MAX_STRING {
			break
		}
	}
	return myres
}

func syscall2_str_helper(h syscall.Handle, func_name string, arg1 string, arg2 string) string {
	proc, e := syscall.GetProcAddress(h, func_name) //One of the functions in the DLL
	if e != nil {
		log.Fatal(e)
	}
	resptr, _, _ := syscall.Syscall9(uintptr(proc), 2, str_to_uintptr(arg1), str_to_uintptr(arg2), 0, 0, 0, 0, 0, 0, 0) //Pay attention to the positioning of the parameter
	return uintptr_to_str_res(resptr)
}
