// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package redis for session provider
//
// depend on github.com/garyburd/redigo/redis
//
// go install github.com/garyburd/redigo/redis
//
// Usage:
// import(
//   _ "github.com/astaxie/beego/session/redis"
//   "github.com/astaxie/beego/session"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("redis", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:7070"}``)
//		go globalSessions.GC()
//	}
//
// more docs: http://beego.me/docs/module/session.md
package redis

import (
	redisclient "GoH/core/redis"
	"GoH/core/session/sess"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"sync"
	"strings"
	"strconv"
)

var redispder = &Provider{}

// SessionStore redis session store
type SessionStore struct {
	pool        *redis.Pool
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
}

// Set value in redis session
func (rs *SessionStore) Set(key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

// Get value in redis session
func (rs *SessionStore) Get(key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in redis session
func (rs *SessionStore) Delete(key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

// Flush clear all values in redis session
func (rs *SessionStore) Flush() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get redis session id
func (rs *SessionStore) SessionID() string {
	return rs.sid
}

// SessionRelease save session values to redis
func (rs *SessionStore) SessionRelease(w http.ResponseWriter) {
	b, err := session.EncodeGob(rs.values)
	if err != nil {
		return
	}
	c := rs.pool.Get()
	defer c.Close()
	c.Do("SETEX", rs.sid, rs.maxlifetime, string(b))
}

// Provider redis session provider
type Provider struct {
	maxlifetime int64
	savePath    string
	dbNum       int
	poollist    *redis.Pool
}

// SessionInit init redis session
func (rp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	rp.maxlifetime = maxlifetime
	configs := strings.Split(savePath, ",")
	if len(configs) > 0 {
		dbNum, err := strconv.Atoi(configs[0])
		if err != nil || dbNum < 0 {
			rp.dbNum = 0
		} else {
			rp.dbNum = dbNum
		}
	}
	rp.poollist = redisclient.RedisClient
	rp.poollist.Get().Do("SELECT", rp.dbNum)
	return rp.poollist.Get().Err()
}

// SessionRead read redis session by sid
func (rp *Provider) SessionRead(sid string) (session.Store, error) {
	c := rp.poollist.Get()
	defer c.Close()
	var kv map[interface{}]interface{}
	kvs, err := redis.String(c.Do("GET", sid))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if len(kvs) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		if kv, err = session.DecodeGob([]byte(kvs)); err != nil {
			return nil, err
		}
	}
	rs := &SessionStore{pool: rp.poollist, sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

// SessionExist check redis session exist by sid
func (rp *Provider) SessionExist(sid string) bool {
	c := rp.poollist.Get()
	defer c.Close()
	if existed, err := redis.Int(c.Do("EXISTS", sid)); err != nil || existed == 0 {
		return false
	}
	return true
}

// SessionRegenerate generate new sid for redis session
func (rp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	c := rp.poollist.Get()
	defer c.Close()
	if existed, _ := redis.Int(c.Do("EXISTS", oldsid)); existed == 0 {
		// oldsid doesn't exists, set the new sid directly
		// ignore error here, since if it return error
		// the existed value will be 0
		c.Do("SET", sid, "", "EX", rp.maxlifetime)
	} else {
		c.Do("RENAME", oldsid, sid)
		c.Do("EXPIRE", sid, rp.maxlifetime)
	}
	return rp.SessionRead(sid)
}

// SessionDestroy delete redis session by id
func (rp *Provider) SessionDestroy(sid string) error {
	c := rp.poollist.Get()
	defer c.Close()
	c.Do("DEL", sid)
	return nil
}

// SessionGC Impelment method, no used.
func (rp *Provider) SessionGC() {
}

// SessionAll return all activeSession
func (rp *Provider) SessionAll() int {
	return 0
}

func init() {
	session.Register("redis", redispder)
}
