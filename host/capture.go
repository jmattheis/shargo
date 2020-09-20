package host

import (
	"bytes"
	"encoding/base64"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"image"
	"image/draw"
	"image/png"
)

const mouse string = `iVBORw0KGgoAAAANSUhEUgAAAA4AAAAPCAQAAAB+HTb/AAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QA/4ePzL8AAAAHdElNRQfkCRQKHQBhkbhHAAABN0lEQVQY02PmPMwv8v2m0vf3DJiAWdYkq1TW8vyXnw+Ffn/HkDaf8OL9/+Xf/TayujCwcqBJsnouf/7///9n/2d+dFrIYMbAyISQk2TgC1zz8////////7/7v/ul2QQGzf9IeoXTT3/7DwOX/9c+0K5jkDODyBUzaFSe+o8Af/6f/p97XSCXgYeJgaGP4camvc8QBv1k+P2flYNRkoEd6iGXhY+g+n79L/8iUcOg8J+BgYGJgYGBQfjvsVVHfjAwMDC8YfjJ4MQhJcjwnBHJUQKJR7/8//o/43HN8+//t/7R7YAZysDCEMygXHT8/6rvAtHcgd2vf/5f90ujnoENoVct7qnrCgaODQyiCdM+/vm/8ptKJ1IoMxcxGDIzMDAwMMrmLPpy+L/BNpRIeAJjsSh3629mkAUAxH+PoY4ST9YAAAAldEVYdGRhdGU6Y3JlYXRlADIwMjAtMDktMjBUMTA6MTk6MTYrMDA6MDAJDUswAAAAJXRFWHRkYXRlOm1vZGlmeQAyMDIwLTA5LTIwVDEwOjE0OjU3KzAwOjAwr5N9cgAAAABJRU5ErkJggg==`

var mouseImg image.Image

func init() {
	b, _ := base64.StdEncoding.DecodeString(mouse)
	mouseImg, _ = png.Decode(bytes.NewReader(b))
}

func ImageStream(idx int, images chan<- image.Image) error {
	for {
		mx, my := robotgo.GetMousePos()
		kek, _ := screenshot.CaptureDisplay(idx)
		rec := image.Rect(mx, my, mx+30, my+30)
		draw.Draw(kek, rec, mouseImg, image.Point{}, draw.Over)
		images <- kek
	}
}
