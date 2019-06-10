package supervisor

import (
	"github.com/Workiva/go-datastructures/queue"
	. "github.com/iotaledger/iota.go/trinary"
	"sync"
	"time"
)

type Supervisor struct {
	Name            string
	queue           *queue.Queue
	idle            bool
	environments    map[string]*environment
	accessLock      *sema // lock for quant. Locks changes in environments and entities: join, affect
	timeout         time.Duration
	quantWG         sync.WaitGroup // released when quant ends
	quantCount      int64
	quantCountMutex sync.RWMutex
}

type Entity struct {
	Supervisor *Supervisor
	Name       string
	inSize     int
	outSize    int
	affecting  []*affectEntData // list of affected environments where effects are sent
	joined     []*environment   // list of environments which are being listened to
	inChan     chan entityMsg   // chan for incoming effects
	core       EntityCore       // function called for each effect
}

type affectEntData struct {
	environment *environment
	delay       int
}

type entityMsg struct {
	effect          Trits
	lastWithinLimit bool
}

type EntityCore interface {
	Call(Trits, Trits) bool
}

type joinEnvData struct {
	entity *Entity
	limit  int
	count  int
}

type environment struct {
	supervisor *Supervisor
	name       string
	invalid    bool
	joins      []*joinEnvData
	affects    []*Entity
	effectChan chan Trits
}
