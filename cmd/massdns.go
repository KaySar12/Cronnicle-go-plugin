package cmd

import (
	"NextDomain-Utils/dto/request"
	"NextDomain-Utils/dto/response"
	model "NextDomain-Utils/model"
	"NextDomain-Utils/utils/cmd"
	"NextDomain-Utils/utils/files"
	"NextDomain-Utils/utils/massdns"
	util "NextDomain-Utils/utils/massdns"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(massdnsCmd)
	massdnsCmd.AddCommand(massLookupCmd)
	massdnsCmd.AddCommand(massLookupCmdDEv)
	massdnsCmd.AddCommand(checkZoneCmd)
	massdnsCmd.AddCommand(checkZonesCmdDEv)
	massLookupCmd.PersistentFlags().String("type", "", "DNS Record type")
	massLookupCmdDEv.PersistentFlags().String("type", "", "DNS Record type")
}

var massdnsCmd = &cobra.Command{
	Use:   "massdns",
	Short: "dns related commands",
}

var massLookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Lookup DNS Record using massdns",
	RunE: func(cmd *cobra.Command, args []string) error {
		recordType, _ := cmd.Flags().GetString("type")
		fmt.Println(recordType)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				continue
			}
			// Parse JSON input
			var job model.Job
			err := json.Unmarshal([]byte(line), &job)
			if err != nil {
				continue
			}
			processMassLookupJob(job, recordType)
		}
		return nil
	},
}
var checkZoneCmd = &cobra.Command{
	Use:   "check-zone",
	Short: "Check Zone ",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				continue
			}
			var job model.Job
			err := json.Unmarshal([]byte(line), &job)
			if err != nil {
				continue
			}
			checkZones(job)
		}
		return nil
	},
}
var massLookupCmdDEv = &cobra.Command{
	Use:   "lookup-dev",
	Short: "Lookup DNS Record using massdns",
	RunE: func(cmd *cobra.Command, args []string) error {
		var perf model.PerfStats
		start := time.Now()
		recordType, _ := cmd.Flags().GetString("type")
		data, err := os.ReadFile("/root/dev/Cronicle/Plugins/go-plugin/NextDomain-Utils/build/massdns.json")
		if err != nil {
			return fmt.Errorf("failed to read lookupdata.json: %w", err)
		}
		var job model.Job
		err = json.Unmarshal([]byte(data), &job)
		if err != nil {
			return fmt.Errorf("failed to read lookupdata.json: %w", err)
		}
		processMassLookupJob(job, recordType)
		perf.ElapsedSec = time.Since(start).Seconds()
		return nil
	},
}
var checkZonesCmdDEv = &cobra.Command{
	Use:   "check-zone-dev",
	Short: "Lookup DNS Record using massdns",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile("/root/dev/Cronicle/Plugins/go-plugin/NextDomain-Utils/build/check-zone.json")
		if err != nil {
			return fmt.Errorf("failed to read lookupdata.json: %w", err)
		}
		var job model.Job
		err = json.Unmarshal([]byte(data), &job)
		if err != nil {
			return fmt.Errorf("failed to read lookupdata.json: %w", err)
		}
		checkZones(job)
		return nil
	},
}

func checkZones(job model.Job) error {
	apikey := job.Params["apikey"].(string)
	server := job.Params["server"].(string)

	zones, err := powerdns.GetZonesPdnsAdmin(server, apikey)
	if err != nil {
		return err
	}

	activeZones := make([]response.GetZonesPdnsAdminResponse, 0)
	deactiveZones := make([]response.GetZonesPdnsAdminResponse, 0)

	for _, zone := range zones {
		if zone.Status == "Active" {
			activeZones = append(activeZones, zone)
		} else if zone.Status == "Deactive" {
			deactiveZones = append(deactiveZones, zone)
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := checkActiveZones(activeZones, job)
		if err != nil {
			// Handle the error (e.g., log it)
			fmt.Println("Error in checkActiveZones:", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := checkDeactiveZones(deactiveZones, job)
		if err != nil {
			// Handle the error (e.g., log it)
			fmt.Println("Error in checkDeactiveZones:", err)
		}
	}()

	wg.Wait()
	return nil
}
func checkActiveZones(activeZones []response.GetZonesPdnsAdminResponse, job model.Job) error {
	fileOp := files.NewFileOp()
	domains := "active-zone/domains.txt"
	resolvers := "active-zone/resolvers.txt"
	results := "active-zone/results.json"
	fileOp.CreateDir("./active-zone", 0755)
	fileOp.CreateFileWithMode(results, 0755)
	fileOp.CreateFileWithMode(domains, 0755)

	content := ""
	for _, zone := range activeZones {
		content += zone.Name + "\n"
	}
	fileOp.SaveFile(domains, content, 0755)
	fileOp.DownloadFile("https://raw.githubusercontent.com/KaySar12/Cronicle-go-Plugins/refs/heads/main/resolvers.txt", resolvers)
	var recordTypes = []string{"NS"}
	err := massdns.BulkLookup(domains, resolvers, results, recordTypes)
	if err != nil {
		fmt.Print(err)
		return err
	}
	queries, err := util.ParseDNSQueries(fileOp, results)
	if err != nil {
		fmt.Print(err)
		return err
	}
	activeZonesMap := make(map[string]response.GetZonesPdnsAdminResponse)
	for _, zone := range activeZones {
		activeZonesMap[fmt.Sprintf("%s.", zone.Name)] = zone
	}

	for _, query := range queries {
		zone := activeZonesMap[query.Name]

		if !checkValidQuery(query, job.Params["assign_zone"].(string), zone) {
			res, err := powerdns.ChangeStatus(job.Params["server"].(string), job.Params["apikey"].(string), zone.Name, "Deactive")
			if err != nil {
				fmt.Print(err)
				return err
			}
			fmt.Println(res)
		}
	}
	// fileOp.DeleteFile(results)
	// fileOp.DeleteFile(domains)
	// fileOp.DeleteFile(resolvers)
	//fileOp.DeleteDir("active-zone")
	return nil
}
func checkDeactiveZones(deactiveZones []response.GetZonesPdnsAdminResponse, job model.Job) error {
	fileOp := files.NewFileOp()
	domains := "deactive-zone/domains.txt"
	resolvers := "deactive-zone/resolvers.txt"
	results := "deactive-zone/results.json"
	fileOp.CreateDir("./deactive-zone", 0755)
	layout := "2006-01-02T15:04:05"
	fileOp.CreateFileWithMode(results, 0755)
	fileOp.CreateFileWithMode(domains, 0755)
	fileOp.WriteFile(domains, strings.NewReader("domains"), 0775)
	now := time.Now()
	for _, zone := range deactiveZones {
		deactivate_age, err := time.Parse(layout, zone.UpdateTimeDeactive)
		if err != nil {
			// Handle the error, e.g., log it or skip this zone
			continue
		}

		duration := now.Sub(deactivate_age)

		// Check if the duration is greater than or equal to a certain number of days
		if duration >= 24*time.Hour*7 { // Example: 7 days
			var req request.ZoneDelete
			req.ZoneId = zone.Name
			err := powerdns.DeleteZone(job.Params["server"].(string), job.Params["apikey"].(string), req)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	content := ""
	for _, zone := range deactiveZones {
		content += zone.Name + "\n"
	}
	fileOp.SaveFile(domains, content, 0755)
	fileOp.DownloadFile("https://raw.githubusercontent.com/KaySar12/Cronicle-go-Plugins/refs/heads/main/resolvers.txt", resolvers)
	var recordTypes = []string{"NS"}
	err := massdns.BulkLookup(domains, resolvers, results, recordTypes)
	if err != nil {
		fmt.Print(err)
		return err
	}
	queries, err := util.ParseDNSQueries(fileOp, results)
	if err != nil {
		fmt.Print(err)
		return err
	}
	activeZonesMap := make(map[string]response.GetZonesPdnsAdminResponse)
	for _, zone := range deactiveZones {
		activeZonesMap[fmt.Sprintf("%s.", zone.Name)] = zone
	}

	for _, query := range queries {
		zone := activeZonesMap[query.Name]

		if checkValidQuery(query, job.Params["assign_zone"].(string), zone) {
			res, err := powerdns.ChangeStatus(job.Params["server"].(string), job.Params["apikey"].(string), zone.Name, "Active")
			if err != nil {
				fmt.Print(err)
				return err
			}
			fmt.Println(res)
		}
	}
	// fileOp.DeleteFile(results)
	// fileOp.DeleteFile(domains)
	// fileOp.DeleteFile(resolvers)
	//fileOp.DeleteDir("deactive-zone")
	return nil
}
func checkValidQuery(query model.DNSQuery, assign_zone string, zone response.GetZonesPdnsAdminResponse) bool {
	if len(query.DataAnswers) == 0 {
		return false
	}
	validNs1 := fmt.Sprintf("ns1.%s.%s", zone.Account.Name, assign_zone)
	validNs2 := fmt.Sprintf("ns2.%s.%s", zone.Account.Name, assign_zone)
	valid1 := false
	valid2 := false
	for _, answer := range query.DataAnswers {
		if answer.Data == validNs1 {
			valid1 = true
		}
		if answer.Data == validNs2 {
			valid2 = true
		}
	}
	return valid1 && valid2
}

func processMassLookupJob(job model.Job, recordType string) error {
	domains := job.Params["domains"].(string)
	resolvers := job.Params["resolvers"].(string)
	results := job.Params["results"].(string)
	err := cmd.ExecCmd(fmt.Sprintf("massdns -r %s -t %s -t A %s > %s", resolvers, recordType, domains, results))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return nil
}
