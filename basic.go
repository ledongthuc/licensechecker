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

// License contains meta data and real license content. All you need here
type License struct {
	LicenseInfo
	LicenseContent
}

// LicenseInfo contains meta data of license. It doesn't contain license content, just info.
type LicenseInfo struct {
	LicenseID    string
	Name         string
	References   []string
	IsDeprecated bool
}

// LicenseContent contains license's content
type LicenseContent struct {
	LicenseID  string
	Content    []byte
	RawContent []byte
}

// All will loads and returns all license that have with their content and me meta data.
func All() ([]License, error) {
	info, err := AllInfo()
	if err != nil {
		return []License{}, err
	}

	result := make([]License, 0, len(info))
	for _, infoItem := range info {
		data, err := infoItem.LoadLicenseContent()
		if err != nil {
			return []License{}, err
		}
		result = append(result, License{
			LicenseInfo:    infoItem,
			LicenseContent: data,
		})
	}

	return result, nil
}

// AllInfo will get all license's information that doesn't contains license's content, just light weight meta data.
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

// TODO: change to receiver method
// LoadLicenseContent will load license content base on their info. It will take care to check license from spdx or custom source
func (licenseInfo LicenseInfo) LoadLicenseContent() (LicenseContent, error) {
	raw, err := data.Asset(licenseInfo.LicenseContentPath())
	if err != nil {
		return LicenseContent{}, errors.Wrap(err, "Error to load data from assets '"+licenseInfo.LicenseContentPath()+"'")
	}

	result := LicenseContent{
		LicenseID: licenseInfo.LicenseID,
	}
	result.Content = raw
	result.RawContent = regexp.MustCompile(`\r?\n`).ReplaceAll(raw, []byte(" "))
	return result, nil
}

// LicenseContentPath compose the path of data assets
func (l LicenseInfo) LicenseContentPath() string {
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
