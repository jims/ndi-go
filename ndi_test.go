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

func TestInitialization(t *testing.T) {
	libDir := os.Getenv("NDI_RUNTIME_DIR_V3")
	if libDir == "" {
		t.Fatal("ndi sdk is not installed")
	}

	if err := LoadAndInitialize(path.Join(libDir, ndiLibName)); err != nil {
		t.Fatal(err)
	}

	t.Logf("Version string is: %s", Version())
	UnloadAndDestroy()
}
