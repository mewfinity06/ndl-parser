package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	VecStartingSize int     `json:"vec_starting_size"`
	VecGrowthFactor float64 `json:"vec_growth_factor"`
	DevMode         bool    `json:"dev_mode"`
	Debug           bool    `json:"debug"`
	OutputFile      string  `json:"output_file"`
}

var configPattern = regexp.MustCompile(`\[(.*?)\](.*)`) // Precompiled regex for efficiency

func parseConfig(filename string) Config {
	config := Config{VecGrowthFactor: 1.5, OutputFile: "main"}
	if filename == "" {
		return config
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
		return config
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}

		matches := configPattern.FindStringSubmatch(line)
		if len(matches) == 3 {
			key := matches[1]
			value := strings.TrimSpace(strings.TrimSuffix(matches[2], ";"))

			switch key {
			case "vec_starting_size":
				if v, err := strconv.Atoi(value); err == nil {
					config.VecStartingSize = v
				}
			case "vec_growth_factor":
				if v, err := strconv.ParseFloat(value, 64); err == nil {
					config.VecGrowthFactor = v
				}
			case "output":
				config.OutputFile = value
			case "dev_mode":
				config.DevMode = (value == "true")
			case "debug":
				config.Debug = (value == "true")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
	}

	return config
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: main <config-directory>")
		os.Exit(1)
	}

	configDir := os.Args[1]
	files, err := filepath.Glob(filepath.Join(configDir, "*.ndl"))
	if err != nil || len(files) == 0 {
		fmt.Fprintln(os.Stderr, "No config file found in the specified directory")
		os.Exit(1)
	}

	config := parseConfig(files[0])
	jsonOutput, err := json.Marshal(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error encoding JSON:", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonOutput))
}
