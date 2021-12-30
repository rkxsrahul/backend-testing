package web

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"git.xenonstack.com/util/continuous-security-backend/config"
)

//checkWorkspaceUser function used for check the user exist in workspace or not.
func checkWorkspaceUser(email, project, token string) error {

	method := "GET"

	client := &http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest(method, config.Conf.Service.Workspace+"/v1/workspaces/"+project+"/user/"+email, nil)

	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	mapd := make(map[string]interface{})
	err = json.Unmarshal(body, &mapd)
	if err != nil {
		log.Println(err)
		return err
	}
	_, ok := mapd["role"].(string)
	if !ok {
		log.Println(mapd)
		return errors.New("no user found")
	}

	return nil
}
