package service

type CronicleService struct{}

type ICronicleService interface {
	GetSchedule(server string, apikey string, offset int64, limit int64)
	GetEvent(server string, apikey string, id string, title string)
	CreateEvent(server string, apikey string)
	UpdateEvent(server string, apikey string)
	DeleteEvent(server string, apikey string)
	GetEventHistory(server string, apikey string)
	GetHistory(server string, apikey string)
	RunEvent(server string, apikey string)
	GetJobStatus(server string, apikey string)
	GetActiveJob(server string, apikey string)
	UpdateJob(server string, apikey string)
	AbortJob(server string, apikey string)
	GetMasterState(server string, apikey string)
	UpdateMasterState(server string, apikey string)
}
