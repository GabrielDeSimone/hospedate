package sql

import "fmt"

const Get_db_current_version = `
    SELECT version
    FROM db_versions
    ORDER BY applied_at DESC
    LIMIT 1;
`
const Update_db_version = `
    INSERT INTO db_versions (version) VALUES ($1);
`

var QueriesV0ToV1 = []string{
    "ALTER TYPE property_status_domain ADD VALUE 'archived';",
    "ALTER TABLE properties ALTER COLUMN airbnb_room_id DROP NOT NULL;",
    "ALTER TABLE properties ADD ex_airbnb_room_id varchar(255);",
}

var QueriesInvertV0ToV1 = []string{
    "CREATE TYPE property_status_domain_old AS ENUM ('loading', 'active');",
    "ALTER TABLE properties ALTER status DROP DEFAULT;",
    "ALTER TABLE properties ALTER COLUMN status TYPE property_status_domain_old USING status::text::property_status_domain_old;",
    "ALTER TABLE properties ALTER status SET DEFAULT 'loading';",
    "DROP TYPE property_status_domain;",
    "ALTER TYPE property_status_domain_old RENAME TO property_status_domain;",
    "ALTER TABLE properties ALTER COLUMN airbnb_room_id SET NOT NULL;",
    "ALTER TABLE properties DROP COLUMN ex_airbnb_room_id;",
}

var QueriesV1ToV2 = []string{
    "CREATE TYPE invitation_kind_domain AS ENUM ('for_traveler', 'for_owner');",
    "ALTER TABLE users_invitations ADD kind invitation_kind_domain DEFAULT 'for_traveler' NOT NULL;",
    Users_credit_ddl,
    "ALTER TABLE properties ADD is_verified bool DEFAULT false NOT NULL",
}

var QueriesInvertV1ToV2 = []string{
    "ALTER TABLE users_invitations DROP COLUMN kind;",
    "DROP TYPE invitation_kind_domain;",
    "DROP TABLE IF EXISTS users_credit;",
    "ALTER TABLE properties DROP COLUMN is_verified;",
}
var QueriesV2ToV3 = []string{
    "UPDATE orders SET total_billed = (total_billed * 100)::integer;",
    "ALTER TABLE orders RENAME total_billed TO total_billed_cents;",
    "ALTER TABLE payments RENAME traveller_amount TO traveler_amount_cents;",
    "ALTER TABLE payments RENAME traveller_currency TO traveler_currency;",
    "ALTER TABLE payments RENAME received_amount TO received_amount_cents;",
    "ALTER TABLE payments RENAME reverted_amount TO reverted_amount_cents;",
    "UPDATE payments SET traveler_amount_cents = (traveler_amount_cents * 100)::integer;",
    "UPDATE payments SET received_amount_cents = (received_amount_cents * 100)::integer;",
    "UPDATE payments SET reverted_amount_cents = (reverted_amount_cents * 100)::integer;",
    "UPDATE owners_earned SET earned_amount = ROUND(earned_amount * 100)::integer;",
    "ALTER TABLE owners_earned RENAME earned_amount TO earned_amount_cents;",
    "ALTER TABLE owners_earned ALTER COLUMN earned_amount_cents TYPE int;",
}

var QueriesInvertV2ToV3 = []string{
    "ALTER TABLE orders RENAME total_billed_cents TO total_billed;",
    "ALTER TABLE payments RENAME traveler_amount_cents TO traveller_amount;",
    "ALTER TABLE payments RENAME traveler_currency TO traveller_currency;",
    "ALTER TABLE payments RENAME received_amount_cents TO received_amount;",
    "ALTER TABLE payments RENAME reverted_amount_cents TO reverted_amount;",
    "ALTER TABLE owners_earned RENAME earned_amount_cents TO earned_amount;",
    "ALTER TABLE owners_earned ALTER COLUMN earned_amount TYPE DECIMAL(10,2);",
    checkAndConvertToDollars("orders", "total_billed"),
    "SELECT update_orders_total_billed_to_dollars();",
    checkAndConvertToDollars("payments", "traveller_amount"),
    "SELECT update_payments_traveller_amount_to_dollars();",
    checkAndConvertToDollars("payments", "received_amount"),
    "SELECT update_payments_received_amount_to_dollars();",
    checkAndConvertToDollars("payments", "reverted_amount"),
    "SELECT update_payments_reverted_amount_to_dollars();",
    "UPDATE owners_earned SET earned_amount = earned_amount / 100;",
    "DROP FUNCTION update_orders_total_billed_to_dollars();",
    "DROP FUNCTION update_payments_traveller_amount_to_dollars();",
    "DROP FUNCTION update_payments_received_amount_to_dollars();",
    "DROP FUNCTION update_payments_reverted_amount_to_dollars();",
}

var QueriesV3ToV4 = []string{
    "UPDATE withdrawals SET reclaimed_amount = ROUND(reclaimed_amount * 100)::integer;",
    "ALTER TABLE withdrawals RENAME reclaimed_amount TO reclaimed_amount_cents;",
    "ALTER TABLE withdrawals ALTER COLUMN reclaimed_amount_cents TYPE int;",
}

var QueriesInvertV3ToV4 = []string{
    "ALTER TABLE withdrawals RENAME reclaimed_amount_cents TO reclaimed_amount;",
    "ALTER TABLE withdrawals ALTER COLUMN reclaimed_amount TYPE DECIMAL(10,2);",
    "UPDATE withdrawals SET reclaimed_amount = reclaimed_amount / 100;",
}

func checkAndConvertToDollars(tableName string, columnName string) string {
    queryTemplate := `
CREATE OR REPLACE FUNCTION update_%[1]s_%[2]s_to_dollars()
RETURNS void AS $$
DECLARE
    row_record record;
BEGIN
    FOR row_record IN (SELECT * FROM %[1]s) LOOP
        IF (row_record.%[2]s/100.0) %% 1 <> 0 THEN
            RAISE EXCEPTION 'Non-integer value in column %[2]s';
        END IF;
    END LOOP;
    
    -- Si todas las filas pasan la verificación, entonces realizamos la actualización.
    EXECUTE 'UPDATE %[1]s SET %[2]s = %[2]s/100';
END;
$$ LANGUAGE plpgsql;
`

    return fmt.Sprintf(queryTemplate, tableName, columnName)
}
