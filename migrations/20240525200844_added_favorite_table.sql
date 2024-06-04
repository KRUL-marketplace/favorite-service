-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE favorite_lists (
    favorite_list_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE favorite_items (
    item_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    favorite_list_id UUID REFERENCES favorite_lists(favorite_list_id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (favorite_list_id, product_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists favorite_lists;
drop table if exists favorite_items;
-- +goose StatementEnd
