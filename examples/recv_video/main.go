/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/diskett-io/ndi-go"
)

const (
	ndiLibName    = "Processing.NDI.Lib.x64.dll"
	ndiSourceName = "ndi-go test"
)

func initializeNDI() {
	libDir := os.Getenv("NDI_RUNTIME_DIR_V3")
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
		for _, source := range findInst.GetCurrentSources() {
			name := source.Name()

			if name == ndiSourceName {
				addr := source.Address()
				recvSettings := ndi.NewRecvCreateSettings()
				recvSettings.SourceToConnectTo = *source

				recvInst = ndi.NewRecvInstanceV2(recvSettings)

				if recvInst == nil {
					log.Printf("unable to connect to %s, %s\n", name, addr)
					continue
				}

				fmt.Printf("Connected to %s, %s\n", name, addr)

				findInst.Destroy()
				pool.Release(findSettings)
				break
			}
		}
	}

	defer recvInst.Destroy()

	if !recvInst.SetTally(&ndi.Tally{OnProgram: true, OnPreview: true}) {
		log.Println("could not set tally")
	}

	for recvInst.GetNumConnections(1000) == 0 {
		fmt.Println("connections..", recvInst.GetNumConnections(1000))
	}

	fmt.Println("Reading video...")

	for {
		var (
			vf ndi.VideoFrameV2
			af ndi.AudioFrameV2
			mf ndi.MetadataFrame
		)

		vf.SetDefault()
		af.SetDefault()
		mf.SetDefault()

		ft, _ := recvInst.CaptureV2(nil, nil, nil, 1000)
		switch ft {
		case ndi.FrameTypeNone:
			fmt.Println("FrameTypeNone")
		case ndi.FrameTypeVideo:
			fmt.Println("FrameTypeVideo")
			recvInst.FreeVideoV2(&vf)
		case ndi.FrameTypeAudio:
			fmt.Println("FrameTypeAudio")
			recvInst.FreeAudioV2(&af)
		case ndi.FrameTypeMetadata:
			fmt.Println("FrameTypeMetadata")
			recvInst.FreeMetadataV2(&mf)
		case ndi.FrameTypeStatusChange:
			fmt.Println("FrameTypeStatusChange")
		default:
			fmt.Println("Unknown frame type!")
		}
	}
}
