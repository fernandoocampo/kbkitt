-- Add up migration script here
CREATE INDEX kb_tags_index ON kbs USING GIN(TAGS);