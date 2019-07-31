package licensechecker

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"

	"github.com/ledongthuc/licensechecker/internal/data"
	"github.com/ledongthuc/licensechecker/internal/toc"
	"github.com/pkg/errors"
)

const (
	listLicenses   = "licenses.json"
	listExceptions = "exceptions.json"
)

var (
	ErrorUninitiatedContainer = errors.New("Can't compose licenses into uninitiated container")
)

// license defines structure of standard license from spdx format
type license struct {
	LicenseListVersion string `json:"licenseListVersion"`
	Licenses           []struct {
		Reference             string   `json:"reference"`
		IsDeprecatedLicenseID bool     `json:"isDeprecatedLicenseId"`
		DetailsURL            string   `json:"detailsUrl"`
		ReferenceNumber       string   `json:"referenceNumber"`
		Name                  string   `json:"name"`
		LicenseID             string   `json:"licenseId"`
		SeeAlso               []string `json:"seeAlso"`
		IsOsiApproved         bool     `json:"isOsiApproved"`
		IsFsfLibre            bool     `json:"isFsfLibre,omitempty"`
	} `json:"licenses"`
	ReleaseDate string `json:"releaseDate"`
}

// contains will check a standard license info is valid and exist
func (e license) contains(info LicenseInfo) bool {
	for _, exception := range e.Licenses {
		if info.LicenseID == exception.LicenseID &&
			info.Name == exception.Name &&
			info.IsDeprecated == exception.IsDeprecatedLicenseID &&
			reflect.DeepEqual(info.References, exception.SeeAlso) {
			return true
		}
	}
	return false
}

// loadStandardLicenses loads all standard licenses from asset resouce
func loadStandardLicenses() (license, error) {
	raw, err := toc.Asset(listLicenses)
	if err != nil {
		return license{}, errors.Wrap(err, "Error when load license info")
	}

	var standardLicenses license
	err = json.Unmarshal(raw, &standardLicenses)
	if err != nil {
		return license{}, errors.Wrap(err, "Error when parsing license info")
	}
	return standardLicenses, nil
}

// exception defines structure of exception license from spdx format
type exception struct {
	LicenseListVersion string `json:"licenseListVersion"`
	ReleaseDate        string `json:"releaseDate"`
	Exceptions         []struct {
		Reference             string   `json:"reference"`
		IsDeprecatedLicenseID bool     `json:"isDeprecatedLicenseId"`
		DetailsURL            string   `json:"detailsUrl"`
		ReferenceNumber       string   `json:"referenceNumber"`
		Name                  string   `json:"name"`
		SeeAlso               []string `json:"seeAlso"`
		LicenseExceptionID    string   `json:"licenseExceptionId"`
	} `json:"exceptions"`
}

// contains will check a exception license info is valid and exist
func (e exception) contains(info LicenseInfo) bool {
	for _, exception := range e.Exceptions {
		if info.LicenseID == exception.LicenseExceptionID &&
			info.Name == exception.Name &&
			info.IsDeprecated == exception.IsDeprecatedLicenseID &&
			reflect.DeepEqual(info.References, exception.SeeAlso) {
			return true
		}
	}
	return false
}

// loadExceptionLicenses loads all standard licenses from asset resouce
func loadExceptionLicenses() (exception, error) {
	var exceptionLicense exception
	raw, err := toc.Asset(listExceptions)
	if err != nil {
		return exception{}, errors.Wrap(err, "Error when load main license info")
	}

	err = json.Unmarshal(raw, &exceptionLicense)
	if err != nil {
		return exception{}, errors.Wrap(err, "Error when parsing license info")
	}
	return exceptionLicense, nil
}

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

// GetLicenseInfo will get all license's info. It's map with key-par value to easier data query.
func GetLicenseInfo() (map[string]LicenseInfo, error) {
	result := make(map[string]LicenseInfo)

	// standard license
	standardLicenses, err := loadStandardLicenses()
	if err != nil {
		return result, err
	}
	err = convertStandardLicenses(result, standardLicenses)
	if err != nil {
		return result, err
	}

	// exception license
	exceptionLicenses, err := loadExceptionLicenses()
	if err != nil {
		return result, err
	}
	err = convertExceptionLicenses(result, exceptionLicenses)
	if err != nil {
		return result, err
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
