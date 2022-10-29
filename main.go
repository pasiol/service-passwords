package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"service-passwords/config"
	"strings"
	"time"

	pq "github.com/pasiol/gopq"
)

var (
	// Version for build
	Version string
	// Build for build
	Build                  string
	jobName                = "service-passwords"
	applicantsImportConfig = "hakija-rekisteri-import"
	studentsImportConfig   = "opiskelija-rekisteri-import"
)

func getApplicantsWithoutPassword() ([]string, error) {
	query := config.PasswordsApplicants()
	c := config.GetPrimusConfig()
	query.Host = c.PrimusHost
	query.Port = c.PrimusPort
	query.User = c.PrimusUser
	query.Pass = c.PrimusPassword
	output, err := pq.ExecuteAndRead(query, 30)
	if err != nil {
		return nil, err
	}
	if output == "" {
		log.Print("No data, nothing to do.")
	}
	return strings.Fields(output), nil
}

func getStudentsWithoutPassword() ([]string, error) {
	query := config.PasswordsStudents()
	c := config.GetPrimusConfig()
	query.Host = c.PrimusHost
	query.Port = c.PrimusPort
	query.User = c.PrimusUser
	query.Pass = c.PrimusPassword
	output, err := pq.ExecuteAndRead(query, 30)
	if err != nil {
		return nil, err
	}
	if output == "" {
		log.Print("No data, nothing to do.")
	}
	return strings.Fields(output), nil
}

func main() {
	c := config.GetPrimusConfig()
	var register string
	if len(os.Args) == 2 {
		if os.Args[1] == "--help" {
			fmt.Printf("Usage: %s [register]\n", jobName)
			os.Exit(0)
		} else {
			register = os.Args[1]
			if register != "opphenk" && register != "hakijat" {
				fmt.Printf("Usage: %s [opphenk|hakijat]\n", jobName)
				os.Exit(0)
			}
		}
	}

	start := time.Now()
	wrt := io.MultiWriter(os.Stdout)
	log.SetOutput(wrt)
	log.Printf("Starttime: %v", start.Format(time.RFC3339))
	log.Printf("Starting job: %s", jobName)
	log.Println("Version: ", Version)
	log.Println("Build Time: ", Build)
	pq.Debug = false

	if register == "hakijat" {
		cards, err := getApplicantsWithoutPassword()
		if err != nil {
			log.Fatalf("getting applicants failed: %s", err)
		}
		for _, applicant := range cards {
			log.Printf("Trying to generate password for applicant: %s", applicant)
			outputFile := config.PasswordXMLApplicants(applicant)
			pq.ExecuteImportQuery(outputFile, c.PrimusHost, c.PrimusPort, c.PrimusUser, c.PrimusPassword, applicantsImportConfig)
		}
		cards = nil
	}

	if register == "opphenk" {
		cards, err := getStudentsWithoutPassword()
		if err != nil {
			log.Fatalf("getting students failed: %s", err)
		}
		for _, student := range cards {
			log.Printf("Trying to generate password for student: %s", student)
			outputFile := config.PasswordXMLApplicants(student)
			pq.ExecuteImportQuery(outputFile, c.PrimusHost, c.PrimusPort, c.PrimusUser, c.PrimusPassword, studentsImportConfig)
		}
		cards = nil
	}

	t := time.Now()
	elapsed := t.Sub(start)

	log.Printf("Ending succesfully %s.", jobName)
	log.Printf("Endtime: %v", t)
	log.Printf("Elapsed processing time %d.", elapsed)
}
