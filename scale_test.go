package fastimagehash

func ExampleScaleNormal() {
	rgba, err := pngToRGBA(testPngImg)
	if err != nil {
		panic(err)
	}

	scaled, err := scaleNormal(rgba, 32, 32)
	if err != nil {
		panic(err)
	}

	path, err := saveImage(scaled)
	if err != nil {
		panic(err)
	}
	println("scale_normal = ", path)

	// Output:
}

func ExampleScaleBox() {
	rgba, err := pngToRGBA(testPngImg)
	if err != nil {
		panic(err)
	}

	scaled, err := scaleBox(rgba, 32, 32)
	if err != nil {
		panic(err)
	}

	path, err := saveImage(scaled)
	if err != nil {
		panic(err)
	}
	println("scale_box = ", path)

	// Output:
}

func ExampleScaleLinear() {
	rgba, err := pngToRGBA(testPngImg)
	if err != nil {
		panic(err)
	}

	scaled, err := scaleLinear(rgba, 32, 32)
	if err != nil {
		panic(err)
	}

	path, err := saveImage(scaled)
	if err != nil {
		panic(err)
	}
	println("scale_linear = ", path)

	// Output:
}

func ExampleScaleGauss() {
	rgba, err := pngToRGBA(testPngImg)
	if err != nil {
		panic(err)
	}

	scaled, err := scaleGauss(rgba, 32, 32)
	if err != nil {
		panic(err)
	}

	path, err := saveImage(scaled)
	if err != nil {
		panic(err)
	}
	println("scale_gauss = ", path)

	// Output:
}
