package script

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"git.xenonstack.com/util/continuous-security-backend/config"
)

func DownloadScripts() bool {

	// download charts from zip
	req, err := http.NewRequest("GET", config.Conf.Service.RepoURL, nil)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println(req.URL)
	req.Header.Set("PRIVATE-TOKEN", config.Conf.Service.RepoPrivateKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		return false
	}

	// create tmp directory
	tmpDir := "tmp" + randomString(10)
	tmpDirPath := "./" + tmpDir
	err = os.MkdirAll(tmpDirPath+"/ws_charts", 0777)
	if err != nil {
		log.Println(err)
		return false
	}
	defer func() {
		os.RemoveAll(tmpDirPath)
	}()

	// create emtpy zip file
	f, err := os.Create(tmpDirPath + "/ws_chart.zip")
	if err != nil {
		log.Println(err)
		return false
	}
	defer f.Close()
	// copy git file to empty file
	num, err := io.Copy(f, resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println(num)
	// unzip file to tmp dir
	err = unzipFile(tmpDirPath+"/ws_chart.zip", tmpDirPath+"/ws_charts")
	if err != nil {
		log.Println(err)
		return false
	}

	// fetch directory name
	dirName := ""
	names, err := ioutil.ReadDir(tmpDirPath + "/ws_charts")
	if err != nil {
		log.Println(err)
		return false
	}
	for _, f := range names {
		if f.IsDir() {
			dirName = f.Name()
			break
		}
	}
	if dirName == "" {
		return false
	}

	//======================
	os.RemoveAll(config.PersistStoragePath)
	os.MkdirAll(config.PersistStoragePath, 0777)
	err = copyDir(tmpDirPath+"/ws_charts/"+dirName, config.PersistStoragePath)
	if err != nil {
		return false
	}

	return true
}
