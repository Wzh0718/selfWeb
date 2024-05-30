package FileUtils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"selfWeb/src/configuration"
	"selfWeb/src/configuration/structs"
	"sort"
	"strconv"
	"strings"
)

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ReadVersionFile(path string) (*structs.VersionConfiguration, error) {
	var versionConfiguration structs.VersionConfiguration
	// 文件不存在的情况下
	if _, err := os.Stat(path); os.IsNotExist(err) {
		message := fmt.Sprintf("%s File does not exist...", path)
		configuration.Logger.Info(message)
		return &versionConfiguration, nil
	}
	message := fmt.Sprintf("File exists, reading file %s", path)
	configuration.Logger.Info(message)

	data, err := os.ReadFile(path)
	if err != nil {
		message := fmt.Sprintf("Error reading file: %s", err.Error())
		configuration.Logger.Error(message)
		return &versionConfiguration, err
	}

	if err := json.Unmarshal(data, &versionConfiguration); err != nil {
		message := fmt.Sprintf("Error parsing file: %s", err.Error())
		configuration.Logger.Error(message)
		return &versionConfiguration, err
	}
	return &versionConfiguration, nil
}

// SaveAsJSON 将 VersionConfiguration 实例保存到指定的文件路径
func SaveAsJSON(vc *structs.VersionConfiguration, filePath string) error {
	// 将结构体序列化为 JSON
	data, err := json.MarshalIndent(vc, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %v", err)
	}

	// 创建/打开文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	// 写入数据到文件
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err.Error())
	}

	return nil
}

// parseFolderName 将文件夹名解析为整数切片
func parseFolderName(name string) ([]int, error) {
	parts := strings.Split(name, "_")
	numbers := make([]int, len(parts))
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		numbers[i] = num
	}
	return numbers, nil
}

// byFolderName 实现 sort.Interface，用于排序
type byFolderName []string

func (f byFolderName) Len() int {
	return len(f)
}

func (f byFolderName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f byFolderName) Less(i, j int) bool {
	a, err := parseFolderName(f[i])
	if err != nil {
		return false
	}
	b, err := parseFolderName(f[j])
	if err != nil {
		return false
	}
	for k := 0; k < len(a) && k < len(b); k++ {
		if a[k] != b[k] {
			return a[k] < b[k]
		}
	}
	return len(a) < len(b)
}

// GetMaxFolderName 获取最大的文件夹名
func GetMaxFolderName(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	var folders []string
	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}

	sort.Sort(sort.Reverse(byFolderName(folders)))

	if len(folders) > 0 {
		return folders[0], nil
	}
	return "", fmt.Errorf("no folders found")
}
