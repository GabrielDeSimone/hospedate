package repositories

import (
    "database/sql"
    "time"

    "github.com/hospedate/backyard/log"
    sqlTmp "github.com/hospedate/backyard/sql"
)

type GlobalRepository interface {
    Migration()
    HardInitDB()
}

type GlobalRepositoryImpl struct {
    db     *sql.DB
    logger log.Logger
}
type dbversion string

type migration struct {
    VersionFrom   dbversion
    VersionTo     dbversion
    Queries       []string
    QueriesInvert []string
}

type MigrationMap map[dbversion]*migration

func NewGlobalRepository(db *sql.DB) GlobalRepository {
    logger := log.NewLogger("GlobalRepository", string(log.INFO_LEVEL))
    return &GlobalRepositoryImpl{db: db, logger: logger}
}

func (r *GlobalRepositoryImpl) getCurrentDBVersion() dbversion {
    result := r.db.QueryRow(sqlTmp.Get_db_current_version)
    var db_version string
    err := result.Scan(&db_version)
    if err != nil {
        r.logger.Fatal("Error getting DB version: ", err.Error())
    }
    return dbversion(db_version)
}

func (r *GlobalRepositoryImpl) executeMigration(migrationAction *migration, reverse bool, tx *sql.Tx) error {

    queries := migrationAction.Queries
    if reverse {
        queries = migrationAction.QueriesInvert
    }

    if err := r.executeQueries(tx, queries); err != nil {
        return err
    }

    if err := r.updateDBVersion(tx, migrationAction.VersionFrom, migrationAction.VersionTo, reverse); err != nil {
        return err
    }

    return nil
}

func (r *GlobalRepositoryImpl) executeQueries(tx *sql.Tx, queries []string) error {
    for _, query := range queries {
        _, err := tx.Exec(query)
        if err != nil {
            r.logger.Error("Cannot execute operation during migration", err.Error())
            return err
        }
    }
    return nil
}

func (r *GlobalRepositoryImpl) updateDBVersion(tx *sql.Tx, fromVersion, toVersion dbversion, reverse bool) error {
    if reverse {
        fromVersion, toVersion = toVersion, fromVersion
    }

    _, err := tx.Exec(sqlTmp.Update_db_version, toVersion)
    if err != nil {
        r.logger.Error("Cannot update db version", err.Error())
        return err
    }
    r.logger.Infof("DB migrated from version %v to %v", fromVersion, toVersion)
    return nil
}

func (r *GlobalRepositoryImpl) ApplyAllMigrations(migrations MigrationMap, currentVersion dbversion) error {

    tx, err := r.db.Begin()
    if err != nil {
        return err
    }

    for migrations[currentVersion] != nil {
        migration := migrations[currentVersion]
        err := r.executeMigration(migration, false, tx)
        if err != nil {
            r.logger.Error("Cannot execute operation during migration, doing rollback", err.Error())
            tx.Rollback()
            return err
        }

        currentVersion = migration.VersionTo
        // set migrations apart in time
        time.Sleep(100 * time.Millisecond)
    }

    err = tx.Commit()
    if err != nil {
        r.logger.Error("Cannot commit transaction to complete migration, doing rollback", err.Error())
        tx.Rollback()
        return err
    }

    return nil
}

func (r *GlobalRepositoryImpl) Migration() {
    finalVersion, migrationMap, err := isValidMigrationList(migrations)
    if err != nil {
        r.logger.Fatal("DB migration list not valid")
    }

    currentVersion := r.getCurrentDBVersion()
    r.logger.Infof("DB migration process is about to start. Final version defined: %v, Current DB version: %v", finalVersion, currentVersion)
    err = r.ApplyAllMigrations(*migrationMap, currentVersion)
    newDBVersion := r.getCurrentDBVersion()
    if err == nil {
        r.logger.Infof("DB migration process finished. Original version: %v. Current version: %v", currentVersion, newDBVersion)
    } else {
        r.logger.Fatal("Error migrating DB: ", err.Error())
    }
}

func (r *GlobalRepositoryImpl) HardInitDB() {
    // In a hard init db, we drop tables and execute the ddls
    r.logger.Info("Hard initializing DB")
    finalVersion, _, err := isValidMigrationList(migrations)
    if err != nil {
        r.logger.Fatal("DB migration list not valid")
    }
    queries := []string{
        sqlTmp.Db_versions_drop,
        sqlTmp.Users_credit_drop,
        sqlTmp.Owners_earned_drop,
        sqlTmp.Withdrawals_drop,
        sqlTmp.Withdrawals_method_categories_drop,
        sqlTmp.Withdrawals_status_categories_drop,
        sqlTmp.Payments_drop,
        sqlTmp.Payment_status_categories_drop,
        sqlTmp.Payment_method_categories_drop,
        sqlTmp.Orders_drop,
        sqlTmp.Order_canceledby_categories_drop,
        sqlTmp.Order_status_categories_drop,
        sqlTmp.Order_type_categories_drop,
        sqlTmp.Blocks_drop,
        sqlTmp.Properties_drop,
        sqlTmp.Property_status_categories_drop,
        sqlTmp.Booking_enum_drop, // New Properties ENUM drops
        sqlTmp.Parking_enum_drop,
        sqlTmp.Availability_enum_drop,
        sqlTmp.TV_enum_drop,
        sqlTmp.Wifi_enum_drop,
        sqlTmp.Location_enum_drop,
        sqlTmp.Accommodation_enum_drop,
        sqlTmp.Users_invitations_drop,
        sqlTmp.Invitation_kind_categories_drop,
        sqlTmp.Users_drop,
        sqlTmp.Pgcrypto_init,
        sqlTmp.Users_ddl,
        sqlTmp.Invitation_kind_categories,
        sqlTmp.Users_invitations_ddl,
        sqlTmp.Property_status_categories,
        sqlTmp.Accommodation_enum, // New Properties ENUMs initialization
        sqlTmp.Location_enum,
        sqlTmp.Wifi_enum,
        sqlTmp.TV_enum,
        sqlTmp.Availability_enum,
        sqlTmp.Parking_enum,
        sqlTmp.Booking_enum,
        sqlTmp.Properties_ddl,
        sqlTmp.Blocks_ddl,
        sqlTmp.Order_status_categories,
        sqlTmp.Order_canceledby_categories,
        sqlTmp.Order_type_categories,
        sqlTmp.Orders_ddl,
        sqlTmp.Payment_status_categories,
        sqlTmp.Payment_method_categories,
        sqlTmp.Payments_ddl,
        sqlTmp.Withdrawals_method_categories,
        sqlTmp.Withdrawals_status_categories,
        sqlTmp.Withdrawals_ddl,
        sqlTmp.Owners_earned_ddl,
        sqlTmp.Users_credit_ddl,
        sqlTmp.Db_versions_ddl,
    }
    r.execSequentially(queries)

    _, err = r.db.Exec(sqlTmp.Update_db_version, finalVersion)
    if err != nil {
        r.logger.Fatal("Error initializing DB: ", err.Error())
    }

    r.logger.Info("DB hard-initialized")
}

func (r *GlobalRepositoryImpl) execSequentially(queries []string) {
    for _, query := range queries {
        r.logger.Debug("Executing query: ", query)
        _, err := r.db.Exec(query)
        if err != nil {
            r.logger.Fatal("Error initializing DB: ", err.Error())
        }
    }
}
