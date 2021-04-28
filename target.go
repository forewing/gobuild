package gobuild

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Target struct {
	Go string

	Source      string
	OutputName  string
	OutputPath  string
	CleanOutput bool

	Cgo bool

	ExtraFlags   []string
	ExtraLdFlags string

	VersionPath string
	HashPath    string

	Compress  CompressType
	Platforms []Platform

	temp    string
	ldflags string
	output  string
}

const (
	PlaceholderVersion = "{Version}"
	PlaceholderArch    = "{Arch}"
	PlaceholderOS      = "{OS}"

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
		t.output = strings.ReplaceAll(t.output, PlaceholderVersion, v)
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
		if !strings.Contains(name, PlaceholderOS) {
			name = fmt.Sprintf("%s-%s", name, p.OS)
		}
		name = strings.ReplaceAll(name, PlaceholderOS, string(p.OS))
	}

	if len(p.Arch) > 0 {
		arch := string(p.Arch) + string(p.GoArm)
		if !strings.Contains(name, PlaceholderArch) {
			name = fmt.Sprintf("%s-%s", name, arch)
		}
		name = strings.ReplaceAll(name, PlaceholderArch, arch)
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
