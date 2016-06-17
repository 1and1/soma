package proto

type Job struct {
	Id           string      `json:"id,omitempty"`
	Status       string      `json:"status,omitempty"`
	Result       string      `json:"result,omitempty"`
	Type         string      `json:"type,omitempty"`
	Serial       int         `json:"serial,omitempty"`
	RepositoryId string      `json:"repositoryId,omitempty"`
	UserId       string      `json:"userId,omitempty"`
	TeamId       string      `json:"teamId,omitempty"`
	TsQueued     string      `json:"queued,omitempty"`
	TsStarted    string      `json:"started,omitempty"`
	TsFinished   string      `json:"finished,omitempty"`
	Error        string      `json:"error,omitempty"`
	Details      *JobDetails `json:"details,omitempty"`
}

type JobFilter struct {
	User   string   `json:"user,omitempty"`
	Team   string   `json:"team,omitempty"`
	Status string   `json:"status,omitempty"`
	Result string   `json:"result,omitempty"`
	Since  string   `json:"since,omitempty"`
	IdList []string `json:"idlist,omitempty"`
}

type JobDetails struct {
	CreatedAt     string `json:"createdAt,omitempty"`
	CreatedBy     string `json:"createdBy,omitempty"`
	Specification string `json:"specification,omitempty"`
}

func NewJobFilter() Request {
	return Request{
		Flags: &Flags{},
		Filter: &Filter{
			Job: &JobFilter{},
		},
	}
}

func NewJobResult() Result {
	return Result{
		Errors: &[]string{},
		Jobs:   &[]Job{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
