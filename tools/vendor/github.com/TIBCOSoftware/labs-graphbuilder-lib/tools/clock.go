/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tools

import (
	"sync"
	"time"
)

//-====================-//
//    Define Clock
//-====================-//

var (
	clock *Clock
	cOnce sync.Once
)

type Clock struct {
	now int64
	mux sync.Mutex
}

func GetClock() *Clock {
	cOnce.Do(func() {
		clock = &Clock{now: 0}
	})
	return clock
}

func (this *Clock) SetCurrentTime(now int64) {
	this.mux.Lock()
	defer this.mux.Unlock()
	//fmt.Println("Try Set Current Time To : ", now)
	if 0 == this.now || now > this.now {
		this.now = now
	}
	//fmt.Println("After Set Current Time : ", this.now)
}

func (this *Clock) GetCurrentTime() int64 {
	//fmt.Println("Get Current Time : ", this.now)
	if 0 == this.now {
		return time.Now().Unix()
	} else {
		return this.now
	}
}
