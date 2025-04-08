package deeplx

type APIVersion uint8

const (
	FreeAPI APIVersion = iota
	ProAPI
	OfficialAPI
)
