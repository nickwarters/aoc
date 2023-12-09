package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"log"
)


var TEMPLATE_PATH = strings.Join([]string{os.Getenv("AOC_PATH"), "templates"}, "/")

func main() {

	var err error = nil

	switch os.Args[1] {
	case "run":
		err = run(os.Args[2:])
	case "setup":
		setup(os.Args[2:])
	case "submit":
		submit(os.Args[2:])
	case "test":
		runTest(os.Args[2:])
	default:
		err = fmt.Errorf("invalid usage")
	}
	

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func getYearDayCombinations(year int, day int) []struct{year int; day int} {
	yearsDays := []struct{
		year int
		day int
		}{}

	var y int
	var d int
	if year == 0 {
		y = 2015
	} else {
		y = year
	}

	if day == 0 {
		d = 1
	} else {
		d = day
	}

	if year == 0 {
		for y < time.Now().Year() + 1 {
			if day == 0 {
				for d < 32 {
					if d > time.Now().Day() && y == time.Now().Year(){
						break
					} else {
						yearsDays = append(yearsDays, struct{year int; day int}{year: y, day: d})
					}
					d++
				}
			} else {
				yearsDays = append(yearsDays, struct{year int; day int}{year: y, day: d})
			}

			y++
		}
	} else {
		if day == 0 {
			for d < 32 {
				if  y == time.Now().Year() && d > time.Now().Day(){
					break
				} else {
					yearsDays = append(yearsDays, struct{year int; day int}{year: y, day: d})
				}
				d++
			}
		} else {
			yearsDays = append(yearsDays, struct{year int; day int}{year: y, day: d})
		}

	}

	return yearsDays
}

func getLangs(lang string) []string {
	var langs []string
	if lang == "all" {
		langs = []string{
			"py",
			"ts",
			"go",
		}
	} else {
		langs = []string{lang}
	}

	return langs
}

func run(args []string) error {
	cmd := flag.NewFlagSet("run", flag.ExitOnError)
	year := cmd.Int("year", time.Now().Year(), "year")
	day := cmd.Int("day", time.Now().Day(), "day")
	// timer := cmd.Bool("time", false, "time")
	lang := cmd.String("lang", "all", "lang")
	cmd.Parse(args)

	if (*year == time.Now().Year() && *day > time.Now().Day()) || *year > time.Now().Year(){
		return fmt.Errorf("can't run a day in the future ... put the kettle on and come back in %d days", (365 * (*year - time.Now().Year() )) + (*day - time.Now().Day()))
	}

	// yearsDays := getYearDayCombinations(*year, *day)
	// langs := getLangs(*lang)

	cmds := []string{}

	var c string
	
	c = createCommand(d, l, true)
	cmds = append(cmds, c)

	c = createCommand(d, l, false)
	cmds = append(cmds, c)
	

	for _, c := range cmds {
		err := runCmd(c); if err != nil {
			return err
		}
		
	}

	return nil
}

func createCommand(d struct{year int; day int}, l string, s bool) string {

	var scriptPath string
	var cmd string

	if s {
		var lang string
		if l == "py" {
			lang = "python"
			
		} else if l == "ts" {
			lang = "typescript"
		} else if l == "go" {
			lang = "go"
		}

		yearDayLang := fmt.Sprintf(">> %d - day %d - %s <<", d.year, d.day, lang)
		startEnd := strings.Repeat("=", 60 - len(yearDayLang) / 2)
		return strings.Join([]string{"clear", "&&", "echo", startEnd, yearDayLang, startEnd}, " ")
	}

	if l == "py" {
		scriptPath = strings.Join([]string{
			os.Getenv("AOC_PATH"), "python", fmt.Sprint(d.year), fmt.Sprintf("day%02d", d.day), "main.py",
		}, "/")
		cmd = fmt.Sprintf("python %s", scriptPath)
	} else if l == "ts" {
		scriptPath = strings.Join([]string{
			os.Getenv("AOC_PATH"), "ts", fmt.Sprint(d.year), fmt.Sprintf("day%02d", d.day), "main.ts",
		}, "/")
		cmd = fmt.Sprintf("bun run %s", scriptPath)
	} else if l == "go" {
		scriptPath = strings.Join([]string{
			os.Getenv("AOC_PATH"), "python", fmt.Sprint(d.year), fmt.Sprintf("day%02d", d.day), "main.py",
		}, "/")
		cmd = fmt.Sprintf("go build %s && %s", scriptPath, strings.Replace(scriptPath, ".go", "", -1))
	}

	return cmd
}

func getInput(year int, day int) (string, error) {
	session_cookie := os.Getenv("AOC_COOKIE")
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.New("could not make request to AOC")
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: session_cookie})
	client := &http.Client{}
    resp, err := client.Do(req)
	if err != nil {
		log.Print("there was an error")
		return "", err
	}
	defer resp.Body.Close()

	fmt.Printf("StatusCode=%d", resp.StatusCode)

	if resp.StatusCode != 200 {
			err = errors.New(url +
					"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
			return "", err
	} 

    input, err := io.ReadAll(resp.Body)
	return string(input), err
}

func setup(args []string) {
	cmd := flag.NewFlagSet("setup", flag.ExitOnError)
	year := cmd.Int("year", time.Now().Year(), "year")
	day := cmd.Int("day", time.Now().Day(), "day")
	cmd.Parse(args)

	input, err := getInput(*year, *day)
	if err != nil {
		fmt.Printf("error=%s\n", err)
	}

	fmt.Println("setting up with args:")
	fmt.Printf("year=%d\n", *year)
	fmt.Printf("day=%d\n", *day)
	fmt.Printf("inout=%s\n", input)
}

func submit(args []string) error {
	cmd := flag.NewFlagSet("submit", flag.ExitOnError)
	year := cmd.Int("year", time.Now().Year(), "year")
	day := cmd.Int("day", time.Now().Day(), "day")
	part := cmd.String("part", "1", "part")
    answer := cmd.String("answer", "", "answer")
	cmd.Parse(args)

    url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/answer", *year, *day)
    session_cookie := os.Getenv("AOC_COOKIE")

    data, err := json.Marshal(map[string]string{"part": *part, "answer": *answer})
    if err != nil {
        fmt.Println("could not package up the answer")
        return err
    }
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        fmt.Println("could not create request")
        return err
    }
    req.AddCookie(&http.Cookie{Name: "session", Value: session_cookie})
    req.Header.Set("data", string(data))
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("failed to post answer")
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        err = fmt.Errorf("request failed, StatusCode=%d", resp.StatusCode)
        return err
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return errors.New("could not read response body")
    }
    
    goodAnswer := "That's the right answer!"
    badAnswer := "That's not the right answer"
    wrongLevel := "You don't seem to be solving the right level"
    tooSoon := "You gave an answer too recently"
    

    if strings.Contains(string(body), goodAnswer) {
        fmt.Println(goodAnswer)
        return nil
    }

    for _, e := range []string{badAnswer, wrongLevel, tooSoon} {
        if strings.Contains(string(body), e) {
            return errors.New(e)
        }
    }

    return errors.New(string(body))
}

func runTest(args []string) {
	cmd := flag.NewFlagSet("test", flag.ExitOnError)
	year := cmd.Int("year", time.Now().Year(), "year")
	day := cmd.Int("day", time.Now().Day(), "day")
	part := cmd.Int("part", 1, "part")
	_ = cmd.Args()
	cmd.Parse(args)

	fmt.Println("submitting with args:")
	fmt.Printf("year=%d\n", *year)
	fmt.Printf("day=%d\n", *day)
	fmt.Printf("part=%d\n", *part)
}


func runCmd(command string) error {
	cmd_arr := strings.Split(command, " ")
	cmd := exec.Command(cmd_arr[0], cmd_arr[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
