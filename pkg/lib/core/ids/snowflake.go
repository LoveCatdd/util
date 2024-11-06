package ids

import (
	"strconv"
	"sync"
	"time"

	"github.com/LoveCatdd/util/pkg/lib/core/log"
)

type snowFlake struct {
	sync.Mutex         // 锁
	timestamp    int64 // 时间戳 ，毫秒
	workerid     int64 // 工作节点
	datacenterid int64 // 数据中心机房id
	sequence     int64 // 序列号
}

var s *snowFlake = new(snowFlake)

func (s *snowFlake) nextId() int64 {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixNano() / ie6 // 转毫秒
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.timestamp {
				now = time.Now().UnixNano() / ie6
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		log.Error("epoch must be between 0 and %d", timestampMax-1)
		return 0
	}
	s.timestamp = now
	r := int64(

		(t << timestampShift) |
			(s.datacenterid << datacenteridShift) |
			(s.workerid << workeridShift) |
			(s.sequence),
	)
	return r
}

func (s *snowFlake) nextStr() string {
	return strconv.FormatInt(s.nextId(), 10)
}

func GenerateId() int64 {
	return s.nextId()
}

func GenerateStr() string {
	return s.nextStr()
}
