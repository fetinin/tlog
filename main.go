package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/andygrunwald/go-jira"
	"github.com/pterm/pterm"
)

func main() {
	conf, err := LoadConfig()
	if err != nil {
		fmt.Printf("cannot load config: %s\n", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: todo <time> <task> [day] [comment]")
		return
	}

	timeLogInput := os.Args[1]
	timeLog, err := convertToTimeLog(timeLogInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	taskInput := os.Args[2]
	jiraID, err := convertToTask(taskInput, conf.DefaultProject, conf.TaskAliases)
	if err != nil {
		fmt.Println(err)
		return
	}

	dayInput := safeGet(os.Args, 3)
	logDay, err := convertToDay(dayInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	logComment := safeGet(os.Args, 4)

	tp := jira.BasicAuthTransport{
		Username: conf.JiraLogin,
		Password: conf.JiraPassword,
	}
	jiraClient, err := jira.NewClient(tp.Client(), conf.JiraURL)
	if err != nil {
		panic(err)
	}

	spinner, _ := pterm.DefaultSpinner.Start("Logging time... (JIRA might be slow🐌)")
	wl, _, err := jiraClient.Issue.AddWorklogRecord(jiraID, &jira.WorklogRecord{
		Comment:          logComment,
		Started:          toPtr(jira.Time(logDay)),
		TimeSpentSeconds: int(timeLog.Seconds()),
	})
	if err != nil {
		spinner.Fail(err.Error())
		return
	}

	spinner.Success(fmt.Sprintf(
		"Created worklog as %s on issue %s for %d munutes: %s",
		wl.Author.Name, jiraID, wl.TimeSpentSeconds/60, wl.Self,
	))
}

func convertToTask(input string, defaultProject string, aliases map[string]string) (string, error) {
	if task, ok := aliases[input]; ok {
		return task, nil
	}

	// if input is number, assume it is issue key
	if _, err := strconv.Atoi(input); err == nil {
		if defaultProject == "" {
			return "", fmt.Errorf("if ussing issue number, set DefaultProject in config")
		}

		return fmt.Sprintf("%s-%s", defaultProject, input), nil
	}

	return input, nil
}

func convertToDay(input string) (time.Time, error) {
	todayStart := time.Now().Truncate(24 * time.Hour).UTC()

	input = strings.ToLower(input)
	if input == "" || input == "today" {
		return todayStart, nil
	}

	if input == "yesterday" {
		return todayStart.Add(-24 * time.Hour), nil
	}

	var day time.Weekday
	switch input {
	case "monday":
		day = time.Monday
	case "tuesday":
		day = time.Tuesday
	case "wednesday":
		day = time.Wednesday
	case "thursday":
		day = time.Thursday
	case "friday":
		day = time.Friday
	case "saturday":
		day = time.Saturday
	case "sunday":
		day = time.Sunday
	default:
		day = -1 // NoneDay
	}

	if day != -1 {
		return todayStart.Add(time.Duration(time.Now().Weekday()-day) * 24 * time.Hour), nil
	}

	if d, err := strconv.Atoi(input); err == nil {
		y, m, _ := time.Now().Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
	}

	return time.Time{}, fmt.Errorf("day of the week, or day number expected")
}

func convertToTimeLog(inputTime string) (time.Duration, error) {
	duration, err := time.ParseDuration(inputTime)
	return duration, err
}

type Config struct {
	JiraURL        string            `toml:"JiraURL"`
	JiraLogin      string            `toml:"JiraLogin"`
	JiraPassword   string            `toml:"JiraPassword"`
	DefaultProject string            `toml:"DefaultProject"`
	TaskAliases    map[string]string `toml:"TaskAliases"`
}

func LoadConfig() (Config, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("cannot obtain home dir: %s\n", err)
	}
	homeConfig := filepath.Join(dirname, ".time_logger_conf.toml")

	if _, err := os.Stat(homeConfig); err != nil {
		cfg := setupConfig()
		err := writeConfig(cfg, homeConfig)
		if err != nil {
			return Config{}, fmt.Errorf("create config: %w", err)
		}
		fmt.Printf("Config created saved at: %s", homeConfig)
	}

	var cfg Config
	if _, err := toml.DecodeFile(homeConfig, &cfg); err != nil {
		return Config{}, fmt.Errorf("cannot decode config file: %s", err)
	}

	return cfg, nil
}

func toPtr[T any](v T) *T {
	return &v
}

func safeGet(arr []string, index int) string {
	if index >= len(arr) {
		return ""
	}
	return arr[index]
}

func setupConfig() Config {
	cfg := Config{}
	area, _ := pterm.DefaultArea.Start()
	area.Update(
		pterm.DefaultSection.Sprint("Hello there 👋"),
		pterm.LightBlue("Let's perform some basic setup. Enter you JIRA username below:"),
	)

	cfg.JiraLogin = readLine()

	area.Update(pterm.LightBlue("Got it.\nYour login is: "), pterm.Yellow(cfg.JiraLogin))
	time.Sleep(2 * time.Second)

	area.Update(pterm.LightBlue("🤫 Now enter your password: "))

	cfg.JiraPassword = readLine()

	area.Update(pterm.LightBlue("Almost done!\n Now enter JIRA url: "))
	cfg.JiraURL = readLine()

	area.Stop()
	return cfg
}

func writeConfig(cfg Config, path string) error {
	tmpl := `
JiraURL = "%s"
JiraLogin = "%s"
JiraPassword = "%s"
DefaultProject = ""

[ TaskAliases ]
`
	tmpl = strings.TrimSpace(tmpl)
	out := fmt.Sprintf(tmpl, cfg.JiraURL, cfg.JiraLogin, cfg.JiraPassword)
	return os.WriteFile(path, []byte(out), 0644)
}

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
