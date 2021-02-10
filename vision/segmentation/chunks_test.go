package segmentation

import (
	"image"
	"image/color"
	"testing"

	"github.com/viamrobotics/robotcore/vision"
)

type chunkImageDebug struct {
}

func (cid *chunkImageDebug) Process(d *vision.MultipleImageTestDebugger, fn string, img vision.Image) error {

	type AShape struct {
		Start      image.Point
		PixelRange []int
	}

	type imgConfig struct {
		Shapes []AShape
	}

	cfg := imgConfig{}
	err := d.CurrentImgConfig(&cfg)
	if err != nil {
		return err
	}

	t1 := .2

	out1 := image.NewGray(img.Bounds())

	for x := 0; x < img.Width(); x += 3 {
		for y := 0; y < img.Height(); y += 3 {

			_, avgDistance := img.AverageColorAndStats(image.Point{x + 1, y + 1}, 1)

			clr := color.Gray{0}
			if avgDistance > t1 {
				clr = color.Gray{255}
			}

			for a := 0; a < 3; a++ {
				for b := 0; b < 3; b++ {
					xx := x + a
					yy := y + b
					out1.SetGray(xx, yy, clr)
				}
			}

		}
	}

	d.GotDebugImage(out1, "t")

	if true {
		out, err := ShapeWalkEntireDebug(img, false)
		if err != nil {
			return err
		}
		d.GotDebugImage(out, "entire")
	}

	if true {
		starts := []image.Point{}

		for _, s := range cfg.Shapes {
			starts = append(starts, s.Start)
		}

		out, err := ShapeWalkMultiple(img, starts, false)
		if err != nil {
			return err
		}

		d.GotDebugImage(out, "shapes")

		for idx, s := range cfg.Shapes {
			numPixels := out.PixelsInSegmemnt(idx + 1)
			if numPixels < s.PixelRange[0] || numPixels > s.PixelRange[1] {
				// run again with debugging on
				_, err := ShapeWalkMultiple(img, starts, true)
				if err != nil {
					return err
				}

				d.T.Errorf("out of pixel range %v %d", s, numPixels)
			}
		}

	}

	return nil
}

func TestChunk1(t *testing.T) {
	d := vision.NewMultipleImageTestDebugger(t, "segmentation/test1", "*.png")
	err := d.Process(&chunkImageDebug{})
	if err != nil {
		t.Fatal(err)
	}

}
