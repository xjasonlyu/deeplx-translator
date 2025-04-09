package deeplx_translator

import (
	"strings"
)

type Version uint8

const (
	VersionV1 Version = iota + 1
	VersionV2
)

func (v Version) IsValid() bool {
	switch v {
	case VersionV1, VersionV2:
		return true
	default:
		return false
	}
}

func inferVersionFromBaseURL(baseURL string) Version {
	switch {
	case strings.HasSuffix(baseURL, "/v1"):
		return VersionV1
	case strings.HasSuffix(baseURL, "/v2"):
		return VersionV2
	default:
		return VersionV1 /* default v1 */
	}
}
