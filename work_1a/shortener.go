package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
)

var (
	//хранилище ссылок
	UrlBaseByShortName map[string]*Url
	UrlBaseByLongName  map[string]*Url
)

type Shortener interface {
	Shorten(url string) string
	Resolve(url string) string
}

//сделано только чтоб грузить в json надо переделать как то по другом
type UrlMap struct {
	List map[string]*Url `json:"url",attr`
}

type Url struct {
	ShortUrl string `json:"short_url",attr`
	LongUrl  string `json:"long_url",attr`
}

func (u *Url) Shorten(longUrl string) string {

	if url, ok := UrlBaseByLongName[longUrl]; ok {
		return url.ShortUrl
	} else {
		urlData := []byte(longUrl)
		for {
			shorUrl := fmt.Sprintf("%x", md5.Sum(urlData))
			if _, ok := UrlBaseByShortName[shorUrl]; ok {
				log.Infof("Short url %s уже есть в базе для long url %s добавляем сиимвол чтобы получить новый хэш", shorUrl, longUrl)
				randomByte := make([]byte, 1)
				rand.Read(randomByte)
				urlData = append(urlData, randomByte[0])
				log.Infof("Новый long url для генерации хэш : %s", urlData)
			} else {
				newUrl := &Url{
					ShortUrl: shorUrl,
					LongUrl:  longUrl,
				}
				UrlBaseByShortName[shorUrl] = newUrl
				UrlBaseByLongName[longUrl] = newUrl
				return shorUrl
			}
		}
	}
}

func (u *Url) Resolve(longUrl string) string {
	if url, ok := UrlBaseByLongName[longUrl]; ok {
		return url.LongUrl
	} else {
		return ""
	}
}

func ReadUrlBaseFromJsonFile(path string) error {
	var list *UrlMap
	if err := readJSON(list, path); err != nil || list.List == nil {
		return err
	}
	log.Debugf("tmap=%v", list.List)
	for _, url := range list.List {
		UrlBaseByShortName[url.ShortUrl] = url
		UrlBaseByLongName[url.LongUrl] = url
		log.Debugf("Прочитано из JSON Short url=%s Long url=%s", url.ShortUrl, url.LongUrl)
	}
	return nil
}

func WriteUrlBaseToJsonFile(path string) error {
	list := &UrlMap{
		List: UrlBaseByLongName,
	}
	if err := writeJSON(list, path); err != nil {
		return err
	}
	return nil
}

func readJSON(interfaceJSON interface{}, path string) error {

	jsfile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Unable to read %v. Err: %v.", path, err)
		return err
	}

	if jerr := json.Unmarshal(jsfile, &interfaceJSON); jerr != nil {
		log.Errorf("jerr:", jerr.Error())
		return jerr
	} else {
		log.Debugf("read JSON from file %s:\n%v", path, interfaceJSON)
	}

	// backupFile := path + time.Now().Format("2006-01-02_15:04:05")
	// log.Infof("Backup json file to %s", backupFile)
	// bk, err := os.Create(backupFile)
	// if err != nil {
	// 	log.Fatalf("Unable to create %v. Err: %v.", backupFile, err)
	// 	return err
	// }
	// defer bk.Close()
	// bk.Write(jsfile)

	return nil
}

func writeJSON(interfaceJSON interface{}, path string) error {

	fp, err := os.Create(path)
	if err != nil {
		log.Errorf("Unable to create %v. Err: %v.", path, err)
		return err
	}
	defer fp.Close()

	j, jerr := json.MarshalIndent(interfaceJSON, "", "  ")
	if jerr != nil {
		log.Errorf("jerr:", jerr.Error())
		return jerr
	}

	fp.Write(j)

	return nil
}

func init() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true})
	log.SetOutput(os.Stderr)
	UrlBaseByShortName = make(map[string]*Url)
	UrlBaseByLongName = make(map[string]*Url)
}
