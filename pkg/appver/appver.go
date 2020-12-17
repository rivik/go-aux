package appver

import (
	"log"
	"strconv"
)

var (
	// link-time string vars
	LDSemVer, LDAltVer, LDRevision string
	LDBuildNo, LDBuildTS           string

	// parsed appver struct
	Version AppVersion
)

/* AppVersion struct

SemVer - semantic version, like '0.8.4'
AltVer - alternative version, like '2020.10-51'
Revision - commit hash or revision number, like '47dc930'
BuildNo - autoincremented build number, like 51
BuildTS - build timestamp (unix epoch), in milliseconds

All but one fields can be empty ('' for text fields, 0 for numeric).
Otherwise, each field MUST uniquely identify this build by itself.
*/

type AppVersion struct {
	SemVer   string `json:"semver"`
	AltVer   string `json:"altver"` // alternative version, like 2020.10
	Revision string `json:"revision"`
	BuildNo  int64  `json:"build_number"`
	BuildTS  int64  `json:"build_timestamp_millis"`
}

func NewAppVersion(semver, altver, revision, buildno, buildts string) (AppVersion, error) {
	var err error
	appVer := AppVersion{}

	if buildts != "" {
		appVer.BuildTS, err = strconv.ParseInt(buildts, 10, 64)
		if err != nil {
			return appVer, err
		}
	}
	if buildno != "" {
		appVer.BuildNo, err = strconv.ParseInt(buildno, 10, 64)
		if err != nil {
			return appVer, err
		}
	}

	appVer.SemVer = semver
	appVer.AltVer = altver
	appVer.Revision = revision

	return appVer, nil
}

func init() {
	appVer, err := NewAppVersion(LDSemVer, LDAltVer, LDRevision, LDBuildNo, LDBuildTS)
	if err != nil {
		log.Panic(err)
	}
	Version = appVer
}
