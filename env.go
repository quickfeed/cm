package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close() // skipcq: GO-S2307

	log.Printf("Loading environment variables from %s", filename)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if ignore(line) {
			continue
		}
		key, val, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		k := strings.TrimSpace(key)
		if os.Getenv(k) != "" {
			// Ignore .env entries already set in the environment.
			continue
		}
		val = os.ExpandEnv(strings.Trim(strings.TrimSpace(val), `"`))
		os.Setenv(k, val)
	}

	return scanner.Err()
}

func ignore(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return trimmedLine == "" || strings.HasPrefix(trimmedLine, "#")
}
