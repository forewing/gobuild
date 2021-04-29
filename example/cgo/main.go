package main

import (
	"fmt"

	"github.com/forewing/gobuild"
)

var (
	target = gobuild.Target{
		Cgo: true,

		Source: "../payload_cgo",
		OutputName: fmt.Sprintf("cgo-%s-%s-%s",
			gobuild.PlaceholderVersion,
			gobuild.PlaceholderOS,
			gobuild.PlaceholderArch),

		OutputPath:  "./output",
		CleanOutput: true,

		ExtraLdFlags: "-s -w",

		VersionPath: "main.Version",
		HashPath:    "main.Hash",

		Compress: gobuild.CompressRaw,
		Platforms: []gobuild.Platform{
			{OS: gobuild.OSLinux, Arch: gobuild.ArchAmd64, CC: "gcc"},
			{OS: gobuild.OSWindows, Arch: gobuild.ArchAmd64, CC: "x86_64-w64-mingw32-gcc"},
		},
	}
)

func main() {
	err := target.Build()
	if err != nil {
		panic(err)
	}
}
