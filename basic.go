package licensechecker

import (
	"regexp"
	"strings"

	"github.com/ledongthuc/licensechecker/internal/data"
	"github.com/pkg/errors"
)

const (
	listLicenses   = "licenses.json"
	listExceptions = "exceptions.json"
)

var (
	ErrorUninitiatedContainer = errors.New("Can't compose licenses into uninitiated container")
)

// LicenseInfo contains meta data of license like ID, nice name, reference urls/resources or license is deprecated
type LicenseInfo struct {
	LicenseID    string
	Name         string
	References   []string
	IsDeprecated bool
}

// LicenseDataPath compose the path of data assets
func (l LicenseInfo) LicenseDataPath() string {
	if l.LicenseID == "" {
		return ""
	}

	// Special cases
	if l.LicenseID == "Nokia-Qt-exception-1.1" {
		return "Nokia-Qt-exception-1.1.txt"
	}

	var pathBuilder strings.Builder
	if l.IsDeprecated {
		pathBuilder.WriteString("deprecated_")
	}
	pathBuilder.WriteString(l.LicenseID)
	pathBuilder.WriteString(".txt")
	return pathBuilder.String()
}

// AllInfo will get all license's info. It's map with key-par value to easier data query.
func AllInfo() ([]LicenseInfo, error) {
	// Load licenses
	standardLicenses, err := loadStandardLicenses()
	if err != nil {
		return []LicenseInfo{}, err
	}
	exceptionLicenses, err := loadExceptionLicenses()
	if err != nil {
		return []LicenseInfo{}, err
	}

	// Merge license
	merging := make(map[string]LicenseInfo)
	err = convertStandardLicenses(merging, standardLicenses)
	if err != nil {
		return []LicenseInfo{}, err
	}

	err = convertExceptionLicenses(merging, exceptionLicenses)
	if err != nil {
		return []LicenseInfo{}, err
	}

	result := make([]LicenseInfo, 0, len(merging))
	for _, m := range merging {
		result = append(result, m)
	}

	return result, nil
}

// convertStandardLicenses will load standard licenses into a map.
func convertStandardLicenses(container map[string]LicenseInfo, ls license) error {
	if container == nil {
		return ErrorUninitiatedContainer
	}
	for _, l := range ls.Licenses {
		container[l.LicenseID] = LicenseInfo{
			LicenseID:    l.LicenseID,
			Name:         l.Name,
			References:   l.SeeAlso,
			IsDeprecated: l.IsDeprecatedLicenseID,
		}
	}
	return nil
}

// convertExceptionLicenses will load all exception licenses into a map.
func convertExceptionLicenses(container map[string]LicenseInfo, ls exception) error {
	if container == nil {
		return ErrorUninitiatedContainer
	}
	for _, l := range ls.Exceptions {
		container[l.LicenseExceptionID] = LicenseInfo{
			LicenseID:    l.LicenseExceptionID,
			Name:         l.Name,
			References:   l.SeeAlso,
			IsDeprecated: l.IsDeprecatedLicenseID,
		}
	}
	return nil
}

// LicenseInfo contains meta data of license like ID, nice name, reference urls/resources or license is deprecated
type LicenseData struct {
	LicenseID  string
	Content    []byte
	RawContent []byte
}

// LoadLicenseData will load license content base on their info. It will take care to check license from spdx or custom source
func LoadLicenseData(licenseInfo LicenseInfo) (LicenseData, error) {
	raw, err := data.Asset(licenseInfo.LicenseDataPath())
	if err != nil {
		return LicenseData{}, errors.Wrap(err, "Error to load data from assets '"+licenseInfo.LicenseDataPath()+"'")
	}

	result := LicenseData{
		LicenseID: licenseInfo.LicenseID,
	}
	result.Content = raw
	result.RawContent = regexp.MustCompile(`\r?\n`).ReplaceAll(raw, []byte(" "))
	return result, nil
}

type License struct {
	LicenseInfo
	LicenseData
}

func All() ([]License, error) {
	info, err := AllInfo()
	if err != nil {
		return []License{}, err
	}

	result := make([]License, 0, len(info))
	for _, infoItem := range info {
		data, err := LoadLicenseData(infoItem)
		if err != nil {
			return []License{}, err
		}
		result = append(result, License{
			LicenseInfo: infoItem,
			LicenseData: data,
		})
	}

	return result, nil
}
