package method

import (
	"log"
	"os/exec"

	"git.xenonstack.com/util/continuous-security-backend/config"
)

func RunBashCommand(url, header string) (string, map[string]interface{}) {
	mapd := make(map[string]interface{})

	//manage the path using header
	path := config.PersistStoragePath + headerFileSlug(header)
	cmd := exec.Command("bash", path, url)
	out, err := cmd.Output()
	if err != nil {
		log.Println(header, " ", err.Error(), cmd.Args)
		mapd["error"] = true
		mapd["header"] = header
		mapd["message"] = "Some error in scanning the URL. Please try after sometime"
		mapd["error_message"] = err.Error()

		return "", mapd
	}
	scriptOut := string(out)
	if scriptOut == " " || scriptOut == "" {
		log.Println(cmd.Args)
	}
	return scriptOut, mapd
}
