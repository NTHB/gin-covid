package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type CovidCaseNode struct {
	CovidCaseNode []CovidData `json:"Data"`
}

type CovidData struct {
	ConfirmDate    string `json:"ConfirmDate"`
	No             int    `json:"No"`
	Age            *int   `json:"Age"`
	Gender         string `json:"Gender"`
	GenderEn       string `json:"GenderEn"`
	Nation         string `json:"Nation"`
	NationEn       string `json:"NationEn"`
	Province       string `json:"Province"`
	ProvinceId     int    `json:"ProvinceId"`
	District       string `json:"District"`
	ProvinceEn     string `json:"ProvinceEn"`
	StatQuarantine int    `json:"StatQuarantine"`
}

type AgeGroup struct {
	GroupName1 int `json:"0-30"`
	GroupName2 int `json:"31-60"`
	GroupName3 int `json:"60+"`
	GroupName4 int `json:"N/A"`
}
type Report struct {
	Province map[string]int `json:"Province"`
	AgeGroup AgeGroup       `json:"AgeGroup"`
}

func main() {
	jsonFile, err := os.Open("covid-cases.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data CovidCaseNode

	json.Unmarshal(byteValue, &data)

	ProvinceMap := make(map[string]int)

	var (
		Group1Count int
		Group2Count int
		Group3Count int
		Group4Count int
	)

	for i := 0; i < len(data.CovidCaseNode); i++ {
		if val, exist := ProvinceMap[data.CovidCaseNode[i].Province]; exist {
			ProvinceMap[data.CovidCaseNode[i].Province] = val + 1
		} else {
			ProvinceMap[data.CovidCaseNode[i].Province] = 1
		}
		if data.CovidCaseNode[i].Age == nil {
			Group4Count++
		} else {
			PatientAge := *data.CovidCaseNode[i].Age

			if PatientAge >= 0 && PatientAge <= 30 {
				Group1Count++
			}
			if PatientAge > 30 && PatientAge <= 60 {
				Group2Count++
			}
			if PatientAge > 60 {
				Group3Count++
			}
		}
	}
	fmt.Println(ProvinceMap)
	fmt.Println(Group1Count, Group2Count, Group3Count, Group4Count)

	delete(ProvinceMap, "")
	fmt.Println(ProvinceMap)

	var Report Report
	Report.Province = ProvinceMap
	Report.AgeGroup.GroupName1 = Group1Count
	Report.AgeGroup.GroupName2 = Group2Count
	Report.AgeGroup.GroupName3 = Group3Count
	Report.AgeGroup.GroupName4 = Group4Count

	fmt.Println(Report)

	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to covid summary project!",
		})
	})

	r.GET("/covid/summary", func(c *gin.Context) {
		c.JSON(http.StatusOK, Report)
	})
	r.Run()
}
