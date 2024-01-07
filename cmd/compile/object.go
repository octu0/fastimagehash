package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const featuresX86 string = `
  features.push_back(Target::AVX);
  features.push_back(Target::AVX2);
  features.push_back(Target::FMA);
  features.push_back(Target::FMA);
  features.push_back(Target::F16C);
  features.push_back(Target::SSE41);
`

const featuresARM string = `
  features.push_back(Target::ARMFp16);
  features.push_back(Target::ARMv81a);
`

const generateRuntimeMainTmpl string = `
#include <Halide.h>
using namespace Halide;
int main() {
  Var x("x");
  Func fn = Func("noop");
  fn(x) = 0;
  std::vector<Argument> args;

  std::vector<Target::Feature> features;
  {{ .Features }}
  {
    Target target;
    target.os = Target::OSX;
    target.arch = {{ .TargetArch }};
    target.bits = 64;
    target.set_features(features);
    fn.compile_to_object("{{ .FileNameDarwin }}", args, "{{ .Name }}", target);
  }
  {
    Target target;
    target.os = Target::Linux;
    target.arch = {{ .TargetArch }};
    target.bits = 64;
    target.set_features(features);
    fn.compile_to_object("{{ .FileNameLinux }}", args, "{{ .Name }}", target);
  }
  fn.compile_to_header("{{ .HeaderName }}", args, "{{ .Name }}");
  return 0;
}
`

const generateGenRunMainTmpl string = `
#include <Halide.h>
#include "{{ .HppFileName }}"
using namespace Halide;
int main() {
  std::vector<Target::Feature> features;
  {{ .Features }}
  features.push_back(Target::Feature::NoRuntime);

  std::tuple<{{ .ExportType }}, std::vector<Argument>> tuple = export_{{ .Name }}();
  {{ .ExportType }} fn = std::get<0>(tuple);
  std::vector<Argument> args =std::get<1>(tuple);
  {
    Target target;
    target.os = Target::OSX;
    target.arch = {{ .TargetArch }};
    target.bits = 64;
    target.set_features(features);
    fn.compile_to_object("{{ .FileNameDarwin }}", args, "{{ .Name }}", target);
  }
  {
    Target target;
    target.os = Target::Linux;
    target.arch = {{ .TargetArch }};
    target.bits = 64;
    target.set_features(features);
    fn.compile_to_object("{{ .FileNameLinux }}", args, "{{ .Name }}", target);
  }
  fn.compile_to_header("{{ .HeaderName }}", args, "{{ .Name }}");
  return 0;
}
`

type GenRun struct {
	FileNameDarwin    string
	FileNameLinux     string
	AsmNameDarwin     string
	AsmNameLinux      string
	LLVMAsmNameDarwin string
	LLVMAsmNameLinux  string
	LLVMBcNameDarwin  string
	LLVMBcNameLinux   string
	HeaderName        string
	Name              string
	HppFileName       string
	ExecFileName      string
	MainTemplate      string
	TargetArch        string
	Features          string
	ExportType        string
}

func main() {
	halidePath := ".halide"

	exportTypeName := os.Args[1]
	funcName := os.Args[2]
	targetFilePath := os.Args[3]
	targetFileBase := filepath.Base(targetFilePath)
	targetFileExt := filepath.Ext(targetFileBase)
	baseName := targetFileBase[0:strings.LastIndex(targetFileBase, targetFileExt)]

	exportType := ""
	switch exportTypeName {
	case "f", "func":
		exportType = "Func"
	case "p", "pipeline":
		exportType = "Pipeline"
	}
	if exportType == "" {
		panic("not support extern type: " + exportTypeName)
	}

	targets := make([]GenRun, 0, 4)
	targets = append(targets, GenRun{
		FileNameDarwin:   fmt.Sprintf("lib/amd64/lib%s_darwin.dylib", "runtime"),
		FileNameLinux:    fmt.Sprintf("lib/amd64/lib%s_linux.o", "runtime"),
		AsmNameDarwin:    fmt.Sprintf("lib/amd64/lib%s_darwin.s", "runtime"),
		AsmNameLinux:     fmt.Sprintf("lib/amd64/lib%s_linux.s", "runtime"),
		LLVMBcNameDarwin: fmt.Sprintf("lib/amd64/lib%s_darwin.bc", "runtime"),
		LLVMBcNameLinux:  fmt.Sprintf("lib/amd64/lib%s_linux.bc", "runtime"),
		HeaderName:       fmt.Sprintf("include/amd64/%s.h", "runtime"),
		Name:             "runtime",
		HppFileName:      "",
		ExecFileName:     fmt.Sprintf("gen/%s_amd64.out", "runtime"),
		MainTemplate:     generateRuntimeMainTmpl,
		TargetArch:       "Target::X86",
		Features:         featuresX86,
		ExportType:       "",
	})
	targets = append(targets, GenRun{
		FileNameDarwin:   fmt.Sprintf("lib/arm64/lib%s_darwin.dylib", "runtime"),
		FileNameLinux:    fmt.Sprintf("lib/arm64/lib%s_linux.o", "runtime"),
		AsmNameDarwin:    fmt.Sprintf("lib/arm64/lib%s_darwin.s", "runtime"),
		AsmNameLinux:     fmt.Sprintf("lib/arm64/lib%s_linux.s", "runtime"),
		LLVMBcNameDarwin: fmt.Sprintf("lib/arm64/lib%s_darwin.bc", "runtime"),
		LLVMBcNameLinux:  fmt.Sprintf("lib/arm64/lib%s_linux.bc", "runtime"),
		HeaderName:       fmt.Sprintf("include/arm64/%s.h", "runtime"),
		Name:             "runtime",
		HppFileName:      "",
		ExecFileName:     fmt.Sprintf("gen/%s_arm64.out", "runtime"),
		MainTemplate:     generateRuntimeMainTmpl,
		TargetArch:       "Target::ARM",
		Features:         featuresARM,
		ExportType:       "",
	})
	targets = append(targets, GenRun{
		FileNameDarwin:   fmt.Sprintf("lib/amd64/lib%s_darwin.dylib", funcName),
		FileNameLinux:    fmt.Sprintf("lib/amd64/lib%s_linux.o", funcName),
		AsmNameDarwin:    fmt.Sprintf("lib/amd64/lib%s_darwin.s", funcName),
		AsmNameLinux:     fmt.Sprintf("lib/amd64/lib%s_linux.s", funcName),
		LLVMBcNameDarwin: fmt.Sprintf("lib/amd64/lib%s_darwin.bc", funcName),
		LLVMBcNameLinux:  fmt.Sprintf("lib/amd64/lib%s_linux.bc", funcName),
		HeaderName:       fmt.Sprintf("include/amd64/%s.h", funcName),
		Name:             funcName,
		HppFileName:      fmt.Sprintf("%s.hpp", baseName),
		ExecFileName:     fmt.Sprintf("gen/%s_amd64.out", funcName),
		MainTemplate:     generateGenRunMainTmpl,
		TargetArch:       "Target::X86",
		Features:         featuresX86,
		ExportType:       exportType,
	})
	targets = append(targets, GenRun{
		FileNameDarwin:   fmt.Sprintf("lib/arm64/lib%s_darwin.dylib", funcName),
		FileNameLinux:    fmt.Sprintf("lib/arm64/lib%s_linux.o", funcName),
		AsmNameDarwin:    fmt.Sprintf("lib/arm64/lib%s_darwin.s", funcName),
		AsmNameLinux:     fmt.Sprintf("lib/arm64/lib%s_linux.s", funcName),
		LLVMBcNameDarwin: fmt.Sprintf("lib/arm64/lib%s_darwin.bc", funcName),
		LLVMBcNameLinux:  fmt.Sprintf("lib/arm64/lib%s_linux.bc", funcName),
		HeaderName:       fmt.Sprintf("include/arm64/%s.h", funcName),
		Name:             funcName,
		HppFileName:      fmt.Sprintf("%s.hpp", baseName),
		ExecFileName:     fmt.Sprintf("gen/%s_arm64.out", funcName),
		MainTemplate:     generateGenRunMainTmpl,
		TargetArch:       "Target::ARM",
		Features:         featuresARM,
		ExportType:       exportType,
	})

	libpng := exec.Command("libpng-config", "--cflags", "--ldflags")
	libpngCfg, err := libpng.Output()
	if err != nil {
		panic(err)
	}
	libpngFlags := strings.TrimSpace(string(libpngCfg))
	libpngFlags = strings.ReplaceAll(libpngFlags, "\n", " ")

	buf := bytes.NewBuffer(nil)
	files := make([]*os.File, 0, len(targets))
	defer func() {
		for _, f := range files {
			os.Remove(f.Name())
		}
	}()

	targetFileStat, err := os.Stat(targetFilePath)
	if err != nil {
		panic(err)
	}
	for _, t := range targets {
		execStat, err := os.Stat(t.ExecFileName)
		if err == nil { // file exists
			if t.Name == "runtime" {
				continue
			}
			sourceMTime := targetFileStat.ModTime()
			if sourceMTime == execStat.ModTime() {
				continue // source code no change
			}
		}
		buf.Reset()

		println("compiling...", t.Name, t.TargetArch)
		tpl, err := template.New(t.Name).Parse(t.MainTemplate)
		if err != nil {
			panic(err)
		}
		if err := tpl.Execute(buf, t); err != nil {
			panic(err)
		}
		mainC, err := os.CreateTemp("", "main*.cpp")
		if err != nil {
			panic(err)
		}
		files = append(files, mainC)
		if _, err := mainC.Write(buf.Bytes()); err != nil {
			panic(err)
		}

		genArgs := []string{
			"clang++",
			"-g",
			"-I.",
			"-I" + halidePath + "/include",
			"-I" + halidePath + "/share/Halide/tools",
			"-L" + halidePath + "/lib",
			libpngFlags,
			"-L/usr/local/opt/jpeg/lib",
			"-I/usr/local/opt/jpeg/include",
			"-ljpeg",
			"-lHalide",
			"-lpthread",
			"-ldl",
			"-lz",
			"-std=c++17",
			"-o", t.ExecFileName,
		}
		if t.Name == "runtime" {
			genArgs = append(genArgs, mainC.Name())
		} else {
			genArgs = append(genArgs, targetFilePath)
			genArgs = append(genArgs, mainC.Name())
		}

		println("compile...", t.Name, t.TargetArch, "cmd:", strings.Join(genArgs, " "))

		cmd := exec.Command("sh", "-c", strings.Join(genArgs, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			println("error", err.Error())
			continue
		}

		println("generate...", t.Name, t.TargetArch)
		gen := exec.Command("sh", "-c", t.ExecFileName)
		gen.Stdout = os.Stdout
		gen.Stderr = os.Stderr
		if err := gen.Run(); err != nil {
			println("error", err.Error())
			continue
		}
		println("done...", t.Name, t.TargetArch)

		if err := os.Chtimes(t.ExecFileName, time.Time{}, targetFileStat.ModTime()); err != nil {
			println("error", err.Error())
			continue
		}
	}
}
