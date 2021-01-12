/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package cache

import (
	"errors"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("dgraph-service")

type CacheEntry struct {
	eId  int
	id   string
	next *CacheEntry
	data interface{}
}

func (this *CacheEntry) Clear() {
	this.id = ""
	this.data = nil
}

func (this *CacheEntry) SetNext(next *CacheEntry) *CacheEntry {
	this.next = next
	return next
}

func NewCacheEntry(eId int) *CacheEntry {
	cacheEntry := &CacheEntry{eId: eId}
	return cacheEntry
}

type Cache struct {
	disable        bool
	size           int
	current        *CacheEntry
	repository     map[string]*CacheEntry
	repositoryRing map[int]*CacheEntry
	mux            sync.Mutex
}

func (this *Cache) Add(id string, data interface{}) error {
	this.mux.Lock()
	defer this.mux.Unlock()

	if this.disable {
		return nil
	}

	/* Before reuse cache entry, remove its association with old key in the map */
	log.Debug("Before reuse, repository =  ", this.repository, ", id =  ", this.current.id, ", entry = ", this.repository[this.current.id])
	if nil != this.current.data {
		log.Debug("Before delete old key association, id =  ", this.current.id, ", entry = ", this.repository[this.current.id])
		delete(this.repository, this.current.id)
	}

	log.Debug("After delete old key association, repository =  ", this.repository)

	this.current.id = id
	this.current.data = data
	this.repository[id] = this.current

	this.current = this.current.next
	log.Debug("After reuse, repository =  ", this.repository)

	log.Debug("len(this.repository) = ", len(this.repository), ", size = ", this.size)
	if len(this.repository) > this.size {
		return errors.New("Cache repository over size!!")
	}

	return nil
}

func (this *Cache) Get(id string) interface{} {
	if this.disable {
		return nil
	}
	log.Debug("Target = ", id, ", Current eid = ", this.current.eId, ", repositryRing = ", this.repositoryRing, ", repositry = ", this.repository)

	target := this.repository[id]
	var data interface{}
	if nil != target {
		log.Debug("Cache HIT !")
		data = target.data
		target.Clear()
		this.Add(id, data)
	} else {
		log.Debug("Cache MISS !")
	}
	return data
}

func NewCache(size int) *Cache {
	cache := &Cache{
		disable:        size < 1,
		size:           size,
		repository:     make(map[string]*CacheEntry),
		repositoryRing: make(map[int]*CacheEntry),
	}

	if !cache.disable {
		entry := NewCacheEntry(0)
		cache.repositoryRing[0] = entry
		cache.current = entry

		for i := 0; i < size-1; i++ {
			entry = entry.SetNext(NewCacheEntry(i + 1))
			cache.repositoryRing[i+1] = entry
		}
		entry.SetNext(cache.current)
	}
	return cache
}
