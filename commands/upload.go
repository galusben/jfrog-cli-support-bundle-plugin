package commands

import (
	"errors"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

func GetUploadSupportBundleCommand() components.Command {
	return components.Command{
		Name:        "upload",
		Description: "Uploads support bundle to supportlogs.",
		Aliases:     []string{"up"},
		Arguments:   getUploadArguments(),
		Flags:       getUploadFlags(),
		EnvVars:     getUploadEnvVar(),
		Action: func(c *components.Context) error {
			return uploadCmd(c)
		},
	}
}

func getUploadArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "filepath",
			Description: "Bundle path on the local file system",
		},
		{
			Name:        "ticket",
			Description: "Ticket number",
		},
	}
}

func getUploadFlags() []components.Flag {
	return []components.Flag{}
}

func getUploadEnvVar() []components.EnvVar {
	return []components.EnvVar{
		{
			Name:        supportLogsUrlEnvVarName,
			Default:     supportLogsUrl,
			Description: "Support logs base url - mostly for debug",
		},
	}
}

type bundleConfig struct {
	filepath           string
	ticket             string
	baseSupportLogsUrl string
}

func uploadCmd(c *components.Context) error {
	if len(c.Arguments) != 2 {
		return errors.New("Wrong number of arguments. Expected: 2, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}
	var conf = new(bundleConfig)
	conf.baseSupportLogsUrl = getSupportLogsUrl()
	conf.filepath = c.Arguments[0]
	conf.ticket = c.Arguments[1]
	err := doUpload(conf)
	if err != nil {
		return err
	}
	return nil
}

func doUpload(conf *bundleConfig) error {
	uploadUrl, err := url.Parse(conf.baseSupportLogsUrl + "/" + conf.ticket + "/" + path.Base(conf.filepath))
	if err != nil {
		return err
	}
	bundle, err := os.Open(conf.filepath)
	if err != nil {
		return err
	}
	defer bundle.Close()
	request := &http.Request{
		Method: http.MethodPut,
		Body:   bundle,
		URL:    uploadUrl,
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		log.Error("error uploading, status code ", res.StatusCode)
		return errors.New("error uploading")
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Output(string(bodyBytes))
	log.Info("sent bundle")
	return nil
}
