package config

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func parseConfig(fileName string) (*dotNotationMap, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	arr := newDotNotationMap()

	lines := strings.Split(string(content), "\n")

	currentKey := ""
	currentValue := ""

	currentTable := ""

	for _, line := range lines {
		line = removeComments(line)

		if len(line) == 0 {
			continue
		}

		if strings.HasSuffix(line, "\"\"\"") && currentKey != "" {
			line = line[:len(line)-3]
			currentValue += line

			arr.Set(
				JoinKeys(currentTable, currentKey),
				parseValue(currentValue),
			)

			currentValue = ""
			currentKey = ""
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentTable = strings.TrimSpace(line[1 : len(line)-1])
			continue
		}

		key, value, ok := strings.Cut(line, "=")

		if !ok {
			if currentKey != "" {
				if !strings.HasSuffix(line, "\\") {
					line += "\n"
				} else {
					line = line[:len(line)-1]
				}

				currentValue += line
				continue
			}

			fmt.Println("Invalid line:", line)
			continue
		}
		key = strings.TrimSpace(key)

		if key == "" {
			fmt.Println("Invalid Line:", line)
			continue
		}

		if k, err := strconv.Unquote(key); err == nil {
			key = k
		}

		value = strings.TrimSpace(value)

		if strings.HasPrefix(value, "\"\"\"") {
			currentValue = ""

			if !strings.HasSuffix(line, "\\") {
				currentValue += "\n"
			}
			currentKey = key
			continue
		}

		arr.Set(
			JoinKeys(currentTable, key),
			parseValue(value),
		)
	}

	return arr, nil
}

func JoinKeys(keys ...string) string {
	out := keys[:0]
	for _, v := range keys {
		if strings.TrimSpace(v) != "" {
			out = append(out, v)
		}
	}
	return strings.Join(out, ".")
}

func removeComments(line string) string {
	before, _, _ := strings.Cut(line, "#")
	return strings.TrimSpace(before)
}

func parseValue(value string) interface{} {

	if len(value) == 0 {
		return nil
	}

	if strings.ToLower(value) == "true" {
		return true
	}

	if strings.ToLower(value) == "false" {
		return false
	}

	if strings.ToLower(value) == "null" {
		return nil
	}

	if strings.ToLower(value) == "+nan" {
		return math.NaN()
	}

	if strings.ToLower(value) == "-nan" {
		return -math.NaN()
	}

	if strings.HasPrefix(value, "0x") {
		v := strings.ReplaceAll(strings.TrimSpace(value), "_", "")

		if v, err := strconv.ParseInt(v, 0, 64); err == nil {
			return v
		}
	}

	if strings.HasPrefix(value, "0o") {
		v := strings.ReplaceAll(strings.TrimSpace(value), "_", "")
		if v, err := strconv.ParseInt(v, 0, 64); err == nil {
			return v
		}
	}

	if strings.HasPrefix(value, "0b") {
		v := strings.ReplaceAll(strings.TrimSpace(value), "_", "")
		if v, err := strconv.ParseInt(v, 0, 64); err == nil {
			return v
		}
	}

	if v, err := strconv.Atoi(value); err == nil {
		return v
	}

	if v, err := strconv.ParseFloat(value, 64); err == nil {
		return v
	}

	if v, err := unquoteMaybe(value); err == nil {
		return v
	}

	return value
}

func unquoteMaybe(s string) (string, error) {
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1], nil
	}
	return strconv.Unquote(s)
}
