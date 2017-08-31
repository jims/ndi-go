/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package ndi

import (
	"reflect"
	"testing"
)

var fieldAlignments = map[string]int{
	"bool":        1,
	"int":         4,
	"int32":       4,
	"int64":       8,
	"uint32":      4,
	"float32":     4,
	"FrameFormat": 4,
}

func fieldAlignmentTest(t *testing.T, v interface{}) {
	value := reflect.ValueOf(v)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		ty := field.Type()

		name := ty.Name()
		if name == "" {
			continue
		}

		n, ok := fieldAlignments[name]
		if !ok {
			t.Fatalf("Invalid field type: %s.", name)
		}

		if n != ty.FieldAlign() {
			t.Errorf("Invalid field alignment for field '%s' in struct '%s'. Expected %d but result is %d.", name, value.Type().Name(), n, ty.FieldAlign())
		}
	}
}

func TestFieldAlignment(t *testing.T) {
	var vf VideoFrameV2
	fieldAlignmentTest(t, vf)

	var scs SendCreateSettings
	fieldAlignmentTest(t, scs)

	var fcs FindCreateSettings
	fieldAlignmentTest(t, fcs)
}

func checkTypeSize(t *testing.T, v interface{}, sz uintptr) {
	if s := reflect.TypeOf(v).Size(); s != sz {
		t.Errorf("Invalid size of struct '%s'. Expected %d but result is %d.", reflect.TypeOf(v).Name(), sz, s)
	}
}

func TestStructSizes(t *testing.T) {
	var vf VideoFrameV2
	checkTypeSize(t, vf, 72)

	var af AudioFrameV2
	checkTypeSize(t, af, 56)

	var scs SendCreateSettings
	checkTypeSize(t, scs, 24)

	var fcs FindCreateSettings
	checkTypeSize(t, fcs, 24)
}
