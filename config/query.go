package config

import (
	"log"
	"regexp"

	"github.com/beevik/etree"
	pq "github.com/pasiol/gopq"
	"github.com/sethvargo/go-password/password"
)

var studentRegistryFilter = ""
var applicantrsRegistryFilter = ""

// PasswordsApplicants func
func PasswordsApplicants() pq.PrimusQuery {
	pq := pq.PrimusQuery{}
	pq.Charset = "UTF-8"
	pq.Database = "hakijat"
	pq.Sort = ""
	pq.Search = applicantrsRegistryFilter
	pq.Data = "#DATA{V1}"
	pq.Footer = ""

	return pq
}

// PasswordsStudents func
func PasswordsStudents() pq.PrimusQuery {
	pq := pq.PrimusQuery{}
	pq.Charset = "UTF-8"
	pq.Database = "opphenk"
	pq.Sort = ""
	pq.Search = studentRegistryFilter
	pq.Data = "#DATA{V1}"
	pq.Footer = ""

	return pq
}

func checkADPassworValitidy(password string) bool {
	if len(password) >= 8 {
		upperCaseLetters, _ := regexp.MatchString(`[A-Z]`, password)
		lowerCaseLetters, _ := regexp.MatchString(`[a-z]`, password)
		digits, _ := regexp.MatchString(`[0-9]`, password)
		if upperCaseLetters && lowerCaseLetters && digits {
			return true
		}
	}
	return false
}

func generatePassword() string {
	notValidated := true
	var pwd string
	for notValidated {
		pwd = password.MustGenerate(10, 5, 0, false, false)
		if checkADPassworValitidy(pwd) {
			notValidated = false
		}
	}
	return pwd
}

// PasswordXMLApplicants generator
func PasswordXMLApplicants(id string) string {
	passwordDoc := etree.NewDocument()
	passwordDoc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	primusquery := passwordDoc.CreateElement("PRIMUSQUERY_IMPORT")
	primusquery.CreateAttr("ARCHIVEMODE", "0")
	primusquery.CreateAttr("CREATEIFNOTFOUND", "0")
	identity := passwordDoc.CreateElement("IDENTITY")
	identity.CreateText("service-password")
	card := passwordDoc.CreateElement("CARD")
	card.CreateAttr("FIND", id)
	passwordElem := card.CreateElement("SALASANA")
	passwordElem.CreateText(generatePassword())
	passwordDoc.Indent(2)
	xmlAsString, _ := passwordDoc.WriteToString()
	filename, err := pq.CreateTMPFile(pq.StringWithCharset(128)+".xml", xmlAsString)
	if err != nil {
		log.Fatalf("Writing output file failed: %s", err)
	}
	return filename
}

// PasswordXMLApplicants generator
func PasswordXMLStudents(id string) string {
	passwordDoc := etree.NewDocument()
	passwordDoc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	primusquery := passwordDoc.CreateElement("PRIMUSQUERY_IMPORT")
	primusquery.CreateAttr("ARCHIVEMODE", "0")
	primusquery.CreateAttr("CREATEIFNOTFOUND", "0")
	identity := passwordDoc.CreateElement("IDENTITY")
	identity.CreateText("service-password")
	card := passwordDoc.CreateElement("CARD")
	card.CreateAttr("FIND", id)
	passwordElem := card.CreateElement("SALASANA")
	res, err := password.Generate(10, 5, 0, false, false)
	if err != nil {
		log.Fatal(err)
	}
	passwordElem.CreateText(res)
	passwordDoc.Indent(2)
	xmlAsString, _ := passwordDoc.WriteToString()
	filename, err := pq.CreateTMPFile(pq.StringWithCharset(128)+".xml", xmlAsString)
	if err != nil {
		log.Fatalf("Writing output file failed: %s", err)
	}
	return filename
}
