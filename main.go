package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/gocolly/colly/v2"
)

var Que1 = []*survey.Question{
	{
		Name:      "Name",
		Prompt:    &survey.Input{Message: "What is your name?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "Usage",
		Prompt: &survey.Select{
			Message: "Laptop Usage (in a day)",
			Options: []string{"5 Hrs", "5 to 10 Hrs", "More than 12 Hrs"},
			Default: "5 Hrs",
		},
	},
	{
		Name: "TouchScreen",
		Prompt: &survey.Select{
			Message: "Do you prefer a touch screen laptop?",
			Options: []string{"yes", "no"},
			Default: "no",
		},
	},

	{
		Name: "Preference",
		Prompt: &survey.Select{
			Message: "Your preference on how much laptop weighs?",
			Options: []string{"on the heavier side", "normal", "on the lighter side"},
			Default: "on the heavier side",
		},
	},
}

var Que2 = []*survey.Question{
	{
		Name: "What do make most out of your laptop?",
		Prompt: &survey.MultiSelect{
			Message: "What do make most out of your laptop? Choose one or more options :",
			Options: []string{
				"Online Shopping & alike",
				"Normal Gaming",
				"Hard-core Gaming",
				"Heavy Softwares",
				"Simple Softwares",
				"Streaming Videos",
				"Watching Movies - theatrical experience - yes",
				"Watching Movies - theatrical experience - not really",
				"PDF reading - simple alike jobs",
				"Video Calling",
			},
		},
	},
}

var LaptopDB []map[string]string

func collectInfo(url string) {
	var LaptopInfo map[string]string
	filePath := "data.json"
	file, err1 := ioutil.ReadFile(filePath)
	if err1 != nil {
		fmt.Printf("// error while reading file %s\n", filePath)
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}

	err2 := json.Unmarshal([]byte(file), &LaptopInfo)
	if err2 != nil {
		log.Fatal(err2)
	}
	var featuresKeys []string
	var features []string
	c := colly.NewCollector()
	c.OnHTML("td", func(e *colly.HTMLElement) {
		if e.Attr("class") == "_1hKmbr col col-3-12" {
			Feature := e.DOM.Text()
			// _, found := LaptopInfo[Feature]
			// if found {
			featuresKeys = append(featuresKeys, Feature)
			// }
		}
	})
	c.OnHTML("td", func(e *colly.HTMLElement) {
		if e.Attr("class") == "URwL2w col col-9-12" {
			Feature := e.DOM.Text()
			features = append(features, Feature)
		}
	})
	c.OnHTML("span", func(e *colly.HTMLElement) {
		if e.Attr("class") == "B_NuCI" {
			LaptopInfo["Laptop"] = e.DOM.Text()
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ...")
	})

	c.Visit(url)
	// fmt.Println("featurekeys: ", featuresKeys)
	// fmt.Println("features: ", features)
	for index, val := range featuresKeys {
		LaptopInfo[val] = features[index]
	}
	// fmt.Println(LaptopInfo)
	LaptopDB = append(LaptopDB, LaptopInfo)
}
func looper(paginatedURL string) {
	c := colly.NewCollector()
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Attr("class") == "_1fQZEK" {
			collectInfo("https://flipkart.com" + e.Attr("href"))
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visited", r.URL)
	})

	c.Visit(paginatedURL)

}
func returnBatteryOrWeight(backup string) float64 {
	var battery float64 = 0
	regex := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	if regex.MatchString(backup) {
		submatchall := regex.FindAllString(backup, -1)
		for _, element := range submatchall {
			batteryStr, err := strconv.ParseFloat(element, 64)
			battery = math.Ceil(float64(batteryStr))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return battery
}
func main() {
	answers := struct {
		Name        string
		Usage       string `survey:"Usage"`
		TouchScreen string `survey:"TouchScreen"`
		Preference  string `survey:"Preference"`
	}{}
	err := survey.Ask(Que1, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	purposes := []string{}

	err3 := survey.Ask(Que2, &purposes)

	if err3 != nil {
		fmt.Println(err.Error())
		return
	}
	splitPurposes := strings.Split(strings.Join(purposes, ", "), ", ")
	mediumPurpose := false
	maxPurpose := false
	for _, v := range splitPurposes {
		if v == "Normal Gaming" || v == "Watching Movies - theatrical experience - yes" || v == "Video Calling" {
			mediumPurpose = true

		} else if v == "Hard-core Gaming" || v == "Heavy Softwares" {
			maxPurpose = true
			break
		}
	}
	// fmt.Println(answers)
	urlInit := "https://www.flipkart.com/search?q=laptop&page="
	for x := 1; x <= 30; x = x + 1 {
		looper(urlInit + strconv.Itoa(x))
	}

	// filePath := "results.json"
	// fmt.Printf("// reading file %s\n", filePath)
	// file, err1 := ioutil.ReadFile(filePath)
	// if err1 != nil {
	// 	fmt.Printf("// error while reading file %s\n", filePath)
	// 	fmt.Printf("File error: %v\n", err1)
	// 	os.Exit(1)
	// }
	// err2 := json.Unmarshal([]byte(file), &LaptopDB)
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }

	// fmt.Println("LaptopDB: ", LaptopDB)

	ProcessorName := []string{"Ryzen 9 Octa Core",
		"Ryzen 7 Octa Core",
		"Ryzen 7 Quad Core",
		"Ryzen 5 Hexa Core",
		"Ryzen 5 Quad Core",
		"Core i9",
		"M1",
		"Core i7",
		"Ryzen 3 Quad Core",
		"Core i5",
		"Ryzen 5 Dual Core",
		"Ryzen 3 Dual Core",
		"MT8183",
		"Core i3",
		"Pentium Quad Core",
		"Pentium Gold",
		"Celeron Dual Core",
		"APU Dual Core A6",
		"Athlon Dual Core",
	}
	var LaptopList []map[string]string
	// Touchscreen, Battery Backup, Weight,
	fmt.Println(returnBatteryOrWeight(answers.Usage))
	for _, val := range LaptopDB {
		if returnBatteryOrWeight(answers.Usage) == 5 {
			_, found := val["Battery"]
			if found {
				if returnBatteryOrWeight(val["Battery"]) <= float64(7) {
					// Battery OK
					if answers.TouchScreen == strings.ToLower(val["Touchscreen"]) {
						// Touchscreen OK
						getWeight := returnBatteryOrWeight(val["Weight"])
						if answers.Preference == "on the heavier side" {
							if getWeight < 5 && getWeight > 2 {
								// Weight OK
								if maxPurpose {
									if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else if mediumPurpose {
									if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else {
									if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								}
							}
						} else if answers.Preference == "normal" {
							if getWeight < 2 && getWeight > 1.5 {
								// Weight OK
								if maxPurpose {
									if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else if mediumPurpose {
									if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else {
									if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								}
							}
						} else {
							if getWeight < 1.5 || getWeight >= 5 {
								// Weight OK
								if maxPurpose {
									if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else if mediumPurpose {
									if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else {
									if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								}
							}
						}
					}
				}
			} else {
				// Battery OK
				if answers.TouchScreen == strings.ToLower(val["Touchscreen"]) {
					// Touchscreen OK
					getWeight := returnBatteryOrWeight(val["Weight"])
					if answers.Preference == "on the heavier side" {
						if getWeight < 5 && getWeight > 2 {
							// Weight OK
							if maxPurpose {
								if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else if mediumPurpose {
								if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else {
								if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							}
						}
					} else if answers.Preference == "normal" {
						if getWeight < 2 && getWeight > 1.5 {
							// Weight OK
							if maxPurpose {
								if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else if mediumPurpose {
								if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else {
								if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							}
						}
					} else {
						if getWeight < 1.5 || getWeight >= 5 {
							// Weight OK
							if maxPurpose {
								if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else if mediumPurpose {
								if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else {
								if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							}
						}
					}
				}

			}

		} else if returnBatteryOrWeight(answers.Usage) == 10 {
			_, found := val["Battery"]
			if found {
				if returnBatteryOrWeight(val["Battery"]) <= float64(12) {
					// Battery OK
					if answers.TouchScreen == strings.ToLower(val["Touchscreen"]) {
						// Touchscreen OK
						getWeight := returnBatteryOrWeight(val["Weight"])
						if answers.Preference == "on the heavier side" {
							if getWeight < 5 && getWeight > 2 {
								// Weight OK
								if maxPurpose {
									if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else if mediumPurpose {
									if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else {
									if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								}
							}
						} else if answers.Preference == "normal" {
							if getWeight < 2 && getWeight > 1.5 {
								// Weight OK
								if maxPurpose {
									if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else if mediumPurpose {
									if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else {
									if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								}
							}
						} else {
							if getWeight < 1.5 || getWeight >= 5 {
								// Weight OK
								if maxPurpose {
									if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else if mediumPurpose {
									if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else {
									if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								}
							}
						}
					}
				}
			} else {
				// Battery OK
				if answers.TouchScreen == strings.ToLower(val["Touchscreen"]) {
					// Touchscreen OK
					getWeight := returnBatteryOrWeight(val["Weight"])
					if answers.Preference == "on the heavier side" {
						if getWeight < 5 && getWeight > 2 {
							// Weight OK
							if maxPurpose {
								if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else if mediumPurpose {
								if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else {
								if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							}
						}
					} else if answers.Preference == "normal" {
						if getWeight < 2 && getWeight > 1.5 {
							// Weight OK
							if maxPurpose {
								if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else if mediumPurpose {
								if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else {
								if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							}
						}
					} else {
						if getWeight < 1.5 || getWeight >= 5 {
							// Weight OK
							if maxPurpose {
								if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else if mediumPurpose {
								if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else {
								if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							}
						}
					}
				}

			}

		} else if returnBatteryOrWeight(answers.Usage) == 12 {
			_, found := val["Battery"]
			if found {
				if returnBatteryOrWeight(val["Battery"]) >= float64(12.1) {
					// Battery OK
					if answers.TouchScreen == strings.ToLower(val["Touchscreen"]) {
						// Touchscreen OK
						getWeight := returnBatteryOrWeight(val["Weight"])
						if answers.Preference == "on the heavier side" {
							if getWeight < 5 && getWeight > 2 {
								// Weight OK
								if maxPurpose {
									if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else if mediumPurpose {
									if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else {
									if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								}
							}
						} else if answers.Preference == "normal" {
							if getWeight < 2 && getWeight > 1.5 {
								// Weight OK
								if maxPurpose {
									if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else if mediumPurpose {
									if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else {
									if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								}
							}
						} else {
							if getWeight < 1.5 || getWeight >= 5 {
								// Weight OK
								if maxPurpose {
									if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else if mediumPurpose {
									if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								} else {
									if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
										// Laptop OK
										LaptopList = append(LaptopList, val)
									}
								}
							}
						}
					}
				}
			} else {
				// Battery OK
				if answers.TouchScreen == strings.ToLower(val["Touchscreen"]) {
					// Touchscreen OK
					getWeight := returnBatteryOrWeight(val["Weight"])
					if answers.Preference == "on the heavier side" {
						if getWeight < 5 && getWeight > 2 {
							// Weight OK
							if maxPurpose {
								if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else if mediumPurpose {
								if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else {
								if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							}
						}
					} else if answers.Preference == "normal" {
						if getWeight < 2 && getWeight > 1.5 {
							// Weight OK
							if maxPurpose {
								if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else if mediumPurpose {
								if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else {
								if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							}
						}
					} else {
						if getWeight < 1.5 || getWeight >= 5 {
							// Weight OK
							if maxPurpose {
								if val["Processor Name"] == ProcessorName[0] || val["Processor Name"] == ProcessorName[1] || val["Processor Name"] == ProcessorName[2] || val["Processor Name"] == ProcessorName[3] || val["Processor Name"] == ProcessorName[4] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else if mediumPurpose {
								if val["Processor Name"] == ProcessorName[5] || val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							} else {
								if val["Processor Name"] == ProcessorName[6] || val["Processor Name"] == ProcessorName[7] || val["Processor Name"] == ProcessorName[8] || val["Processor Name"] == ProcessorName[9] || val["Processor Name"] == ProcessorName[10] || val["Processor Name"] == ProcessorName[11] || val["Processor Name"] == ProcessorName[12] || val["Processor Name"] == ProcessorName[13] || val["Processor Name"] == ProcessorName[14] || val["Processor Name"] == ProcessorName[15] || val["Processor Name"] == ProcessorName[16] || val["Processor Name"] == ProcessorName[17] || val["Processor Name"] == ProcessorName[18] {
									// Laptop OK
									LaptopList = append(LaptopList, val)
								}
							}
						}
					}
				}

			}

		}
	}

	// fmt.Println(ProcessorName)
	fmt.Println("Hey " + answers.Name + "Please find your laptop recommendations at recommendation.json ")
	fmt.Println("==============================================================")
	fmt.Println("Hey " + answers.Name + "Please find all laptops information at results.json ")
	file2, _ := json.MarshalIndent(LaptopList, "", " ")
	// fmt.Println(LaptopList)
	_ = ioutil.WriteFile("recommendation.json", file2, 0644)
	file3, _ := json.MarshalIndent(LaptopDB, "", " ")
	_ = ioutil.WriteFile("results.json", file3, 0644)
}
