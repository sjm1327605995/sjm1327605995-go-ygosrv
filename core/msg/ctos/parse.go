package ctos

import (
	"go-ygosrv/core/msg/host"
	"go-ygosrv/utils"
)

type PlayerInfo struct {
	Name string
}

const (
	StrLimit = 40
)

func (p *PlayerInfo) Parse(b *utils.BitReader) (err error) {
	// 将二进制数组转换为字符串
	p.Name, err = utils.UTF16ToStr(b.Next(StrLimit))
	return
}

type TPResult struct {
	res uint8
}

func (h *TPResult) Parse(b *utils.BitReader) (err error) {
	return utils.GetData(b, &h.res)

}

type CreateGame struct {
	Info host.HostInfo
	Name string
	Pass string
}

func (h *CreateGame) Parse(b *utils.BitReader) (err error) {
	err = h.Info.Parse(b)
	if err != nil {
		return err
	}
	h.Name, err = utils.UTF16ToStr(b.Next(StrLimit))
	if err != nil {
		return
	}
	h.Pass, err = utils.UTF16ToStr(b.Next(StrLimit))
	if err != nil {
		return
	}
	return
}

type JoinGame struct {
	Version uint16
	Align   uint16
	GameId  uint32
	Pass    string
}

// Pass: [40] - 房间密码
func (h *JoinGame) Parse(b *utils.BitReader) (err error) {
	err = utils.GetData(b, &h.Version, &h.Align, &h.GameId)
	if err != nil {
		return
	}
	h.Pass, err = utils.UTF16ToStr(b.Next(StrLimit))
	if err != nil {
		return
	}
	return
}

type Kick struct {
	Pos uint16
}

func (h *Kick) Parse(b *utils.BitReader) (err error) {
	return utils.GetData(b, &h.Pos)

}