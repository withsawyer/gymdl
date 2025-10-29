package utils

import (
    "time"

    "github.com/go-ping/ping"
)

// http相关工具

func CheckWrapperConnection(containerName string) bool {
    // 你必须启用wrapper
    pinger, err := ping.NewPinger(containerName)
    if err != nil {
        ErrorWithFormat("创建 pinger 失败: %v\n", err)
        return false
    }

    pinger.Count = 3              // 发送 3 个包
    pinger.Timeout = 5 * time.Second
    pinger.SetPrivileged(true)    // Linux 下需要 root 权限发送 ICMP

    InfoWithFormat("开始 ping %s...\n", containerName)
    err = pinger.Run() // 阻塞
    if err != nil {
        ErrorWithFormat("Ping 出错: %v\n", err)
        return false
    }

    stats := pinger.Statistics() // 获取统计信息
    if stats.PacketsRecv > 0 {
        InfoWithFormat("%s 可达, 丢包率: %.2f%%\n", containerName, stats.PacketLoss)
    } else {
        ErrorWithFormat("%s 不可达\n", containerName)
        return false
    }
    return true
}