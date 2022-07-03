package cmn

import (
	"github.com/importcjj/sensitive"
	"sync"
)

var (
	_SensitiveFilter *sensitive.Filter
	_Once            sync.Once
)

func SensitiveFilter() *sensitive.Filter {
	return _SensitiveFilter
}

func InitSensitiveWords(url string) {
	if _SensitiveFilter != nil {
		return
	}
	_SensitiveFilter = sensitive.New()
	_SensitiveFilter.LoadNetWordDict(url)
}
