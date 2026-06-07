package util

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
)

/**
* check if the string is empty
*
* @param val: the string to check
* @return: true if the string is empty, false otherwise
* @example:
* - IsEmptyStr("") -> true
* - IsEmptyStr("hello") -> false
 */
func IsEmptyStr(val string) bool {
	return len(strings.TrimSpace(val)) == 0
}

/**
* generate a new server ID string
*
* @param prefix: the prefix of the server ID
* @return: the generated server ID
* @example:
* - NewServerID("svr") -> "svr-123e4567-e89b-12d3-a456-426614174000"
 */
func NewServerID(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, ksuid.New().String())
}

/**
* generate a new UUID string
*
* @param slash: if true, the UUID will be generated with a slash
* @return: the generated UUID
* @example:
* - NewUUID(true) -> "123e4567-e89b-12d3-a456-426614174000"
* - NewUUID(false) -> "123e4567e89b12d3a456426614174000"
 */
func NewUUID(slash bool) string {
	if !slash {
		return strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return uuid.New().String()
}
