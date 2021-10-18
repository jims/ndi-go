/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/FlowingSPDG/ndi-go"
)

const (
	ndiLibName  = "Processing.NDI.Lib.x64.dll"
	scanTimeout = 5000
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

	currentSources := make(map[string]string)
	for {
		inst.WaitForSources(scanTimeout)
		sources := inst.GetCurrentSources()

		var newListing bool
		if len(currentSources) != len(sources) {
			newListing = true
		} else {
			for _, source := range sources {
				name := source.Name()
				addr := source.Address()

				if n, ok := currentSources[addr]; !ok || n != name {
					newListing = true
					break
				}
			}
		}

		if newListing {
			fmt.Printf("%d available source(s):\n", len(sources))
			currentSources = make(map[string]string)

			for _, source := range sources {
				name := source.Name()
				addr := source.Address()

				currentSources[addr] = name
				fmt.Printf("Name: %s, Address: %s\n", name, addr)
			}
		}
	}
}
