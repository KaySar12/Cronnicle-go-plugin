package cronicle

import (
	cron "NextDomain-Utils/constant"
	"NextDomain-Utils/model"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func Report(job model.Job, status string, args ...any) {
	var response model.Response
	var err error

	// Handle arguments
	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			arg := args[i]
			switch v := arg.(type) {
			case model.PerfStats:
				response.Perf = v
			case model.Table:
				response.Table = v
			default:
				if err == nil {
					err = v.(error)
				}
			}
		}
	}

	// Handle status cases
	switch status {
	case "Success":
		response.HTML = model.HTMLReport{
			Title:   "Report MassDNS CheckZone Result",
			Content: "",
			Caption: "",
		}
		response.Code = cron.SUCCESS
		response.Description = ""
		response.Complete = 1
	case "Expire":
		response.HTML = model.HTMLReport{
			Title:   "Report MassDNS CheckZone Result",
			Content: fmt.Sprintf("<pre>%s is not valid domain managed by NextDomain</pre>", job.Params["zone"].(string)),
			Caption: "",
		}
		response.Code = cron.EXPIRE
		response.Description = "Expire!"
		response.Complete = 1
	case "Error":
		response.HTML = model.HTMLReport{
			Title:   "Report MassDNS CheckZone Result",
			Content: fmt.Sprintf("<pre>Error Report:%s</pre>", err),
			Caption: "",
		}
		response.Code = cron.ERROR
		response.Description = "Error!"
		response.Complete = 1
	}

	OutputJSON(response)
}
func OutputJSON(response model.Response) {
	encoder := json.NewEncoder(os.Stdout)
	err := encoder.Encode(response)
	if err != nil {
		log.Printf("Failed to encode response: %v\n", err)
	}
}
