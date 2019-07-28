package licensechecker

import "github.com/ledongthuc/licensechecker/data"

func Detect(licenseContent []byte) (data.LicenseInfo, error) {
	return data.LicenseInfo{}
}

func DetectFromPath(localPath string) (data.LicenseInfo, error) {
	return data.LicenseInfo{}
}

func DetectFromURL(URL string) (data.LicenseInfo, error) {
	return data.LicenseInfo{}
}

func Add(licenseContent []byte, pathOfSource string) error {
}

func AddWithOption(licenseContent []byte, pathOfSource string, options string) error {
}
