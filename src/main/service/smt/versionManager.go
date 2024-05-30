package smt

import (
	"archive/zip"
	"bytes"
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
	// spiderFile 文件所在上级文件夹
	spiderFileDir := ""
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
		spiderFileDir = configPath + "\\spider\\"
		versionConfigurationDir = spiderVersionDir + "\\config.json"
		fmt.Println(runSpiderVersion)
	} else {
		runSpiderVersion = "/home/dav/runSpiderVersionHome/" + file.Filename
		spiderVersionDir = "/home/dav/runSpiderVersionHome/" + spiderVersion
		versionConfigurationDir = spiderVersionDir + "/config.json"
		spiderFileDir = "/home/dav/runSpiderVersionHome/spider/"
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

	// 保证版本只保留5个
	err = FileUtils.RetainLatestFolders(spiderFileDir, 5)
	if err != nil {
		errMessage := fmt.Sprintf("retain Latest Folders err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}

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

func DownLoadSpiderFile(context *gin.Context) *gin.Context {
	mainDir := ""
	if runtime.GOOS == "windows" {
		configPath, _ := filepath.Abs(".")
		mainDir = configPath
	} else {
		mainDir = "/home/dav/runSpiderVersionHome/"
	}
	// spiderFile 文件所在上级文件夹
	spiderFileDir := filepath.Join(mainDir, "\\spider\\")
	maxFolderName, err := FileUtils.GetMaxFolderName(spiderFileDir)
	if err != nil {
		errMessage := fmt.Sprintf("get max file name err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}
	filesDir := filepath.Join(spiderFileDir, maxFolderName)
	// 创建一个buffer来存放zip文件的内容
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	err = filepath.Walk(filesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			f, err := zipWriter.Create(info.Name())
			if err != nil {
				return err
			}
			_, err = f.Write(fileContent)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create zip file"})
		return context
	}
	// 关闭zip文件
	err = zipWriter.Close()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close zip writer"})
		return context
	}

	// 设置响应头
	context.Header("Content-Type", "application/octet-stream")
	context.Header("Content-Disposition", "attachment; filename=files.zip")
	context.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))

	// 写入响应
	context.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
	return context
}
