package systemguard

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kod/kod/interceptor"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
)

type Config struct {
	Enable     bool `json:"enable"`
	SystemLoad uint `json:"system_load"`
	CpuUsage   uint `json:"cpu_usage" default:"80"`
}

func (c Config) Build() interceptor.Interceptor {
	if !c.Enable {
		return nil
	}

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) (err error) {
		if c.SystemLoad > 0 {
			// 获取系统负载
			load, err := getSystemLoad()
			if err != nil {
				return err
			}

			if load > c.SystemLoad {
				return fmt.Errorf("system load exceeded")
			}
		}

		if c.CpuUsage > 0 {
			// 获取CPU使用率
			cpuUsage, err := getCPUUsage()
			if err != nil {
				return err
			}

			if cpuUsage > c.CpuUsage {
				return fmt.Errorf("cpu usage exceeded")
			}
		}

		return invoker(ctx, info, req, reply)
	}
}

// getSystemLoad 获取系统1分钟负载值的百分比
func getSystemLoad() (uint, error) {
	// 读取系统负载
	info, err := load.Avg()
	if err != nil {
		return 0, fmt.Errorf("get load avg failed: %w", err)
	}

	return uint(info.Load1), nil
}

// getCPUUsage 获取CPU使用率百分比
func getCPUUsage() (uint, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, fmt.Errorf("get cpu percent failed: %w", err)
	}

	if len(percent) == 0 {
		return 0, fmt.Errorf("no cpu usage data")
	}

	return uint(percent[0]), nil
}
