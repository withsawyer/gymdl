package core

import "testing"

// 测试解析链接返回对应的音乐平台处理器
func TestParseLink(t *testing.T) {
	text := "https://www.youtube.com/watch?v=YJ7IrWCEPyo&list=RD91zzlZbwATA&index=2"

	url, handler := ParseLink(text)
	if handler == nil {
		t.Error("handler is nil")
	}
	handler.HandlerMusic(url)
}
