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
func (inst *RecvInstance) SetTally(tally *Tally) bool {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibRecvSetTally, 2, uintptr(unsafe.Pointer(inst)), uintptr(unsafe.Pointer(tally)), 0)
	if eno != 0 {
		panic(eno)
	}
	return ret != 0
}

//This function will send a meta message to the source that we are connected too. This returns FALSE if we are
//not currently connected to anything.
func (inst *RecvInstance) SendMetadata(mf *MetadataFrame) bool {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibRecvSendMetadata, 2, uintptr(unsafe.Pointer(inst)), uintptr(unsafe.Pointer(mf)), 0)
	if eno != 0 {
		panic(eno)
	}
	return ret != 0
}

func (inst *RecvInstance) CaptureV2(vf *VideoFrameV2, af *AudioFrameV2, mf *MetadataFrame, timeoutInMs uint32) FrameType {
	ret, _, _ := syscall.Syscall6(
		funcPtrs.NDIlibRecvCaptureV2,
		5,
		uintptr(unsafe.Pointer(inst)),
		uintptr(unsafe.Pointer(vf)),
		uintptr(unsafe.Pointer(af)),
		uintptr(unsafe.Pointer(mf)),
		uintptr(timeoutInMs),
		0,
	)

	return FrameType(ret)
}

func (inst *RecvInstance) FreeVideoV2(vf *VideoFrameV2) {
	if _, _, eno := syscall.Syscall(funcPtrs.NDIlibRecvFreeVideoV2, 2, uintptr(unsafe.Pointer(inst)), uintptr(unsafe.Pointer(vf)), 0); eno != 0 {
		panic(eno)
	}
}

func (inst *RecvInstance) FreeAudioV2(af *AudioFrameV2) {
	if _, _, eno := syscall.Syscall(funcPtrs.NDIlibRecvFreeAudioV2, 2, uintptr(unsafe.Pointer(inst)), uintptr(unsafe.Pointer(af)), 0); eno != 0 {
		panic(eno)
	}
}

func (inst *RecvInstance) FreeMetadataV2(mf *MetadataFrame) {
	if _, _, eno := syscall.Syscall(funcPtrs.NDIlibRecvFreeMetadata, 2, uintptr(unsafe.Pointer(inst)), uintptr(unsafe.Pointer(mf)), 0); eno != 0 {
		panic(eno)
	}
}

//Is this receiver currently connected to a source on the other end, or has the source not yet been found or is no longe ronline.
//This will normally return 0 or 1.
func (inst *RecvInstance) GetNumConnections(timeoutInMs uint32) (int, error) {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibRecvGetNoConnections, 2, uintptr(unsafe.Pointer(inst)), uintptr(timeoutInMs), 0)
	if eno != 0 {
		return 0, Error{eno}
	}
	return int(ret), nil
}
