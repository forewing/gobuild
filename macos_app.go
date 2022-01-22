package gobuild

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/jackmordaunt/icns/v2"
)

const (
	infoTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleName</key>
	<string>%v</string>
	<key>CFBundleExecutable</key>
	<string>%v</string>
	<key>CFBundleIdentifier</key>
	<string>%v</string>
	<key>CFBundleIconFile</key>
	<string>icon.icns</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0.0</string>
	<key>CFBundleSupportedPlatforms</key>
	<array>
		<string>MacOSX</string>
	</array>
	<key>CFBundleVersion</key>
	<string>1</string>
	<key>NSHighResolutionCapable</key>
	<true/>
	<key>NSSupportsAutomaticGraphicsSwitching</key>
	<true/>
	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>LSApplicationCategoryType</key>
	<string>public.app-category.</string>
	<key>LSMinimumSystemVersion</key>
	<string>10.11</string>
</dict>
</plist>`
)

func BuildMacOSApp(output, name, exe, id, icon string, convertIcon bool) error {
	appContents := filepath.Join(output, name+".app", "Contents")
	appResources := filepath.Join(appContents, "Resources")
	appMacOS := filepath.Join(appContents, "MacOS")
	os.MkdirAll(appResources, os.ModePerm)
	os.MkdirAll(appMacOS, os.ModePerm)

	exeName := filepath.Base(exe)

	err := os.WriteFile(filepath.Join(appContents, "Info.plist"),
		[]byte(fmt.Sprintf(infoTemplate, name, exeName, id)), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Rename(exe, filepath.Join(appMacOS, exeName))
	if err != nil {
		return err
	}

	targetIcon := filepath.Join(appResources, "icon.icns")
	if convertIcon {
		err = convertIcns(icon, targetIcon)
	} else {
		err = os.Rename(icon, targetIcon)
	}
	if err != nil {
		return err
	}

	return nil
}

func convertIcns(src, dst string) error {
	imgFp, err := os.Open(src)
	if err != nil {
		return err
	}
	defer imgFp.Close()

	img, _, err := image.Decode(imgFp)
	if err != nil {
		return err
	}

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	if err := icns.Encode(output, img); err != nil {
		return err
	}

	return nil
}
