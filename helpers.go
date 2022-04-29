package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func getConf() *configStruct {
	configData, err := ioutil.ReadFile(configLocation)
	if err != nil {
		fmt.Print("Config not found. \n\nPlease copy config-sample.json to ~/.config/qcard/config.json and modify it accordingly.\n\n")
		log.Fatal(err)
	}

	conf := configStruct{}
	err = json.Unmarshal(configData, &conf)
	//fmt.Println(conf)
	if err != nil {
		log.Fatal(err)
	}

	return &conf
}

func getAbProps(calNo int, p *[]calProps, wg *sync.WaitGroup) {
	req, err := http.NewRequest("PROPFIND", config.Addressbooks[calNo].Url, nil)
	req.SetBasicAuth(config.Addressbooks[calNo].Username, config.Addressbooks[calNo].Password)

	/*tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli := &http.Client{Transport: tr}*/
	cli := &http.Client{}
	resp, err := cli.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	xmlContent, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	xmlProps := xmlProps{}
	err = xml.Unmarshal(xmlContent, &xmlProps)
	if err != nil {
		log.Fatal(err)
	}
	displayName := xmlProps.DisplayName

	thisCal := calProps{
		calNo:       calNo,
		displayName: displayName,
		url:         config.Addressbooks[calNo].Url,
	}
	*p = append(*p, thisCal)

	wg.Done()
}

func getAbList() {
	p := []calProps{}

	var wg sync.WaitGroup
	wg.Add(len(config.Addressbooks)) // waitgroup length = num calendars

	for i := range config.Addressbooks {
		go getAbProps(i, &p, &wg)
	}
	wg.Wait()

	sort.Slice(p, func(i, j int) bool {
		return p[i].calNo < (p[j].calNo)
	})

	for i := range p {
		u, err := url.Parse(config.Addressbooks[i].Url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(`[` + fmt.Sprintf("%v", i) + `] - ` + Colors[i] + colorBlock + ColDefault +
			` ` + p[i].displayName + ` (` + u.Hostname() + `)`)
	}
}

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func splitAfter(s string, re *regexp.Regexp) (r []string) {
	re.ReplaceAllStringFunc(s, func(x string) string {
		s = strings.Replace(s, x, "::"+x, -1)
		return s
	})
	for _, x := range strings.Split(s, "::") {
		if x != "" {
			r = append(r, x)
		}
	}
	return
}

func (e contactStruct) fancyOutput() {
	if showColor {
		fmt.Print(e.Color + colorBlock + ColDefault + ` `)
	}
	fmt.Println(e.fullName)
	//fmt.Println(e.name)
	//fmt.Printf(`%5s`, ` `)
	//fmt.Println(` M: ` + e.phoneCell)
	if showDetails {
		if e.title != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`T: ` + e.title)
		}
		if e.organisation != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`O: ` + e.organisation)
		}
		if e.phoneCell != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`M: ` + e.phoneCell)
		}
		if e.phoneHome != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`P: ` + e.phoneHome)
		}
		if e.emailHome != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`E: ` + e.emailHome)
		}
		if e.emailWork != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`e: ` + e.emailWork)
		}
		if e.addressHome != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`A: ` + e.addressHome)
		}
		if e.addressWork != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`a: ` + e.addressWork)
		}
		if e.phoneWork != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`p: ` + e.phoneWork)
		}
		if e.birthday != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`B: ` + e.birthday)
		}
		if e.name != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`N: ` + e.name)
		}
		if e.note != "" {
			fmt.Printf(`%2s`, ` `)
			fmt.Println(`n: ` + e.note)
		}
	}

	if showFilename {
		if e.Href != "" {
			fmt.Println(path.Base(e.Href))
		}
	}
	//fmt.Println()
}
func (e contactStruct) vcfOutput() {
	// whole day or greater
	fmt.Println(`Contact
=======`)
	//fmt.Printf(`Summary:%6s`, ` `)
	//fmt.Print(e.Summary)
	fmt.Printf(`Full Name:%2s`+e.fullName, ` `)
	fmt.Println(``)
	fmt.Printf(`Cell:%5s`+e.phoneCell, ` `)
	fmt.Println(``)
}

func genUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}

func strToInt(str string) (int, error) {
	nonFractionalPart := strings.Split(str, ".")
	return strconv.Atoi(nonFractionalPart[0])
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func filterMatch(fullName string) bool {
	re, _ := regexp.Compile(`(?i)` + filter)
	return re.FindString(fullName) != ""
}

func deleteContact(abNumber string, contactFilename string) (status string) {
	if contactFilename == "" {
		log.Fatal("No contact filename given")
	}

	abNo, _ := strconv.ParseInt(abNumber, 0, 64)

	req, _ := http.NewRequest("DELETE", config.Addressbooks[abNo].Url+contactFilename, nil)
	req.SetBasicAuth(config.Addressbooks[abNo].Username, config.Addressbooks[abNo].Password)

	cli := &http.Client{}
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Status)

	return
}

func dumpContact(abNumber string, contactFilename string, toFile bool) (status string) {
	abNo, _ := strconv.ParseInt(abNumber, 0, 64)
	//fmt.Println(config.Addressbooks[calNo].Url + eventFilename)

	req, _ := http.NewRequest("GET", config.Addressbooks[abNo].Url+contactFilename, nil)
	req.SetBasicAuth(config.Addressbooks[abNo].Username, config.Addressbooks[abNo].Password)

	cli := &http.Client{}
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(resp.Status)
	xmlContent, _ := ioutil.ReadAll(resp.Body)

	if toFile {
		// create cache dir if not exists
		os.MkdirAll(cacheLocation, os.ModePerm)
		err := ioutil.WriteFile(cacheLocation+"/"+contactFilename, xmlContent, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return contactFilename + " written"
	} else {
		fmt.Println(string(xmlContent))
		return
	}
}

func uploadVCF(abNumber string, contactFilePath string, contactEdit bool) (status string) {
	abNo, _ := strconv.ParseInt(abNumber, 0, 64)
	//fmt.Println(config.Calendars[calNo].Url + eventFilePath)

	var vcfData string
	var contactVCF string
	var contactFileName string

	if contactFilePath == "-" {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			vcfData += scanner.Text() + "\n"
		}
		//eventICS, _ = explodeEvent(&icsData)
		contactVCF = vcfData
		contactFileName = genUUID() + `.ics`
		fmt.Println(contactVCF)

	} else {
		//eventICS, err := ioutil.ReadFile(cacheLocation + "/" + eventFilename)
		contactVCFByte, err := ioutil.ReadFile(contactFilePath)
		if err != nil {
			log.Fatal(err)
		}

		contactVCF = string(contactVCFByte)
		if contactEdit == true {
			contactFileName = path.Base(contactFilePath) // use old filename again
		} else {
			contactFileName = genUUID() + `.ics` // no edit, so new filename
		}
	}
	req, _ := http.NewRequest("PUT", config.Addressbooks[abNo].Url+contactFileName, strings.NewReader(contactVCF))
	req.SetBasicAuth(config.Addressbooks[abNo].Username, config.Addressbooks[abNo].Password)
	req.Header.Add("Content-Type", "text/calendar; charset=utf-8")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Status)

	return
}

func displayVCF() {
	scanner := bufio.NewScanner(os.Stdin)

	var vcfData string

	for scanner.Scan() {
		vcfData += scanner.Text() + "\n"
	}

	parseMain(&vcfData, &contactsSlice, "none", "none")
	for _, e := range contactsSlice {
		e.vcfOutput()
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

}

func editContact(abNumber string, contactFilename string) (status string) {
	toFile = true
	contactEdit := true
	dumpContact(abNumber, contactFilename, toFile)
	//fmt.Println(appointmentEdit)
	filePath := cacheLocation + "/" + contactFilename
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	beforeMTime := fileInfo.ModTime()

	shell := exec.Command(editor, filePath)
	shell.Stdout = os.Stdin
	shell.Stdin = os.Stdin
	shell.Stderr = os.Stderr
	shell.Run()

	fileInfo, err = os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	afterMTime := fileInfo.ModTime()

	if beforeMTime.Before(afterMTime) {
		uploadVCF(abNumber, filePath, contactEdit)
	} else {
		log.Fatal("no changes")
	}

	return
}
