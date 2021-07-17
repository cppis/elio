package main

import (
	"github.com/cppis/elio"
)

// LogonRequest logon request
type LogonRequest struct {
	UID elio.UID `struc:"int64,little"`
}

// LogonResponse logon response
type LogonResponse struct {
	Result int32 `struc:"int32,little"`
}

// LogoffRequest logoff request
type LogoffRequest struct {
}

// LogoffResponse logoff response
type LogoffResponse struct {
	Result int32    `struc:"int32,little"`
	UID    elio.UID `struc:"int64,little"`
}

// EnterRequest enter request
type EnterRequest struct {
	LenRoomID int16 `struc:"int16,little,sizeof=RoomID"`
	RoomID    string
}

// EnterResponse enter response
type EnterResponse struct {
	UID    elio.UID `struc:"int64,little"`
	Result int32    `struc:"int32,little"`
}

// LeaveRequest leave request
type LeaveRequest struct {
	UID       elio.UID `struc:"int64,little"`
	LenRoomID int16    `struc:"int16,little,sizeof=RoomID"`
	RoomID    string
}

// LeaveResponse leave response
type LeaveResponse struct {
	Result    int32    `struc:"int32,little"`
	UID       elio.UID `struc:"int64,little"`
	LenRoomID int16    `struc:"int16,little,sizeof=RoomID"`
	RoomID    string
}

// Message message
type Message struct {
	LenMessage int16 `struc:"int16,little,sizeof=Message"`
	Message    []byte
}
