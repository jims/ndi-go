/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package ndi

import (
	"syscall"
	"unsafe"
)

type RecvInstance struct{}

func NewRecvInstanceV2(settings *RecvCreateSettings) *RecvInstance {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibRecvCreateV2, 1, uintptr(unsafe.Pointer(settings)), 0, 0)
	if eno != 0 {
		panic(eno)
	}
	return (*RecvInstance)(unsafe.Pointer(ret))
}

func (inst *RecvInstance) Destroy() {
	if _, _, eno := syscall.Syscall(funcPtrs.NDIlibRecvDestroy, 1, uintptr(unsafe.Pointer(inst)), 0, 0); eno != 0 {
		panic(eno)
	}
}

//Set the up-stream tally notifications. This returns FALSE if we are not currently connected to anything. That
//said, the moment that we do connect to something it will automatically be sent the tally state.
func (inst *RecvInstance) SetTally(tally *Tally) {
	if _, _, eno := syscall.Syscall(funcPtrs.NDIlibRecvSetTally, 2, uintptr(unsafe.Pointer(inst)), uintptr(unsafe.Pointer(tally)), 0); eno != 0 {
		panic(eno)
	}
}
