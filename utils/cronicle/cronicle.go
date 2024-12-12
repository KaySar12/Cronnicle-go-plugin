package cronicle

import (
	"NextDomain-Utils/model"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func Report(job model.Job, status string, perf interface{}, table model.Table, err ...error) {

	switch status {
	case "Success":
		var html = model.HTMLReport{
			Title:   "Report NextDomain lookup Result",
			Content: "",
			Caption: "",
		}
		OutputJSON(model.Response{
			Complete:    1,
			Perf:        perf,
			Table:       table,
			Code:        0,
			Description: "Success!",
			HTML:        html,
		})
		return
	case "Expire":
		var html = model.HTMLReport{
			Title:   "Report NextDomain lookup Result",
			Content: fmt.Sprintf("<pre>%s is not valid domain managed by NextDomain</pre>", job.Params["zone"].(string)),
			Caption: "",
		}
		OutputJSON(model.Response{
			Complete:    1,
			Perf:        perf,
			Code:        999,
			Table:       table,
			Description: "Expire!",
			HTML:        html,
		})
		return
	case "Error":
		var html = model.HTMLReport{
			Title:   "Report NextDomain lookup Result",
			Content: fmt.Sprintf("<pre>Error Report:%s</pre>", err),
			Caption: "",
		}
		OutputJSON(model.Response{
			Complete:    1,
			Perf:        perf,
			Table:       table,
			Code:        500,
			Description: "Error!",
			HTML:        html,
		})
	}

}

func OutputJSON(response model.Response) {
	encoder := json.NewEncoder(os.Stdout)
	err := encoder.Encode(response)
	if err != nil {
		log.Printf("Failed to encode response: %v\n", err)
	}
}
