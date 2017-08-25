/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/diskett-io/ndi-go"
)

const (
	ndiLibName  = "Processing.NDI.Lib.x64.dll"
	scanTimeout = 5000
)

func main() {
	runtime.LockOSThread()

	libDir := os.Getenv("NDI_RUNTIME_DIR_V3")
	if libDir == "" {
		log.Fatalln("ndi sdk is not installed")
	}

	if err := ndi.LoadAndInitialize(path.Join(libDir, ndiLibName)); err != nil {
		log.Fatalln(err)
	}

	pool := ndi.NewObjectPool()
	settings := pool.NewFindCreateSettings(true, "", "")
	inst := ndi.NewFindInstanceV2(settings)
	if inst == nil {
		log.Fatalln("could not create finder")
	}

	defer func() {
		inst.Destroy()
		ndi.DestroyAndUnload()
	}()

	fmt.Println("Searching for NDI sources...")

	listed := make(map[string]struct{})
	for {
		inst.WaitForSources(scanTimeout)
		source := inst.GetCurrentSources()

		for _, source := range source {
			name := source.Name()
			addr := source.Address()

			key := name + addr
			if _, ok := listed[key]; !ok {
				fmt.Printf("Source: %s, Address: %s\n", name, addr)
				listed[key] = struct{}{}
			}
		}
	}
}
