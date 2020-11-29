package common

type STATUS uint8

const (
	SUCCESS        STATUS = 100
	PARAM_ERROR    STATUS = 101
	INTERNAL_ERROR STATUS = 102
)
