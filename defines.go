package main

import (
	"encoding/xml"
	"os"
	"time"
)

var err string
var homedir string = os.Getenv("HOME")
var editor string = os.Getenv("EDITOR")
var configLocation string = (homedir + "/" + ConfigDir + "/config.json")
var cacheLocation string = (homedir + "/" + CacheDir)
var versionLocation string = (cacheLocation + "/version.json")
var timezone, _ = time.Now().Zone()
var xmlContent []byte
var showDetails bool
var showFilename bool
var showEmailOnly *bool
var displayFlag bool
var toFile bool
var filter string
var searchterm string

//var colorBlock string = "█"
var colorBlock string = "|"
var contactsSlice []contactStruct
var Colors = [10]string{"\033[0;31m", "\033[0;32m", "\033[1;33m", "\033[1;34m", "\033[1;35m", "\033[1;36m", "\033[1;37m", "\033[1;38m", "\033[1;39m", "\033[1;40m"}
var showColor bool = true
var qcardversion string = "0.6.0"

const (
	ConfigDir  = ".config/qcard"
	CacheDir   = ".cache/qcard"
	IcsFormat  = "20060102T150405Z"
	ColWhite   = "\033[1;37m"
	ColDefault = "\033[0m"
	ColGreen   = "\033[0;32m"
	ColYellow  = "\033[1;33m"
	ColBlue    = "\033[1;34m"
)

type configStruct struct {
	Addressbooks []struct {
		Url      string
		Username string
		Password string
	}
	DetailThreshold int
	SortByLastname  bool
}

type contactStruct struct {
	Href         string
	Color        string
	fullName     string
	name         string
	title        string
	organisation string
	phoneCell    string
	phoneHome    string
	phoneWork    string
	emailHome    string
	emailWork    string
	addressHome  string
	addressWork  string
	birthday     string
	note         string
}

type xmlProps struct {
	calNo        string
	Url          string
	XMLName      xml.Name `xml:"multistatus"`
	Href         string   `xml:"response>href"`
	DisplayName  string   `xml:"response>propstat>prop>displayname"`
	Color        string   `xml:"response>propstat>prop>calendar-color"`
	CTag         string   `xml:"response>propstat>prop>getctag"`
	ETag         string   `xml:"response>propstat>prop>getetag"`
	LastModified string   `xml:"response>propstat>prop>getlastmodified"`
}

type calProps struct {
	calNo       int
	displayName string
	url         string
	color       string
}

type XmlDataStruct struct {
	XMLName  xml.Name          `xml:"multistatus"`
	Elements []xmlDataElements `xml:"response"`
}

type xmlDataElements struct {
	XMLName xml.Name `xml:"response"`
	Href    string   `xml:"href"`
	ETag    string   `xml:"propstat>prop>getetag"`
	Data    string   `xml:"propstat>prop>address-data"`
}
