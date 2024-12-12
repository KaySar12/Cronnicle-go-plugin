package model

type Job struct {
	Params        map[string]interface{} `json:"params"`
	Timeout       int                    `json:"timeout"`
	CatchUp       int                    `json:"catch_up"`
	QueueMax      int                    `json:"queue_max"`
	Timezone      string                 `json:"timezone"`
	Category      string                 `json:"category"`
	Plugin        string                 `json:"plugin"`
	Target        string                 `json:"target"`
	Algo          string                 `json:"algo"`
	Multiplex     int                    `json:"multiplex"`
	Stagger       int                    `json:"stagger"`
	Retries       int                    `json:"retries"`
	RetryDelay    int                    `json:"retry_delay"`
	Detached      int                    `json:"detached"`
	Queue         int                    `json:"queue"`
	Chain         string                 `json:"chain"`
	ChainError    string                 `json:"chain_error"`
	NotifySuccess string                 `json:"notify_success"`
	NotifyFail    string                 `json:"notify_fail"`
	WebHook       string                 `json:"web_hook"`
	CPULimit      int                    `json:"cpu_limit"`
	CPUSustain    int                    `json:"cpu_sustain"`
	MemoryLimit   int                    `json:"memory_limit"`
	MemorySustain int                    `json:"memory_sustain"`
	LogMaxSize    int                    `json:"log_max_size"`
	Notes         string                 `json:"notes"`
	CategoryTitle string                 `json:"category_title"`
	GroupTitle    string                 `json:"group_title"`
	PluginTitle   string                 `json:"plugin_title"`
	Now           int                    `json:"now"`
	Source        string                 `json:"source"`
	ID            string                 `json:"id"`
	TimeStart     float64                `json:"time_start"`
	Hostname      string                 `json:"hostname"`
	Event         string                 `json:"event"`
	EventTitle    string                 `json:"event_title"`
	NiceTarget    string                 `json:"nice_target"`
	Command       string                 `json:"command"`
	LogFile       string                 `json:"log_file"`
	Pid           int                    `json:"pid"`
}

type Response struct {
	Complete    int         `json:"complete"`
	Code        int         `json:"code"`
	Description string      `json:"description"`
	Perf        interface{} `json:"perf,omitempty"` // Now dynamic
	Table       Table       `json:"table,omitempty"`
	HTML        HTMLReport  `json:"html,omitempty"`
	Progress    float64     `json:"progress,omitempty"`
}

type PerfStats struct {
	ElapsedSec float64 `json:"elapsed_sec"`
}
type Table struct {
	Title   string          `json:"title"`
	Header  []string        `json:"header"`
	Rows    [][]interface{} `json:"rows"`
	Caption string          `json:"caption"`
}
type HTMLReport struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Caption string `json:"caption"`
}
