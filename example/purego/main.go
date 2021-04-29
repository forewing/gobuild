package main

import (
	"fmt"

	"github.com/forewing/gobuild"
)

var (
	target = gobuild.Target{
		Source: "../payload",
		OutputName: fmt.Sprintf("purego-%s-%s-%s",
			gobuild.PlaceholderVersion,
			gobuild.PlaceholderOS,
			gobuild.PlaceholderArch),
		OutputPath:  "./output",
		CleanOutput: true,

		ExtraLdFlags: "-s -w",

		VersionPath: "main.Version",
		HashPath:    "main.Hash",

		Compress:  gobuild.CompressAuto,
		Platforms: gobuild.PlatformCommon,
	}
)

func main() {
	err := target.Build()
	if err != nil {
		panic(err)
	}
}
