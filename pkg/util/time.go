package util

import "time"

/**
* get the current time string
*
* @return: the current time string
* @example:
* - NowString() -> "2026-01-06 15:04:05"
 */
func NowString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
