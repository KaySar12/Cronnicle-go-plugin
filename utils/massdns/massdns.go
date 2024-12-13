package massdns

import (
	"NextDomain-Utils/model"
	"NextDomain-Utils/utils/cmd"
	"NextDomain-Utils/utils/files"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func ParseDNSQueries(fileOp files.FileOp, baseFilename string, processors int) ([]model.DNSQuery, error) {
	if processors <= 1 {
		return parseSingleFile(fileOp, baseFilename)
	}

	// Merge multiple files into a single content buffer
	var mergedContent strings.Builder
	for i := 0; i <= processors; i++ {
		filename := baseFilename
		if i > 0 {
			filename = fmt.Sprintf("%s%d", baseFilename, i-1)
		}

		content, err := fileOp.GetContent(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to get content from file %s: %w", filename, err)
		}
		if len(content) == 0 {
			if err := fileOp.Fs.Remove(filename); err != nil {
				return nil, fmt.Errorf("failed to remove file %s: %w", filename, err)
			}
			continue // Skip empty files
		}
		mergedContent.WriteString(string(content))
		mergedContent.WriteString("\n")

		// Remove the individual file
		if err := fileOp.Fs.Remove(filename); err != nil {
			return nil, fmt.Errorf("failed to remove file %s: %w", filename, err)
		}
	}

	// Write merged content back to the base file
	err := fileOp.SaveFile(baseFilename, mergedContent.String(), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to save merged file %s: %w", baseFilename, err)
	}

	// Parse the merged file
	return parseSingleFile(fileOp, baseFilename)
}
func parseSingleFile(fileOp files.FileOp, filename string) ([]model.DNSQuery, error) {
	content, err := fileOp.GetContent(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get content from file %s: %w", filename, err)
	}

	reader := strings.NewReader(string(content))
	decoder := json.NewDecoder(reader)

	var queries []model.DNSQuery
	queryMap := make(map[string]bool)
	for decoder.More() {
		var query model.DNSQuery
		if err := decoder.Decode(&query); err != nil {
			return nil, fmt.Errorf("failed to decode JSON in file %s: %w", filename, err)
		}
		// Check for duplicates
		queryKey := fmt.Sprintf("%s:%s", query.Name, query.Type)
		if !queryMap[queryKey] {
			queries = append(queries, query)
			queryMap[queryKey] = true
		}
	}

	return queries, nil
}

func BulkLookup(domains string, resolvers string, results string, recordTypes []string, processor int, logpath string) error {
	var recordString strings.Builder
	for _, recordType := range recordTypes {
		recordString.WriteString(fmt.Sprintf("-t %s ", recordType))
	}
	err := cmd.ExecCmdWithOutput(fmt.Sprintf("massdns -r %s %s -o J -w %s --processes %d --status-format ansi %s --error-log %s", resolvers, recordString.String(), results, processor, domains, logpath))
	if err != nil {
		return err
	}
	return nil
}
func ParseInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}
