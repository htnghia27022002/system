-- Elasticsearch was removed; admin search now queries Postgres directly,
-- so the search indexing outbox is no longer used.
DROP TABLE IF EXISTS search_outbox;
