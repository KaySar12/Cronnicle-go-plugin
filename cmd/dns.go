package cmd

import (
	"NextDomain-Utils/dto/request"
	model "NextDomain-Utils/model"
	"NextDomain-Utils/service"
	"NextDomain-Utils/utils/cronicle"
	"NextDomain-Utils/utils/dnsclient"
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var (
	powerdns = service.NewIPowerDNSService()
)

func init() {
	RootCmd.AddCommand(dnsCmd)
	dnsCmd.AddCommand(lookupCmd)
	dnsCmd.AddCommand(lookupCmdDev)
}

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "dns related commands",
}

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Lookup DNS record",
	RunE: func(cmd *cobra.Command, args []string) error {
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
			processLookupJob(job)
		}
		return nil
	},
}
var lookupCmdDev = &cobra.Command{
	Use:   "lookup-dev",
	Short: "Lookup zone type ns",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		var perf model.PerfStats
		var table = model.Table{
			Title:  "NextDomain lookup stats",
			Header: []string{"DNS Lookup", "NS Records", "IP Address"},
			Rows: [][]interface{}{
				{"unknown", "unknown", "unknown"},
			},
		}
		data, err := os.ReadFile("/root/dev/Cronicle/Plugins/go-plugin/NextDomain-Utils/build/lookupdata.json")
		if err != nil {
			return fmt.Errorf("failed to read lookupdata.json: %w", err)
		}
		// Parse JSON input
		var job model.Job
		err = json.Unmarshal([]byte(data), &job)
		if err != nil {
			slog.Error(fmt.Sprintf("Error parsing 'job' input: %v", err))
			perf.ElapsedSec = time.Since(start).Seconds()
			cronicle.Report(job, "Error", perf, table, err) // Report error and exit early continue
		}
		processLookupJob(job)
		return nil
	},
}

func processLookupJob(job model.Job) {
	var perf model.PerfStats
	var table = model.Table{
		Title:  "NextDomain lookup stats",
		Header: []string{"DNS Lookup", "NS Records", "IP Address"},
		Rows: [][]interface{}{
			{job.Params["zone"], "unknown", "unknown"},
		},
	}
	duration := float64(job.Timeout)
	start := time.Now()
	server := job.Params["server"].(string)
	server_id := job.Params["server_id"].(string)
	apikey := job.Params["apikey"].(string)
	// Use ParseInt instead of Atoi for better error handling and to allow different bases
	maxAttempts, err := strconv.ParseInt(job.Params["retry"].(string), 10, 64)
	if err != nil {
		slog.Error(fmt.Sprintf("Error parsing 'retry' param: %v", err))
		perf.ElapsedSec = time.Since(start).Seconds()
		cronicle.Report(job, "Error", perf, table, err) // Report error and exit early
		return
	}

	retryDelay, err := strconv.ParseInt(job.Params["retry_delay"].(string), 10, 64)
	if err != nil {
		slog.Error(fmt.Sprintf("Error parsing 'retry_delay' param: %v", err))
		perf.ElapsedSec = time.Since(start).Seconds()
		cronicle.Report(job, "Error", perf, table, err) // Report error and exit early
		return
	}
	interval := time.Duration(retryDelay) * time.Second
	slog.Info(fmt.Sprintf("Duration : %f", duration))
	slog.Info(fmt.Sprintf("Max Attempts : %d", maxAttempts))
	slog.Info(fmt.Sprintf("Domain to check : %s", job.Params["zone"]))
	slog.Info(fmt.Sprintf("Assign Zone : %s", job.Params["assign_zone"]))
	attempt := 1
	elapsed := 0.0
	timer := time.NewTicker(150 * time.Millisecond)
	slog.Info(fmt.Sprintf("Begin to check %s......", job.Params["zone"]))
	for elapsed < duration {
		select {
		case <-timer.C:
			elapsed = time.Since(start).Seconds()
			progress := math.Min(elapsed/float64(duration), 1.0)
			cronicle.OutputJSON(model.Response{Progress: progress})

			if lookup(job, &table) {
				var request request.ZoneChangeStatus
				request.Name = job.Params["zone"].(string)
				response, err := powerdns.ChangeStatus(server, apikey, job.Params["zone"].(string), "Active")
				if err != nil {
					slog.Error(fmt.Sprintf("Error changing zone status: %v", err))
					perf.ElapsedSec = time.Since(start).Seconds()
					cronicle.Report(job, "Error", perf, table, err)
					return // Exit on error
				}
				slog.Info(response.Status)
				perf.ElapsedSec = time.Since(start).Seconds()
				cronicle.Report(job, "Success", perf, table)
				return // Exit on success
			}
			if int64(attempt) <= maxAttempts {
				slog.Info(fmt.Sprintf("Progress : %f percent", progress*100))
				slog.Info(fmt.Sprintf("Current Attempt : %d ", attempt))
				attempt++
				time.Sleep(interval)
			} else {
				slog.Info("All attempts failed, returning expired")
				// Delete Zone
				var request request.ZoneDelete
				request.ZoneId = job.Params["zone"].(string)
				slog.Info(fmt.Sprintf("Deleting Zone %s", job.Params["zone"].(string)))
				err := powerdns.DeleteZone(server, apikey, request, server_id)
				if err != nil {
					slog.Error(fmt.Sprintf("Problem with deleting zone %v", err))
					perf.ElapsedSec = time.Since(start).Seconds()
					cronicle.Report(job, "Error", perf, table, err)
					return
				}
				perf.ElapsedSec = time.Since(start).Seconds()
				cronicle.Report(job, "Expire", perf, table)
				return // Exit if time is up
			}
		}
	}
}

func lookup(job model.Job, table *model.Table) bool {
	// Access job.Params safely with type checks and error handling
	zone, ok := job.Params["zone"].(string)
	if !ok {
		slog.Error("zone parameter is missing or not a string")
		return false
	}
	assignZone, ok := job.Params["assign_zone"].(string)
	if !ok {
		slog.Error("assign_zone parameter is missing or not a string")
		return false
	}
	account, ok := job.Params["account"].(string)
	if !ok {
		slog.Error("account parameter is missing or not a string")
		return false
	}
	zone = dnsclient.FormatZone(zone)
	assignZone = dnsclient.FormatZone(assignZone)
	// Load the system's default DNS server from resolv.conf
	config, err := dnsclient.LoadDnsServer()
	if err != nil {
	}
	response, err := dnsclient.Lookup(config, dns.TypeNS, zone, 10)
	if err != nil {
	}
	var nsrecords []string
	for _, answer := range response.Answer {
		if nsRecord, ok := answer.(*dns.NS); ok {
			nsrecords = append(nsrecords, fmt.Sprintf("%s.%s", account, assignZone))
			if !strings.HasSuffix(nsRecord.Ns, assignZone) {
				fmt.Printf("Found NS record: %s different with assign Zone: %s\n", nsRecord.Ns, fmt.Sprintf("%s.%s", account, assignZone))
				return false
			}
		}
	}

	if len(nsrecords) == 0 {
		slog.Error("No NS records found for the zone")
		return false
	} else {
		table.Rows[0][1] = fmt.Sprintf("%v", nsrecords)
		slog.Info(fmt.Sprintf("NS Records for the zone: %v", nsrecords))
	}

	// --- Get A records (IP addresses) for the zone ---
	response, err = dnsclient.Lookup(config, dns.TypeA, zone, 10)
	if err != nil {

	}

	var ipAddresses []string
	for _, answer := range response.Answer {
		if aRecord, ok := answer.(*dns.A); ok {
			ipAddresses = append(ipAddresses, aRecord.A.String())
		}
	}

	if len(ipAddresses) == 0 {
		slog.Warn("No IP addresses (A records) found for the zone") // Log as a warning
	} else {
		slog.Info(fmt.Sprintf("IP addresses for the zone: %v", ipAddresses))
		table.Rows[0][2] = fmt.Sprintf("%v", ipAddresses)
	}

	slog.Info("Zone is Valid")
	return true
}
