package main

import (
	"strconv"
	"time"
)

func formatTimeRelativeLongest(tm time.Time) string {
	tmSince := time.Since(tm)
	hrs := tmSince.Hours()
	yrs := int(hrs / (365 * 24))
	if yrs > 0 {
		if yrs == 1 {
			return "1 year ago"
		}
		return strconv.Itoa(yrs) + " years ago"
	}
	months := int(hrs / (30 * 24))
	if months > 0 {
		if months == 1 {
			return "1 month ago"
		}
		return strconv.Itoa(months) + " months ago"
	}
	days := int(hrs / 24)
	if days > 0 {
		if days == 1 {
			return "1 day ago"
		}
		return strconv.Itoa(days) + " days ago"
	}
	hrsInt := int(hrs)
	if hrsInt > 0 {
		if hrsInt == 1 {
			return "1 hour ago"
		}
		return strconv.Itoa(hrsInt) + " hours ago"
	}
	mins := int(hrs*60) % 60
	if mins > 0 {
		if mins == 1 {
			return "1 min ago"
		}
		return strconv.Itoa(mins) + " mins ago"
	}
	secs := int(hrs*60*60) % 60
	if secs > 0 {
		if secs == 1 {
			return "1 sec ago"
		}
		return strconv.Itoa(secs) + " secs ago"
	}
	return "now"
}

func formatTimeRelativeLong(tm time.Time) string {
	tmSince := time.Since(tm)
	hrs := tmSince.Hours()
	yrs := int(hrs / (365 * 24))
	if yrs > 0 {
		if yrs == 1 {
			return "1 year"
		}
		return strconv.Itoa(yrs) + " years"
	}
	months := int(hrs / (30 * 24))
	if months > 0 {
		if months == 1 {
			return "1 month"
		}
		return strconv.Itoa(months) + " months"
	}
	days := int(hrs / 24)
	if days > 0 {
		if days == 1 {
			return "1 day"
		}
		return strconv.Itoa(days) + " days"
	}
	hrsInt := int(hrs)
	if hrsInt > 0 {
		if hrsInt == 1 {
			return "1 hour"
		}
		return strconv.Itoa(hrsInt) + " hours"
	}
	mins := int(hrs*60) % 60
	if mins > 0 {
		if mins == 1 {
			return "1 min"
		}
		return strconv.Itoa(mins) + " mins"
	}
	secs := int(hrs*60*60) % 60
	if secs > 0 {
		if secs == 1 {
			return "1 sec"
		}
		return strconv.Itoa(secs) + " secs"
	}
	return "now"
}

func formatTimeMixedLongest(tm time.Time) string {
	tmSince := time.Since(tm)
	hrs := tmSince.Hours()
	yrs := int(hrs / (365 * 24))
	if yrs > 0 {
		if yrs == 1 {
			return "1 year ago"
		}
		return strconv.Itoa(yrs) + " years ago"
	}
	months := int(hrs / (30 * 24))
	if months > 0 {
		if months == 1 {
			return "1 month ago"
		}
		return strconv.Itoa(months) + " months ago"
	}
	days := int(hrs / 24)
	if days > 0 {
		if days == 1 {
			return "1 day ago"
		}
		return strconv.Itoa(days) + " days ago"
	}
	hrsInt := int(hrs)
	mins := int(hrs*60) % 60
	return strconv.Itoa(hrsInt) + ":" + strconv.Itoa(mins)
}

func formatTimeRelativeShort(tm time.Time) string {
	tmSince := time.Since(tm)
	hrs := tmSince.Hours()
	yrs := int(hrs / (365 * 24))
	if yrs > 0 {
		return strconv.Itoa(yrs) + " Y"
	}
	months := int(hrs / (30 * 24))
	if months > 0 {
		return strconv.Itoa(months) + " M"
	}
	days := int(hrs / 24)
	if days > 0 {
		return strconv.Itoa(days) + " D"
	}
	hrsInt := int(hrs)
	if hrsInt > 0 {
		return strconv.Itoa(hrsInt) + " h"
	}
	mins := int(hrs*60) % 60
	if mins > 0 {
		return strconv.Itoa(mins) + " m"
	}
	secs := int(hrs*60*60) % 60
	if secs > 0 {
		return strconv.Itoa(secs) + " s"
	}
	return "now"
}

func formatTimeMixedShort(tm time.Time) string {
	tmSince := time.Since(tm)
	hrs := tmSince.Hours()
	yrs := int(hrs / (365 * 24))
	if yrs > 0 {
		return strconv.Itoa(yrs) + " Y"
	}
	months := int(hrs / (30 * 24))
	if months > 0 {
		return strconv.Itoa(months) + " M"
	}
	days := int(hrs / 24)
	if days > 0 {
		return strconv.Itoa(days) + " D"
	}
	hrsInt := int(hrs)
	mins := int(hrs*60) % 60
	return strconv.Itoa(hrsInt) + ":" + strconv.Itoa(mins)
}
