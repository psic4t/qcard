package main

import (
	// 	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
)

var config = getConf()

func fetchAbData(abNo int, wg *sync.WaitGroup) {
	var xmlBody string
	/*var xmlFilter string
	if filter != "" {
		xmlFilter = `<c:filter><c:prop-filter name="FN">
				<c:text-match collation="i;unicode-casemap" match-type="contains">` + filter + `</c:text-match>
			     </c:prop-filter></c:filter>`
	}*/

	xmlBody = `<c:addressbook-query xmlns:d="DAV:" xmlns:c="urn:ietf:params:xml:ns:carddav"><d:prop>
            		<d:getetag /><c:address-data />
		   </d:prop></c:addressbook-query>`

	//fmt.Println(xmlBody)
	req, err := http.NewRequest("REPORT", config.Addressbooks[abNo].Url, strings.NewReader(xmlBody))
	req.SetBasicAuth(config.Addressbooks[abNo].Username, config.Addressbooks[abNo].Password)
	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	req.Header.Add("Depth", "1") // needed for SabreDAV
	req.Header.Add("Prefer", "return-minimal")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	xmlContent, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	//fmt.Println(string(xmlContent))
	xmlData := XmlDataStruct{}
	err = xml.Unmarshal(xmlContent, &xmlData)
	if err != nil {
		log.Fatal(err)
	}

	for i := range xmlData.Elements {
		contactData := xmlData.Elements[i].Data
		contactHref := xmlData.Elements[i].Href
		ABColor := Colors[abNo]
		//fmt.Println(contactData)
		parseMain(&contactData, &contactsSlice, contactHref, ABColor)
	}

	wg.Done()
}

func showAddresses(singleAB string) {
	var wg sync.WaitGroup            // use waitgroups to fetch calendars in parallel
	wg.Add(len(config.Addressbooks)) // waitgroup length = num calendars
	for i := range config.Addressbooks {
		if singleAB == fmt.Sprintf("%v", i) || singleAB == "all" { // sprintf because convert int to string
			go fetchAbData(i, &wg)
		} else {
			wg.Done()
		}
	}
	wg.Wait()

	// TODO: Allow sort by first and last name
	sort.Slice(contactsSlice, func(i, j int) bool {
		return contactsSlice[i].fullName < contactsSlice[j].fullName
		//return contactsSlice[i].name < contactsSlice[j].name
	})

	if len(contactsSlice) == 0 {
		log.Fatal("no contacts") // get out if nothing found
	}
	if len(contactsSlice) <= config.DetailThreshold {
		showDetails = true
	}

	for _, e := range contactsSlice {
		e.fancyOutput() // pretty print
	}
}

func main() {
	toFile := false

	flag.StringVar(&filter, "s", "", "Search term")
	//flag.BoolVar(&showInfo, "i", false, "Show additional info like description and location for contacts")
	flag.BoolVar(&showFilename, "f", false, "Show contact filename for editing or deletion")
	flag.BoolVar(&displayFlag, "p", false, "Print VCF file piped to qcard (for CLI mail tools like mutt)")
	abNumber := flag.String("a", "all", "Show only single addressbook (number)")
	version := flag.Bool("v", false, "Show version")
	showAddressbooks := flag.Bool("l", false, "List configured addressbooks with their corresponding numbers (for \"-c\")")
	contactFile := flag.String("u", "", "Upload contact file. Provide filename and use with \"-c\"")
	contactDelete := flag.String("d", "", "Delete contact. Get filename with \"-f\" and use with \"-c\"")
	contactDump := flag.String("dump", "", "Dump raw contact data. Get filename with \"-f\" and use with \"-c\"")
	contactEdit := flag.String("edit", "", "Edit + upload contact data. Get filename with \"-f\" and use with \"-c\"")
	flag.Parse()
	flagset := make(map[string]bool) // map for flag.Visit. get bools to determine set flags
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })

	if *showAddressbooks {
	}
	if flagset["l"] {
		getAbList()
	} else if flagset["d"] {
		deleteContact(*abNumber, *contactDelete)
	} else if flagset["dump"] {
		dumpContact(*abNumber, *contactDump, toFile)
	} else if flagset["p"] {
		displayVCF()
	} else if flagset["edit"] {
		editContact(*abNumber, *contactEdit)
	} else if flagset["u"] {
		contactEdit := false
		uploadVCF(*abNumber, *contactFile, contactEdit)
	} else if *version {
		fmt.Print("qcard ")
		fmt.Println(qcardversion)
	} else {
		//filter = os.Args[1]
		showAddresses(*abNumber)
	}
}
