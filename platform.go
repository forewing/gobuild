package gobuild

// Platform set the target platform
type Platform struct {
	// Arch sets GOARCH.
	Arch PlatformArch

	// OS sets GOOS.
	OS PlatformOS

	// GoArm sets GOARM when Arch is `arm`.
	GoArm PlatformGoArm

	// CC sets C compiler when cgo is enabled.
	CC string
}

// PlatformArch is all available GOARCH.
type PlatformArch string

// PlatformOS is all available GOOS.
type PlatformOS string

// PlatformGoArm is all available GOARM.
type PlatformGoArm string

const (
	Arch386         PlatformArch = "386"
	ArchAmd64       PlatformArch = "amd64"
	ArchAmd64p32    PlatformArch = "amd64p32"
	ArchArm         PlatformArch = "arm"
	ArchArmbe       PlatformArch = "armbe"
	ArchArm64       PlatformArch = "arm64"
	ArchArm64be     PlatformArch = "arm64be"
	ArchPpc64       PlatformArch = "ppc64"
	ArchPpc64le     PlatformArch = "ppc64le"
	ArchMips        PlatformArch = "mips"
	ArchMipsle      PlatformArch = "mipsle"
	ArchMips64      PlatformArch = "mips64"
	ArchMips64le    PlatformArch = "mips64le"
	ArchMips64p32   PlatformArch = "mips64p32"
	ArchMips64p32le PlatformArch = "mips64p32le"
	ArchPpc         PlatformArch = "ppc"
	ArchRiscv       PlatformArch = "riscv"
	ArchRiscv64     PlatformArch = "riscv64"
	ArchS390        PlatformArch = "s390"
	ArchS390x       PlatformArch = "s390x"
	ArchSparc       PlatformArch = "sparc"
	ArchSparc64     PlatformArch = "sparc64"
	ArchWasm        PlatformArch = "wasm"
)

const (
	OSAix       PlatformOS = "aix"
	OSAndroid   PlatformOS = "android"
	OSDarwin    PlatformOS = "darwin"
	OSDragonfly PlatformOS = "dragonfly"
	OSFreebsd   PlatformOS = "freebsd"
	OSHurd      PlatformOS = "hurd"
	OSIllumos   PlatformOS = "illumos"
	OSIos       PlatformOS = "ios"
	OSJs        PlatformOS = "js"
	OSLinux     PlatformOS = "linux"
	OSNacl      PlatformOS = "nacl"
	OSNetbsd    PlatformOS = "netbsd"
	OSOpenbsd   PlatformOS = "openbsd"
	OSPlan9     PlatformOS = "plan9"
	OSSolaris   PlatformOS = "solaris"
	OSWindows   PlatformOS = "windows"
	OSZos       PlatformOS = "zos"
)

const (
	GoArm5 PlatformGoArm = "5"
	GoArm6 PlatformGoArm = "6"
	GoArm7 PlatformGoArm = "7"
)

var (
	PlatformWindows386   = Platform{OS: OSWindows, Arch: Arch386}
	PlatformWindowsAmd64 = Platform{OS: OSWindows, Arch: ArchAmd64}
	PlatformWindowsArm5  = Platform{OS: OSWindows, Arch: ArchArm, GoArm: GoArm5}
	PlatformWindowsArm6  = Platform{OS: OSWindows, Arch: ArchArm, GoArm: GoArm6}
	PlatformWindowsArm7  = Platform{OS: OSWindows, Arch: ArchArm, GoArm: GoArm7}

	PlatformLinux386   = Platform{OS: OSLinux, Arch: Arch386}
	PlatformLinuxAmd64 = Platform{OS: OSLinux, Arch: ArchAmd64}
	PlatformLinuxArm64 = Platform{OS: OSLinux, Arch: ArchArm64}
	PlatformLinuxArm5  = Platform{OS: OSLinux, Arch: ArchArm, GoArm: GoArm5}
	PlatformLinuxArm6  = Platform{OS: OSLinux, Arch: ArchArm, GoArm: GoArm6}
	PlatformLinuxArm7  = Platform{OS: OSLinux, Arch: ArchArm, GoArm: GoArm7}

	PlatformDarwinAmd64 = Platform{OS: OSDarwin, Arch: ArchAmd64}
	PlatformDarwinArm64 = Platform{OS: OSDarwin, Arch: ArchArm64}
)

var (
	// PlatformCommon includes the most used (~99.9%) platforms.
	//
	// Including OS: Windows, Linux, Darwin(MacOS).
	// Including Arch: 386, amd64, arm32(v5, v6, v7), arm64.
	PlatformCommon = []Platform{
		PlatformWindows386,
		PlatformWindowsAmd64,
		PlatformWindowsArm5,
		PlatformWindowsArm6,
		PlatformWindowsArm7,

		PlatformLinux386,
		PlatformLinuxAmd64,
		PlatformLinuxArm64,
		PlatformLinuxArm5,
		PlatformLinuxArm6,
		PlatformLinuxArm7,

		PlatformDarwinAmd64,
		PlatformDarwinArm64,
	}

	// PlatformNative contains the config for native platform
	PlatformNative = []Platform{
		{},
	}
)
