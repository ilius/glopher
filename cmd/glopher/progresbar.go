package main

import (
	"github.com/cheggaaa/pb"
)

func NewCmdProgressBar() *CmdProgressBar {
	pbar := &CmdProgressBar{
		pb: pb.New(0),
	}
	return pbar
}

type CmdProgressBar struct {
	pb *pb.ProgressBar
}

func (p *CmdProgressBar) SetTotal(total int) {
	p.pb.Total = int64(total)
}
func (p *CmdProgressBar) Start(msg string) {
	p.pb.Start()
}

func (p *CmdProgressBar) Update(index int) {
	p.pb.Set(index)
}
func (p *CmdProgressBar) SetMessage(msg string) {
	// p.pb.Prefix(msg)
}
