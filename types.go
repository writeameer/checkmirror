package main

type Config struct {
	SqlServerHost string
	SqlServerPort uint16
	ListenPort    uint16
}

type JsonResponse struct {
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Time    string      `json:"time"`
	Version string      `json:"version"`
}

type JsonDbMirroring struct {
	Name          string `json:"name"`
	MirroringRole string `json:"mirroring_role"`
}

type JsonMirroring struct {
	OverallMirroringRole string             `json:"overall_mirroring_role"`
	DatabasesMirroring   []*JsonDbMirroring `json:"databases_mirroring"`
}
