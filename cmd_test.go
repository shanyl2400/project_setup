package main

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestListFilesInDirectory(t *testing.T) {
	config.Ignore = []string{".git", "web"}
	out, err := listFilesInDirectory("D:\\code\\go\\able\\sipserver\\")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(out)
}
func TestFilterForCandidates(t *testing.T) {
	config.Ignore = []string{".git", "web"}
	out, err := listFilesInDirectory("D:\\code\\go\\able\\sipserver\\")
	if err != nil {
		t.Fatal(err)
	}
	out["sip\\sipserver_test.go"] = struct{}{}
	out["sipserver.exe"] = struct{}{}
	out["deploy\\docker-compose.yaml"] = struct{}{}
	out["deploy\\data\\logs"] = struct{}{}
	out["config.ini"] = struct{}{}
	// t.Log("out:", out)
	InitRegex([]string{"^.git*", "^sipserver.exe$", "^config.ini$"}, []string{"_test.go$", "^deploy*"})

	moved, removed := filterForCandidates(out)
	t.Log("moved:", moved)
	t.Log("removed:", removed)
}

func TestCopyFile(t *testing.T) {
	src := "D:\\tmp\\p1\\sp1\\1.txt"
	dist := "D:\\tmp\\p3\\sp3\\spsp333\\1.txt"
	err := DeepCopyFile(src, dist)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("done2")
}

func TestMarshalYAML(t *testing.T) {
	config := &Config{
		Remove: []string{
			"*.exe$",
			"^.gitignore$",
			"^go.sum$",
		},
		Move: []string{
			"*_test.go$",
			"*config.go$",
			"^config.ini$",
		},

		Source: "D:\\tmp\\sipserver",
		Backup: "D:\\tmp\\sipserver_backup",
		Setup:  "D:\\tmp\\sipserver_test",

		MoveCpy:    "D:\\tmp\\sipserver_mov",
		ReleaseCpy: "D:\\tmp\\sipserver_release",

		Ignore: []string{".git", "web"},
	}

	yamlBytes, _ := yaml.Marshal(config)
	fmt.Println(string(yamlBytes))
}

func TestRestore(t *testing.T) {
	restore("D:\\tmp\\1", "D:\\tmp\\2")
}
