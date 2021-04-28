package main

import (
	"fmt"

	"github.com/forewing/gobuild"
)

var (
	target = gobuild.Target{
		Source: "./payload",
		OutputName: fmt.Sprintf("purego-%s-%s-%s",
			gobuild.PlaceholderVersion,
			gobuild.PlaceholderOS,
			gobuild.PlaceholderArch),
		OutputPath: "./output",

		ExtraLdFlags: "-s -w",

		VersionPath: "main.Version",
		HashPath:    "main.Hash",

		Compress: gobuild.CompressAuto,
		Platforms: []gobuild.Platform{
			{OS: gobuild.OSWindows, Arch: gobuild.ArchAmd64},
			{OS: gobuild.OSWindows, Arch: gobuild.Arch386},
			{OS: gobuild.OSLinux, Arch: gobuild.ArchAmd64},
			{OS: gobuild.OSLinux, Arch: gobuild.ArchArm, GoArm: "5"},
		},
	}
)

func main() {
	fmt.Println(target.Build())
}
