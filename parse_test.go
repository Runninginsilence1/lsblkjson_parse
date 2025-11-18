package lsblkjson_parse

import (
	"fmt"
	"testing"

	"github.com/Runninginsilence1/lsblkjson_parse/internal/utils"
)

func TestAll(t *testing.T) {
	disk, err := ReadForensicDisk()
	if err != nil {
		t.Error(err)
		return
	}
	if len(disk) == 0 {
		t.Error("No disk found")
		return
	}
}

func TestOutput(t *testing.T) {
	disk, err := ReadForensicDisk()
	if err != nil {
		t.Error(err)
		return
	}
	pretty, err := utils.JSONMarshalPretty(disk)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(pretty))
}
