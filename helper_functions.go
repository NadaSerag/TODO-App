package main

import "time"

func PriorityValid(priority string) bool {
	if priority != "High" && priority != "Medium" && priority != "Low" {
		return false
	}
	return true
}

func CategoryValid(cat string) bool {
	if cat != "Work" && cat != "Study" && cat != "Personal" {
		return false
	}
	return true
}

func DueDateValid(date *time.Time) bool {
	if (*date).Before(time.Now()) {
		return false
	}
	return true
}
