package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/ttacon/chalk"
)

func showOpenTasks(filter string) {
	searchDir := tc.TaskDir

	// check filter for project
	if strings.HasPrefix(filter, "+") {
		searchDir = filepath.Join(tc.TaskDir, filter[1:])
	}

	fmt.Printf("  %s : %-10s : %-50s : %s\n", "ID", "Project", "Task", "Age")
	fmt.Println("-----:------------:----------------------------------------------------:---------------")
	filepath.Walk(searchDir, func(fn string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Warn("Open tasks walk", err)
			return err
		}
		if !fi.IsDir() {
			if filepath.Ext(fn) == ".toml" {
				if !strings.Contains(fn, ".done") {
					displayTaskFromFile(fn)
				}
			}
		}
		return nil
	})
}

func showTask(taskID int) {
	task, err := getTaskByID(taskID)
	log.FatalErrNotNil(err, "Task not found")
	fmt.Print(getColorForProject(task.Project))
	fmt.Println("Name         :", task.Name)
	fmt.Println("Project      :", task.Project)
	fmt.Println("Created On   :", task.CreationDate.Format("Jan 2, 2006"))
	fmt.Print(chalk.Reset)

	if !task.CompletionDate.IsZero() {
		fmt.Print(chalk.Green)
		fmt.Println("Completed On :", task.CompletionDate.Format("Jan 2, 2006"))
		fmt.Print(chalk.Reset)
	}
	if len(task.Notes) > 0 {
		fmt.Println("\nNotes:")
		for _, note := range task.Notes {
			fmt.Println("    Date:", note.CreationDate.Format("Jan 2, 2006"))
			fmt.Println("   ", note.Entry)
			fmt.Println("----")
		}
	}
}

// displayTaskFromFile reads a task file and displays an entry
func displayTaskFromFile(filename string) {
	task, err := readTaskFromFilename(filename)
	if err == nil {
		fmt.Print(getColorForProject(task.Project))
		fmt.Printf("%4d : %-10s : %-50s : %s\n", task.ID, trunc(task.Project, 10), trunc(task.Name, 48), humanize.Time(task.CreationDate))
		fmt.Print(chalk.Reset)
	}
}

// displayCompletedTaskFromFile reads a task file and displays an entry
func displayCompletedTaskFromFile(filename string) {
	task, err := readTaskFromFilename(filename)
	if err == nil {
		fmt.Print(getColorForProject(task.Project))
		fmt.Printf("%4d : %-10s : %-50s : %s\n", task.ID, trunc(task.Project, 10), trunc(task.Name, 48), humanize.Time(task.CompletionDate))
		fmt.Print(chalk.Reset)
	}
}

// showReport displays completed tasks
func showCompletedReport(filter string) {
	searchDir := tc.TaskDir

	// check filter for project
	if strings.HasPrefix(filter, "+") {
		searchDir = filepath.Join(tc.TaskDir, filter[1:])
	}

	fmt.Printf("  %s : %-10s : %-50s : %s\n", "ID", "Project", "Task", "Completed")
	fmt.Println("-----:------------:----------------------------------------------------:---------------")
	filepath.Walk(searchDir, func(fn string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Warn("Open tasks walk", err)
			return err
		}
		if !fi.IsDir() {
			if filepath.Ext(fn) == ".toml" {
				if strings.Contains(fn, ".done") {
					displayCompletedTaskFromFile(fn)
				}
			}
		}
		return nil
	})
}

func trunc(s string, l int) string {
	if len(s) > l {
		return s[0:l-1] + "..."
	}
	return s
}

func getColorForProject(project string) chalk.Color {
	val := 0
	colors := []chalk.Color{
		chalk.Green, chalk.Yellow, chalk.Blue, chalk.Magenta, chalk.Cyan,
	}

	for _, s := range project {
		val = val + int(s)
	}

	index := val % len(colors)
	return colors[index]
}
