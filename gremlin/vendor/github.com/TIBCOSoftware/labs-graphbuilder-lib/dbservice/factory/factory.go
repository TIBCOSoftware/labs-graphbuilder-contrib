/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package factory

import (
	"sync"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
	gremlin "github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/gremlin"
)

var (
	instances map[dbservice.DBType]dbservice.DBServiceFactory
	once      sync.Once
)

func GetFactory(graphDB dbservice.DBType) dbservice.DBServiceFactory {
	once.Do(func() {
		instances = make(map[dbservice.DBType]dbservice.DBServiceFactory)
		base := dbservice.BaseDBServiceFactory{DBServices: make(map[string]dbservice.DBService)}
		switch graphDB {
		case dbservice.Gremlin:
			instances[graphDB] = &gremlin.GremlinServiceFactory{
				BaseDBServiceFactory: base,
			}
		}

	})
	instances[graphDB].Initialize()
	return instances[graphDB]
}
