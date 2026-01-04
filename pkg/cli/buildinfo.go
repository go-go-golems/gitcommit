package cli

import "fmt"

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

var buildInfo = BuildInfo{
	Version: "dev",
	Commit:  "none",
	Date:    "unknown",
}

func SetBuildInfo(version, commit, date string) {
	if version != "" {
		buildInfo.Version = version
	}
	if commit != "" {
		buildInfo.Commit = commit
	}
	if date != "" {
		buildInfo.Date = date
	}

	rootCmd.Version = buildInfo.Version
	rootCmd.SetVersionTemplate(fmt.Sprintf("{{with .Name}}{{printf \"%%s \" .}}{{end}}{{printf \"%%s\\n\" .Version}}commit: %s\\ndate: %s\\n", buildInfo.Commit, buildInfo.Date))
}
