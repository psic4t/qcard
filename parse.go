package main

import (
	"fmt"
	"regexp"
	"strings"
)

func trimField(field, cutset string) string {
	re, _ := regexp.Compile(cutset)
	cutsetRem := re.ReplaceAllString(field, "")
	return strings.TrimRight(cutsetRem, "\r\n")
}

func parseContactFullName(contactData *string) string {
	re, _ := regexp.Compile(`\nFN:.*\n`)
	result := re.FindString(*contactData)
	return trimField(result, "\nFN:")
}

func parseContactName(contactData *string) string {
	re, _ := regexp.Compile(`\nN:.*?\n`)
	result := re.FindString(*contactData)
	result = strings.Replace(result, ";;;", "", -1) // remove triple semicola
	return trimField(result, "\nN:")
}

func parseContactPhoneCell(contactData *string) string {
	re, _ := regexp.Compile(`(?i)TEL;TYPE=CELL:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)TEL;TYPE=CELL:")
}

func parseContactPhoneHome(contactData *string) string {
	re, _ := regexp.Compile(`(?i)TEL;TYPE=HOME:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)TEL;TYPE=HOME:")
}

func parseContactPhoneWork(contactData *string) string {
	re, _ := regexp.Compile(`(?i)TEL;TYPE=WORK:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)TEL;TYPE=WORK:")
}

func parseContactEmail(contactData *string) string {
	var emailType string
	re, _ := regexp.Compile(`(?i)EMAIL(;TYPE=(.*))?:.*?\n`)
	parts := re.FindStringSubmatch(*contactData)

	if len(parts) > 1 {
		types := strings.Split(parts[2], ",")
		for _, i := range types {
			//fmt.Println(i)
			switch strings.ToLower(i) {
			case "internet":
				emailType = "home"
			case "home":
				emailType = "home"
			case "pref":
				emailType = "home"
			case "work":
				emailType = "work"
			}
		}
	}

	fmt.Println(emailType)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)EMAIL(;TYPE=(.*))?:")
}

func parseContactEmailHome(contactData *string) string {
	/*workre, _ := regexp.Compile(`(?i)EMAIL;TYPE=(.*WORK.*):.*?\n`) // check, if work email
	if workre.FindString(*contactData) != "" {
		fmt.Println("yes")
		return "" // if work email, get out
	}*/
	re, _ := regexp.Compile(`(?i)EMAIL(;TYPE=(HOME|INTERNET|PREF|INTERNET,HOME|HOME,INTERNET))?:.*?\n`)
	//re, _ := regexp.Compile(`(?i)EMAIL(;TYPE=(.*))?:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)EMAIL(;TYPE=(.*))?:")

}

func parseContactEmailWork(contactData *string) string {
	//re, _ := regexp.Compile(`(?i)EMAIL;TYPE=(WORK|INTERNET,WORK):.*?\n`)
	re, _ := regexp.Compile(`(?i)EMAIL;TYPE=(.*WORK.*):.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)EMAIL;TYPE=(WORK|INTERNET,WORK):")
}

func parseContactAddressHome(contactData *string) string {
	re, _ := regexp.Compile(`(?i)ADR;TYPE=HOME:.*?\n`)
	result := re.FindString(*contactData)
	result = strings.Replace(result, ";;;", "", -1) // remove triple semicola
	return trimField(result, "(?i)ADR;TYPE=HOME:")
}

func parseContactAddressWork(contactData *string) string {
	re, _ := regexp.Compile(`(?i)ADR;TYPE=WORK:.*?\n`)
	result := re.FindString(*contactData)
	result = strings.Replace(result, ";;;", "", -1) // remove triple semicola
	return trimField(result, "(?i)ADR;TYPE=WORK:")
}

func parseContactBirthday(contactData *string) string {
	re, _ := regexp.Compile(`(?i)BDAY:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)BDAY:")
}

func parseContactNote(contactData *string) string {
	re, _ := regexp.Compile(`(?i)NOTE:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)NOTE:")
}

func parseContactOrg(contactData *string) string {
	re, _ := regexp.Compile(`(?i)ORG:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)ORG:")
}

func parseContactTitle(contactData *string) string {
	re, _ := regexp.Compile(`(?i)TITLE:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)TITLE:")
}

func parseContactRole(contactData *string) string {
	re, _ := regexp.Compile(`(?i)ROLE:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)ROLE:")
}

func parseContactNickname(contactData *string) string {
	re, _ := regexp.Compile(`(?i)NICKNAME:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)NICKNAME:")
}

func parseMain(contactData *string, contactsSlice *[]contactStruct, href, color string) {
	//fmt.Println(parseContactName(contactData))
	fullName := parseContactFullName(contactData)
	organisation := parseContactOrg(contactData)

	//if flagset["so"] {
	// TODO: How to filter with name and org?
	if (filter == "") || ((filterMatch(fullName) == true) || (filterOrgMatch(organisation) == true)) {
		data := contactStruct{
			Href:         href,
			Color:        color,
			fullName:     fullName,
			name:         parseContactName(contactData),
			title:        parseContactTitle(contactData),
			role:         parseContactRole(contactData),
			organisation: organisation,
			phoneCell:    parseContactPhoneCell(contactData),
			phoneHome:    parseContactPhoneHome(contactData),
			phoneWork:    parseContactPhoneWork(contactData),
			emailHome:    parseContactEmailHome(contactData),
			emailWork:    parseContactEmailWork(contactData),
			addressHome:  parseContactAddressHome(contactData),
			addressWork:  parseContactAddressWork(contactData),
			birthday:     parseContactBirthday(contactData),
			nickname:     parseContactNickname(contactData),
			note:         parseContactNote(contactData),
		}
		*contactsSlice = append(*contactsSlice, data)
	}
}
