package pool

import (
	"errors"
	"sync"
	"time"
)

type Coon struct {
	Client   interface{}
	LeftTime time.Time
}

type Pool struct {
	MaxOpenConns    int
	ConnMaxLifeTime int
	ConnTimeOut     int
	mu              sync.RWMutex

	CreateClient  func() interface{}
	DestroyClient func(coon *Coon)
	Pools         chan *Coon
}

func (p *Pool) Create() {
	client := p.CreateClient()
	lefttime := time.Now().Add(time.Duration(p.ConnMaxLifeTime) * time.Second)

	coon := new(Coon)
	coon.Client = client
	coon.LeftTime = lefttime
	p.mu.RLock()
	p.Pools <- coon
	p.mu.RUnlock()
}

func (p *Pool) Init() error {
	if p.MaxOpenConns > 0 {
		p.Pools = make(chan *Coon, p.MaxOpenConns)
		return nil
	} else {
		return errors.New("parameter:MaxOpenConns must be greater than zero.")
	}
}
func (p *Pool) Get() (*Coon, error) {
	if len(p.Pools) < 1 {
		p.Create()
		return <-p.Pools, nil
	}
	for {
		select {
		case <-time.After(time.Duration(p.ConnTimeOut / 1000)):
			if len(p.Pools) < p.MaxOpenConns {
				p.Create()
			}
		case <-time.After(time.Duration(p.ConnTimeOut)):
			return nil, errors.New("Get connection time out")
		case client := <-p.Pools:
			if client.LeftTime.Unix() < time.Now().Unix() {
				continue
			} else {
				return client, nil
			}
		}
	}
}
func (p *Pool) Put(c *Coon) {
	if len(p.Pools) > p.MaxOpenConns {
		p.DestroyClient(c)
	} else {
		p.mu.RLock()
		p.Pools <- c
		p.mu.RUnlock()
	}
}
func (p *Pool)Close(){
	p.mu.Lock()
	pools := p.Pools
	p.Pools = nil
	p.mu.Unlock()
	
	if pools != nil {
		close(pools)
		for conn := range pools {
			p.DestroyClient(conn)
		}
	}
}