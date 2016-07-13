package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/siky/model"
)

func GetTags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := r.FormValue("name")
	group := r.FormValue("group")
	log.Println("group:" + group + ",name:" + name)
	var tagListUrl string
	if group != "" {
		tagListUrl = "http://10.3.15.35:5000/v2/" + group + "/" + name + "/tags/list"
	} else {
		tagListUrl = "http://10.3.15.35:5000/v2/" + name + "/tags/list"
	}
	resp, err := http.Get(tagListUrl)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

	}
	var tagList map[string]interface{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		//todo
		log.Println("Unmarshal error")
		log.Println(err)

	}
	var wg sync.WaitGroup
	tags := tagList["tags"]
	wg.Add(len(tags.([]interface{})))
	var tagsInfo = make(model.ImageList, 0)
	for i, tag := range tags.([]interface{}) {
		log.Println(i)
		go func(t string) {
			detail := getTag(group, name, t)
			tagsInfo = append(tagsInfo, detail)
			//tagsInfo = jsonDetail
			//tagsInfo["tag"] = t
			defer wg.Done()
		}(tag.(string))
	}
	wg.Wait()
	log.Println("tasks done")
	log.Println(tagsInfo)
	sort.Sort(tagsInfo)
	json, err := json.Marshal(&tagsInfo)
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
	detail := getTag(group, name, tag)

	json, err := json.Marshal(&detail)
	if err != nil {
		//todo
	}
	//w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(json))
}
func getTag(group, name, version string) model.Image {
	log.Println("tag:" + version)
	// http request for get mainfest
	var detailUrl string
	if group != "" {
		detailUrl = "http://10.3.15.35:5000/v2/" + group + "/" + name + "/manifests/" + version
	} else {
		detailUrl = "http://10.3.15.35:5000/v2/" + name + "/manifests/" + version
	}
	resp, err := http.Get(detailUrl)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {

	}
	info := make(map[string]interface{})
	err = json.Unmarshal(data, &info)
	if err != nil {
		// todo
	}
	detail := info["history"].([]interface{})[0].(map[string]interface{})["v1Compatibility"]

	//json unmarshal
	//jsonDetail := make(map[string]interface{});
	tag := &model.Image{}
	strs, ok := detail.(string)
	if ok {
		log.Println("IS STRING")
		log.Println(strs)
		err = json.Unmarshal([]byte(strs), tag)
		if err != nil {
			//todo
		}
	}
	tag.Version = version
	return *tag
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
		log.Printf("v:%v", v)
		imgPath := strings.Split(v, "/")
		if len(imgPath) > 1 {
			log.Printf("image gourp:%s\n", imgPath[0])
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
