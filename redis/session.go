package redis

import (
	"fmt"
	"sync"

	"github.com/gomodule/redigo/redis"
)

//default redis write and read buf is 4k

//Session session
type Session struct {
	mutex sync.Mutex
	pool  *redis.Pool
	conn  redis.Conn
	nums  int
}

//NewSession return redis session
func (c *Client) NewSession() *Session {
	return &Session{
		pool: c.pool,
		nums: -1,
	}
}

//Begin init nums and pool
func (s *Session) Begin() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.nums != -1 {
		return fmt.Errorf("must use NewSession to init! ")
	}

	s.nums = 0
	s.conn = s.pool.Get()

	return

}

//Exec nums++ and fill in redis.conn.write buf
func (s *Session) Exec(commandName string, args ...interface{}) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.nums == -1 {
		return fmt.Errorf("must use Begin first! ")
	}

	s.nums++

	return s.conn.Send(commandName, args[:]...)
}

//Commit commit write buf to redis server
func (s *Session) Commit() error {
	return s.conn.Flush()
}

//Receive receive data from redis.conn.read buf
func (s *Session) Receive() (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.nums <= 0 {
		return nil, fmt.Errorf("no data can be received")
	}

	s.nums--

	return s.conn.Receive()

}

//Close clear conn read buf and close net.conn
func (s *Session) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := 0; i < s.nums; i++ {
		s.conn.Receive()
	}

	return s.conn.Close()
}
