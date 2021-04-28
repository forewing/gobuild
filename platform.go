package gobuild

type Platform struct {
	Arch  PlatformArch
	OS    PlatformOS
	GoArm PlatformGoArm

	CC string
}

type PlatformArch string
type PlatformOS string
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
