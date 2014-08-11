// Copyright 2013 Abram Hindle <abram.hindle@softwareprocess.es>
// Copyright 2011 <chaishushan@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main
//	"opencv"

import (
	"fmt"
	"os"
	"github.com/chai2010/go-opencv/opencv"
	"math/rand"
	"time"
	"unsafe"
	"flag"
)
//#include <opencv/cv.h>
//#include <opencv/highgui.h>
//#include <opencv2/photo/photo_c.h>
//#include <opencv2/imgproc/imgproc_c.h>
//#cgo linux  pkg-config: opencv
//#cgo darwin pkg-config: opencv
//#cgo windows LDFLAGS: -lopencv_core242.dll -lopencv_imgproc242.dll -lopencv_photo242.dll -lopencv_highgui242.dll -lstdc++
import "C"
func flipImage(src *opencv.IplImage, dst *opencv.IplImage, flags int) {
        C.cvFlip(unsafe.Pointer(src), unsafe.Pointer(dst), C.int(flags))
}

type kali struct {
    ww int
    frames int
}

func kaliMain(kal kali, canvas *opencv.IplImage, images []*opencv.IplImage, mask *opencv.IplImage, altMask *opencv.IplImage,  vw *opencv.VideoWriter) {

	i := 0
	flip_or_not := (rand.Int() % 2 == 1)
	w := canvas.Width()
	opencv.Zero(canvas)
	ww := kal.ww
	wmul := ww
	wn := w / ww
	frames := kal.frames
	for j := 0 ; j < frames; j++ {
		var image = images[(j/w) % len(images)] //(i / wn) % len(images)]//rand.Int() % len(images)]
		for k := 0 ; k < wn; k++ {	

			//subimage := new(opencv.Mat)
			var r1 = opencv.Rect{}		
			mh := canvas.Height()
			if (mh > image.Height()) {
				mh = image.Height()
			}
			r1.Init((wmul*k) % (canvas.Width()),0,ww,mh)
			var r2 = opencv.Rect{}
			r2.Init(j % (image.Width()-ww), 0,ww, mh)
			canvas.SetROI(r1)
			image.SetROI(r2)
			//subimage = opencv.GetSubRect(opencv.Arr(canvas.ImageData()), subimage, r)
			//opencv.Copy(image, canvas, nil)
			//flipImage(image, canvas, 1 * ((1+k)%2))
			if (flip_or_not) {
				if (k%2 == 0) {
					flipImage(image, canvas, 1)// * ((1+k)%2))
				} else {
					opencv.Copy(image, canvas, nil)
				}
			} else {
				// was
				flipImage(image, canvas, -1 * (k%2))
			}
			
			image.ResetROI()
			canvas.ResetROI()
			i = i + 1
		}
		vw.WriteFrame(canvas)
	}
}



func main() {
	rand.Seed( time.Now().UTC().UnixNano())
	filenames := []string{}//"1.png"}
	frames := flag.Int("frames", 100, "n frames")
	outfile := flag.String("out", "out.mkv", "output movie")
	flag.Parse()
	fmt.Printf("%#v\n", *frames)

	if len(os.Args) >= 2 {
		//filenames = os.Args[1:]
		filenames = flag.Args()
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
	mw := images[0].Width()
	mh := images[0].Height()
	for _, image := range images {
		iw := image.Width()
		ih := image.Height()
		if (iw != mw) {
			panic("Not all images are the same size")
		}
		if (ih != mh) {
			panic("Not all images are the same size")
		}
	}
	fmt.Printf("Created Canvas")
	canvas := opencv.CreateImage(mw, mh, 8, 3)
	fmt.Printf("Blit")
	opencv.Copy(images[0], canvas, nil)
	fmt.Printf("Make Mask")
	mask := opencv.CreateImage(canvas.Width(), canvas.Height(), 8, 1)
	altMask := opencv.CreateImage(canvas.Width(), canvas.Height(), 8, 1)
	opencv.Zero(mask)


	//vw := opencv.NewVideoWriter("out.mkv", int(opencv.FOURCC('X','V','I','D')), 30.0,canvas.Width(),canvas.Height(),1)
	//vw := opencv.NewVideoWriter("out.flv", int(opencv.FOURCC('H','2','6','4')), 30.0,canvas.Width(),canvas.Height(),1)
	//vw := opencv.NewVideoWriter("out.flv", int(opencv.FOURCC('A','V','C','1')), 30.0,canvas.Width(),canvas.Height(),1)
	vw := opencv.NewVideoWriter(*outfile, int(opencv.FOURCC('F','M','P','4')), 30.0,canvas.Width(),canvas.Height(),1)
	if vw == nil {
		panic("No video writer!")
	}
	sizes := []int{10,20,40,80,160}
	var k = kali{ frames: *frames, ww: sizes[ rand.Int() % len(sizes) ] }
	kaliMain(k, canvas, images, mask, altMask,  vw);
	vw.Release();
	os.Exit(0);	


}
