-- Shared geo, places, news, and HTTP ingest tables (no map_ prefix).
-- Closed enums persist as SMALLINT. Audit columns on every table.

CREATE TABLE IF NOT EXISTS countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(2) NOT NULL,
    code3 VARCHAR(3) NOT NULL,
    name TEXT NOT NULL,
    name_local TEXT NOT NULL DEFAULT '',
    phone_code TEXT NOT NULL DEFAULT '',
    currency VARCHAR(3) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT countries_code_key UNIQUE (code),
    CONSTRAINT countries_code3_key UNIQUE (code3),
    CONSTRAINT countries_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT countries_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS administrative_divisions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_id UUID NOT NULL,
    parent_id UUID,
    level SMALLINT NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    name_en TEXT NOT NULL DEFAULT '',
    full_name TEXT NOT NULL DEFAULT '',
    type TEXT NOT NULL,
    path TEXT NOT NULL,
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT administrative_divisions_country_code_key UNIQUE (country_id, code),
    CONSTRAINT administrative_divisions_country_id_fkey FOREIGN KEY (country_id) REFERENCES countries (id) ON DELETE RESTRICT,
    CONSTRAINT administrative_divisions_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES administrative_divisions (id) ON DELETE RESTRICT,
    CONSTRAINT administrative_divisions_level_check CHECK (level >= 1),
    CONSTRAINT administrative_divisions_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT administrative_divisions_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_admin_div_country_parent ON administrative_divisions (country_id, parent_id);
CREATE INDEX IF NOT EXISTS idx_admin_div_country_level ON administrative_divisions (country_id, level);
CREATE INDEX IF NOT EXISTS idx_admin_div_path ON administrative_divisions (path);
CREATE INDEX IF NOT EXISTS idx_admin_div_path_prefix ON administrative_divisions (path text_pattern_ops);

CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    location_key TEXT,
    country_id UUID,
    admin_division_id UUID,
    street TEXT NOT NULL DEFAULT '',
    postal_code TEXT NOT NULL DEFAULT '',
    formatted TEXT NOT NULL DEFAULT '',
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT locations_country_id_fkey FOREIGN KEY (country_id) REFERENCES countries (id) ON DELETE RESTRICT,
    CONSTRAINT locations_admin_division_id_fkey FOREIGN KEY (admin_division_id) REFERENCES administrative_divisions (id) ON DELETE RESTRICT,
    CONSTRAINT locations_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT locations_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS locations_location_key_key
    ON locations (location_key)
    WHERE location_key IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_locations_formatted_norm
    ON locations (lower(btrim(formatted)))
    WHERE formatted <> '';
CREATE INDEX IF NOT EXISTS idx_locations_street_norm
    ON locations (lower(btrim(street)))
    WHERE street <> '';
CREATE INDEX IF NOT EXISTS idx_locations_coords
    ON locations (lat, lng)
    WHERE lat IS NOT NULL AND lng IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_locations_country_id ON locations (country_id);
CREATE INDEX IF NOT EXISTS idx_locations_admin_division_id ON locations (admin_division_id);

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT NOT NULL,
    name TEXT NOT NULL,
    name_local TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT categories_key_key UNIQUE (key),
    CONSTRAINT categories_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT categories_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_categories_active_sort ON categories (is_active, sort_order, name);

CREATE TABLE IF NOT EXISTS places (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL,
    category_id UUID NOT NULL,
    place_key TEXT NOT NULL,
    name TEXT NOT NULL,
    unit TEXT NOT NULL DEFAULT '',
    status SMALLINT NOT NULL DEFAULT 1,
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT places_location_place_key_key UNIQUE (location_id, place_key),
    CONSTRAINT places_location_id_fkey FOREIGN KEY (location_id) REFERENCES locations (id) ON DELETE RESTRICT,
    CONSTRAINT places_category_id_fkey FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE RESTRICT,
    CONSTRAINT places_status_check CHECK (status IN (1, 2, 3)),
    CONSTRAINT places_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT places_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_places_pinnable
    ON places (category_id, status)
    WHERE status = 2 AND lat IS NOT NULL AND lng IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_places_location_id ON places (location_id);
CREATE INDEX IF NOT EXISTS idx_places_category_id ON places (category_id);
CREATE INDEX IF NOT EXISTS idx_places_name ON places (lower(name));

CREATE TABLE IF NOT EXISTS data_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    http_method SMALLINT NOT NULL DEFAULT 1,
    url TEXT NOT NULL,
    headers JSONB NOT NULL DEFAULT '{}'::jsonb,
    query_params JSONB NOT NULL DEFAULT '{}'::jsonb,
    body TEXT NOT NULL DEFAULT '',
    field_mapping JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT data_sources_http_method_check CHECK (http_method IN (1, 2, 3, 4)),
    CONSTRAINT data_sources_url_check CHECK (url ~* '^https?://'),
    CONSTRAINT data_sources_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT data_sources_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_data_sources_enabled ON data_sources (enabled);

CREATE TABLE IF NOT EXISTS news (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    place_id UUID NOT NULL,
    category_id UUID NOT NULL,
    source_id UUID,
    source_name TEXT NOT NULL,
    original_url TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT news_original_url_key UNIQUE (original_url),
    CONSTRAINT news_place_id_fkey FOREIGN KEY (place_id) REFERENCES places (id) ON DELETE CASCADE,
    CONSTRAINT news_category_id_fkey FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE RESTRICT,
    CONSTRAINT news_source_id_fkey FOREIGN KEY (source_id) REFERENCES data_sources (id) ON DELETE SET NULL,
    CONSTRAINT news_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT news_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_news_place_created ON news (place_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_news_category_id ON news (category_id);

CREATE TABLE IF NOT EXISTS data_ingest_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status SMALLINT NOT NULL DEFAULT 1,
    triggered_by UUID NOT NULL,
    error_message TEXT NOT NULL DEFAULT '',
    sources_total INTEGER NOT NULL DEFAULT 0,
    success INTEGER NOT NULL DEFAULT 0,
    error INTEGER NOT NULL DEFAULT 0,
    items_success INTEGER NOT NULL DEFAULT 0,
    items_error INTEGER NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT data_ingest_runs_status_check CHECK (status IN (1, 2, 3, 4)),
    CONSTRAINT data_ingest_runs_triggered_by_fkey FOREIGN KEY (triggered_by) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT data_ingest_runs_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT data_ingest_runs_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_data_ingest_runs_created ON data_ingest_runs (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_data_ingest_runs_active ON data_ingest_runs (status) WHERE status IN (1, 2);

CREATE TABLE IF NOT EXISTS data_ingest_run_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL,
    source_id UUID,
    source_name TEXT NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    error_message TEXT NOT NULL DEFAULT '',
    items_success INTEGER NOT NULL DEFAULT 0,
    items_error INTEGER NOT NULL DEFAULT 0,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT data_ingest_run_sources_status_check CHECK (status IN (1, 2, 3, 4, 5)),
    CONSTRAINT data_ingest_run_sources_run_id_fkey FOREIGN KEY (run_id) REFERENCES data_ingest_runs (id) ON DELETE CASCADE,
    CONSTRAINT data_ingest_run_sources_source_id_fkey FOREIGN KEY (source_id) REFERENCES data_sources (id) ON DELETE SET NULL,
    CONSTRAINT data_ingest_run_sources_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT data_ingest_run_sources_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_data_ingest_run_sources_run ON data_ingest_run_sources (run_id);
