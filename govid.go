//
// Copyright 2011 <chaishushan@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"opencv"
)

func main() {
	filename := "1.png"
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}
	img0 := opencv.LoadImage(filename)
	if img0 == nil {
		panic("LoadImage fail")
	}
	defer img0.Release()

	os.Exit(0);
}
