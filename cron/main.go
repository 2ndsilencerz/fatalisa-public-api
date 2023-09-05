package main

import (
	praySchedulePkpuSvc "fatalisa-public-api/service/pray-schedule/pkpu"
)

// This is a onetime run to download all pray schedule data
// use this with cron
func main() {
	praySchedulePkpuSvc.PrayScheduleDownload()
}
