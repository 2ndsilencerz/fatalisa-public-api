package main

import (
	prayScheduleSvc "fatalisa-public-api/service/pray-schedule"
)

// This is a onetime run to download all pray schedule data
// use this with cron
func main() {
	prayScheduleSvc.PrayScheduleDownload()
}
