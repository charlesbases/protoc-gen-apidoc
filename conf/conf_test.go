package conf

import (
	"testing"
)

func Test(t *testing.T) {
	Parse("configfile=../apidoc.yaml")
}
