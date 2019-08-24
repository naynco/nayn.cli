package main

import (
	"fmt"

	"log"

	"github.com/blang/semver"
	"github.com/fatih/color"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"github.com/ungerik/go-rss"
)

const version = "1.0.0"

func doSelfUpdate() {
	v := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(v, "naynco/nayn.cli")
	if err != nil {
		log.Println("Binary update failed:", err)
		return
	}
	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		log.Println("Current binary is the latest version", version)
	} else {
		log.Println("Successfully updated to version", latest.Version)
		log.Println("Release note:\n", latest.ReleaseNotes)
	}
}

func main() {

	var cmdAll = &cobra.Command{
		Use:   "all ",
		Short: "Read all news",
		Long:  `read read read`,
		Run: func(cmd *cobra.Command, args []string) {
			nayn := color.New(color.FgWhite, color.Bold, color.Underline).SprintFunc()
			co := color.New(color.FgYellow, color.Bold).SprintFunc()
			ntime := color.New(color.FgWhite, color.Bold).SprintFunc()
			narrow := color.New(color.FgRed).SprintFunc()
			ntitle := color.New(color.FgGreen).SprintFunc()

			channel, err := rss.Read("https://nayn.co/feed/")
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("\n" + nayn("NAYN") + co(".CO") + "\n")

			for _, item := range channel.Item {
				s, err := item.PubDate.Parse()
				if err != nil {
					fmt.Print(err)
				}

				fmt.Println(ntime(s.Format("15:04")), narrow(">"), ntitle(item.Title))
			}

			lbd, err := channel.LastBuildDate.Parse()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("\n", lbd.Format("2006-01-02 15:04:05"))

		},
	}

	var rootCmd = &cobra.Command{Use: "nayn", Version: version}
	rootCmd.AddCommand(cmdAll)
	rootCmd.Execute()
}
