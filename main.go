package main

import (
	"fmt"
	"log"
	"time"

	"github.com/blang/semver"
	"github.com/fatih/color"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"github.com/ungerik/go-rss"

	mex "github.com/monocash/exchange-rates/pkg/exchanger"
	"github.com/monocash/exchange-rates/pkg/swap"
)

const version = "1.0.4"

func main() {
	var cmdUpdate = &cobra.Command{
		Use:   "update ",
		Short: "Update ",
		Long:  `up up up`,
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

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
			nlink := color.New(color.FgBlue).SprintFunc()

			channel, err := rss.Read("https://nayn.co/feed/")
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("\n" + nayn("NAYN") + co(".CO") + "\n")

			loc, _ := time.LoadLocation("Europe/Istanbul")
			for _, item := range channel.Item {
				s, err := item.PubDate.Parse()
				if err != nil {
					fmt.Print(err)
				}

				fmt.Println(ntime(s.In(loc).Format("15:04")), narrow(">"), ntitle(item.Title), nlink(item.GUID))
			}

			lbd, err := channel.LastBuildDate.Parse()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("\n", narrow("Son güncelleme :"), lbd.In(loc).Format("2006-01-02 15:04:05"))

			Swap1 := swap.NewSwap()
			Swap1.AddExchanger(mex.NewyahooAPI(nil)).Build()
			usdToTryRate := Swap1.Latest("USD/TRY")

			Swap2 := swap.NewSwap()
			Swap2.AddExchanger(mex.NewyahooAPI(nil)).Build()
			eurToTryRate := Swap2.Latest("EUR/TRY")

			fmt.Println("\n", "USD", narrow(usdToTryRate.GetRateValue()), nlink("EUR"), narrow(eurToTryRate.GetRateValue()))

			fmt.Println("\n", "sürüm", version)

		},
	}

	var rootCmd = &cobra.Command{Use: "nayn", Version: version}
	rootCmd.AddCommand(cmdAll)
	rootCmd.AddCommand(cmdUpdate)
	rootCmd.Execute()
}
