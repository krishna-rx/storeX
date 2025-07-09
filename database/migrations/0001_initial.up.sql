BEGIN;
DROP TYPE IF EXISTS role_types, user_types, asset_types, asset_statuses CASCADE;
CREATE TYPE role_types AS ENUM ('admin', 'asset_manager', 'employee_manager', 'employee');
CREATE TYPE user_types AS ENUM ('full_time', 'intern', 'freelancer');
CREATE TYPE asset_types AS ENUM ('laptop', 'mouse', 'keyboard', 'charger', 'phone');
CREATE TYPE asset_statuses AS ENUM ('available', 'assigned', 'waiting_for_repair', 'service', 'damaged');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone_no TEXT UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by UUID,
    updated_at TIMESTAMPTZ,
    user_type user_types NOT NULL,
    archived_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email) WHERE archived_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone_no ON users (phone_no) WHERE archived_at IS NULL;
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) NOT NULL,
    role role_types NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    serial_no TEXT NOT NULL,
    model TEXT NOT NULL,
    brand TEXT NOT NULL,
    asset_status asset_statuses DEFAULT 'available',
    asset_type asset_types NOT NULL,
    purchased_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by TEXT,
    updated_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS asset_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID REFERENCES assets(id) NOT NULL,
    assigned_to UUID REFERENCES users(id),
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    retrieved_at TIMESTAMPTZ,
    description TEXT
);

CREATE TABLE IF NOT EXISTS laptops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assets_id UUID REFERENCES assets(id) NOT NULL,
    name TEXT NOT NULL,
    IEMI_no_1 TEXT NOT NULL,
    IEMI_no_2 TEXT NOT NULL,
    price TEXT NOT NULL,
    ram INTEGER NOT NULL,
    os TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS mouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assets_id UUID REFERENCES assets(id) NOT NULL,
    name TEXT NOT NULL,
    IEMI_no_1 TEXT NOT NULL,
    IEMI_no_2 TEXT NOT NULL,
    price TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS keyboards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assets_id UUID REFERENCES assets(id) NOT NULL,
    name TEXT NOT NULL,
    IEMI_no_1 TEXT NOT NULL,
    IEMI_no_2 TEXT NOT NULL,
    type TEXT NOT NULL,
    price TEXT NOT NULL,
    is_wired BOOLEAN NOT NULL,
    keys TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS mobiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assets_id UUID REFERENCES assets(id) NOT NULL,
    name TEXT NOT NULL,
    IEMI_no_1 TEXT NOT NULL,
    IEMI_no_2 TEXT NOT NULL,
    price TEXT NOT NULL,
    ram INTEGER NOT NULL,
    os TEXT NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ
);

COMMIT;