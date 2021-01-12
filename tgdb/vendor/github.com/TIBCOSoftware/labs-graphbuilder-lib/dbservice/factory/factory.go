/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package factory

import (
	"sync"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice"
	tgdb2 "github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/tgdb/v2"
	tgdb3 "github.com/TIBCOSoftware/labs-graphbuilder-lib/dbservice/tgdb/v3"
)

var (
	instances map[dbservice.DBType]dbservice.DBServiceFactory
	once      sync.Once
)

func GetFactory(graphDB dbservice.DBType, version string) dbservice.DBServiceFactory {
	once.Do(func() {
		instances = make(map[dbservice.DBType]dbservice.DBServiceFactory)
		base := dbservice.BaseDBServiceFactory{DBServices: make(map[string]dbservice.DBService)}
		switch version {
		case "v2":
			instances[graphDB] = &tgdb2.TGDBServiceFactory{
				BaseDBServiceFactory: base,
			}
		case "v3":
			instances[graphDB] = &tgdb3.TGDBServiceFactory{
				BaseDBServiceFactory: base,
			}
		}
	})
	instances[graphDB].Initialize()
	return instances[graphDB]
}
