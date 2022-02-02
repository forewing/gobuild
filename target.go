package gobuild

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v4"
)

// Target describe build configs
type Target struct {
	// Go executeable path.
	//
	// If empty, `go` will be used.
	Go string

	// Source is directory of source main package.
	Source string

	// OutputName is the base name of output.
	//
	// It may contains `Placeholder`s, which will be replaced when building.
	OutputName string

	// OutputPath is the directory where all outputs will be stored.
	OutputPath string

	// If CleanOutput is true, all contents of OutputPath will be removed before building.
	CleanOutput bool

	// If Cgo is true, enabled cgo, and Platform.CC is used as C compiler.
	//
	// If false, cgo is disabled.
	Cgo bool

	// ExtraFlags to be passed to go compiler.
	ExtraFlags []string

	// ExtraLdFlags to be passed to loader.
	//
	// Example: to behave like `go build -ldflags "-s -w"`, this field
	// should be set to "-s -w".
	ExtraLdFlags string

	// VersionPath is the path of a variable of your source package
	// where you want it to be set to output of `git describe --tags` when building.
	//
	// Example: you have a variable `Version` in package `main`,
	// set VersionPath to `main.Version`, and Version will be set to your git tag.
	VersionPath string

	// HashPath is the path of a variable where you want it to be the current git hash.
	HashPath string

	// Compress set the compress methods.
	Compress CompressType

	// Platforms is the target platforms.
	Platforms []Platform

	// Envs set extra environment variables
	Envs map[string]string

	temp    string
	ldflags string
	output  string
}

// Placeholder will be replaced when building.
type Placeholder string

const (
	// PlaceholderVersion will be replaced by output of `git describe --tags` on success
	PlaceholderVersion Placeholder = "{Version}"

	// PlaceholderArch will be replaced by GOARCH.
	PlaceholderArch Placeholder = "{Arch}"

	// PlaceholderOS will be replaced by GOOS.
	PlaceholderOS Placeholder = "{OS}"

	tempDirPattern = "go-build*"

	envCgo    = "CGO_ENABLED"
	envCC     = "CC"
	envGoOS   = "GOOS"
	envGoArch = "GOARCH"
	envGoArm  = "GOARM"

	defaultGoExec = "go"
)

var (
	OutputMode = 0755

	defaultEnvs = map[string]string{
		envCgo: "0",
	}
)

// Build the target
func (t *Target) Build() error {
	if t.CleanOutput {
		if err := os.RemoveAll(t.OutputPath); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(t.OutputPath, os.ModePerm); err != nil {
		return err
	}

	if err := t.init(); err != nil {
		return err
	}
	defer os.RemoveAll(t.temp)

	for i := range t.Platforms {
		bin, err := t.build(i)
		if err != nil {
			return err
		}
		err = t.pack(i, bin)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Target) init() error {
	t.output = t.OutputName

	realSource, err := evalSymlinksAbs(t.Source)
	if err != nil {
		return err
	}
	t.Source = realSource

	realOutput, err := evalSymlinksAbs(t.OutputPath)
	if err != nil {
		return err
	}
	t.OutputPath = realOutput

	temp, err := os.MkdirTemp("", tempDirPattern)
	if err != nil {
		return err
	}
	t.temp = temp

	ldflags := []string{t.ExtraLdFlags}
	if v, err := GetGitVersion(t.Source); err == nil && len(v) > 0 {
		if len(t.VersionPath) > 0 {
			ldflags = append(ldflags,
				fmt.Sprintf("-X '%s=%s'", t.VersionPath, v),
			)
		}
		t.output = strings.ReplaceAll(t.output, string(PlaceholderVersion), v)
	}
	if len(t.HashPath) > 0 {
		if h, err := GetGitHash(t.Source); err == nil && len(h) > 0 {
			ldflags = append(ldflags,
				fmt.Sprintf("-X '%s=%s'", t.HashPath, h),
			)
		}
	}
	t.ldflags = strings.Join(ldflags, " ")

	return nil
}

func (t *Target) build(id int) (string, error) {
	p := t.Platforms[id]

	if p.Arch == ArchUniversal {
		return t.buildUniversal(id)
	}

	output := filepath.Join(t.temp, fmt.Sprintf("output-%v", id))
	args := []string{
		"build",
		"-ldflags",
		t.ldflags,
		"-o",
		output,
	}
	args = append(args, t.ExtraFlags...)
	args = append(args, t.Source)

	goexec := defaultGoExec
	if len(t.Go) > 0 {
		goexec = t.Go
	}
	cmd := exec.Command(goexec, args...)
	cmd.Dir = t.Source

	envs := make(map[string]string)
	for k, v := range defaultEnvs {
		envs[k] = v
	}

	if t.Cgo {
		envs[envCgo] = "1"
	}
	if len(p.CC) > 0 {
		envs[envCC] = p.CC
	}
	envs[envGoOS] = string(p.OS)
	envs[envGoArch] = string(p.Arch)
	envs[envGoArm] = string(p.GoArm)

	for k := range t.Envs {
		envs[k] = t.Envs[k]
	}
	for k := range p.Envs {
		envs[k] = p.Envs[k]
	}

	cmd.Env = append(cmd.Env, os.Environ()...)
	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	combined, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("err: %v, output: %v", err, string(combined))
	}

	os.Chmod(output, fs.FileMode(OutputMode))
	return output, nil
}

func (t *Target) pack(id int, input string) error {
	p := t.Platforms[id]

	name := t.output

	forceReplace := false
	if len(t.Platforms) > 1 {
		forceReplace = true
	}

	osReplace := runtime.GOOS
	if len(p.OS) > 0 {
		if forceReplace && !strings.Contains(name, string(PlaceholderOS)) {
			name = fmt.Sprintf("%s-%s", name, p.OS)
		}
		osReplace = string(p.OS)
	} else if runtime.GOOS == string(OSWindows) {
		p.OS = OSWindows
	}
	name = strings.ReplaceAll(name, string(PlaceholderOS), osReplace)

	archReplace := runtime.GOARCH
	if len(p.Arch) > 0 {
		arch := string(p.Arch) + string(p.GoArm)
		if forceReplace && !strings.Contains(name, string(PlaceholderArch)) {
			name = fmt.Sprintf("%s-%s", name, arch)
		}
		archReplace = arch
	}
	name = strings.ReplaceAll(name, string(PlaceholderArch), archReplace)

	binary := name
	if p.OS == OSWindows {
		binary += ".exe"
	}

	files := map[string]string{input: binary}

	switch t.Compress {
	case CompressTarGz:
		return Compress(filepath.Join(t.OutputPath, name+".tar.gz"),
			files, archiver.CompressedArchive{
				Compression: archiver.Gz{},
				Archival:    archiver.Tar{},
			})
	case CompressZip:
		return Compress(filepath.Join(t.OutputPath, name+".zip"), files, archiver.Zip{})
	case CompressRaw:
		fallthrough
	default:
		return moveWithoutCompress(filepath.Join(t.OutputPath, binary), input)
	}
}

func (t *Target) buildUniversal(id int) (string, error) {
	p := t.Platforms[id]

	if p.OS != OSDarwin {
		return "", fmt.Errorf("%v does not support universal arch", p.OS)
	}

	t2 := *t
	t2.temp = filepath.Join(t.temp, "u")
	os.MkdirAll(t2.temp, os.ModePerm)

	t2.Platforms = []Platform{p, p}
	t2.Platforms[0].Arch = ArchAmd64
	t2.Platforms[1].Arch = ArchArm64

	o0, err := t2.build(0)
	if err != nil {
		return "", err
	}
	o1, err := t2.build(1)
	if err != nil {
		return "", err
	}

	output := filepath.Join(t.temp, fmt.Sprintf("output-%v", id))
	lipo := exec.Command("lipo", "-create", "-output", output, o0, o1)
	err = lipo.Run()
	if err != nil {
		return "", err
	}

	return output, nil
}
