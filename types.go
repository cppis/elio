package elio

import "strconv"

// UIDInvalid uid invalid
var UIDInvalid = UID(0)

// UID uid
// 사용자 ID 정의 타입
// KBO(CPB)의 uid, MLB의 vid 값
type UID int64

// ToInt uid to int
func (u *UID) ToInt() int64 {
	return int64(*u)
}

// FromInt uid from int
func (u *UID) FromInt(i int64) {
	*u = UID(i)
}

// ToString uid to string
func (u *UID) ToString() string {
	return strconv.FormatInt(int64(*u), 10)
}

// FromString uid from string
func (u *UID) FromString(s string) {
	n, _ := strconv.ParseInt(s, 10, 64)
	*u = UID(n)
}

// IsInvalid uid is invalid
func (u UID) IsInvalid() bool {
	if UIDInvalid == u {
		return true
	}

	return false
}
