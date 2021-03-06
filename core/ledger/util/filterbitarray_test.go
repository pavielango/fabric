/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// To run this test by itself, enter `go test -run FilterBitArray` from this folder
package util

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestFilterBitArrayFixed(t *testing.T) {
	ba := NewFilterBitArray(25)
	// Note that capacity may be greater than 25
	t.Logf("FilterBitArray capacity: %d\n", ba.Capacity())
	var i uint
	for i = 0; i < ba.Capacity(); i++ {
		ba.Set(i)
	}
	// All bits must be set
	for i = 0; i < ba.Capacity(); i++ {
		if !ba.IsSet(i) {
			t.FailNow()
		}
	}
	for i = 0; i < ba.Capacity(); i++ {
		ba.Unset(i)
	}
	// All bits must be unset
	for i = 0; i < ba.Capacity(); i++ {
		if ba.IsSet(i) {
			t.FailNow()
		}
	}
}

func TestFilterBitArraySparse(t *testing.T) {
	ba := new(FilterBitArray)
	// test byte boundary
	ba.Unset(0)
	ba.Set(8)
	ba.Unset(9)
	ba.Set(116)
	if ba.IsSet(0) {
		t.FailNow()
	}
	if !ba.IsSet(8) {
		t.FailNow()
	}
	if ba.IsSet(9) {
		t.FailNow()
	}
	if !ba.IsSet(116) {
		t.FailNow()
	}
}

func TestFilterBitArrayIO(t *testing.T) {
	ba := NewFilterBitArray(20)
	var i uint
	for i = 0; i < 20; i++ {
		if i%2 == 0 {
			ba.Set(i)
		} else {
			ba.Unset(i)
		}
	}
	b := ba.ToBytes()
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, b); err != nil {
		t.Fatalf("binary.Write failed: %s", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, b); err != nil {
		t.Fatalf("binary.Read failed: %s", err)
	}
	ba.FromBytes(b)
	for i = 0; i < 20; i++ {
		if i%2 == 0 {
			if !ba.IsSet(i) {
				t.FailNow()
			}
		} else {
			if ba.IsSet(i) {
				t.FailNow()
			}
		}
	}
}
