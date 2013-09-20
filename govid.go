//
// Copyright 2011 <chaishushan@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
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
	//vw := opencv.NewVideoWriter("out.mkv", int(opencv.FOURCC('X','V','I','D')), 30.0,1456,1000,1)
	vw := opencv.NewVideoWriter("out.mkv", int(opencv.FOURCC('X','V','I','D')), 30.0,1456,1000,1)
	if vw == nil {
		panic("No video writer!")
	}
	defer vw.Release()
	for i := 0 ; i < 120; i++ {
		vw.WriteFrame( img0 )
		rgb := opencv.ScalarAll(255.0)
		pt1 := opencv.Point{1456/2-10*i, 500-200}
		pt2 := opencv.Point{1456/2+10*i, 500+200}
		opencv.Line(img0, pt1, pt2, rgb, 5, 8, 0)
		fmt.Printf("%d\n", i)
	}
	os.Exit(0);
}
