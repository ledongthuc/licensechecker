package licensechecker

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/ledongthuc/licensechecker/internal/data"
)

const exampleLicenseRaw = `{
  "licenseListVersion": "3.6",
  "licenses": [
    {
      "reference": "./0BSD.html",
      "isDeprecatedLicenseId": false,
      "detailsUrl": "http://spdx.org/licenses/0BSD.json",
      "referenceNumber": "319",
      "name": "BSD Zero Clause License",
      "licenseId": "0BSD",
      "seeAlso": [
        "http://landley.net/toybox/license.html"
      ],
      "isOsiApproved": true
    },
    {
      "reference": "./AAL.html",
      "isDeprecatedLicenseId": false,
      "detailsUrl": "http://spdx.org/licenses/AAL.json",
      "referenceNumber": "21",
      "name": "Attribution Assurance License",
      "licenseId": "AAL",
      "seeAlso": [
        "https://opensource.org/licenses/attribution"
      ],
      "isOsiApproved": true
    },
	{
      "reference": "./AGPL-3.0.html",
      "isDeprecatedLicenseId": true,
      "isFsfLibre": true,
      "detailsUrl": "http://spdx.org/licenses/AGPL-3.0.json",
      "referenceNumber": "229",
      "name": "GNU Affero General Public License v3.0",
      "licenseId": "AGPL-3.0",
      "seeAlso": [
        "https://www.gnu.org/licenses/agpl.txt",
        "https://opensource.org/licenses/AGPL-3.0"
      ],
      "isOsiApproved": true
    }
  ],
  "releaseDate": "2019-07-10"
}`

const exampleExceptionLicenseRaw = `{
  "licenseListVersion": "3.6",
  "releaseDate": "2019-07-10",
  "exceptions": [
    {
      "reference": "./Libtool-exception.html",
      "isDeprecatedLicenseId": false,
      "detailsUrl": "http://spdx.org/licenses/Libtool-exception.json",
      "referenceNumber": "1",
      "name": "Libtool Exception",
      "seeAlso": [
        "http://git.savannah.gnu.org/cgit/libtool.git/tree/m4/libtool.m4"
      ],
      "licenseExceptionId": "Libtool-exception"
    },
    {
      "reference": "./Classpath-exception-2.0.html",
      "isDeprecatedLicenseId": false,
      "detailsUrl": "http://spdx.org/licenses/Classpath-exception-2.0.json",
      "referenceNumber": "13",
      "name": "Classpath exception 2.0",
      "seeAlso": [
        "http://www.gnu.org/software/classpath/license.html",
        "https://fedoraproject.org/wiki/Licensing/GPL_Classpath_Exception"
      ],
      "licenseExceptionId": "Classpath-exception-2.0"
    },
    {
      "reference": "./Nokia-Qt-exception-1.1.html",
      "isDeprecatedLicenseId": true,
      "detailsUrl": "http://spdx.org/licenses/Nokia-Qt-exception-1.1.json",
      "referenceNumber": "23",
      "name": "Nokia Qt LGPL exception 1.1",
      "seeAlso": [
        "https://www.keepassx.org/dev/projects/keepassx/repository/revisions/b8dfb9cc4d5133e0f09cd7533d15a4f1c19a40f2/entry/LICENSE.NOKIA-LGPL-EXCEPTION"
      ],
      "licenseExceptionId": "Nokia-Qt-exception-1.1"
    }
  ]
}`

func TestAllInfo(t *testing.T) {
	m, err := AllInfo()
	if err != nil {
		t.Errorf("AllInfo() got error: %v", err)
		return
	}

	standardLicenses, err := loadStandardLicenses()
	if err != nil {
		t.Errorf("LoadStandardLicenses() got error: %v", err)
		return
	}

	exceptionLicenses, err := loadExceptionLicenses()
	if err != nil {
		t.Errorf("TestLoadExceptionLicenses() got error: %v", err)
		return
	}

	if len(m) != len(standardLicenses.Licenses)+len(exceptionLicenses.Exceptions) {
		t.Errorf("Expect number of exception: %d, got: %d", len(m), len(standardLicenses.Licenses)+len(exceptionLicenses.Exceptions))
		return
	}
	for _, licenseInfo := range m {
		if standardLicenses.contains(licenseInfo) || exceptionLicenses.contains(licenseInfo) {
			continue
		}
		t.Errorf("License '%s' doesn't exist in standard licenses and exception licenses", licenseInfo.LicenseID)
	}
}

func TestLicenseInfo_LicenseContentPath(t *testing.T) {
	tests := []struct {
		name        string
		licenseInfo LicenseInfo
		want        string
	}{
		{
			name: "Emtpy id",
			licenseInfo: LicenseInfo{
				LicenseID: "",
			},
			want: "",
		},
		{
			name: "Full name",
			licenseInfo: LicenseInfo{
				LicenseID: "AAL",
			},
			want: "AAL.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.licenseInfo.LicenseContentPath(); got != tt.want {
				t.Errorf("LicenseInfo.LicenseContentPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLicenseInfo_LoadLicenseContent(t *testing.T) {
	tests := []struct {
		name    string
		argInfo LicenseInfo
		want    LicenseContent
		wantErr bool
	}{
		{
			name: "standard license",
			argInfo: LicenseInfo{
				LicenseID: "0BSD",
			},
			want: LicenseContent{
				LicenseID: "0BSD",
				Content: []byte(`Copyright (C) 2006 by Rob Landley <rob@landley.net>

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE
OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THIS SOFTWARE.
`),
				RawContent: []byte(`Copyright (C) 2006 by Rob Landley <rob@landley.net>  Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted.  THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE. `),
			},
			wantErr: false,
		},
		{
			name: "exception license",
			argInfo: LicenseInfo{
				LicenseID: "Libtool-exception",
			},
			want: LicenseContent{
				LicenseID: "Libtool-exception",
				Content: []byte(`As a special exception to the GNU General Public License, if you distribute this file as part of a program or library that is built using GNU Libtool, you may include this file under the same distribution terms that you use for the rest of that program.
`),
				RawContent: []byte(`As a special exception to the GNU General Public License, if you distribute this file as part of a program or library that is built using GNU Libtool, you may include this file under the same distribution terms that you use for the rest of that program. `),
			},
			wantErr: false,
		},
		{
			name: "wrong license id",
			argInfo: LicenseInfo{
				LicenseID: "it's not real",
			},
			want:    LicenseContent{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.argInfo.LoadLicenseContent()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadLicenseContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadLicenseContent() = %v, want %v", got, tt.want)
			}
		})
	}
	licenseInfo, err := AllInfo()
	if err != nil {
		t.Errorf("GetLicenseInfo() error = %v", err)
		return
	}
	for _, licenseI := range licenseInfo {
		t.Run(licenseI.LicenseID, func(t *testing.T) {
			actual, err := licenseI.LoadLicenseContent()
			if err != nil {
				t.Errorf("LoadLicenseContent() error = %v", err)
				return
			}
			if actual.LicenseID != licenseI.LicenseID {
				t.Errorf("GetLicenseInfo() = %v, want %v", actual.LicenseID, licenseI.LicenseID)
				return
			}

			path := licenseI.LicenseContentPath()
			expectedData, _ := data.Asset(path)
			if !reflect.DeepEqual(actual.Content, expectedData) {
				t.Errorf("GetLicenseInfo() = %v, want %v", string(actual.Content), string(expectedData))
			}

			expectedRawData := regexp.MustCompile(`\r?\n`).ReplaceAll(expectedData, []byte(" "))
			if !reflect.DeepEqual(actual.RawContent, expectedRawData) {
				t.Errorf("GetLicenseInfo() raw = %v, want raw %v", string(actual.Content), string(expectedData))
			}
		})
	}
}

func TestAll(t *testing.T) {
	licenses, err := All()
	if err != nil {
		t.Errorf("All() error = %v", err)
		return
	}
	for _, license := range licenses {
		expectedLicenseContent, err := license.LicenseInfo.LoadLicenseContent()
		if err != nil {
			t.Errorf("LoadLicenseContent(%s) error = %v", license.LicenseInfo.Name, err)
			return
		}
		if !reflect.DeepEqual(license.LicenseContent, expectedLicenseContent) {
			t.Errorf("All() = %v, want %v", license.LicenseContent, expectedLicenseContent)
		}
	}
}
