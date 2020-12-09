package commands

import (
	"encoding/json"
	"errors"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io"
	"net/http"
	"os"
	"strconv"
)

func GetGenerateSupportBundleCommand() components.Command {
	return components.Command{
		Name:        "generate",
		Description: "Generates support bundle to supportlogs.",
		Aliases:     []string{"up"},
		Arguments:   getGenerateArguments(),
		Flags:       getGenerateFlags(),
		EnvVars:     getGenerateEnvVar(),
		Action: func(c *components.Context) error {
			return generateCmd(c)
		},
	}
}

func getGenerateArguments() []components.Argument {
	return []components.Argument{}
}

func getGenerateFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:        "server-id",
			Description: "Artifactory server ID configured using the config command.",
		},
		components.BoolFlag{
			Name:         "send-to-support",
			Description:  "Rather to upload the support bundle to JFrog support or not",
			DefaultValue: false,
		},
		components.StringFlag{
			Name:        "ticket",
			Description: "Ticket identifier for JFrog support team - must be provided when send-to-support = true",
		},
		components.StringFlag{
			Name:        "name",
			Description: "Support bundle name - when empty will be auto generated",
		},
		components.StringFlag{
			Name:        "description",
			Description: "Support bundle description",
		},
		components.BoolFlag{
			Name:         "config",
			Description:  "Include service configuration",
			DefaultValue: true,
		},
		components.BoolFlag{
			Name:         "system",
			Description:  "Include service system information",
			DefaultValue: true,
		},
		components.BoolFlag{
			Name:        "logs",
			Description: "Include logs",
		},
		components.BoolFlag{
			Name:        "dumps",
			Description: "Include thread dumps",
		},
		components.StringFlag{
			Name:        "dumps-count",
			Description: "number of times to collect thread dump. Default:1",
		},
		components.StringFlag{
			Name:        "dumps-interval",
			Description: "Interval between times of collection in milliseconds. Default:0",
		},
		components.StringFlag{
			Name:        "start",
			Description: "start date from which to fetch the logs. pattern: YYYY-MM-DD",
		},
		components.StringFlag{
			Name:        "end",
			Description: "end date until which to fetch the logs. pattern: YYYY-MM-DD",
		},
	}
}

func getGenerateEnvVar() []components.EnvVar {
	return []components.EnvVar{
		{
			Name:        supportLogsUrlEnvVarName,
			Default:     supportLogsUrl,
			Description: "Support logs base url - mostly for debug",
		},
	}
}

type generateConfig struct {
	sendToSupportlogs bool
	ticket            string
}

func generateCmd(c *components.Context) error {
	var conf = new(generateConfig)
	conf.sendToSupportlogs = c.GetBoolFlagValue("send-to-support")
	conf.ticket = c.GetStringFlagValue("ticket")
	if conf.sendToSupportlogs && conf.ticket == "" {
		return errors.New("when providing send-to-support you must provide ticket")
	}
	rtDetails, err := getRtDetails(c)
	if err != nil {
		return err
	}
	auth, err := rtDetails.CreateArtAuthConfig()
	if err != nil {
		return err
	}
	rtClient, err := httpclient.ArtifactoryClientBuilder().
		SetCertificatesPath(auth.GetClientCertKeyPath()).
		SetInsecureTls(rtDetails.InsecureTls).SetServiceDetails(&auth).
		Build()
	if err != nil {
		return err
	}
	model, err := getCreateBundleModel(c)
	if err != nil {
		return err
	}
	clientDetails := auth.CreateHttpClientDetails()
	responseModel, err := createBundle(clientDetails, rtDetails, rtClient, model)
	if err != nil {
		return err
	}
	log.Info("Downloading bundle")
	resBody, err := downloadSupportBundle(responseModel, rtDetails, rtClient, clientDetails)
	if err != nil {
		return err
	}
	defer resBody.Close()
	filename := responseModel.Id + "_" + responseModel.Artifactory.ServiceId + ".zip"
	log.Info("Writing bundle to file:", filename)
	err = writeBundleToFile(err, filename, resBody)
	if err != nil {
		return err
	}
	if conf.sendToSupportlogs {
		uploadConf := &bundleConfig{
			filepath:           filename,
			ticket:             conf.ticket,
			baseSupportLogsUrl: getSupportLogsUrl(),
		}
		return doUpload(uploadConf)
	}
	return nil
}

func getCreateBundleModel(c *components.Context) (*supportBundleApiModel, error) {
	dumpsInterval := 0
	dumpsCount := 1
	var err error
	if c.GetStringFlagValue("dumps-interval") != "" {
		dumpsInterval, err = strconv.Atoi(c.GetStringFlagValue("dumps-interval"))
		if err != nil {
			return nil, err
		}
	}
	if c.GetStringFlagValue("dumps-count") != "" {
		dumpsCount, err = strconv.Atoi(c.GetStringFlagValue("dumps-count"))
		if err != nil {
			return nil, err
		}
	}
	model := supportBundleApiModel{
		Name:        c.GetStringFlagValue("name"),
		Description: c.GetStringFlagValue("description"),
		Parameters: supportBundleParametersApiModel{
			Configuration: c.GetBoolFlagValue("config"),
			System:        c.GetBoolFlagValue("system"),
			Logs: supportBundleParametersLogs{
				Include:   c.GetBoolFlagValue("logs"),
				StartDate: c.GetStringFlagValue("start"),
				EndDate:   c.GetStringFlagValue("end"),
			},
			ThreadDump: supportBundleParametersThreadDump{
				Count:    dumpsInterval,
				Interval: dumpsCount,
			},
		},
	}
	return &model, nil
}

func writeBundleToFile(err error, fileName string, bundleReader io.ReadCloser) error {
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, bundleReader)
	if err != nil {
		return err
	}
	return nil
}

func downloadSupportBundle(responseModel supportBundleResponseModel, rtDetails *config.ArtifactoryDetails, rtClient *httpclient.ArtifactoryHttpClient, clientDetails httputils.HttpClientDetails) (io.ReadCloser, error) {
	bundlePath := "api/system/support/bundle/" + responseModel.Id + "/archive"
	url, err := utils.BuildArtifactoryUrl(rtDetails.GetUrl(), bundlePath, nil)
	if err != nil {
		return nil, err
	}
	res, _, _, err := rtClient.Send(http.MethodGet, url, nil, true, false, &clientDetails)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		defer res.Body.Close()
		log.Error("error downloading release bundle, status code: ", res.StatusCode, "url:", url)
		return nil, errors.New("error downloading release bundle")
	}
	return res.Body, nil
}

func createBundle(clientDetails httputils.HttpClientDetails, rtDetails *config.ArtifactoryDetails, rtClient *httpclient.ArtifactoryHttpClient, msg *supportBundleApiModel) (supportBundleResponseModel, error) {
	clientDetails.Headers["Content-Type"] = "application/json"
	url, err := utils.BuildArtifactoryUrl(rtDetails.GetUrl(), "api/system/support/bundle", nil)
	if err != nil {
		return supportBundleResponseModel{}, err
	}
	postBody, err := json.Marshal(msg)
	if err != nil {
		return supportBundleResponseModel{}, err
	}
	res, respBody, err := rtClient.SendPost(url, postBody, &clientDetails)
	if err != nil {
		return supportBundleResponseModel{}, err
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		log.Error("error creating release bundle, status code ", res.StatusCode)
		return supportBundleResponseModel{}, errors.New("error creating release bundle")
	}
	log.Output(string(respBody))
	log.Info("Created bundle")

	responseModel := supportBundleResponseModel{}
	err = json.Unmarshal(respBody, &responseModel)
	if err != nil {
		return supportBundleResponseModel{}, err
	}
	return responseModel, nil
}

func getRtDetails(c *components.Context) (*config.ArtifactoryDetails, error) {
	serverId := c.GetStringFlagValue("server-id")
	details, err := commands.GetConfig(serverId, false)
	if err != nil {
		return nil, err
	}
	if details.Url == "" {
		return nil, errors.New("no server-id was found, or the server-id has no url")
	}
	details.Url = clientutils.AddTrailingSlashIfNeeded(details.Url)
	err = config.CreateInitialRefreshableTokensIfNeeded(details)
	if err != nil {
		return nil, err
	}
	return details, nil
}
