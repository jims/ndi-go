/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package ndi

import (
	"os"
	"path"
	"testing"
)

const ndiLibName = "Processing.NDI.Lib.x64.dll"

func doInit(t *testing.T) {
	libDir := os.Getenv("NDI_RUNTIME_DIR_V3")
	if libDir == "" {
		t.Fatal("ndi sdk is not installed")
	}

	if err := LoadAndInitialize(path.Join(libDir, ndiLibName)); err != nil {
		t.Fatal(err)
	}
}

func TestInitialization(t *testing.T) {
	doInit(t)
	t.Logf("Version string is: %s", Version())
	DestroyAndUnload()
}

func TestFrame(t *testing.T) {
	doInit(t)

	pool := NewObjectPool()
	settings := pool.NewSendCreateSettings("ndi-go test", "", false, false)
	inst := SendCreate(settings)

	frameData := make([]byte, 1920*1080*4)
	frame := VideoFrameV2{
		FourCC:     FourCCTypeBGRA,
		Xres:       1920,
		Yres:       1080,
		LineStride: 1920 * 4,
		Data:       &frameData[0],
	}

	SendSendVideoV2(inst, &frame)

	SendDestroy(inst)
	DestroyAndUnload()
}
