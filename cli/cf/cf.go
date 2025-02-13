package cf // import "gitlab.com/clburlison/cfprefs"

import (
	"unsafe"
)

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"

// Convert a Go string to a CFString
// Make sure to release the CFString when finished
func stringToCFString(s string) C.CFStringRef {
	return C.CFStringCreateWithCString(C.kCFAllocatorDefault, C.CString(s), C.kCFStringEncodingUTF8)
}

// Convert a CFString to a Go string
func cfstringToString(s C.CFStringRef) string {
	return C.GoString(C.CFStringGetCStringPtr(s, C.kCFStringEncodingUTF8))
}

// Convert a CFBoolean to a Go bool
func cfbooleanToBoolean(s C.CFBooleanRef) bool {
	if C.CFBooleanGetValue(s) == 1 {
		return true
	}
	return false
}

// Convert a CFData to a Go byte
func cfdataToData(s C.CFDataRef) []uint8 {
	d := C.GoBytes(unsafe.Pointer(C.CFDataGetBytePtr(s)), C.int(C.CFDataGetLength(s)))
	return d
}

// CFPreferencesCopyAppValue - Return a value from a preference
func CFPreferencesCopyAppValue(key string, domain string) interface{} {
	k := stringToCFString(key)
	defer release(C.CFTypeRef(k))

	d := stringToCFString(domain)
	defer release(C.CFTypeRef(d))
	
	if ret := C.CFPreferencesCopyAppValue(k, d); ret != 0 && C.CFGetTypeID(ret) == C.CFStringGetTypeID() {
		defer release(ret)
		return cfstringToString(C.CFStringRef(ret))
	}
	if ret := C.CFPreferencesCopyAppValue(k, d); ret != 0 && C.CFGetTypeID(ret) == C.CFBooleanGetTypeID() {
		defer release(ret)
		return cfbooleanToBoolean(C.CFBooleanRef(ret))
	}
	if ret := C.CFPreferencesCopyAppValue(k, d); ret != 0 && C.CFGetTypeID(ret) == C.CFDataGetTypeID() {
		defer release(ret)
		return cfdataToData(C.CFDataRef(ret))
	}
	return nil
}

func release(ref C.CFTypeRef) {
	if ref != 0 {
		C.CFRelease(ref)
	}
}
