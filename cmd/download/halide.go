package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"runtime"
)

func main() {
	switch runtime.GOOS {
	case "darwin":
		dir := dirnameDarwin()
		hiddenDir := "." + dir
		if exists(hiddenDir) {
			return
		}

		downloadDarwin()
		mustRename(dir, hiddenDir)
		mustSymlink(hiddenDir, ".halide")
	case "linux":
		dir := dirnameLinux()
		hiddenDir := "." + dir
		if exists(hiddenDir) {
			return
		}

		downloadLinux()
		mustRename(dir, hiddenDir)
		mustSymlink(hiddenDir, ".halide")
	default:
		panic("does not generate OS type: " + runtime.GOOS)
	}
	println("complete")
}

func exists(p string) bool {
	s, err := os.Stat(p)
	if err != nil {
		return false
	}
	if s.IsDir() != true {
		return false
	}
	return true
}

func mustRename(prev, next string) {
	if err := os.Rename(prev, next); err != nil {
		panic(err)
	}
}

func mustSymlink(src, name string) {
	if err := os.Symlink(src, name); err != nil {
		panic(err)
	}
}

func mustDownload(url string) {
	println("download Halide...")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	println("extract tar gz...")
	if err := tarxzf(resp.Body); err != nil {
		panic(err)
	}
}

func tarxzf(r io.Reader) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()

	t := tar.NewReader(gz)
	for {
		h, err := t.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if err := tarf(h, t); err != nil {
			return err
		}
	}
	return nil
}

func tarf(h *tar.Header, t *tar.Reader) error {
	switch h.Typeflag {
	case tar.TypeDir:
		if err := os.MkdirAll(h.Name, os.FileMode(h.Mode)); err != nil {
			return err
		}
	case tar.TypeReg:
		f, err := os.Create(h.Name)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := f.Chmod(os.FileMode(h.Mode)); err != nil {
			return err
		}
		if _, err := io.Copy(f, t); err != nil {
			return err
		}
	}
	return nil
}
