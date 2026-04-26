package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadDotEnv reads a KEY=VAL file; it only sets a key when it is not already in the environment.
func LoadDotEnv(path string) error { return loadDotenv(path, false) }

// LoadDotEnvOverride reads a KEY=VAL file; it sets every key from the file, overriding existing.
func LoadDotEnvOverride(path string) error { return loadDotenv(path, true) }

func loadDotenv(path string, override bool) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(line[len("export "):])
		}
		i := strings.IndexByte(line, '=')
		if i < 0 {
			return fmt.Errorf("%s:%d: no = in line", path, lineNo)
		}
		key := strings.TrimSpace(line[:i])
		val := strings.TrimSpace(line[i+1:])
		val = unquoteValue(val)
		if key == "" {
			return fmt.Errorf("%s:%d: empty key", path, lineNo)
		}
		if !override {
			if _, set := os.LookupEnv(key); set {
				continue
			}
		}
		_ = os.Setenv(key, val)
	}
	return sc.Err()
}

func unquoteValue(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return strings.Trim(s, `"`)
		}
		if s[0] == '\'' && s[len(s)-1] == '\'' {
			return strings.Trim(s, `'`)
		}
	}
	return s
}
