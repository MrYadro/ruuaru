package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"github.com/generaltso/vibrant"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func getImageColors(fileName string) (fgColorVibrant, bgColorVibrant, fgColorMuted, bgColorMuted string) {

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	palette, err := vibrant.NewPaletteFromImage(img)
	if err != nil {
		panic(err)
	}

	for name, swatch := range palette.ExtractAwesome() {
		if name == "Vibrant" {
			fgColorVibrant = swatch.Color.TitleTextColor().RGBHex()
			bgColorVibrant = swatch.Color.RGBHex()
		}
		if name == "Muted" {
			fgColorMuted = swatch.Color.TitleTextColor().RGBHex()
			bgColorMuted = swatch.Color.RGBHex()
		}
	}
	return fgColorVibrant, bgColorVibrant, fgColorMuted, bgColorMuted
}

func resizeImageToFill(mw *imagick.MagickWand, w, h uint) {
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	wOut := width
	hOut := height

	if width < height {
		if width != w {
			mult := float64(w) / float64(width)
			wOut = uint(float64(width) * mult)
			hOut = uint(float64(height) * mult)
		}
	} else {
		if height != h {
			mult := float64(h) / float64(height)
			wOut = uint(float64(width) * mult)
			hOut = uint(float64(height) * mult)
		}
	}

	err := mw.ResizeImage(wOut, hOut, imagick.FILTER_CUBIC)
	if err != nil {
		panic(err)
	}
}

func resizeImageToFit(mw *imagick.MagickWand, w, h uint) {
	height := mw.GetImageHeight()
	width := mw.GetImageWidth()
	var wOut uint
	var hOut uint

	if width > height {
		if width != w {
			mult := float64(w) / float64(width)
			wOut = uint(float64(width) * mult)
			hOut = uint(float64(height) * mult)
		}
	} else {
		if height != h {
			mult := float64(h) / float64(height)
			wOut = uint(float64(width) * mult)
			hOut = uint(float64(height) * mult)
		}
	}

	err := mw.ResizeImage(wOut, hOut, imagick.FILTER_CUBIC)
	if err != nil {
		panic(err)
	}
}

func makeRatingStar(color string) *imagick.MagickWand {
	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	pw := imagick.NewPixelWand()
	defer pw.Destroy()

	pw.SetColor(color)
	err := mw.NewImage(512, 512, pw)
	if err != nil {
		panic(err)
	}

	starShapeImage := imagick.NewMagickWand()
	err = starShapeImage.ReadImageBlob(starBin)
	if err != nil {
		panic(err)
	}

	err = mw.CompositeImage(starShapeImage, imagick.COMPOSITE_OP_COPY_ALPHA, true, 0, 0)
	if err != nil {
		panic(err)
	}

	starShapeImage.Destroy()

	return mw
}

func makeRatingStars(rating int, colorStars, colorBack string, width, height uint) *imagick.MagickWand {
	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	pw := imagick.NewPixelWand()
	defer pw.Destroy()

	pw.SetColor("none")
	err := mw.NewImage(228, 90, pw)
	if err != nil {
		panic(err)
	}

	err = dw.SetFont("TT-Commons-Bold")
	if err != nil {
		panic(err)
	}
	dw.SetFontSize(80)

	starFull := makeRatingStar(colorStars)

	err = starFull.ResizeImage(60, 60, imagick.FILTER_CUBIC)
	if err != nil {
		panic(err)
	}

	err = dw.SetFont("TT-Commons-Bold")
	if err != nil {
		panic(err)
	}
	dw.SetFontSize(90)

	pw.SetColor(colorBack)
	dw.SetFillColor(pw)
	dw.RoundRectangle(0, 0, 228, 90, 10, 10)

	pw.SetColor(colorStars)
	dw.SetFillColor(pw)
	dw.SetTextAlignment(imagick.ALIGN_RIGHT)
	dw.Annotation(104, 73, strconv.Itoa(rating))

	err = mw.DrawImage(dw)
	if err != nil {
		panic(err)
	}

	err = mw.CompositeImage(starFull, imagick.COMPOSITE_OP_OVER, true, 124, 15)
	if err != nil {
		panic(err)
	}

	starFull.Destroy()

	return mw
}

func makeStoryText(fgColor, bgColor, text string) *imagick.MagickWand {
	mw := imagick.NewMagickWand()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	pw := imagick.NewPixelWand()
	defer pw.Destroy()

	pw.SetColor("none")
	err := mw.NewImage(1080, 103, pw)
	if err != nil {
		panic(err)
	}

	err = dw.SetFont("TT-Commons-Bold")
	if err != nil {
		panic(err)
	}
	dw.SetFontSize(80)

	textMetrics := mw.QueryFontMetrics(dw, text)

	textWidth := textMetrics.TextWidth

	imageWidth := textWidth + 40

	err = mw.ResizeImage(uint(imageWidth), 103, imagick.FILTER_CUBIC)
	if err != nil {
		panic(err)
	}

	pw.SetColor(bgColor)
	dw.SetFillColor(pw)

	dw.RoundRectangle(0, 0, imageWidth, 103, 10, 10)

	pw.SetColor(fgColor)
	dw.SetFillColor(pw)
	dw.SetTextAlignment(imagick.ALIGN_CENTER)
	dw.Annotation(imageWidth/2, 93-17, text)

	err = mw.DrawImage(dw)
	if err != nil {
		panic(err)
	}

	return mw
}

func makeStory(title, posterPath, backgroundPath, reviewLine1, reviewLine2, reviewLine3 string, rating int) string {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	pw := imagick.NewPixelWand()
	defer pw.Destroy()

	fontColor, backgroundColor, darkFontColor, darkBackgroundColor := getImageColors(posterPath)

	backgroundImage := imagick.NewMagickWand()
	err := backgroundImage.ReadImage(backgroundPath)
	if err != nil {
		panic(err)
	}

	posterImage := imagick.NewMagickWand()
	err = posterImage.ReadImage(posterPath)
	if err != nil {
		panic(err)
	}

	pw.SetColor(backgroundColor)
	err = mw.NewImage(storyWidth, storyHeight, pw)
	if err != nil {
		panic(err)
	}

	resizeImageToFill(backgroundImage, storyWidth, storyHeight)
	resizeImageToFit(posterImage, posterMaxWidth, posterMaxHeight)

	backgroundImageWidth := backgroundImage.GetImageWidth()

	err = backgroundImage.CropImage(storyWidth, storyHeight, int((backgroundImageWidth-storyWidth)/2), 0)
	if err != nil {
		panic(err)
	}

	err = backgroundImage.SetImageColorspace(imagick.COLORSPACE_GRAY)
	if err != nil {
		panic(err)
	}

	err = mw.CompositeImage(backgroundImage, imagick.COMPOSITE_OP_COPY, true, 0, 0)
	if err != nil {
		panic(err)
	}

	backgroundImage.Destroy()

	pw.SetColor("#000")

	pwo := imagick.NewPixelWand()
	pwo.SetColor("rgb(50%,50%,50%)")

	err = mw.ColorizeImage(pw, pwo)
	if err != nil {
		panic(err)
	}

	posterWidth := posterImage.GetImageWidth()

	posterShapeImage := imagick.NewMagickWand()
	err = posterShapeImage.ReadImageBlob(posterBin)
	if err != nil {
		panic(err)
	}

	err = posterImage.CompositeImage(posterShapeImage, imagick.COMPOSITE_OP_COPY_ALPHA, true, 0, 0)
	if err != nil {
		panic(err)
	}

	posterShapeImage.Destroy()

	err = mw.CompositeImage(posterImage, imagick.COMPOSITE_OP_OVER, true, int(storyWidth-posterWidth)/2, 300)
	if err != nil {
		panic(err)
	}

	posterImage.Destroy()

	starRating := makeRatingStars(rating, darkFontColor, darkBackgroundColor, 400, 600)

	err = mw.CompositeImage(starRating, imagick.COMPOSITE_OP_OVER, true, 655, 253)
	if err != nil {
		panic(err)
	}

	starRating.Destroy()

	movieTitle := makeStoryText(fontColor, backgroundColor, title)

	if movieTitle.GetImageWidth() > 1000 {
		resizeImageToFit(movieTitle, 1000, 1000)
	}

	err = mw.CompositeImage(movieTitle, imagick.COMPOSITE_OP_OVER, true, int(storyWidth-movieTitle.GetImageWidth())/2, 1312)
	if err != nil {
		panic(err)
	}

	movieTitle.Destroy()

	if len(reviewLine1) != 0 {
		err = dw.SetFont("TT-Commons-Bold-Italic")
		if err != nil {
			panic(err)
		}
		dw.SetFontSize(70)
		pw.SetColor("#fff")
		dw.SetFillColor(pw)
		dw.SetTextAlignment(imagick.ALIGN_CENTER)
		dw.Annotation(1080/2, 1560, reviewLine1)

		textMetrics := mw.QueryFontMetrics(dw, reviewLine1)

		pw.SetColor("#9a9a9a")
		dw.SetFillColor(pw)
		dw.SetTextAlignment(imagick.ALIGN_RIGHT)
		dw.Annotation((1080-textMetrics.TextWidth)/2, 1560, "“")
		if len(reviewLine2) != 0 {
			pw.SetColor("#fff")
			dw.SetFillColor(pw)
			dw.SetTextAlignment(imagick.ALIGN_CENTER)
			dw.Annotation(1080/2, 1660, reviewLine2)

			if len(reviewLine3) != 0 {
				pw.SetColor("#fff")
				dw.SetFillColor(pw)
				dw.SetTextAlignment(imagick.ALIGN_CENTER)
				dw.Annotation(1080/2, 1760, reviewLine3)

				textMetrics := mw.QueryFontMetrics(dw, reviewLine3)
				pw.SetColor("#9a9a9a")
				dw.SetFillColor(pw)
				dw.SetTextAlignment(imagick.ALIGN_LEFT)
				dw.Annotation((1080+textMetrics.TextWidth)/2, 1760, "”")
			} else {
				textMetrics := mw.QueryFontMetrics(dw, reviewLine2)
				pw.SetColor("#9a9a9a")
				dw.SetFillColor(pw)
				dw.SetTextAlignment(imagick.ALIGN_LEFT)
				dw.Annotation((1080+textMetrics.TextWidth)/2, 1660, "”")
			}
		} else {
			textMetrics := mw.QueryFontMetrics(dw, reviewLine1)
			pw.SetColor("#9a9a9a")
			dw.SetFillColor(pw)
			dw.SetTextAlignment(imagick.ALIGN_LEFT)
			dw.Annotation((1080+textMetrics.TextWidth)/2, 1560, "”")
		}
	}

	err = mw.DrawImage(dw)
	if err != nil {
		panic(err)
	}

	filename := RandStringBytes(5)

	err = mw.WriteImage("tmp/" + filename + ".png")
	if err != nil {
		panic(err)
	}
	return filename + ".png"
}
