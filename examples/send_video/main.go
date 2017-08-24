/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"crypto/rand"
	"log"
	"os"
	"path"
	"time"

	ndi "github.com/diskett-io/ndi-go"
)

const ndiLibName = "Processing.NDI.Lib.x64.dll"

func main() {
	libDir := os.Getenv("NDI_RUNTIME_DIR_V3")
	if libDir == "" {
		log.Fatalln("ndi sdk is not installed")
	}

	if err := ndi.LoadAndInitialize(path.Join(libDir, ndiLibName)); err != nil {
		log.Fatalln(err)
	}

	pool := ndi.NewObjectPool()
	settings := pool.NewSendCreateSettings("ndi-go test", "", false, false)
	inst := ndi.SendCreate(settings)

	frameData := make([]byte, 1920*1080*4)
	frame := ndi.VideoFrameV2{
		FourCC:     ndi.FourCCTypeBGRA,
		Xres:       1920,
		Yres:       1080,
		LineStride: 1920 * 4,
		Data:       &frameData[0],
	}

	log.Println("Streaming video frame...")

	ticker := time.NewTicker(time.Millisecond * (1000 / 25))
	for range ticker.C {
		if _, err := rand.Read(frameData); err != nil {
			log.Fatalln(err)
		}
		ndi.SendSendVideoV2(inst, &frame)
	}
	ticker.Stop()

	ndi.SendDestroy(inst)
	ndi.DestroyAndUnload()
}
