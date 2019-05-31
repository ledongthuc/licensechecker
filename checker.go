package licensechecker

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/parnurzeal/gorequest"
)

type Result struct {
	Source      string `json:"source"`
	Name        string `json:"name"`
	LicenseName string `json:"license_name"`
	LicenseURL  string `json:"license_url"`
}

type Results []Result

func FromURL(url string) (Results, error) {
	for _, config := range DefaultConfigs {
		var validPattern = regexp.MustCompile(config.Pattern)
		if valid := validPattern.MatchString(url); !valid {
			continue
		}

		for _, licenseURL := range config.LicenseURLs {
			// Change it to ReplaceAll on 1.12
			licenseURL = ReplaceLisenceURL(config, licenseURL, url)
			resp, body, errs := gorequest.New().Get(licenseURL).End()
			if len(errs) > 0 {
				continue
			}
			if resp.StatusCode != http.StatusOK {
				continue
			}

			fmt.Printf(OUTPUT, config.GetName(url), GetLicenseName(string(body)), licenseURL)
			goto NextURL
		}
	}
	fmt.Printf(OUTPUT, url, "FAIL TO CHECK", "")
	return
}
