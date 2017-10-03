/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package ndi

import (
	"syscall"
	"unsafe"
)

type Source struct {
	name, address *byte
}

func (s *Source) Name() string {
	if s.name == nil {
		return ""
	}
	return goStringFromCString(uintptr(unsafe.Pointer(s.name)))
}

func (s *Source) Address() string {
	if s.address == nil {
		return ""
	}
	return goStringFromCString(uintptr(unsafe.Pointer(s.address)))
}

type FindInstance struct{}

func NewFindInstanceV2(settings *FindCreateSettings) *FindInstance {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibFindCreateV2, 1, uintptr(unsafe.Pointer(settings)), 0, 0)
	if eno != 0 {
		panic(eno)
	}
	return (*FindInstance)(unsafe.Pointer(ret))
}

func (inst *FindInstance) Destroy() {
	if _, _, eno := syscall.Syscall(funcPtrs.NDIlibFindDestroy, 1, uintptr(unsafe.Pointer(inst)), 0, 0); eno != 0 {
		panic(eno)
	}
}

type FindError struct {
	syscall.Errno
}

//This will allow you to wait until the number of online sources have changed.
func (inst *FindInstance) WaitForSources(timeoutInMs uint32) (int, error) {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibFindWaitForSources, 2, uintptr(unsafe.Pointer(inst)), uintptr(timeoutInMs), 0)
	if eno != 0 {
		return int(ret), FindError{eno}
	}
	return int(ret), nil
}

//This function will recover the current set of sources (i.e. the ones that exist right this second).
func (inst *FindInstance) GetCurrentSources() []*Source {
	var numSources uint32
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibFindGetCurrentSources, 2, uintptr(unsafe.Pointer(inst)), uintptr(unsafe.Pointer(&numSources)), 0)
	if eno != 0 {
		panic(eno)
	}

	sources := make([]*Source, numSources)
	for i, s := range sources {
		sources[i] = (*Source)(unsafe.Pointer(ret))
		ret += unsafe.Sizeof(*s)
	}
	return sources
}
