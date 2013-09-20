// Copyright 2013 Abram Hindle <abram.hindle@softwareprocess.es>
// Copyright 2011 <chaishushan@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"opencv"
	"math/rand"
)
// func fillConvexPoly(image *IplImage, pts, pt2 Point, color Scalar, thickness, line_type, shift int) {
// 	C.cvLine(
// 		unsafe.Pointer(image),
// 		C.cvPoint(C.int(pt1.X), C.int(pt1.Y)),
// 		C.cvPoint(C.int(pt2.X), C.int(pt2.Y)),
// 		(C.CvScalar)(color),
// 		C.int(thickness), C.int(line_type), C.int(shift),
// 	)
// 	//Scalar
// }
// 
// //CVAPI(void) cvFillConvexPoly(CvArr* img, const CvPoint* pts, int npts, CvScalar color, int line_type=8, int shift=0 )
// 

type WLine struct {
	xoff int
	x int
	rate int	
}
func (wl *WLine) update(move int) {
	wl.x = wl.x + move
}
func (wl *WLine) next() {
	wl.x = wl.x + wl.rate
}

func (wl WLine) line(w, h int) (opencv.Point, opencv.Point) {
	pt1 := opencv.Point{ wl.x,	0}
	pt2 := opencv.Point{ wl.x - wl.xoff, h}
	return pt1, pt2	
}
func (wl *WLine) jumble() {
	wl.xoff = rand.Int() % 400 - 200
	wl.x = wl.xoff
	wl.rate = 1 + (rand.Int() % 10)
	if (wl.xoff > 0) {
		wl.x = -1 * wl.xoff
	}
}
func NewWLine() *WLine {
	w := new(WLine)
	w.jumble()
	return w
}


func genline(canvas *opencv.IplImage) (opencv.Point, opencv.Point) {
	pt1 := opencv.Point{ 
		(rand.Int() % canvas.Width()),
		0}
	pt2 := opencv.Point{ 
		(rand.Int() % canvas.Width()),
		canvas.Height()}
	return pt1, pt2
}


func main() {
	filenames := []string{}//"1.png"}
	if len(os.Args) >= 2 {
		filenames = os.Args[1:]
	}
	images := make([]*opencv.IplImage, len(filenames) )
	for i, filename := range filenames {
		images[i] = opencv.LoadImage( filename )
		fmt.Printf("%#v\n",filename)
		if images[i] == nil {
			// how to insert the filename into that string?
			panic("LoadImage fail")
		}
	}
	canvas := images[0]
	mask := opencv.CreateImage(canvas.Width(), canvas.Height(), 8, 1)
	opencv.Zero(mask)

	vw := opencv.NewVideoWriter("out.mkv", int(opencv.FOURCC('X','V','I','D')), 30.0,canvas.Width(),canvas.Height(),1)
	if vw == nil {
		panic("No video writer!")
	}

	rgb := opencv.ScalarAll(255.0)
	wl := NewWLine()
	imin := rand.Int() % len(images)
	for i := 0 ; i < 4738*30; i++ {
		var image = images[imin]
		//vw.WriteFrame( images[i % len(images)]  )
		opencv.Copy(image, canvas, mask)
		vw.WriteFrame(canvas)
		pt1, pt2 := wl.line(canvas.Width(),canvas.Height())

		//pt1 := opencv.Point{ (rand.Int() % canvas.Width()),
		//	(rand.Int() % canvas.Height())}
		//pt2 := opencv.Point{ (rand.Int() % canvas.Width()),
		//	(rand.Int() % canvas.Height())}
		fmt.Println("%#v", wl)
		opencv.Line(mask, pt1, pt2, rgb, 200, 4, 0)
		wl.next()
		if (rand.Int() % 100 == 0 || pt1.X >= canvas.Width() && pt2.X >= canvas.Width()) {
			opencv.Zero(mask)
			wl.jumble()			
			opencv.Copy(image, canvas, nil)
			imin = rand.Int() % len(images)
		}
		if (rand.Int() % 15 == 0) {
			imin = rand.Int() % len(images)
		}
		fmt.Printf("%d\n", i)		
	}
	vw.Release();
	os.Exit(0);	
}
