package egoja

import (
	"golang.org/x/sync/singleflight"
	"sync"
)

const (
	mainFuncScript = ` (function() {
 				%s
        })()`
)

var (
	cacheProgram = sync.Map{}
	single       singleflight.Group
)
