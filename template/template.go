package template

import (
	"math"
	"regexp"

	"github.com/agnivade/levenshtein"
)

var Templates = map[string]string{
	"Apache License 1.0": Apache10,
	"Apache License 1.1": Apache11,
	"Apache License 2.0": Apache20,

	"BSD 2-Clause 'Simplified' License":       bsd2Clause,
	"BSD 3-Clause 'New' or 'Revised' License": bsd3Clause,
	"BSD 3-Clause Clear License":              bsd3ClauseClear,

	"Creative Commons Attribution 4.0": ccBy40,

	"GNU General Public License v2.0": gnu20,
	"GNU General Public License v3.0": gnu30,

	"GNU Lesser General Public License v2.0": lgpl20,
	"GNU Lesser General Public License v2.1": lgpl21,
	"GNU Lesser General Public License v3.0": lgpl30,

	"MIT License": mit,

	"ISC License": isc,
}

func FindMatchName(license string) (string, int) {
	// var matchName string
	// var matchPercent float64
	// for name, template := range Templates {
	// 	p := smetrics.JaroWinkler(license, template, 0.7, 4)
	// 	if p > matchPercent {
	// 		matchName = name
	// 		matchPercent = p
	// 	}
	// }
	// return matchName, matchPercent

	// dmp := diffmatchpatch.New()
	// var matchName string
	// noDiff := int(math.MaxInt32)
	// for name, template := range Templates {
	// 	diffs := dmp.DiffMain(license, template, false)
	// 	if len(diffs) < noDiff {
	// 		matchName = name
	// 		noDiff = len(diffs)
	// 	}
	// }
	// return matchName, 1.0

	// Need to replace by preformated
	license = regexp.MustCompile(`\r?\n`).ReplaceAllString(license, " ")

	var matchName string
	var distance = int(math.MaxInt32)
	for name, template := range Templates {

		// Need to replace by preformated
		template = regexp.MustCompile(`\r?\n`).ReplaceAllString(template, " ")

		d := levenshtein.ComputeDistance(license, template)
		if d < distance {
			matchName = name
			distance = d
		}
	}
	return matchName, distance
}
