package licensechecker

func Detect(licenseContent []byte) (LicenseInfo, error) {
	return LicenseInfo{}, nil
}

func DetectFromPath(localPath string) (LicenseInfo, error) {
	return LicenseInfo{}, nil
}

func DetectFromURL(URL string) (LicenseInfo, error) {
	return LicenseInfo{}, nil
}

func Add(licenseContent []byte, pathOfSource string) error {
	return nil
}

func AddWithOption(licenseContent []byte, pathOfSource string, options string) error {
	return nil
}
