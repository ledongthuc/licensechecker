package data

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pkg/errors"
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

func Test_loadLicenses(t *testing.T) {
	var exampleLicense license
	err := json.Unmarshal([]byte(exampleLicenseRaw), &exampleLicense)
	if err != nil {
		t.Errorf("Can't pase example license for testing: %s", err)
	}

	type args struct {
		container map[string]LicenseInfo
		ls        license
	}
	tests := []struct {
		name          string
		args          args
		expected      map[string]LicenseInfo
		expectedError error
	}{
		{
			name: "Test example license",
			args: args{
				container: make(map[string]LicenseInfo),
				ls:        exampleLicense,
			},
			expected: map[string]LicenseInfo{
				"0BSD": LicenseInfo{
					ID:   "0BSD",
					Name: "BSD Zero Clause License",
					References: []string{
						"http://landley.net/toybox/license.html",
					},
					IsDeprecated: false,
				},
				"AAL": LicenseInfo{
					ID:   "AAL",
					Name: "Attribution Assurance License",
					References: []string{
						"https://opensource.org/licenses/attribution",
					},
					IsDeprecated: false,
				},
				"AGPL-3.0": LicenseInfo{
					ID:   "AGPL-3.0",
					Name: "GNU Affero General Public License v3.0",
					References: []string{
						"https://www.gnu.org/licenses/agpl.txt",
						"https://opensource.org/licenses/AGPL-3.0",
					},
					IsDeprecated: true,
				},
			},
		},
		{
			name: "Test nil map",
			args: args{
				container: nil,
				ls:        exampleLicense,
			},
			expected:      nil,
			expectedError: errors.New("Container shouldn't nil"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loadLicenses(tt.args.container, tt.args.ls)
			if err != nil && tt.expectedError == nil {
				t.Errorf("not match error, expected=%s, actual=nil", err)
				return
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("not match error, expected=nil, actual=%s", tt.expectedError)
				return
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("not match error, expected=%s, actual=%s", err, tt.expectedError)
				return
			}
			if len(tt.args.container) != len(tt.expected) {
				t.Errorf("not match loading licenses, expected len=%d, actual len=%d", len(tt.expected), len(tt.args.container))
				return
			}
			for id, expectedLicense := range tt.expected {
				actualLicense, exist := tt.args.container[id]
				if !exist {
					t.Errorf("not match license: %s, expected: exist, actual=not exist", id)
				}
				if !reflect.DeepEqual(expectedLicense, actualLicense) {
					t.Errorf("not match license: %s, expected: %+v, actual=%+v", id, expectedLicense, actualLicense)
				}
			}
		})
	}
}

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

func Test_loadExceptionLicenses(t *testing.T) {
	var exampleExceptionLicense exception
	err := json.Unmarshal([]byte(exampleExceptionLicenseRaw), &exampleExceptionLicense)
	if err != nil {
		t.Errorf("Can't pase example exception license for testing: %s", err)
	}

	type args struct {
		container map[string]LicenseInfo
		ls        exception
	}
	tests := []struct {
		name          string
		args          args
		expected      map[string]LicenseInfo
		expectedError error
	}{
		{
			name: "Test example license",
			args: args{
				container: make(map[string]LicenseInfo),
				ls:        exampleExceptionLicense,
			},
			expected: map[string]LicenseInfo{
				"Libtool-exception": LicenseInfo{
					ID:   "Libtool-exception",
					Name: "Libtool Exception",
					References: []string{
						"http://git.savannah.gnu.org/cgit/libtool.git/tree/m4/libtool.m4",
					},
					IsDeprecated: false,
				},
				"Classpath-exception-2.0": LicenseInfo{
					ID:   "Classpath-exception-2.0",
					Name: "Classpath exception 2.0",
					References: []string{
						"http://www.gnu.org/software/classpath/license.html",
						"https://fedoraproject.org/wiki/Licensing/GPL_Classpath_Exception",
					},
					IsDeprecated: false,
				},
				"Nokia-Qt-exception-1.1": LicenseInfo{
					ID:   "Nokia-Qt-exception-1.1",
					Name: "Nokia Qt LGPL exception 1.1",
					References: []string{
						"https://www.keepassx.org/dev/projects/keepassx/repository/revisions/b8dfb9cc4d5133e0f09cd7533d15a4f1c19a40f2/entry/LICENSE.NOKIA-LGPL-EXCEPTION",
					},
					IsDeprecated: true,
				},
			},
		},
		{
			name: "Test nil map",
			args: args{
				container: nil,
				ls:        exampleExceptionLicense,
			},
			expected:      nil,
			expectedError: errors.New("Container shouldn't nil"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loadExceptionLicenses(tt.args.container, tt.args.ls)
			if err != nil && tt.expectedError == nil {
				t.Errorf("not match error, expected=%s, actual=nil", err)
				return
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("not match error, expected=nil, actual=%s", tt.expectedError)
				return
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("not match error, expected=%s, actual=%s", err, tt.expectedError)
				return
			}
			if len(tt.args.container) != len(tt.expected) {
				t.Errorf("not match loading licenses, expected len=%d, actual len=%d", len(tt.expected), len(tt.args.container))
				return
			}
			for id, expectedLicense := range tt.expected {
				actualLicense, exist := tt.args.container[id]
				if !exist {
					t.Errorf("not match license: %s, expected: exist, actual=not exist", id)
				}
				if !reflect.DeepEqual(expectedLicense, actualLicense) {
					t.Errorf("not match license: %s, expected: %+v, actual=%+v", id, expectedLicense, actualLicense)
				}
			}
		})
	}
}
