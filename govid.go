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

// renders onto canvas
// mask gets wiped
// canvas and mask are mutated, images are not
func renderNewFrame( canvas *opencv.IplImage, images []*opencv.IplImage, mask *opencv.IplImage) {
	opencv.Zero(mask)
	opencv.Copy(images[rand.Int() % len(images)], canvas, nil)
	rgb := opencv.ScalarAll(255.0)
	for i := 0 ; i < canvas.Width()/2; i+=20 {
		pt1 := opencv.Point{ ((-2 * canvas.Width() / 4) + i), 0 }
		pt2 := opencv.Point{ i, canvas.Height() }
		pt3 := opencv.Point{ (6 * canvas.Width() / 4) - i, 0 }
		pt4 := opencv.Point{ canvas.Width() - i, canvas.Height() }
		opencv.Line(mask, pt1, pt2, rgb, 200, 8, 0)		
		opencv.Line(mask, pt3, pt4, rgb, 200, 8, 0)		
	}
	opencv.Copy(images[rand.Int() % len(images)], canvas, mask)
}

var masks [4000]*opencv.IplImage

func getMask(base int, width int, height int) *opencv.IplImage {
	if (masks[base] == nil) {
		masks[base] = opencv.CreateImage(width, height, 8, 1)
		opencv.Zero(masks[base])
		//opencv.Copy(image, canvas, nil)
		rgb := opencv.ScalarAll(255.0)
		for i := 0 ; i < width/2 + base; i+=5 {
			pt1 := opencv.Point{ ((-2 * width / 4) + i), height }
			pt2 := opencv.Point{ i, 0 }
			pt3 := opencv.Point{ (6 * width / 4) - i, height }
			pt4 := opencv.Point{ width - i, 0 }
			opencv.Line(masks[base], pt1, pt2, rgb, 200, 8, 0)		
			opencv.Line(masks[base], pt3, pt4, rgb, 200, 8, 0)		
		}
	}
	return masks[base]
}




func renderTriangle(base int, canvas *opencv.IplImage, image *opencv.IplImage, mask *opencv.IplImage) {
	// opencv.Zero(mask)
	// //opencv.Copy(image, canvas, nil)
	// rgb := opencv.ScalarAll(255.0)
	// for i := 0 ; i < canvas.Width()/2 + base; i+=5 {
	// 	pt1 := opencv.Point{ ((-2 * canvas.Width() / 4) + i), canvas.Height() }
	// 	pt2 := opencv.Point{ i, 0 }
	// 	pt3 := opencv.Point{ (6 * canvas.Width() / 4) - i, canvas.Height() }
	// 	pt4 := opencv.Point{ canvas.Width() - i, 0 }
	// 	opencv.Line(mask, pt1, pt2, rgb, 200, 8, 0)		
	// 	opencv.Line(mask, pt3, pt4, rgb, 200, 8, 0)		
	// }
	ourMask := getMask(base, canvas.Width(), canvas.Height())
	opencv.Copy(image, canvas, ourMask)
}


func lineMain(canvas *opencv.IplImage, images []*opencv.IplImage, mask *opencv.IplImage, altMask *opencv.IplImage,  vw *opencv.VideoWriter) {

	rgb := opencv.ScalarAll(255.0)
	wl := NewWLine()
	imin := rand.Int() % len(images)
	//for i := 0 ; i < 4738*30; i++ {
	for i := 0 ; i < 300; i++ {
		var image = images[imin]
		opencv.Copy(image, canvas, mask)
		vw.WriteFrame(canvas)
		pt1, pt2 := wl.line(canvas.Width(),canvas.Height())
		fmt.Println("%#v", wl)
		opencv.Line(mask, pt1, pt2, rgb, 200, 4, 0)
		wl.next()
		if (rand.Int() % 25 == 0 || pt1.X >= canvas.Width() && pt2.X >= canvas.Width()) {
			opencv.Zero(mask)
			wl.jumble()
			opencv.Copy(image, canvas, nil)
			renderNewFrame( canvas, images, altMask )			
			imin = rand.Int() % len(images)
		}
		fmt.Printf("%d\n", i)		
	}
}


func triangleMain(canvas *opencv.IplImage, images []*opencv.IplImage, mask *opencv.IplImage, altMask *opencv.IplImage,  vw *opencv.VideoWriter) {


	imin := rand.Int() % len(images)
	count := 0
	for i := 0 ; i < 600; i++ {
		var image = images[imin]
		opencv.Copy(image, canvas, mask)
		vw.WriteFrame(canvas)
		renderTriangle(count, canvas, image, altMask)
		count += 5
		if (rand.Int() % 100 == 0 || count > 200) {
			opencv.Zero(mask)
			opencv.Copy(image, canvas, nil)
			renderNewFrame( canvas, images, altMask )			
			imin = rand.Int() % len(images)
			count = 0
		}
		fmt.Printf("%d\n", i)		
	}
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
	canvas := opencv.CreateImage(images[0].Width(), images[0].Height(), 8, 3)
	opencv.Copy(images[0], canvas, nil)
	mask := opencv.CreateImage(canvas.Width(), canvas.Height(), 8, 1)
	altMask := opencv.CreateImage(canvas.Width(), canvas.Height(), 8, 1)
	opencv.Zero(mask)


	vw := opencv.NewVideoWriter("out.mkv", int(opencv.FOURCC('X','V','I','D')), 30.0,canvas.Width(),canvas.Height(),1)
	if vw == nil {
		panic("No video writer!")
	}

	triangleMain(canvas, images, mask, altMask,  vw);
	//lineMain(canvas, images, mask, altMask,  vw);
	vw.Release();
	os.Exit(0);	


}
