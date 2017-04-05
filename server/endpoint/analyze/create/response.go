package create

import (
	"github.com/xh3b4sd/wafer/server/endpoint/analyze/create/response"
)

type Response struct {
	Body    response.Body    `json:"body"`
	Analyze response.Analyze `json:"analyze"`
}
