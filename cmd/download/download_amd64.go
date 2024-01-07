package main

const (
	darwinDir string = "Halide-14.0.0-x86-64-osx"
	linuxDir  string = "Halide-14.0.0-x86-64-osx"
	darwinUrl string = "https://github.com/halide/Halide/releases/download/v14.0.0/Halide-14.0.0-x86-64-osx-6b9ed2afd1d6d0badf04986602c943e287d44e46.tar.gz"
	linuxUrl  string = "https://github.com/halide/Halide/releases/download/v14.0.0/Halide-14.0.0-x86-64-linux-6b9ed2afd1d6d0badf04986602c943e287d44e46.tar.gz"
	//darwinDir string = "Halide-16.0.0-x86-64-osx"
	//linuxDir  string = "Halide-16.0.0-x86-64-linux"
	//darwinUrl string = "https://github.com/halide/Halide/releases/download/v16.0.0/Halide-16.0.0-x86-64-osx-1e963ff817ef0968cc25d811a25a7350c8953ee6.tar.gz"
	//linuxUrl  string = "https://github.com/halide/Halide/releases/download/v16.0.0/Halide-16.0.0-x86-64-linux-1e963ff817ef0968cc25d811a25a7350c8953ee6.tar.gz"
)

func dirnameDarwin() string {
	return darwinDir
}

func dirnameLinux() string {
	return linuxDir
}

func downloadDarwin() {
	mustDownload(darwinUrl)
}

func downloadLinux() {
	mustDownload(linuxUrl)
}
