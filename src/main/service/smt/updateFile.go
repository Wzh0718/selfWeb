package smt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"selfWeb/src/configuration"
	"selfWeb/src/tools/FileUtils"
	"selfWeb/src/tools/HttpsUtil/response"
	utils "selfWeb/src/tools/Unzip"
	"strings"
	"time"
)

// DownloadSpiderFile 解析zip --> 下载到指定的zip
func DownloadSpiderFile(context *gin.Context) *gin.Context {

	// 爬虫zip缓存数据地址
	runSpiderVersion := ""
	// 存储spiderVersion地址
	spiderVersionDir := ""
	// spider配置文件地址
	versionConfigurationDir := ""
	// 获取参数
	file, err := context.FormFile("data")

	// 解析版本号
	spiderVersion := context.PostForm("version")
	spiderVersion = strings.ReplaceAll(spiderVersion, ".", "_")

	if err != nil {
		errMessage := fmt.Sprintf("get form err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}

	if runtime.GOOS == "windows" {
		configPath, _ := filepath.Abs(".")
		runSpiderVersion = configPath + "\\spider\\" + file.Filename
		spiderVersionDir = configPath + "\\spider\\" + spiderVersion
		versionConfigurationDir = spiderVersionDir + "\\config.json"
		fmt.Println(runSpiderVersion)
	} else {
		runSpiderVersion = "/home/dav/runSpiderVersionHome/" + file.Filename
		spiderVersionDir = "/home/dav/runSpiderVersionHome/" + spiderVersion
		versionConfigurationDir = spiderVersionDir + "/config.json"
	}

	versionConfiguration, err := FileUtils.ReadVersionFile(versionConfigurationDir)
	if err != nil {
		errMessage := fmt.Sprintf("read version file err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}

	versionConfiguration.Version = spiderVersion
	versionConfiguration.UpdateDate = time.Now().String()
	versionConfiguration.Id = versionConfiguration.Id + 1

	// 缓存文件
	if err := context.SaveUploadedFile(file, runSpiderVersion); err != nil {
		errMessage := fmt.Sprintf("save file err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}

	z := &utils.ZipUtil{}
	zipErr := z.UnzipFile(runSpiderVersion, spiderVersionDir)
	if zipErr != nil {
		errMessage := fmt.Sprintf("unzip file is err: %s", zipErr.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}

	// After successfully unzipping, delete the temporary zip file
	err = os.Remove(runSpiderVersion)
	if err != nil {
		errMessage := fmt.Sprintf("delete file err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}

	if err := FileUtils.SaveAsJSON(versionConfiguration, versionConfigurationDir); err != nil {
		errMessage := fmt.Sprintf("save file err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}

	configuration.Logger.Info("VersionConfiguration saved to " + versionConfigurationDir)

	body := response.ReturnContextNoBody(context, http.StatusOK, "成功")
	return body
}

func GetSpiderVersion(context *gin.Context) *gin.Context {
	spiderFileDir := ""
	if runtime.GOOS == "windows" {
		configPath, _ := filepath.Abs(".")
		spiderFileDir = configPath + "\\spider\\"
	} else {
		spiderFileDir = "/home/dav/runSpiderVersionHome/spider/"
	}

	maxFolderName, err := FileUtils.GetMaxFolderName(spiderFileDir)
	if err != nil {
		errMessage := fmt.Sprintf("get max file name err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}

	versionConfiguration, err := FileUtils.ReadVersionFile(spiderFileDir + maxFolderName + "\\config.json")
	if err != nil {
		errMessage := fmt.Sprintf("read version file err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}
	returnContext := response.ReturnContext(context, http.StatusOK, "成功", versionConfiguration)
	return returnContext
}
