package util

type SlsConf struct {
	FrameRate int `json:"frame_rate"`
	CahkSize  int `json:"chak_size"`
}

type Container struct {
	Data []string `json:"data"`
}
