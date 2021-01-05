package units

import (
	"Caesar/internal/library/director"
	"testing"

	"Caesar/internal/pkg/finger"
)

// 指纹识别单元测试
func TestFinger(t *testing.T) {
	f3 := finger.NewLoads("../../assets/fingerprint/apps.json")

	_, _, body, _ := director.GenerateNormalGet("http://www.qsbanks.cc:5050/back/ls")

	println(f3.CheckFinger(string(body)))

}
