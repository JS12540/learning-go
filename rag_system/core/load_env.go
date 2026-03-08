package core

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// LoadEnv reads a .env file and sets each key=value as an environment variable.
// It is safe to call multiple times; existing env vars are NOT overwritten.
// Call this once at the very start of main() before anything else.
func LoadEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		// .env is optional — if it doesn't exist, rely on real env vars
		log.Printf("No .env file found (%s), using system environment variables", filename)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Strip surrounding quotes if present: "value" or 'value'
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// Only set if not already set in the real environment
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading .env file: %v", err)
	}
}
