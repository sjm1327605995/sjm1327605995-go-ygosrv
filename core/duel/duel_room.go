package duel

import "sync"

var allDuelRoom sync.Map

type DuelRoom struct {
	locker     sync.RWMutex
	Players    map[string]*DuelPlayer
	HostPlayer *DuelPlayer
}

func JoinOrCreateRoom(dp *DuelPlayer) *DuelRoom {
	var room = &DuelRoom{}
	val, has := allDuelRoom.LoadOrStore(dp.Pass, &room)
	if has {
		room = val.(*DuelRoom)
	}
	room.locker.Lock()
	defer room.locker.Unlock()
	if room.Players == nil {
		room.Players = make(map[string]*DuelPlayer, 2)
		room.HostPlayer = dp
	}
	oldDp, has := room.Players[dp.Name]
	if has {
		if oldDp == room.HostPlayer {
			oldDp = dp
		}
		_ = oldDp.Conn.Close()
	} else {
		dp.Pos = uint8(len(room.Players))
	}
	room.Players[dp.Name] = dp
	return room
}
func (r *DuelRoom) TypeChange(dp *DuelPlayer) uint8 {
	if r.HostPlayer == dp {
		return 0x10
	}
	return 0
}
