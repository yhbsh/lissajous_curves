package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"net/http"
)

var palette = []color.Color{color.RGBA{0, 0, 0, 255}, color.RGBA{170, 220, 150, 255}, color.RGBA{60, 120, 150, 255}, color.RGBA{60, 90, 240, 255}}

func main() {
	fmt.Printf("[INFO]: server is listening at port %d\n", 8080)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[INFO]: new request from %s\n", r.RemoteAddr)
	lissajous(w)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			if int(y+x)%3 == 0 {
				img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), 1)
			} else if int(y+x)%3 == 1 {
				img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), 2)
			} else {
				img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), 3)
			}

		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
