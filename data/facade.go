package data

import (
	"encoding/json"
	"reflect"

	"github.com/ledongthuc/licensechecker/data/toc"
	"github.com/pkg/errors"
)

var (
	ErrorUninitiatedContainer = errors.New("Can't compose licenses into uninitiated container")
)

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

func (e license) contains(info LicenseInfo) bool {
	for _, exception := range e.Licenses {
		if info.ID == exception.LicenseID &&
			info.Name == exception.Name &&
			info.IsDeprecated == exception.IsDeprecatedLicenseID &&
			reflect.DeepEqual(info.References, exception.SeeAlso) {
			return true
		}
	}
	return false
}

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

func (e exception) contains(info LicenseInfo) bool {
	for _, exception := range e.Exceptions {
		if info.ID == exception.LicenseExceptionID &&
			info.Name == exception.Name &&
			info.IsDeprecated == exception.IsDeprecatedLicenseID &&
			reflect.DeepEqual(info.References, exception.SeeAlso) {
			return true
		}
	}
	return false
}

// LicenseInfo contains meta data of license like ID, nice name, reference urls/resources or license is deprecated
type LicenseInfo struct {
	ID           string
	Name         string
	References   []string
	IsDeprecated bool
}

const (
	listLicenses   = "licenses.json"
	listExceptions = "exceptions.json"
)

// GetLicenseInfo will get all license's info. It's map with key-par value to easier data query.
func GetLicenseInfo() (map[string]LicenseInfo, error) {
	result := make(map[string]LicenseInfo)

	// standard license
	standardLicenses, err := LoadStandardLicenses()
	if err != nil {
		return result, err
	}
	err = convertStandardLicenses(result, standardLicenses)
	if err != nil {
		return result, err
	}

	// exception license
	exceptionLicenses, err := LoadExceptionLicenses()
	if err != nil {
		return result, err
	}
	err = convertExceptionLicenses(result, exceptionLicenses)
	if err != nil {
		return result, err
	}

	return result, nil
}

// LoadStandardLicenses loads all standard licenses from asset resouce
func LoadStandardLicenses() (license, error) {
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

// convertStandardLicenses will load standard licenses into a map.
func convertStandardLicenses(container map[string]LicenseInfo, ls license) error {
	if container == nil {
		return ErrorUninitiatedContainer
	}
	for _, l := range ls.Licenses {
		container[l.LicenseID] = LicenseInfo{
			ID:           l.LicenseID,
			Name:         l.Name,
			References:   l.SeeAlso,
			IsDeprecated: l.IsDeprecatedLicenseID,
		}
	}
	return nil
}

// LoadExceptionLicenses loads all standard licenses from asset resouce
func LoadExceptionLicenses() (exception, error) {
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

// convertExceptionLicenses will load all exception licenses into a map.
func convertExceptionLicenses(container map[string]LicenseInfo, ls exception) error {
	if container == nil {
		return ErrorUninitiatedContainer
	}
	for _, l := range ls.Exceptions {
		container[l.LicenseExceptionID] = LicenseInfo{
			ID:           l.LicenseExceptionID,
			Name:         l.Name,
			References:   l.SeeAlso,
			IsDeprecated: l.IsDeprecatedLicenseID,
		}
	}
	return nil
}
