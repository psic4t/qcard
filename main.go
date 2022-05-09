package main

import (
	// 	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	//"strconv"
	"strings"
	"sync"
	"time"
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

func showAddresses(abNo int) {
	var wg sync.WaitGroup            // use waitgroups to fetch calendars in parallel
	wg.Add(len(config.Addressbooks)) // waitgroup length = num calendars
	for i := range config.Addressbooks {
		if abNo == i || abNo == 1000 {
			go fetchAbData(i, &wg)
		} else {
			wg.Done()
		}
	}
	wg.Wait()

	// TODO: Allow sort by first and last name
	sort.Slice(contactsSlice, func(i, j int) bool {
		if config.SortByLastname {
			return contactsSlice[i].name < contactsSlice[j].name
		} else {
			return contactsSlice[i].fullName < contactsSlice[j].fullName
		}
	})

	if len(contactsSlice) == 0 {
		log.Fatal("no contacts") // get out if nothing found
	}
	if len(contactsSlice) <= config.DetailThreshold {
		showDetails = true
	}

	if *showEmailOnly {
		fmt.Println("Searching...")
	}

	for _, e := range contactsSlice {
		e.fancyOutput() // pretty print
	}
}

func createContact(abNo int, contactData string) {
	curTime := time.Now()
	d := regexp.MustCompile(`\s[a-z,A-Z]:`)
	dataArr := splitAfter(contactData, d) // own function, splitAfter is not supported by regex module

	var fullName string
	var name string
	var phoneCell string
	var phoneHome string
	var phoneWork string
	var emailHome string
	var emailWork string
	var addressHome string
	var addressWork string
	var note string
	var birthday string
	var organisation string
	var title string
	var role string
	var nickname string

	newUUID := genUUID()

	for i, e := range dataArr {
		if i == 0 {
			fullName = e
			lastInd := strings.LastIndex(e, " ")              // split name at last space
			name = e[lastInd+1:] + ";" + e[0:lastInd] + ";;;" // lastname, givenname1 givenname2
		} else {
			attr := strings.Split(e, ":")

			switch attr[0] {
			case " M":
				phoneCell = "\nTEL;TYPE=CELL:" + attr[1]
			case " P":
				phoneHome = "\nTEL;TYPE=HOME:" + attr[1]
			case " p":
				phoneWork = "\nTEL;TYPE=WORK:" + attr[1]
			case " E":
				emailHome = "\nEMAIL;TYPE=HOME:" + attr[1]
			case " e":
				emailWork = "\nEMAIL;TYPE=WORK:" + attr[1]
			case " A":
				if strings.Contains(attr[1], ";") == false {
					log.Fatal("Address must be splitted in Semicola")
				}
				addressHome = "\nADR;TYPE=HOME:" + attr[1]
			case " a":
				addressWork = "\nADR;TYPE=WORK:" + attr[1]
			case " O":
				organisation = "\nORG:" + attr[1]
			case " B":
				birthday = "\nBDAY:" + attr[1]
			case " n":
				note = "\nNOTE:" + attr[1]
			case " T":
				title = "\nTITLE:" + attr[1]
			case " R":
				role = "\nROLE:" + attr[1]
			case " I":
				role = "\nNICKNAME:" + attr[1]
			}
		}
	}

	var contactSkel = `BEGIN:VCARD
VERSION:3.0
PRODID:-//qcard
UID:` + newUUID +
		emailHome +
		phoneCell +
		phoneHome +
		phoneWork +
		emailHome +
		emailWork +
		addressHome +
		addressWork +
		birthday +
		organisation +
		note +
		title +
		role +
		nickname + `
FN:` + fullName + `
N:` + name + `
REV:` + curTime.UTC().Format(IcsFormat) + `
END:VCARD`
	//fmt.Println(contactSkel)
	//os.Exit(3)

	newElem := newUUID + `.vcf`

	req, _ := http.NewRequest("PUT", config.Addressbooks[abNo].Url+newElem, strings.NewReader(contactSkel))
	req.SetBasicAuth(config.Addressbooks[abNo].Username, config.Addressbooks[abNo].Password)
	req.Header.Add("Content-Type", "application/xml; charset=utf-8")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Status)
}

func main() {
	toFile := false

	if len(os.Args[1:]) > 0 {
		searchterm = os.Args[1]
	}
	flag.StringVar(&filter, "s", "", "Search term")
	//flag.BoolVar(&showInfo, "i", false, "Show additional info like description and location for contacts")
	flag.BoolVar(&showFilename, "f", false, "Show contact filename for editing or deletion")
	flag.BoolVar(&displayFlag, "p", false, "Print VCF file piped to qcard (for CLI mail tools like mutt)")
	abNo := flag.Int("a", 0, "Show only single addressbook (number).")
	version := flag.Bool("v", false, "Show version")
	showAddressbooks := flag.Bool("l", false, "List configured addressbooks with their corresponding numbers (for \"-a\")")
	contactFile := flag.String("u", "", "Upload contact file. Provide filename and use with \"-c\"")
	contactDelete := flag.String("delete", "", "Delete contact. Get filename with \"-f\" and use with \"-a\"")
	contactDump := flag.String("d", "", "Dump raw contact data. Get filename with \"-f\" and use with \"-a\"")
	contactEdit := flag.String("edit", "", "Edit + upload contact data. Get filename with \"-f\" and use with \"-a\"")
	contactNew := flag.String("n", "", "Add a new contact. Check README.md for syntax")
	showEmailOnly = flag.Bool("emailonly", false, "Show only email addresses and names without further formatting (for CLI mail tools like mutt)")
	flag.Parse()
	flagset := make(map[string]bool) // map for flag.Visit. get bools to determine set flags
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })

	if *showAddressbooks {
	}
	if flagset["l"] {
		getAbList()
	} else if flagset["delete"] {
		deleteContact(*abNo, *contactDelete)
	} else if flagset["d"] {
		dumpContact(*abNo, *contactDump, toFile)
	} else if flagset["p"] {
		displayVCF()
	} else if flagset["n"] {
		createContact(*abNo, *contactNew)
	} else if flagset["edit"] {
		editContact(*abNo, *contactEdit)
	} else if flagset["u"] {
		contactEdit := false
		uploadVCF(*abNo, *contactFile, contactEdit)
	} else if *version {
		fmt.Println("qcard " + qcardversion)
	} else if flagset["a"] {
		showAddresses(*abNo)
	} else {
		showAddresses(1000)
	}
}
