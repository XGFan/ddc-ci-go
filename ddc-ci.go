package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	MOUSE_X     = 1000
	MOUSE_Y     = 1000
	M_POINT     = (MOUSE_X & 0xFFFFFFFF) | (MOUSE_Y << 32)
	InPutSource = '\x60'
	Brightness  = '\x10'
	Contrast    = '\x12'
	Volume      = '\x62'
	DP          = 15
	HDMI        = 17
)

var (
	user32, _                          = syscall.LoadLibrary("User32.dll")
	dxva2, _                           = syscall.LoadLibrary("dxva2.dll")
	monitorFromPoint, _                = syscall.GetProcAddress(user32, "MonitorFromPoint")
	GetPhysicalMonitorsFromHMONITOR, _ = syscall.GetProcAddress(dxva2, "GetPhysicalMonitorsFromHMONITOR")
	SetVCPFeature, _                   = syscall.GetProcAddress(dxva2, "SetVCPFeature")
	GetVCPFeatureAndVCPFeatureReply, _ = syscall.GetProcAddress(dxva2, "GetVCPFeatureAndVCPFeatureReply")
	DestroyPhysicalMonitor, _          = syscall.GetProcAddress(dxva2, "DestroyPhysicalMonitor")
)

func getMonitorHandle() (result uintptr) {
	ret, _, callErr := syscall.Syscall(
		monitorFromPoint,
		2,
		uintptr(M_POINT),
		uintptr(1),
		0)

	if callErr != 0 {
		abort("Call getMonitorHandle", callErr)
	}
	res := ret
	result = getPhysicalMonitor(res)
	return
}

func getPhysicalMonitor(handle uintptr) (result uintptr) {
	b := make([]byte, 256)
	_, _, callErr := syscall.Syscall(
		GetPhysicalMonitorsFromHMONITOR,
		3,
		handle,
		uintptr(1),
		uintptr(unsafe.Pointer(&b[0])))

	if callErr != 0 {
		abort("Call getPhysicalMonitor", callErr)
	}
	result = uintptr(b[0])
	return
}

func GetMonitorValue(key int32) int {
	mHandle := getMonitorHandle()
	var pvct uint32 = 0
	var curr uint32 = 0
	var max uint32 = 0

	_, _, callErr := syscall.Syscall6(
		GetVCPFeatureAndVCPFeatureReply,
		5,
		mHandle,
		uintptr(key),
		uintptr(unsafe.Pointer(&pvct)), //pvct
		uintptr(unsafe.Pointer(&curr)), //pdwCurrentValue
		uintptr(unsafe.Pointer(&max)),  //pdwMaximumValue
		0,                              //unused
	)
	if callErr != 0 {
		abort("Call GetVCPFeatureAndVCPFeatureReply", callErr)
	}
	return int(curr)
}

func SetMonitorValue(key int32, count int) {
	mHandle := getMonitorHandle()
	_, _, callErr := syscall.Syscall(
		SetVCPFeature,
		3,
		mHandle,
		uintptr(key),
		uintptr(count))

	if callErr != 0 {
		abort("Call SetVCPFeature", callErr)
	}
	destroyPhysicalMonitor(mHandle)
}

func destroyPhysicalMonitor(monitor uintptr) {
	_, _, callErr := syscall.Syscall(
		DestroyPhysicalMonitor,
		1,
		monitor,
		0,
		0)
	if callErr != 0 {
		abort("Call destroyPhysicalMonitor", callErr)
	}
}

func abort(funcName string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcName, err))
}
