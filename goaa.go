package goaa

import (
	"errors"
	"image"
	"image/color"

	"github.com/aybabtme/rgbterm"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

var sampleChar = map[uint16]string{
	0x0000: "`",
	0x0001: ".",
	0x0002: "*",
	0x0010: ".",
	0x0020: "*",
	0x0011: "_",
	0x0012: "o",
	0x0021: "o",
	0x0022: "=",
	0x0100: ".",
	0x0200: "*",
	0x0101: "|",
	0x0102: "l",
	0x0201: "t",
	0x0202: "I",
	0x0110: "/",
	0x0120: "/",
	0x0210: "/",
	0x0220: "<",
	0x0111: "d",
	0x0112: "d",
	0x0121: "d",
	0x0122: "d",
	0x0211: "d",
	0x0212: "d",
	0x0221: "d",
	0x0222: "d",
	0x1000: "^",
	0x2000: "^",
	0x1001: "\\",
	0x1002: "\\",
	0x2001: "\\",
	0x2002: "\\",
	0x1010: "l",
	0x1020: "l",
	0x2010: "l",
	0x2020: "l",
	0x1011: "b",
	0x1012: "b",
	0x1021: "b",
	0x1022: "b",
	0x2011: "b",
	0x2012: "b",
	0x2021: "b",
	0x2022: "b",
	0x1100: "\"",
	0x1200: "\"",
	0x2100: "\"",
	0x2200: "\"",
	0x1101: "T",
	0x1102: "T",
	0x1201: "q",
	0x1202: "q",
	0x2101: "q",
	0x2102: "q",
	0x2201: "q",
	0x2202: "q",
	0x1110: "T",
	0x1120: "T",
	0x1210: "p",
	0x1220: "p",
	0x2110: "p",
	0x2120: "p",
	0x2210: "p",
	0x2220: "p",
	0x1111: "#",
	0x1112: "#",
	0x1121: "#",
	0x1122: "#",
	0x1211: "#",
	0x1212: "#",
	0x1221: "#",
	0x1222: "#",
	0x2111: "#",
	0x2112: "#",
	0x2121: "#",
	0x2122: "#",
	0x2211: "#",
	0x2212: "#",
	0x2221: "#",
	0x2222: "#",
}

func ConvertGray(img image.Image) (*image.Gray16, error) {
	bounds := img.Bounds()
	dest := image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.Gray16Model.Convert(img.At(x, y))
			gray, ok := c.(color.Gray16)
			if !ok {
				return nil, errors.New("color is miss match")
			}
			dest.Set(x, y, gray)
		}
	}
	return dest, nil
}

func ConvertASCII(img image.Image) ([]string, error) {
	if img.Bounds().Max.X > 800 {
		img = resize.Resize(800, 0, img, resize.Lanczos3)
	}
	if img.Bounds().Max.Y > 800 {
		img = resize.Resize(0, 800, img, resize.Lanczos3)
	}
	ret := []string{}
	gray, err := ConvertGray(imaging.AdjustContrast(imaging.Sharpen(img, 50), 10))
	if err != nil {
		return nil, err
	}

	var buffer = [2][2]uint16{}

	for y := gray.Bounds().Min.Y; y < gray.Bounds().Max.Y; y += 2 {
		line := ""
		for x := gray.Bounds().Min.X; x < gray.Bounds().Max.X; x++ {
			buffer[0] = [2]uint16{gray.At(x, y).(color.Gray16).Y, gray.At(x+1, y).(color.Gray16).Y}
			buffer[1] = [2]uint16{gray.At(x, y+1).(color.Gray16).Y, gray.At(x+1, y+1).(color.Gray16).Y}
			r, g, b, _ := img.At(x, y).RGBA()
			line += rgbterm.FgString(ConvertString(buffer), uint8(r>>8), uint8(g>>8), uint8(b>>8))
		}
		ret = append(ret, line)
	}
	return ret, nil
}

func ConvertString(p [2][2]uint16) string {
	var pattern uint16
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			pattern |= checkColorGroup(p[y][x])
			pattern <<= 4
		}
	}
	return sampleChar[pattern]
}

//輝度を調べる
func checkColorGroup(n uint16) uint16 {
	if n > ((0xffff * 2) / 3) {
		return 2
	} else if n > ((0xffff * 1) / 3) {
		return 1
	}
	return 0
}
