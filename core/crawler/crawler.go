package crawler

import "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Crawler는 newneek.co에서 기사 목록을 가져와 data에 추가하는 역할을 합니다.
// 저장 형식은 (월)(일).htm으로 저장이 됩니다.
type Crawler interface {
	Get() ([]byte, error)
}
