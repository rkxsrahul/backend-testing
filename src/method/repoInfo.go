package method

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Language struct {
	Language string `json:"language"`
}

type Name struct {
	Name string `json:"name"`
}

//FetchLanguage function used for fetch the github specific language
func FetchLanguage(projectName, repoName string) (Language, string, error) {
	client := &http.Client{}
	var language Language

	//fetch the repository related information
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+projectName+"/"+repoName, nil)
	if err != nil {
		log.Println(err)
		return language, projectName, err
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return language, projectName, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return language, projectName, err
	}

	err = json.Unmarshal(body, &language)
	if err != nil {
		log.Println(err)
		return language, projectName, err
	}

	//fetch the use related information
	var name Name
	req, err = http.NewRequest("GET", "https://api.github.com/users/"+projectName, nil)
	if err != nil {
		log.Println(err)
		return language, name.Name, err
	}
	res, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return language, name.Name, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return language, name.Name, err
	}

	err = json.Unmarshal(body, &name)
	if err != nil {
		log.Println(err)
		return language, name.Name, err
	}

	return language, name.Name, err
}
