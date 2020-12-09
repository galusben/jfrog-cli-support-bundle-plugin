package commands

import "os"

const supportLogsUrl string = "https://supportlogs.jfrog.com/logs"
const supportLogsUrlEnvVarName string = "SUPPORT_LOGS_URL"

func getSupportLogsUrl() string {
	url := os.Getenv(supportLogsUrlEnvVarName)
	if url == "" {
		url = supportLogsUrl
	}
	return url
}
