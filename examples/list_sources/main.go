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
	maxRetries  = 10
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

	fmt.Println("Searching for NDI sources...")

	for n := 0; n < maxRetries; n++ {
		inst.WaitForSources(scanTimeout)
		source := inst.GetCurrentSources()

		if len(source) > 0 {
			for _, source := range source {
				fmt.Printf("Source: %s, Address: %s\n", source.Name(), source.Address())
			}
			break
		}
	}

	inst.Destroy()
	ndi.DestroyAndUnload()
}
