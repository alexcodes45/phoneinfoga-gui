package dto

type ScanRequest struct {
	Number    string `json:"number"`
	Mode      string `json:"mode"` // auto|serve|cli
	Proxy     string `json:"proxy"`
	TimeoutMs int    `json:"timeoutMs"`
	CaseID    *int   `json:"caseId"`
}

type ScanResponse struct {
	JobID  string `json:"jobId"`
	Status string `json:"status"`
	Number string `json:"number"`
}
