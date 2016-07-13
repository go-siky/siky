package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Tag struct {
	Author        string    `json:"author"`
	Created       time.Time `json:"created"`
	DockerVersion string    `json:"docker_version"`
	OS            string    `json:"os"`
	ID            string    `json:"id"`
	Tag           string    `json:"tag"`
}

type TagList []Tag

func (p TagList) Len() int {
	return len(p)
}

func (p TagList) Less(i, j int) bool {
	return p[i].Created.After(p[j].Created)
}

func (p TagList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Repository struct {
	Group  string   `json:"group"`
	Images []string `json:"images"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
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
	var tagsInfo = make(TagList, 0)
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
func getTag(group, name, version string) Tag {
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
	tag := &Tag{}
	strs, ok := detail.(string)
	if ok {
		log.Println("IS STRING")
		log.Println(strs)
		err = json.Unmarshal([]byte(strs), tag)
		if err != nil {
			//todo
		}
	}
	tag.Tag = version
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

func catalog(repos map[string][]string) []*Repository {
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
	var repositories []*Repository
	var sortedKeys []string
	for k, _ := range catalog {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, v := range sortedKeys {
		images := catalog[v]
		sort.Strings(images)
		group := &Repository{v, images}
		repositories = append(repositories, group)
	}
	log.Printf("After catalog:%v", catalog)
	return repositories
}

//My Handles xxx
type Handlers []http.Handler

// Implement the ServerHTTP method on our new type
func (handles Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, handler := range handles {
		if strings.HasPrefix(r.URL.Path, "/api/v") {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
		}
		handler.ServeHTTP(w, r)
	}
}
func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	//router.GET("/api/v2/tag/:group", GetTag)
	router.GET("/api/v2/repository/tag/:id", GetTag)
	router.GET("/api/v2/tags", GetTags)
	router.GET("/api/v2/repositories", GetRepositories)

	mHandlers := make(Handlers, 0)
	mHandlers = append(mHandlers, router)
	log.Fatal(http.ListenAndServe(":8080", mHandlers))
}
