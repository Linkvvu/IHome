package model

// image capatcha record
type ImageCaptchaPair struct {
	Uuid string `cache:"key"`
	Code string `cache:"val"`
}

type SmsCaptchaPair struct {
	PhoneNum string `cache:"key"`
	Code     string `cache:"val"`
}
