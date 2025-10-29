package processor

import "github.com/nichuanfang/gymdl/core/domain"

// 顶级接口定义

type Processor interface {
	Handle(link string) (string, error) // 处理链接并返回结果
	Category() domain.ProcessorCategory //所属分类
}
