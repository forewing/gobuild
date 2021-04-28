package gobuild

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
	defaultEnvs = map[string]string{
		envCgo:    "0",
		envCC:     "gcc",
		envGoOS:   "",
		envGoArch: "",
		envGoArm:  "",
	}
)

// Build the target
func (t *Target) Build() error {
	if t.CleanOutput {
		if err := cleanDirectory(t.OutputPath); err != nil {
			return nil
		}
	}
	if err := os.MkdirAll(t.OutputPath, os.ModePerm); err != nil {
		return err
	}

	if err := t.init(); err != nil {
		return err
	}
	defer cleanDirectory(t.temp)

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

	absSource, err := filepath.Abs(t.Source)
	if err != nil {
		return err
	}
	t.Source = absSource

	absOutput, err := filepath.Abs(t.OutputPath)
	if err != nil {
		return err
	}
	t.OutputPath = absOutput

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

	cmd.Env = append(cmd.Env, os.Environ()...)
	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	combined, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("err: %v, output: %v", err, string(combined))
	}

	return output, nil
}

func (t *Target) pack(id int, input string) error {
	p := t.Platforms[id]

	name := t.output

	if len(p.OS) > 0 {
		if !strings.Contains(name, string(PlaceholderOS)) {
			name = fmt.Sprintf("%s-%s", name, p.OS)
		}
		name = strings.ReplaceAll(name, string(PlaceholderOS), string(p.OS))
	} else if runtime.GOOS == string(OSWindows) {
		p.OS = OSWindows
	}

	if len(p.Arch) > 0 {
		arch := string(p.Arch) + string(p.GoArm)
		if !strings.Contains(name, string(PlaceholderArch)) {
			name = fmt.Sprintf("%s-%s", name, arch)
		}
		name = strings.ReplaceAll(name, string(PlaceholderArch), arch)
	}

	binary := name
	if p.OS == OSWindows {
		binary += ".exe"
	}

	outputTarGz := filepath.Join(t.OutputPath, name+".tar.gz")
	outputZip := filepath.Join(t.OutputPath, name+".zip")
	outputRaw := filepath.Join(t.OutputPath, binary)

	switch t.Compress {
	case CompressAllTarGz:
		return compressTarGz(outputTarGz, input, binary)
	case CompressAllZip:
		return compressZip(outputZip, input, binary)
	case CompressAuto:
		if p.OS == OSWindows {
			return compressZip(outputZip, input, binary)
		}
		return compressTarGz(outputTarGz, input, binary)
	case CompressRaw:
		return compressRaw(outputRaw, input)
	default:
		return compressRaw(outputRaw, input)
	}
}

func cleanDirectory(path string) error {
	return os.RemoveAll(path)
}
