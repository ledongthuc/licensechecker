package licensechecker

import (
	"encoding/json"
	"reflect"

	"github.com/ledongthuc/licensechecker/internal/toc"
	"github.com/pkg/errors"
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
