# Go Build Tools

[![Go Report Card](https://goreportcard.com/badge/github.com/forewing/gobuild?style=flat-square)](https://goreportcard.com/report/github.com/forewing/gobuild)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/forewing/gobuild?style=flat-square)](https://github.com/forewing/gobuild/releases/latest)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/forewing/gobuild)](https://pkg.go.dev/github.com/forewing/gobuild)

Tools for building and distributing Go executables

## Examples

### Directly Use

[example/purego](./example/purego)

```golang
target := gobuild.Target{
    Source: "./example/payload",

    OutputName: fmt.Sprintf("payload-%s-%s-%s",
        gobuild.PlaceholderVersion,
        gobuild.PlaceholderOS,
        gobuild.PlaceholderArch),

    OutputPath:  "./output",
    CleanOutput: true,

    ExtraLdFlags: "-s -w",

    VersionPath: "main.Version",
    HashPath:    "main.Hash",

    Compress:  gobuild.CompressZip,
    Platforms: gobuild.PlatformCommon,
}

err := target.Build()
if err != nil {
    panic(err)
}
```

### Config File

[example/config](./example/config)

```json
{
    "Go": "",
    "Source": "./example/payload",
    "OutputName": "payload-{Version}-{OS}-{Arch}",
    "OutputPath": "./output",
    "CleanOutput": true,
    "Cgo": false,
    "ExtraFlags": null,
    "ExtraLdFlags": "-s -w",
    "VersionPath": "main.Version",
    "HashPath": "main.Hash",
    "Compress": "auto",
    "PlatformShortcut": "common",
    "Platforms": [
        {
            "Arch": "riscv",
            "OS": "linux"
        }
    ]
}
```

```golang
target, err := gobuild.GetTargetFromJson("./config.json")
if err != nil {
    panic(err)
}
err = target.Build()
if err != nil {
    panic(err)
}
```

### CGO

```golang
target = gobuild.Target{
    Platforms: []gobuild.Platform{
        {OS: gobuild.OSLinux, Arch: gobuild.ArchAmd64, CC: "gcc"},
        {OS: gobuild.OSWindows, Arch: gobuild.ArchAmd64, CC: "x86_64-w64-mingw32-gcc"},
    },
}
```
