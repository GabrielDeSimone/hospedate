package sql

const Pgcrypto_init = `
CREATE EXTENSION IF NOT EXISTS pgcrypto;
`

const Users_drop = `
DROP TABLE IF EXISTS users;
`

const Users_ddl = `
CREATE TABLE IF NOT EXISTS users (
    id varchar(255) NOT NULL PRIMARY KEY,
    email varchar(255) NOT NULL UNIQUE,
    password_enc varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    phone_number varchar (25) NOT NULL UNIQUE,
    is_host bool DEFAULT false NOT NULL
  );
`

const Users_invitations_drop = `
DROP TABLE IF EXISTS users_invitations;
`

const Invitation_kind_categories = `
DO $$ BEGIN CREATE TYPE invitation_kind_domain AS ENUM ('for_traveler', 'for_owner'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`
const Invitation_kind_categories_drop = `
DROP TYPE IF EXISTS invitation_kind_domain;
`

const Users_invitations_ddl = `
CREATE TABLE IF NOT EXISTS users_invitations (
    id varchar(255) NOT NULL PRIMARY KEY,
    used_by varchar(255),
    kind invitation_kind_domain DEFAULT 'for_traveler' NOT NULL,
    generated_by varchar(255) NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    CONSTRAINT fk_user_id
      FOREIGN KEY(used_by)
      REFERENCES users(id)
      ON DELETE CASCADE,
    CONSTRAINT fk_generated_by
      FOREIGN KEY(generated_by)
      REFERENCES users(id)
      ON DELETE CASCADE
)
`

const Property_status_categories = `
DO $$ BEGIN CREATE TYPE property_status_domain AS ENUM ('loading', 'active','archived'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`
const Property_status_categories_drop = `
DROP TYPE IF EXISTS property_status_domain;
`

const Accommodation_enum = `
DO $$ BEGIN CREATE TYPE accommodation_domain AS ENUM ('house', 'apartment', 'private_room', 'shared_room'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`

const Location_enum = `
DO $$ BEGIN CREATE TYPE location_domain AS ENUM ('city_center', 'near_beach', 'residential_area', 'countryside', 'mountain'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`

const Wifi_enum = `
DO $$ BEGIN CREATE TYPE wifi_domain AS ENUM ('shared', 'private', 'not_available'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`

const TV_enum = `
DO $$ BEGIN CREATE TYPE tv_domain AS ENUM ('available', 'available_cable_or_streaming', 'not_available'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`

const Availability_enum = `
DO $$ BEGIN CREATE TYPE availability_domain AS ENUM ('available', 'not_available'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`

const Parking_enum = `
DO $$ BEGIN CREATE TYPE parking_domain AS ENUM ('available_in_public_area', 'available_private_uncovered', 'available_private_covered', 'not_available'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`

const Booking_enum = `
DO $$ BEGIN CREATE TYPE booking_domain AS ENUM ('owner_directly', 'in_platform', 'both'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`

const Properties_ddl = `
CREATE TABLE IF NOT EXISTS properties (
    id varchar(255) NOT NULL PRIMARY KEY,
    title varchar(255) DEFAULT '' NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    max_guests int NOT NULL,
    airbnb_room_id varchar(255) UNIQUE,
    ex_airbnb_room_id varchar(255),
    price int NOT NULL,
    city varchar(255) NOT NULL,
    status property_status_domain DEFAULT 'loading' NOT NULL,
    user_id varchar(255) NOT NULL,
    is_verified bool DEFAULT false NOT NULL,
    images JSON default '[]' NOT NULL,
    booking_options booking_domain DEFAULT 'both' NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    accommodation accommodation_domain,
    location location_domain,
    wifi wifi_domain,
    tv tv_domain,
    microwave availability_domain,
    oven availability_domain,
    kettle availability_domain,
    toaster availability_domain,
    coffee_machine availability_domain,
    air_conditioning availability_domain,
    heating availability_domain,
    parking parking_domain,
    pool availability_domain,
    gym availability_domain,
    half_bathrooms int,
    bedrooms int,
    CONSTRAINT fk_owner
      FOREIGN KEY(user_id) 
      REFERENCES users(id)
      ON DELETE CASCADE
  );
`
const Accommodation_enum_drop = `
DROP TYPE IF EXISTS accommodation_domain;
`

const Location_enum_drop = `
DROP TYPE IF EXISTS location_domain;
`

const Wifi_enum_drop = `
DROP TYPE IF EXISTS wifi_domain;
`

const TV_enum_drop = `
DROP TYPE IF EXISTS tv_domain;
`

const Availability_enum_drop = `
DROP TYPE IF EXISTS availability_domain;
`

const Parking_enum_drop = `
DROP TYPE IF EXISTS parking_domain;
`

const Booking_enum_drop = `
DROP TYPE IF EXISTS booking_domain;
`

const Properties_drop = `
DROP TABLE IF EXISTS properties;
`

const Blocks_drop = `
DROP TABLE IF EXISTS blocks;
`

const Blocks_ddl = `
CREATE TABLE IF NOT EXISTS blocks (
    id varchar(255) NOT NULL PRIMARY KEY,
    property_id varchar(255) NOT NULL,
    block_init DATE NOT NULL,
    block_end DATE NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    CONSTRAINT fk_property
      FOREIGN KEY(property_id) 
      REFERENCES properties(id)
      ON DELETE CASCADE
  );
`

const Orders_drop = `
DROP TABLE IF EXISTS orders;
`
const Order_status_categories = `
DO $$ BEGIN CREATE TYPE order_status_domain AS ENUM ('pending', 'confirmed','completed','in_progress','canceled','ephemeral','discarded'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`
const Order_status_categories_drop = `
DROP TYPE IF EXISTS order_status_domain;
`

const Order_type_categories = `
DO $$ BEGIN CREATE TYPE order_type_domain AS ENUM ('in_platform', 'owner_directly'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`
const Order_type_categories_drop = `
DROP TYPE IF EXISTS order_type_domain;
`

const Order_canceledby_categories = `
DO $$ BEGIN CREATE TYPE order_canceledby_domain AS ENUM ('visitor', 'owner'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`
const Order_canceledby_categories_drop = `
DROP TYPE IF EXISTS order_canceledby_domain;
`

const Orders_ddl = `
CREATE TABLE IF NOT EXISTS orders (
    id varchar(255) NOT NULL PRIMARY KEY,
    user_id varchar(255) NOT NULL,
    status order_status_domain DEFAULT 'pending' NOT NULL,
    property_id varchar(255) NOT NULL,
    checkin_date DATE NOT NULL,
    checkout_date DATE NOT NULL,
    number_guests int NOT NULL,
    price int NOT NULL,
    price_currency varchar(255) DEFAULT 'USDT' NOT NULL,
    total_billed_cents int NOT NULL,
    canceled_by order_canceledby_domain,
    order_type order_type_domain NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    wallet_address varchar(255) DEFAULT '' NOT NULL
  );
`

const Payment_status_categories = `
DO $$ BEGIN CREATE TYPE payment_status_domain AS ENUM ('pending', 'confirmed','reverted'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`
const Payment_status_categories_drop = `
DROP TYPE IF EXISTS payment_status_domain;
`
const Payment_method_categories = `
DO $$ BEGIN CREATE TYPE payment_method_domain AS ENUM ('crypto_wallet'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`
const Payment_method_categories_drop = `
DROP TYPE IF EXISTS payment_method_domain;
`

const Payments_ddl = `
CREATE TABLE IF NOT EXISTS payments (
    id varchar(255) NOT NULL PRIMARY KEY,
    order_id varchar(255) NOT NULL UNIQUE,
    method payment_method_domain DEFAULT 'crypto_wallet' NOT NULL,
    status payment_status_domain DEFAULT 'pending' NOT NULL,
    traveler_amount_cents int NOT NULL,
    traveler_currency varchar(255) DEFAULT 'USDT' NOT NULL,
    received_amount_cents int,
    received_currency varchar(255),
    reverted_amount_cents int,
    reverted_currency varchar(255),
    created_at timestamp without time zone default (now() at time zone 'utc'),
    private_key varchar(255) DEFAULT '' NOT NULL
  );
`
const Payments_drop = `
DROP TABLE IF EXISTS payments;
`
const Withdrawals_status_categories = `
DO $$ BEGIN CREATE TYPE withdrawal_status_domain AS ENUM ('pending','completed','aborted'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`
const Withdrawals_status_categories_drop = `
DROP TYPE IF EXISTS withdrawal_status_domain;
`

const Withdrawals_method_categories = `
DO $$ BEGIN CREATE TYPE withdrawal_method_domain AS ENUM ('crypto_wallet'); EXCEPTION WHEN duplicate_object THEN null; END $$;
`
const Withdrawals_method_categories_drop = `
DROP TYPE IF EXISTS withdrawal_method_domain;
`

const Withdrawals_ddl = `
CREATE TABLE IF NOT EXISTS withdrawals (
    id varchar(255) NOT NULL PRIMARY KEY,
    user_id varchar(255) NOT NULL,
    method withdrawal_method_domain DEFAULT 'crypto_wallet' NOT NULL,
    status withdrawal_status_domain DEFAULT 'pending' NOT NULL,
    reclaimed_amount_cents int NOT NULL,
    reclaimed_currency varchar(255) DEFAULT 'USDT' NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    processed_at timestamp without time zone,
    wallet_address varchar(255) NOT NULL
  );
`
const Withdrawals_drop = `
DROP TABLE IF EXISTS withdrawals;
`

const Owners_earned_ddl = `
CREATE TABLE IF NOT EXISTS owners_earned (
    id varchar(255) NOT NULL PRIMARY KEY,
    order_id varchar(255) NOT NULL UNIQUE,
    user_id varchar(255) NOT NULL, 
    earned_amount_cents int NOT NULL,
    earned_currency varchar(255) DEFAULT 'USDT' NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc')
  );
`
const Owners_earned_drop = `
DROP TABLE IF EXISTS owners_earned;
`
const Db_versions_ddl = `
CREATE TABLE db_versions (
    version varchar(255) NOT NULL,
    applied_at timestamp without time zone default (now() at time zone 'utc') NOT NULL
);
`
const Db_versions_drop = `
DROP TABLE IF EXISTS db_versions;
`
const Users_credit_ddl = `
CREATE TABLE IF NOT EXISTS users_credit (
    id varchar(255) NOT NULL PRIMARY KEY,
    user_id varchar(255) NOT NULL, 
    invitation_id varchar(255) NOT NULL,
    earned_amount DECIMAL(10,2) NOT NULL,
    earned_currency varchar(255) DEFAULT 'USDT' NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc')
  );
`
const Users_credit_drop = `
DROP TABLE IF EXISTS users_credit;
`
