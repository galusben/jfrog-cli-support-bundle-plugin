package commands

type supportBundleResponseModel struct {
	Id          string                                `json:"id,omitempty"`
	Artifactory artifactorySupportBundleResponseModel `json:"artifactory,omitempty"`
}

type artifactorySupportBundleResponseModel struct {
	BundleUrl string `json:"bundle_url,omitempty"`
	ServiceId string `json:"service_id,omitempty"`
}

type supportBundleApiModel struct {
	Name        string                          `json:"name,omitempty"`
	Description string                          `json:"description,omitempty"`
	Parameters  supportBundleParametersApiModel `json:"parameters,omitempty"`
}

type supportBundleParametersApiModel struct {
	Configuration bool                              `json:"configuration,omitempty"`
	System        bool                              `json:"system,omitempty"`
	Logs          supportBundleParametersLogs       `json:"logs,omitempty"`
	ThreadDump    supportBundleParametersThreadDump `json:"thread_dump,omitempty"`
}

type supportBundleParametersLogs struct {
	Include   bool   `json:"include,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

type supportBundleParametersThreadDump struct {
	Count    int `json:"count,omitempty"`
	Interval int `json:"interval,omitempty"`
}
