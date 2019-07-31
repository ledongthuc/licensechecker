package licensechecker

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/parnurzeal/gorequest"

	"github.com/ledongthuc/licensechecker/template"
)

const OUTPUT = "%q,%q,%q\n"

type URLType int

const (
	URLGithub URLType = iota
	URLUnknow
)

type Config struct {
	Pattern     string
	URLType     URLType
	LicenseURLs []string
}

var (
	DefaultConfigs = []Config{
		Config{
			Pattern: "github.com", // will be replace to raw.githubusercontent.com to get license file
			URLType: URLGithub,
			LicenseURLs: []string{
				"{{URL}}/master/LICENSE",
				"{{URL}}/master/LICENSE.md",
				"{{URL}}/master/LICENSE.txt",
				"{{URL}}/master/License",
				"{{URL}}/master/License.md",
				"{{URL}}/master/License.txt",
				"{{url}}/master/license",
				"{{url}}/master/license.md",
				"{{url}}/master/license.txt",
				"{{URL}}/master/README",
				"{{URL}}/master/README.md",
				"{{URL}}/master/README.txt",
				"{{URL}}/master/Readme",
				"{{URL}}/master/Readme.md",
				"{{URL}}/master/Readme.txt",
				"{{url}}/master/readme",
				"{{url}}/master/readme.md",
				"{{url}}/master/readme.txt",
			},
		},
		Config{
			Pattern: "githubusercontent.com",
			URLType: URLGithub,
			LicenseURLs: []string{
				"{{URL}}/master/LICENSE",
				"{{URL}}/master/LICENSE.md",
				"{{URL}}/master/LICENSE.txt",
				"{{URL}}/master/License",
				"{{URL}}/master/License.md",
				"{{URL}}/master/License.txt",
				"{{url}}/master/license",
				"{{url}}/master/license.md",
				"{{url}}/master/license.txt",
				"{{URL}}/master/README",
				"{{URL}}/master/README.md",
				"{{URL}}/master/README.txt",
				"{{URL}}/master/Readme",
				"{{URL}}/master/Readme.md",
				"{{URL}}/master/Readme.txt",
				"{{url}}/master/readme",
				"{{url}}/master/readme.md",
				"{{url}}/master/readme.txt",
			},
		},
		Config{
			Pattern: "golang.org", // will be replace to raw.githubusercontent.com to get license file
			URLType: URLGithub,
			LicenseURLs: []string{
				"{{URL}}/master/LICENSE",
				"{{URL}}/master/LICENSE.md",
				"{{URL}}/master/LICENSE.txt",
				"{{URL}}/master/License",
				"{{URL}}/master/License.md",
				"{{URL}}/master/License.txt",
				"{{url}}/master/license",
				"{{url}}/master/license.md",
				"{{url}}/master/license.txt",
				"{{URL}}/master/README",
				"{{URL}}/master/README.md",
				"{{URL}}/master/README.txt",
				"{{URL}}/master/Readme",
				"{{URL}}/master/Readme.md",
				"{{URL}}/master/Readme.txt",
				"{{url}}/master/readme",
				"{{url}}/master/readme.md",
				"{{url}}/master/readme.txt",
			},
		},
		Config{
			Pattern: "",
			URLType: URLUnknow,
			LicenseURLs: []string{
				"{{URL}}",
			},
		},
	}
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Need to input lib path")
		return
	}

	urls := os.Args[1:]
	duplicating := make(map[string]struct{})

	for _, url := range urls {
		_, ok := duplicating[url]
		if ok {
			continue
		}
		duplicating[url] = struct{}{}

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

	NextURL:
	}
	return
}

func GetLicenseName(license string) string {
	if strings.Contains(strings.ToLower(string(license)), strings.ToLower("MIT license")) {
		return "MIT license"
	}
	if strings.Contains(strings.ToLower(string(license)), strings.ToLower("Mozilla Public License")) {
		return "Mozilla Public License"
	}
	if strings.Contains(strings.ToLower(string(license)), strings.ToLower("Apache License")) && strings.Contains(strings.ToLower(string(license)), strings.ToLower("Version 2.0")) {
		return "Apache License 2.0"
	}

	name, distance := template.FindMatchName(license)
	if distance < 500 {
		return name
	}

	re := regexp.MustCompile(`(?i)Copyright(.+)`)
	return re.FindString(license)
}

func ReplaceLisenceURL(config Config, licenseURL, url string) string {
	licenseURL = strings.Replace(licenseURL, "{{URL}}", url, -1)
	if config.URLType == URLGithub {
		licenseURL = strings.Replace(licenseURL, "github.com", "raw.githubusercontent.com", -1)
		licenseURL = strings.Replace(licenseURL, "golang.org", "raw.githubusercontent.com", -1)
	}
	if !strings.HasPrefix("http", licenseURL) {
		licenseURL = fmt.Sprintf("http://%s", licenseURL)
	}
	return licenseURL
}

func (config Config) GetName(url string) string {
	if config.URLType == URLGithub {
		return config.GetGithubLibraryName(url)
	}
	return url
}

func (config Config) GetGithubLibraryName(url string) string {
	re := regexp.MustCompile(`(github.com|githubusercontent.com)/(.+)/(.+)`)
	parts := re.FindAllStringSubmatch(url, -1)
	if len(parts) == 0 {
		return url
	}
	if len(parts[0]) < 2 {
		return url
	}
	if len(parts[0]) == 2 {
		return parts[0][1]
	}
	return parts[0][len(parts[0])-1]
}
