package bootstrap

import (
	"charites/global"
	"errors"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

const (
	_defaultStartTime = "2021-12-31"
)

func setupSnowflake(startTime string, machineId int64) error {
	if machineId < 0 {
		return errors.New("snowflake need machineId")
	}
	if len(startTime) == 0 {
		startTime = _defaultStartTime
	}
	var st time.Time
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	sf.Epoch = st.UnixNano() / 100_0000          // 时间戳，开始时间 69年
	global.SnowNode, err = sf.NewNode(machineId) // 机器编号，1024
	if err != nil {
		return err
	}
	return nil
}
