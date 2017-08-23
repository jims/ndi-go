/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package ndi

import "testing"

func TestInitialization(t *testing.T) {
	if err := LoadAndInitialize("C:\\Program Files\\NewTek\\NewTek NDI SDK\\Bin\\x64\\Processing.NDI.Lib.x64.dll"); err != nil {
		t.Fatal(err)
	}
	UnloadAndDestroy()
}
