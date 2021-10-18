/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/FlowingSPDG/ndi-go"
)

const (
	ndiLibName    = "Processing.NDI.Lib.x64.dll"
	ndiSourceName = "FL-9900K (Test Pattern)"
)

func initializeNDI() {
	libDir := os.Getenv("NDI_RUNTIME_DIR_V5")
	if libDir == "" {
		log.Fatalln("ndi sdk is not installed")
	}

	if err := ndi.LoadAndInitialize(path.Join(libDir, ndiLibName)); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	initializeNDI()
	defer ndi.DestroyAndUnload()

	pool := ndi.NewObjectPool()
	findSettings := pool.NewFindCreateSettings(true, "", "")
	findInst := ndi.NewFindInstanceV2(findSettings)
	if findInst == nil {
		log.Fatalln("could not create finder")
	}

	var recvInst *ndi.RecvInstance

	fmt.Println("Searching for NDI sources...")

	for recvInst == nil {
		srcs := findInst.GetCurrentSources()
		log.Printf("Got %d sources\n", len(srcs))
		for i, source := range srcs {
			name := source.Name()
			log.Printf("srcs[%d] : %s %s", i, name, source.Address())

			if name == ndiSourceName {
				addr := source.Address()
				recvSettings := ndi.NewRecvCreateSettings()
				recvSettings.SourceToConnectTo = *source

				recvInst = ndi.NewRecvInstanceV2(recvSettings)

				if recvInst == nil {
					log.Printf("unable to connect to %s, %s\n", name, addr)
					continue
				}

				log.Printf("Connected to %s, %s\n", name, addr)

				findInst.Destroy()
				pool.Release(findSettings)
				break
			}
		}
		time.Sleep(time.Second)
	}

	defer recvInst.Destroy()

	if !recvInst.SetTally(&ndi.Tally{OnProgram: true, OnPreview: true}) {
		log.Println("could not set tally")
	}

	go func() {
		for {
			i, err := recvInst.GetNumConnections(1000)
			if err != nil {
				log.Println("Failed to get numconnections:", err)
				return
			}
			if i != 0 {
				fmt.Println("connections..", i)
			}
			time.Sleep(time.Second)
		}

	}()

	fmt.Println("Reading NDI...")

	for {
		var (
			vf ndi.VideoFrameV2
			af ndi.AudioFrameV2
			mf ndi.MetadataFrame
		)

		vf.SetDefault()
		af.SetDefault()
		mf.SetDefault()

		ft := recvInst.CaptureV2(&vf, &af, &mf, 1000)
		switch ft {
		case ndi.FrameTypeNone:
			log.Println("FrameTypeNone")
		case ndi.FrameTypeVideo:
			log.Println("FrameTypeVideo")
			log.Printf("VideoFrame : %#v\n", vf)
			recvInst.FreeVideoV2(&vf)
		case ndi.FrameTypeAudio:
			log.Println("FrameTypeAudio")
			log.Printf("AudioFrame : %#v\n", af)
			recvInst.FreeAudioV2(&af)
		case ndi.FrameTypeMetadata:
			log.Println("FrameTypeMetadata")
			log.Printf("Metadata : %#v\n", mf)
			recvInst.FreeMetadataV2(&mf)
		case ndi.FrameTypeStatusChange:
			log.Println("FrameTypeStatusChange")
		default:
			log.Println("Unknown frame type!")
		}
	}
}
