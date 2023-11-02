package dbmodel

var (
	CreateBlockHeadersStmt = `CREATE TABLE IF NOT EXISTS block_headers (
		block_hash CHAR(66) NOT NULL,
		parent_hash CHAR(66) NOT NULL,
		block_number BIGINT NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		PRIMARY KEY (block_hash)
		);`

	CreateIndexBlockHeaderStmt = `CREATE INDEX IF NOT EXISTS idx_block_headers ON block_headers(block_hash,parent_hash,block_number,timestamp)`

	CreateSyncEventsStmt = `CREATE TABLE IF NOT EXISTS sync_events (
			id uuid NOT NULL,
			topic varchar(255) NOT NULL,
			key varchar(255) NOT NULL,
			sequence SERIAL NOT NULL,
			created_at TIMESTAMP NOT NULL,
			PRIMARY KEY (id)
			);`
)
