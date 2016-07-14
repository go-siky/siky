package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"

	"github.com/qiniu/log"
	"gopkg.in/siky/config"
	"gopkg.in/siky/model"
	"gopkg.in/siky/util"
)

//GetAllImages from docker registry
func GetAllImages(group, name string) *model.ImageList {
	log.Debug("group:", group, ",name:", name)

	tagListURL := util.JoinURL(config.Get().Registry.URL, config.Get().Registry.Version)

	if group != "" {
		tagListURL = util.JoinURL(tagListURL, group, name, "tags", "list")
	} else {
		tagListURL = util.JoinURL(tagListURL, name, "tags", "list")
	}
	log.Debug("Tag URL:", tagListURL)
	resp, err := http.Get(tagListURL)
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
		log.Debug("Unmarshal error")
		log.Debug(err)

	}
	var wg sync.WaitGroup
	tags := tagList["tags"]
	wg.Add(len(tags.([]interface{})))
	var images = make(model.ImageList, 0)
	for i, tag := range tags.([]interface{}) {
		log.Debug(i)
		go func(t string) {
			detail := GetTag(group, name, t)
			images = append(images, detail)
			//tagsInfo = jsonDetail
			//tagsInfo["tag"] = t
			defer wg.Done()
		}(tag.(string))
	}
	wg.Wait()
	log.Debug("tasks done")
	log.Debug(images)
	sort.Sort(images)
	// json, err := json.Marshal(&tagsInfo)
	// if err != nil {
	// 	//todo
	// }
	// fmt.Fprintf(w, string(json))
	return &images
}
func GetTag(group, name, version string) model.Image {
	log.Debug("tag:" + version)
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
