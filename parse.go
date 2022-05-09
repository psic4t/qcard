package main

import (
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

func parseContactEmailHome(contactData *string) string {
	re, _ := regexp.Compile(`(?i)EMAIL(;TYPE=(HOME|INTERNET|PREF|INTERNET,HOME))?:.*?\n`)
	result := re.FindString(*contactData)
	return trimField(result, "(?i)EMAIL(;TYPE=(HOME|INTERNET|PREF|INTERNET,HOME))?:")
}

func parseContactEmailWork(contactData *string) string {
	re, _ := regexp.Compile(`(?i)EMAIL;TYPE=(WORK|INTERNET,WORK):.*?\n`)
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

func parseMain(contactData *string, contactsSlice *[]contactStruct, href, color string) {
	//fmt.Println(parseContactName(contactData))
	fullName := parseContactFullName(contactData)

	if (filter == "") || (filterMatch(fullName) == true) {
		data := contactStruct{
			Href:         href,
			Color:        color,
			fullName:     fullName,
			name:         parseContactName(contactData),
			title:        parseContactTitle(contactData),
			role:         parseContactRole(contactData),
			organisation: parseContactOrg(contactData),
			phoneCell:    parseContactPhoneCell(contactData),
			phoneHome:    parseContactPhoneHome(contactData),
			phoneWork:    parseContactPhoneWork(contactData),
			emailHome:    parseContactEmailHome(contactData),
			emailWork:    parseContactEmailWork(contactData),
			addressHome:  parseContactAddressHome(contactData),
			addressWork:  parseContactAddressWork(contactData),
			birthday:     parseContactBirthday(contactData),
			note:         parseContactNote(contactData),
		}
		*contactsSlice = append(*contactsSlice, data)
	}
}
