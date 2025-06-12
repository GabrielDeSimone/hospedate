package repositories

import sqlTmp "github.com/hospedate/backyard/sql"

const INITIAL_VERSION = dbversion("v0")

var migrations = []migration{
    {
        VersionFrom:   INITIAL_VERSION,
        VersionTo:     dbversion("v1"),
        Queries:       sqlTmp.QueriesV0ToV1,
        QueriesInvert: sqlTmp.QueriesInvertV0ToV1,
    },
    {
        VersionFrom:   dbversion("v1"),
        VersionTo:     dbversion("v2"),
        Queries:       sqlTmp.QueriesV1ToV2,
        QueriesInvert: sqlTmp.QueriesInvertV1ToV2,
    },
    {
        VersionFrom:   dbversion("v2"),
        VersionTo:     dbversion("v3"),
        Queries:       sqlTmp.QueriesV2ToV3,
        QueriesInvert: sqlTmp.QueriesInvertV2ToV3,
    },
    {
        VersionFrom:   dbversion("v3"),
        VersionTo:     dbversion("v4"),
        Queries:       sqlTmp.QueriesV3ToV4,
        QueriesInvert: sqlTmp.QueriesInvertV3ToV4,
    },
}

func isValidMigrationList(migrations []migration) (dbversion, *MigrationMap, error) {
    migrationMap, err := createMigrationMap(migrations)
    if err != nil {
        return "", nil, err
    }

    finalVersion, count := countMovementsInChain(migrationMap)

    if count != len(migrationMap) {
        return "", nil, ErrWrongMigrationChain
    }

    return finalVersion, &migrationMap, nil
}

func createMigrationMap(migrations []migration) (MigrationMap, error) {
    migrationMap := make(map[dbversion]*migration, len(migrations))

    for _, m := range migrations {
        if _, ok := migrationMap[m.VersionFrom]; ok {
            return nil, ErrWrongMigrationChain
        }

        for _, value := range migrationMap {
            if value.VersionTo == m.VersionTo {
                return nil, ErrWrongMigrationChain
            }
        }
        mCopy := m
        migrationMap[m.VersionFrom] = &mCopy
    }

    return migrationMap, nil
}

func countMovementsInChain(migrationMap MigrationMap) (finalVersion dbversion, count int) {
    currVersion := INITIAL_VERSION
    count = 0
    for {
        if m, ok := migrationMap[currVersion]; ok {
            currVersion = m.VersionTo
            count++
        } else {
            break
        }
    }

    return currVersion, count
}
