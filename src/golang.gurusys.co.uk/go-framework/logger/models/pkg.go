package models

type LogAppDef struct {
	Status       string
	Appname      string
	Repository   string
	Groupname    string
	Namespace    string
	DeploymentID string
	StartupID    string
}

type LogLine struct {
	Time int64
	Line string
}

type LogRequest struct {
	AppDef *LogAppDef
	Lines  []*LogLine
}

type LogResponse struct{}

type LogFilter struct {
	Host     string
	UserName string
	AppDef   *AppDef
}

type GetLogRequest struct {
	LogFilter    []*LogFilter
	MinimumLogID int64
}

type LogEntry struct {
	ID       uint64
	Host     string
	UserName string
	Occured  uint64
	AppDef   *LogAppDef
	Line     string
}

type GetLogResponse struct {
	Entries []*LogEntry
}
