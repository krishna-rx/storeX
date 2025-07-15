CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    start_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_date TIMESTAMPTZ,
    cost INTEGER CHECK (cost >= 0),
    service_type TEXT,
    description TEXT
);
CREATE INDEX IF NOT EXISTS idx_assets_asset_type ON assets(asset_type);
CREATE INDEX IF NOT EXISTS idx_assets_asset_status ON assets(asset_status);
CREATE INDEX IF NOT EXISTS idx_assets_archived_at ON assets(archived_at);

CREATE INDEX IF NOT EXISTS idx_assignments_assigned_by ON asset_assignments(assigned_by);
CREATE INDEX IF NOT EXISTS idx_assignments_assigned_to ON asset_assignments(assigned_to);
CREATE INDEX IF NOT EXISTS idx_assignments_asset_id ON asset_assignments(asset_id);

CREATE INDEX IF NOT EXISTS idx_laptops_assets_id ON laptops(assets_id);
CREATE INDEX IF NOT EXISTS idx_mouses_assets_id ON mouses(assets_id);
CREATE INDEX IF NOT EXISTS idx_keyboards_assets_id ON keyboards(assets_id);
CREATE INDEX IF NOT EXISTS idx_mobiles_assets_id ON mobiles(assets_id);
