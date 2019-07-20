package data

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ledongthuc/licensechecker/data/toc"
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

func Test_convertStandardLicenses(t *testing.T) {
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
			expectedError: ErrorUninitiatedContainer,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := convertStandardLicenses(tt.args.container, tt.args.ls)
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

func Test_convertExceptionLicenses(t *testing.T) {
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
			expectedError: ErrorUninitiatedContainer,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := convertExceptionLicenses(tt.args.container, tt.args.ls)
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

func TestGetLicenseInfo(t *testing.T) {
	m, err := GetLicenseInfo()
	if err != nil {
		t.Errorf("GetLicenseInfo() got error: %v", err)
		return
	}

	standardLicenses, err := LoadStandardLicenses()
	if err != nil {
		t.Errorf("LoadStandardLicenses() got error: %v", err)
		return
	}

	exceptionLicenses, err := LoadExceptionLicenses()
	if err != nil {
		t.Errorf("TestLoadExceptionLicenses() got error: %v", err)
		return
	}

	if len(m) != len(standardLicenses.Licenses)+len(exceptionLicenses.Exceptions) {
		t.Errorf("Expect number of exception: %d, got: %d", len(m), len(standardLicenses.Licenses)+len(exceptionLicenses.Exceptions))
		return
	}
	for licenseKey, licenseInfo := range m {
		if licenseKey != licenseInfo.ID {
			t.Errorf("LicenseKey '%s' != licenseInfo.ID '%s'", licenseKey, licenseInfo.ID)
		}
		if standardLicenses.contains(licenseInfo) || exceptionLicenses.contains(licenseInfo) {
			continue
		}
		t.Errorf("License '%s' doesn't exist in standard licenses and exception licenses", licenseKey)
	}
}

func TestLoadStandardLicenses(t *testing.T) {
	l, err := LoadStandardLicenses()
	if err != nil {
		t.Errorf("LoadStandardLicenses() got error: %v", err)
		return
	}
	reversedLicense, err := json.Marshal(l)
	if err != nil {
		t.Errorf("LoadStandardLicenses() output can't reverse")
		return
	}

	raw, err := toc.Asset(listLicenses)
	if err != nil {
		t.Errorf("Can't load raw data to check from '%s'", listLicenses)
		return
	}
	var rawLicenses license
	err = json.Unmarshal(raw, &rawLicenses)
	if err != nil {
		t.Errorf(errors.Wrapf(err, "Can't load raw data to check from '%s'", listLicenses).Error())
		return
	}
	raw, err = json.Marshal(rawLicenses)
	if err != nil {
		t.Errorf(errors.Wrapf(err, "Can't load raw data to check from '%s'", listLicenses).Error())
		return
	}

	if !reflect.DeepEqual(reversedLicense, raw) {
		t.Errorf("LoadStandardLicenses(): output is not match with raw data in: %s", listLicenses)
	}
}

func TestLoadExceptionLicenses(t *testing.T) {
	l, err := LoadExceptionLicenses()
	if err != nil {
		t.Errorf("LoadExceptionLicenses() got error: %v", err)
		return
	}
	reversedLicense, err := json.Marshal(l)
	if err != nil {
		t.Errorf("LoadExceptionLicenses() output can't reverse")
		return
	}

	raw, err := toc.Asset(listExceptions)
	if err != nil {
		t.Errorf("Can't load raw data to check from '%s'", listExceptions)
		return
	}
	var rawLicenses exception
	err = json.Unmarshal(raw, &rawLicenses)
	if err != nil {
		t.Errorf(errors.Wrapf(err, "Can't load raw data to check from '%s'", listLicenses).Error())
		return
	}
	raw, err = json.Marshal(rawLicenses)
	if err != nil {
		t.Errorf(errors.Wrapf(err, "Can't load raw data to check from '%s'", listLicenses).Error())
		return
	}

	if !reflect.DeepEqual(reversedLicense, raw) {
		t.Errorf("LoadExceptionLicenses(): output is not match with raw data in: %s", listLicenses)
	}
}

func Test_license_contains(t *testing.T) {
	var exampleLicense license
	err := json.Unmarshal([]byte(exampleLicenseRaw), &exampleLicense)
	if err != nil {
		t.Errorf("Can't pase example license for testing: %s", err)
	}

	tests := []struct {
		name string
		l    license
		info LicenseInfo
		want bool
	}{
		{
			name: "Contains",
			l:    exampleLicense,
			info: LicenseInfo{
				ID:           "AGPL-3.0",
				Name:         "GNU Affero General Public License v3.0",
				IsDeprecated: true,
				References: []string{
					"https://www.gnu.org/licenses/agpl.txt",
					"https://opensource.org/licenses/AGPL-3.0",
				},
			},
			want: true,
		},
		{
			name: "Wrong ID",
			l:    exampleLicense,
			info: LicenseInfo{
				ID:           "AGPL-3.1",
				Name:         "GNU Affero General Public License v3.0",
				IsDeprecated: true,
				References: []string{
					"https://www.gnu.org/licenses/agpl.txt",
					"https://opensource.org/licenses/AGPL-3.0",
				},
			},
			want: false,
		},
		{
			name: "Wrong Name",
			l:    exampleLicense,
			info: LicenseInfo{
				ID:           "AGPL-3.0",
				Name:         "GNU Affero General Public License v3.1",
				IsDeprecated: true,
				References: []string{
					"https://www.gnu.org/licenses/agpl.txt",
					"https://opensource.org/licenses/AGPL-3.0",
				},
			},
			want: false,
		},
		{
			name: "Wrong Deprecated",
			l:    exampleLicense,
			info: LicenseInfo{
				ID:           "AGPL-3.0",
				Name:         "GNU Affero General Public License v3.0",
				IsDeprecated: false,
				References: []string{
					"https://www.gnu.org/licenses/agpl.txt",
					"https://opensource.org/licenses/AGPL-3.0",
				},
			},
			want: false,
		},
		{
			name: "Wrong References - different",
			l:    exampleLicense,
			info: LicenseInfo{
				ID:           "AGPL-3.0",
				Name:         "GNU Affero General Public License v3.0",
				IsDeprecated: true,
				References: []string{
					"https://www.gnu.org/licenses/agpl.json",
					"https://opensource.org/licenses/AGPL-3.0",
				},
			},
			want: false,
		},
		{
			name: "Wrong References - miss",
			l:    exampleLicense,
			info: LicenseInfo{
				ID:           "AGPL-3.0",
				Name:         "GNU Affero General Public License v3.0",
				IsDeprecated: true,
				References: []string{
					"https://opensource.org/licenses/AGPL-3.0",
				},
			},
			want: false,
		},
		{
			name: "Wrong References - more",
			l:    exampleLicense,
			info: LicenseInfo{
				ID:           "AGPL-3.0",
				Name:         "GNU Affero General Public License v3.0",
				IsDeprecated: true,
				References: []string{
					"https://www.gnu.org/licenses/agpl.txt",
					"https://opensource.org/licenses/AGPL-3.0",
					"https://spin.atomicobject.com/2014/11/23/record-vim-macros/",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.contains(tt.info); got != tt.want {
				t.Errorf("license.contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_exception_contains(t *testing.T) {
	var exampleExceptionLicense exception
	err := json.Unmarshal([]byte(exampleExceptionLicenseRaw), &exampleExceptionLicense)
	if err != nil {
		t.Errorf("Can't pase example exception license for testing: %s", err)
	}

	tests := []struct {
		name string
		e    exception
		info LicenseInfo
		want bool
	}{
		{
			name: "Contains",
			e:    exampleExceptionLicense,
			info: LicenseInfo{
				ID:           "Libtool-exception",
				Name:         "Libtool Exception",
				IsDeprecated: false,
				References: []string{
					"http://git.savannah.gnu.org/cgit/libtool.git/tree/m4/libtool.m4",
				},
			},
			want: true,
		},
		{
			name: "Wrong ID",
			e:    exampleExceptionLicense,
			info: LicenseInfo{
				ID:           "Libtool-exception-12",
				Name:         "Libtool Exception",
				IsDeprecated: false,
				References: []string{
					"http://git.savannah.gnu.org/cgit/libtool.git/tree/m4/libtool.m4",
				},
			},
			want: false,
		},
		{
			name: "Wrong Name",
			e:    exampleExceptionLicense,
			info: LicenseInfo{
				ID:           "Libtool-exception",
				Name:         "Libtool 123 Exception",
				IsDeprecated: false,
				References: []string{
					"http://git.savannah.gnu.org/cgit/libtool.git/tree/m4/libtool.m4",
				},
			},
			want: false,
		},
		{
			name: "Wrong Deprecated",
			e:    exampleExceptionLicense,
			info: LicenseInfo{
				ID:           "Libtool-exception",
				Name:         "Libtool Exception",
				IsDeprecated: true,
				References: []string{
					"http://git.savannah.gnu.org/cgit/libtool.git/tree/m4/libtool.m4",
				},
			},
			want: false,
		},
		{
			name: "Wrong References - different",
			e:    exampleExceptionLicense,
			info: LicenseInfo{
				ID:           "Libtool-exception",
				Name:         "Libtool Exception",
				IsDeprecated: false,
				References: []string{
					"http://git.savannah.gnu.org/cgit/libtool.git/tree/m4/libtool.m5",
				},
			},
			want: false,
		},
		{
			name: "Wrong References - miss",
			e:    exampleExceptionLicense,
			info: LicenseInfo{
				ID:           "Libtool-exception",
				Name:         "Libtool Exception",
				IsDeprecated: true,
				References:   []string{},
			},
			want: false,
		},
		{
			name: "Wrong References - more",
			e:    exampleExceptionLicense,
			info: LicenseInfo{
				ID:           "Libtool-exception",
				Name:         "Libtool Exception",
				IsDeprecated: true,
				References: []string{
					"http://git.savannah.gnu.org/cgit/libtool.git/tree/m4/libtool.m4",
					"http://git.savannah.gnu.org/cgit/libtool.git/tree/m4/libtool.m5",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.contains(tt.info); got != tt.want {
				t.Errorf("exception.contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
