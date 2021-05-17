-- +goose Up
CREATE TABLE IF NOT EXISTS  list (
	id SERIAL, 
	cidr varchar(100) NOT NULL DEFAULT '',  
	list varchar(36) NOT NULL DEFAULT '', 
	PRIMARY KEY (id)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE list;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
