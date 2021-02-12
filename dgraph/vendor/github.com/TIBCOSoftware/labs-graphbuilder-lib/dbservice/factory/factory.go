/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package factory

import (
	"sync"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
	dgraph "github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/dgraph/factory"
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
		case dbservice.Dgraph:
			instances[graphDB] = &dgraph.DgraphServiceFactory{
				BaseDBServiceFactory: base,
			}
		}
	})
	instances[graphDB].Initialize()
	return instances[graphDB]
}
