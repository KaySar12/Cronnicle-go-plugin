package massdns

import (
	"NextDomain-Utils/model"
	"NextDomain-Utils/utils/cmd"
	"NextDomain-Utils/utils/files"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParseDNSQueries(fileOp files.FileOp, filename string) ([]model.DNSQuery, error) {
	var queries []model.DNSQuery
	var query model.DNSQuery

	content, err := fileOp.GetContent(filename)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(string(content))
	decoder := json.NewDecoder(reader)

	for {
		err := decoder.Decode(&query)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		queries = append(queries, query)
	}

	return queries, nil
}
func BulkLookup(domains string, resolvers string, results string, recordTypes []string) error {
	var recordString strings.Builder
	for _, recordType := range recordTypes {
		recordString.WriteString(fmt.Sprintf("-t %s ", recordType))
	}
	err := cmd.ExecCmdWithOutput(fmt.Sprintf("massdns -r %s %s -o J -w %s --status-format ansi %s", resolvers, recordString.String(), results, domains))
	if err != nil {
		return err
	}
	return nil
}
func ParseInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}
