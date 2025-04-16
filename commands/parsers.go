package commands

import (
  "fmt"
  "strings"
  "bufio"
  "strconv"
)

type KeyspaceInfo struct {
  Keys int64
}

func parseDBLine(line string) (*KeyspaceInfo, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format: expected ':' separator in %q", line)
	}
	dbPart := parts[0]
	dataPart := parts[1]

	info := &KeyspaceInfo{}

	if !strings.HasPrefix(dbPart, "db0") {
		return nil, fmt.Errorf("invalid format: expected dbPart %q to start with 'db0'", dbPart)
	}

	kvPairs := strings.Split(dataPart, ",")
	for _, pair := range kvPairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid format: expected key=value pair in %q", pair)
		}
		key := kv[0]
		valueStr := kv[1]

		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid format: could not parse value %q for key %q: %w", valueStr, key, err)
		}

		switch key {
		case "keys":
			info.Keys = value
		default:
			// fmt.Printf("Warning: unknown key %q encountered\n", key)
		}
	}

	return info, nil
}


func parseKeyspaceDBLines(inputText string) ([]KeyspaceInfo, error) {
	var dbLines []KeyspaceInfo
	var inKeyspaceSection bool = false

	// Use a scanner to efficiently read the input text line by line.
	scanner := bufio.NewScanner(strings.NewReader(inputText))

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "# Keyspace" {
			inKeyspaceSection = true 
			continue                
		}

		if inKeyspaceSection {
			if strings.HasPrefix(trimmedLine, "db0") {
        info, err := parseDBLine(trimmedLine)
        if err != nil {
          return dbLines, err
        }
        dbLines = append(dbLines, *info)

			} else if strings.HasPrefix(trimmedLine, "#") {
				inKeyspaceSection = false
			} else if trimmedLine == "" {
				// Skip empty lines.
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning input text: %w", err)
	}

	return dbLines, nil
}

