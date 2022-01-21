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
		ExtraFlags:   []string{"-trimpath"},

		Compress: gobuild.CompressRaw,
		Platforms: []gobuild.Platform{
			{OS: gobuild.OSDarwin, Arch: gobuild.ArchUniversal},
		},
	}
)

func main() {
	err := target.Build()
	if err != nil {
		panic(err)
	}
}
