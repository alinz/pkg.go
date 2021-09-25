package sqlite_test

import (
	"strings"
	"testing"

	"github.com/alinz/pkg.go/sqlite"
)

func TestGenerateInValues(t *testing.T) {
	sqlSegment, params := sqlite.GenerateInValues("sample", 2)

	if sqlSegment != "$sample0,$sample1" {
		t.Errorf("expect %s, but got %s", "$sample0,$sample1", sqlSegment)
	}

	if strings.Join(params, ",") != "$sample0,$sample1" {
		t.Errorf("expect %s, but got %s", "$sample0,$sample1", sqlSegment)
	}
}
