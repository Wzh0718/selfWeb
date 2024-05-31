package smt

import (
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"selfWeb/src/configuration"
	"selfWeb/src/tools/Cache"
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
	if err != nil {
		errMessage := fmt.Sprintf("get form err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}

	// 解析版本号
	spiderVersion := context.PostForm("version")
	spiderVersion = strings.ReplaceAll(spiderVersion, ".", "_")

	if runtime.GOOS == "windows" {
		configPath, _ := filepath.Abs(".")
		runSpiderVersion = configPath + "\\spider\\" + file.Filename
		spiderVersionDir = configPath + "\\spider\\" + spiderVersion
		spiderFileDir = configPath + "\\spider\\"
		versionConfigurationDir = spiderVersionDir + "\\config.json"
		fmt.Println(runSpiderVersion)
	} else {
		spiderFileDir = "/home/dav/runSpiderVersionHome/spider/"
		runSpiderVersion = spiderFileDir + file.Filename
		spiderVersionDir = spiderFileDir + spiderVersion
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
	configDir := filepath.Join(spiderFileDir, maxFolderName, "config.json")
	versionConfiguration, err := FileUtils.ReadVersionFile(configDir)
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
	spiderFileDir := filepath.Join(mainDir, "spider")
	maxFolderName, err := FileUtils.GetMaxFolderName(spiderFileDir)
	if err != nil {
		errMessage := fmt.Sprintf("get max file name err: %s", err.Error())
		body := response.ReturnContextNoBody(context, http.StatusBadRequest, errMessage)
		configuration.Logger.Error(errMessage)
		return body
	}
	filesDir := filepath.Join(spiderFileDir, maxFolderName)

	files, err := GetCachedOrNewFiles(filesDir)
	if err != nil {
		return nil
	}
	return sendZipFile(context, files)
}

// GetCachedOrNewFiles 判断cache里面是否有指定数据，如果没有则添加
func GetCachedOrNewFiles(filesDir string) (map[string][]byte, error) {
	cacheKey := filesDir

	filesData, ok := Cache.GetFromCache(cacheKey)
	if ok {
		return filesData, nil
	}

	filesData, cacheErr := Cache.CacheFilesInMemory(filesDir)
	if cacheErr != nil {
		return nil, fmt.Errorf("failed to cache files in memory: %v", cacheErr)
	}

	Cache.AddToCache(cacheKey, filesData)

	return filesData, nil
}

func sendZipFile(context *gin.Context, filesData map[string][]byte) *gin.Context {
	context.Header("Content-Type", "application/zip")
	context.Header("Content-Disposition", "attachment; filename=files.zip")

	zipWriter := zip.NewWriter(context.Writer)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			return
		}
	}(zipWriter)

	for fileName, fileContent := range filesData {
		f, err := zipWriter.Create(fileName)
		// 防止重复写入头信息和状态码
		if err != nil {
			if !context.Writer.Written() {
				message := fmt.Sprintf("Failed to create zip entry for file %s: %v", fileName, err)
				response.ReturnContextNoBody(context, http.StatusInternalServerError, message)
				configuration.Logger.Error(message)
			}

			return context
		}
		_, err = f.Write(fileContent)
		if err != nil {
			// 防止重复写入头信息和状态码
			if !context.Writer.Written() {
				message := fmt.Sprintf("Failed to write file %s to zip: %v", fileName, err)
				response.ReturnContextNoBody(context, http.StatusInternalServerError, message)
				configuration.Logger.Error(message)
			}
			return context
		}
	}
	return context
}
