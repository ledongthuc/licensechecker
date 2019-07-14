package data

import (
	"encoding/json"

	"github.com/ledongthuc/licensechecker/data/toc"
	"github.com/pkg/errors"
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

	// Load Main license
	raw, err := toc.Asset(listLicenses)
	if err != nil {
		return result, errors.Wrap(err, "Error when load license info")
	}

	var mainLicense license
	err = json.Unmarshal(raw, &mainLicense)
	if err != nil {
		return result, errors.Wrap(err, "Error when parsing license info")
	}
	loadLicenses(result, mainLicense)

	// Load Exception license
	var exceptionLicense exception
	raw, err = toc.Asset(listExceptions)
	if err != nil {
		return result, errors.Wrap(err, "Error when load main license info")
	}
	err = json.Unmarshal(raw, &exceptionLicense)
	if err != nil {
		return result, errors.Wrap(err, "Error when parsing license info")
	}
	loadExceptionLicenses(result, exceptionLicense)

	return result, nil
}

func loadLicenses(container map[string]LicenseInfo, ls license) error {
	if container == nil {
		return errors.New("Container shouldn't nil")
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

func loadExceptionLicenses(container map[string]LicenseInfo, ls exception) error {
	if container == nil {
		return errors.New("Container shouldn't nil")
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
