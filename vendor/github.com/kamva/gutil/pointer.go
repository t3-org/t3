package gutil

import "time"

func NewBool(b bool) *bool           { return &b }
func NewInt(i int) *int              { return &i }
func NewInt32(i int32) *int32        { return &i }
func NewInt64(i int64) *int64        { return &i }
func NewString(s string) *string     { return &s }
func NewTime(t time.Time) *time.Time { return &t }
