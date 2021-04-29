package gobuild

import (
	"encoding/json"
	"fmt"
	"os"
)

// TargetConfig is used to parse configs.
//
// Some shortcuts are added here, so you will not
// have to write many redundant data.
type TargetConfig struct {
	// Target is the real target.
	Target

	// PlatformShortcut is shortcut for platform settings
	PlatformShortcut PlatformShortcut
}

// GetTargetFromJson parse json file as TargetConfig.
//
// Modifies set by Shortcuts will be applied to TargetConfig.Target.
func GetTargetFromJson(path string) (*Target, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result TargetConfig
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	target := result.Target
	if len(result.PlatformShortcut) > 0 {
		v, ok := platformShorcutMap[result.PlatformShortcut]
		if !ok {
			return nil, fmt.Errorf("invalid PlatformShortcut: %s", result.PlatformShortcut)
		}
		target.Platforms = append(target.Platforms, v...)
	}

	fmt.Printf("%+v", target)
	return &target, nil
}
