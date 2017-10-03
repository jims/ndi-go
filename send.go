/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package ndi

import (
	"syscall"
	"unsafe"
)

type SendInstance struct{}

func NewSendInstance(settings *SendCreateSettings) *SendInstance {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibSendCreate, 1, uintptr(unsafe.Pointer(settings)), 0, 0)
	if eno != 0 {
		panic(eno)
	}
	return (*SendInstance)(unsafe.Pointer(ret))
}

func (inst *SendInstance) Destroy() {
	if _, _, eno := syscall.Syscall(funcPtrs.NDIlibSendDestroy, 1, uintptr(unsafe.Pointer(inst)), 0, 0); eno != 0 {
		panic(eno)
	}
}

//This will add a video frame.
func (inst *SendInstance) SendVideoV2(frame *VideoFrameV2) {
	if _, _, eno := syscall.Syscall(funcPtrs.NDIlibSendSendVideoV2, 2, uintptr(unsafe.Pointer(inst)), uintptr(unsafe.Pointer(frame)), 0); eno != 0 {
		panic(eno)
	}
}

type SendError struct {
	syscall.Errno
}

//Get the current number of receivers connected to this source. This can be used to avoid even rendering when nothing is connected to the video source.
//which can significantly improve the efficiency if you want to make a lot of sources available on the network. If you specify a timeout that is not
//0 then it will wait until there are connections for this amount of time.
func (inst *SendInstance) GetNumConnections(timeoutInMs uint32) (int, error) {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibSendGetNoConnections, 2, uintptr(unsafe.Pointer(inst)), uintptr(timeoutInMs), 0)
	if eno != 0 {
		return int(ret), SendError{eno}
	}
	return int(ret), nil
}
