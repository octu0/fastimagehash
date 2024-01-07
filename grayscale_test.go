package fastimagehash

func ExampleGrayscale() {
	rgba, err := pngToRGBA(testPngImg)
	if err != nil {
		panic(err)
	}

	gray, err := grayscale(rgba)
	if err != nil {
		panic(err)
	}

	path, err := saveImage(gray)
	if err != nil {
		panic(err)
	}
	println("grayscale = ", path)

	// Output:
}
