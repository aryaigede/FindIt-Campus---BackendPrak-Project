-- 1. Roles
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Seed default roles
INSERT INTO roles (name)
VALUES ('admin'),
    ('security'),
    ('user') ON CONFLICT (name) DO NOTHING;
-- 2. Users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    role_id INT NOT NULL REFERENCES roles(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    name VARCHAR(120) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    student_id VARCHAR(30),
    phone VARCHAR(25),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- 3. Categories
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- 4. Locations
CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(120) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- 5. Item Reports
CREATE TABLE IF NOT EXISTS item_reports (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    category_id INT NOT NULL REFERENCES categories(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    location_id INT NOT NULL REFERENCES locations(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    report_type VARCHAR(10) NOT NULL CHECK (report_type IN ('lost', 'found')),
    title VARCHAR(180) NOT NULL,
    description TEXT NOT NULL,
    custody_note VARCHAR(150),
    event_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (
        status IN ('open', 'claimed', 'returned', 'closed')
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- 6. Item Images
CREATE TABLE IF NOT EXISTS item_images (
    id SERIAL PRIMARY KEY,
    report_id INT NOT NULL REFERENCES item_reports(id) ON UPDATE CASCADE ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- 7. Claims
CREATE TABLE IF NOT EXISTS claims (
    id SERIAL PRIMARY KEY,
    report_id INT NOT NULL REFERENCES item_reports(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    claimant_id INT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    reviewed_by INT REFERENCES users(id) ON UPDATE CASCADE ON DELETE
    SET NULL,
        proof_description TEXT NOT NULL,
        proof_image_url TEXT,
        status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
        review_note TEXT,
        reviewed_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- 8. Item Returns (mendukung serah terima laporan 'found' via claim maupun laporan 'lost' langsung)
CREATE TABLE IF NOT EXISTS item_returns (
    id SERIAL PRIMARY KEY,
    report_id INT UNIQUE NOT NULL REFERENCES item_reports(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    claim_id INT UNIQUE REFERENCES claims(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    handed_over_by INT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    received_by INT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    returned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes TEXT
);
-- 9. Report Histories
CREATE TABLE IF NOT EXISTS report_histories (
    id SERIAL PRIMARY KEY,
    report_id INT NOT NULL REFERENCES item_reports(id) ON UPDATE CASCADE ON DELETE CASCADE,
    actor_id INT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    old_status VARCHAR(20),
    new_status VARCHAR(20) NOT NULL,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Business Rule Constraints & Performance Indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_claims_one_approved_per_report ON claims (report_id)
WHERE status = 'approved';
CREATE UNIQUE INDEX IF NOT EXISTS idx_claims_one_pending_per_user ON claims (report_id, claimant_id)
WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_item_reports_filters ON item_reports (status, report_type, category_id, location_id);
CREATE INDEX IF NOT EXISTS idx_item_reports_user_id ON item_reports (user_id);
CREATE INDEX IF NOT EXISTS idx_claims_report_id ON claims (report_id);
CREATE INDEX IF NOT EXISTS idx_report_histories_report_id ON report_histories (report_id);