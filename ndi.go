/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package ndi

import (
	"errors"
	"log"
	"syscall"
	"unsafe"
)

var (
	alreadyLoadedErr     = errors.New("library is already loaded")
	loadProcsErr         = errors.New("failed to load library procs")
	initializeLibraryErr = errors.New("unable to initialize library")
)

var (
	ndiSharedLibrary syscall.Handle
	funcPtrs         *ndiLIBv3
)

type ObjectPool struct {
	objects map[interface{}]struct{}
}

func NewObjectPool() *ObjectPool {
	return &ObjectPool{make(map[interface{}]struct{})}
}

func (p *ObjectPool) Register(o interface{}) {
	if _, ok := p.objects[o]; ok {
		log.Fatalln("object is already in the object pool")
	}
	p.objects[o] = struct{}{}
}

func (p *ObjectPool) Release(o interface{}) {
	if _, ok := p.objects[o]; !ok {
		log.Fatalln("object was not found in the object pool")
	}
	delete(p.objects, o)
}

type SendCreateSettings struct {
	ndiName, groups        *byte
	clockVideo, clockAudio bool
}

func (p *ObjectPool) NewSendCreateSettings(name, groups string, clockVideo, clockAudio bool) *SendCreateSettings {
	var bNamePtr *byte
	if name != "" {
		bName := make([]byte, len(name)+1)
		copy(bName, name)
		bNamePtr = &bName[0]
	}

	var bGroupsPtr *byte
	if groups != "" {
		bGroups := make([]byte, len(groups)+1)
		copy(bGroups, groups)
		bGroupsPtr = &bGroups[0]
	}

	o := &SendCreateSettings{bNamePtr, bGroupsPtr, clockVideo, clockAudio}
	p.Register(o)
	return o
}

func LoadAndInitialize(path string) error {
	if ndiSharedLibrary != 0 {
		return alreadyLoadedErr
	}

	var err error
	if ndiSharedLibrary, err = syscall.LoadLibrary(path); err != nil {
		return err
	}

	var ndiLoadProc uintptr
	if ndiLoadProc, err = syscall.GetProcAddress(ndiSharedLibrary, "NDIlib_v3_load"); err != nil {
		syscall.FreeLibrary(ndiSharedLibrary)
		ndiSharedLibrary = 0
		return err
	}

	var (
		ret uintptr
		eno syscall.Errno
	)

	if ret, _, eno = syscall.Syscall(ndiLoadProc, 0, 0, 0, 0); eno != 0 {
		syscall.FreeLibrary(ndiSharedLibrary)
		ndiSharedLibrary = 0
		return eno
	}

	funcPtrs = (*ndiLIBv3)(unsafe.Pointer(ret))
	if funcPtrs == nil {
		syscall.FreeLibrary(ndiSharedLibrary)
		ndiSharedLibrary = 0
		return loadProcsErr
	}

	if ret, _, eno = syscall.Syscall(funcPtrs.NDIlibInitialize, 0, 0, 0, 0); eno != 0 {
		syscall.FreeLibrary(ndiSharedLibrary)
		ndiSharedLibrary = 0
		return eno
	}

	if ret == 0 {
		syscall.FreeLibrary(ndiSharedLibrary)
		ndiSharedLibrary = 0
		return initializeLibraryErr
	}
	return nil
}

func DestroyAndUnload() {
	if ndiSharedLibrary == 0 {
		return
	}

	if funcPtrs != nil {
		if _, _, eno := syscall.Syscall(funcPtrs.NDIlibDestroy, 0, 0, 0, 0); eno != 0 {
			panic(eno)
		}
	}

	syscall.FreeLibrary(ndiSharedLibrary)
	ndiSharedLibrary = 0
}

func Version() string {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibVersion, 0, 0, 0, 0)
	if eno != 0 {
		panic(eno)
	}
	return goStringFromConst(ret)
}

func IsSupportedCPU() bool {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibIsSupportedCPU, 0, 0, 0, 0)
	if eno != 0 {
		panic(eno)
	}
	return ret != 0
}

type SendInstance int

func SendCreate(settings *SendCreateSettings) *SendInstance {
	ret, _, eno := syscall.Syscall(funcPtrs.NDIlibSendCreate, 1, uintptr(unsafe.Pointer(settings)), 0, 0)
	if eno != 0 {
		panic(eno)
	}
	return (*SendInstance)(unsafe.Pointer(ret))
}

func SendDestroy(inst *SendInstance) {
	if _, _, eno := syscall.Syscall(funcPtrs.NDIlibSendDestroy, 1, uintptr(unsafe.Pointer(inst)), 0, 0); eno != 0 {
		panic(eno)
	}
}

func SendSendVideoV2(inst *SendInstance, frame *VideoFrameV2) {
	if _, _, eno := syscall.Syscall(funcPtrs.NDIlibSendSendVideoV2, 2, uintptr(unsafe.Pointer(inst)), uintptr(unsafe.Pointer(frame)), 0); eno != 0 {
		panic(eno)
	}
}
