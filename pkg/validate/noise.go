package validate

import (
	"path/filepath"
	"strings"
)

type NoiseFinding struct {
	Path   string
	Reason string
}

func FindNoise(paths []string) []NoiseFinding {
	var findings []NoiseFinding
	for _, p := range paths {
		p = filepath.ToSlash(strings.TrimSpace(p))
		if p == "" {
			continue
		}

		base := filepath.Base(p)
		ext := strings.ToLower(filepath.Ext(base))

		if strings.HasPrefix(p, "node_modules/") || p == "node_modules" {
			findings = append(findings, NoiseFinding{Path: p, Reason: "node_modules/"})
			continue
		}
		if strings.HasPrefix(p, "vendor/") || p == "vendor" {
			findings = append(findings, NoiseFinding{Path: p, Reason: "vendor/"})
			continue
		}
		if strings.HasPrefix(p, "dist/") || p == "dist" {
			findings = append(findings, NoiseFinding{Path: p, Reason: "dist/ build output"})
			continue
		}
		if strings.HasPrefix(p, "build/") || p == "build" || strings.HasPrefix(p, "out/") || p == "out" {
			findings = append(findings, NoiseFinding{Path: p, Reason: "build output"})
			continue
		}
		if strings.HasPrefix(p, "tmp/") || p == "tmp" || strings.HasPrefix(p, "temp/") || p == "temp" {
			findings = append(findings, NoiseFinding{Path: p, Reason: "temporary files"})
			continue
		}
		if strings.HasPrefix(p, "coverage/") || p == "coverage" || ext == ".cover" || strings.HasSuffix(base, ".cover") {
			findings = append(findings, NoiseFinding{Path: p, Reason: "coverage output"})
			continue
		}
		if strings.HasPrefix(p, ".idea/") || p == ".idea" || strings.HasPrefix(p, ".vscode/") || p == ".vscode" {
			findings = append(findings, NoiseFinding{Path: p, Reason: "IDE config"})
			continue
		}
		if base == ".DS_Store" || base == "Thumbs.db" {
			findings = append(findings, NoiseFinding{Path: p, Reason: "OS junk"})
			continue
		}
		if base == ".env" || strings.HasSuffix(base, ".env") || strings.HasPrefix(base, ".env.") {
			findings = append(findings, NoiseFinding{Path: p, Reason: "env/secrets"})
			continue
		}
		if ext == ".log" {
			findings = append(findings, NoiseFinding{Path: p, Reason: "log file"})
			continue
		}
		if ext == ".exe" || ext == ".bin" || ext == ".dll" || ext == ".so" || ext == ".dylib" {
			findings = append(findings, NoiseFinding{Path: p, Reason: "binary artifact"})
			continue
		}
		if ext == ".pyc" || strings.Contains(p, "__pycache__/") {
			findings = append(findings, NoiseFinding{Path: p, Reason: "python cache"})
			continue
		}
		if ext == ".o" || ext == ".a" {
			findings = append(findings, NoiseFinding{Path: p, Reason: "compiled object"})
			continue
		}
	}

	return findings
}
