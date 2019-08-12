package collector

import (
	"testing"
	"io/ioutil"
	"bufio"
	"bytes"
	"os"
)

func TestStructureRawStats(t *testing.T) {
  rawStats := new(bytes.Buffer)
	dataRoot := "../testdata/"
	files, err := ioutil.ReadDir(dataRoot)

  if err != nil {
      t.Error(err)
  }

  for _, statFile := range files {
		statFileContent, err := os.Open(dataRoot + statFile.Name())
		if err != nil {
			t.Error(err)
			os.Exit(1)
			t.Fail()
		}
		statReader := bufio.NewReader(statFileContent)
		rawStats.ReadFrom(statReader)
	}

	structuredStats := (*structureRawStats(rawStats))

	if want, got := 7, len((*structuredStats["df"])); want != got {
		t.Errorf("want %d df elements, got %d", want, got)
	}

	if want, got := 8, len((*structuredStats["diskstat"])); want != got {
		t.Errorf("want %d diskstat elements, got %d", want, got)
	}

	if want, got := 2, len((*structuredStats["mounts"])); want != got {
		t.Errorf("want %d mount elements, got %d", want, got)
	}

}
