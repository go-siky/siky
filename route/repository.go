package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/qiniu/log"
	"gopkg.in/siky/model"
	"gopkg.in/siky/service"
)

func GetTags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := r.FormValue("name")
	group := r.FormValue("group")
	log.Debug("group:" + group + ",name:" + name)
	images := service.GetAllImages(group, name)
	json, err := json.Marshal(images)
	if err != nil {
		//todo
	}
	fmt.Fprintf(w, string(json))

}
func GetTag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	var group, name, tag string

	params := strings.Split(id, "---")
	if len(params) == 2 {
		name = params[0]
		tag = params[1]
	}
	if len(params) == 3 {
		group = params[0]
		name = params[1]
		tag = params[2]
	}
	log.Printf("group:%s\n", group)
	log.Printf("name:%s\n", name)
	log.Printf("tag:%s\n", tag)
	detail := service.GetTag(group, name, tag)

	json, err := json.Marshal(&detail)
	if err != nil {
		//todo
	}
	//w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(json))
}

func GetRepositories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := "http://10.3.15.35:5000/v2/_catalog"
	resp, err := http.Get(url)
	if err != nil {
		//todo
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

	}
	log.Printf("body:%s", body)
	repos := make(map[string][]string)

	err = json.Unmarshal(body, &repos)
	if err != nil {
		//todo
		log.Println(err)
	}

	log.Printf("repos:%v", repos)
	json, err := json.Marshal(catalog(repos))
	if err != nil {
		//todo
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(json))
}

func catalog(repos map[string][]string) []*model.Repository {
	catalog := make(map[string][]string)
	for _, v := range repos["repositories"] {
		log.Debugf("v:%v", v)
		imgPath := strings.Split(v, "/")
		if len(imgPath) > 1 {
			log.Debugf("image gourp:%s\n", imgPath[0])
			content := catalog[imgPath[0]]
			content = append(content, imgPath[1])
			catalog[imgPath[0]] = content
		} else {
			content := catalog["libs"]
			content = append(content, imgPath[0])
			catalog["libs"] = content
		}
	}
	var repositories []*model.Repository
	var sortedKeys []string
	for k := range catalog {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, key := range sortedKeys {
		images := catalog[key]
		sort.Strings(images)
		repository := &model.Repository{Group: key, Images: images}
		repositories = append(repositories, repository)
	}
	log.Printf("After catalog:%v", catalog)
	return repositories
}
