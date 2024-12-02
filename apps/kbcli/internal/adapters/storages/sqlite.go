package storages

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type SQLiteSetup struct {
	DB *sql.DB
}

// SQLite implements logic to store data into sqlite repository.
type SQLite struct {
	db *sql.DB
}

const (
	sqliteVersion = "sqlite3"

	createKBSQL = `INSERT INTO kbs
	(KB_ID, KB_KEY, KB_VALUE, NOTES, CATEGORY, TAG_VALUES, REFERENCE, NAMESPACE, CREATED_ON)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?)`

	updateKBSQL = `UPDATE kbs
SET KB_KEY = ?, KB_VALUE = ?, NOTES = ?, CATEGORY = ?, TAG_VALUES = ?, REFERENCE = ?, NAMESPACE = ?
WHERE KB_ID = ?`

	queryAKBByIDSQL            = "SELECT INTERNAL_ID, KB_ID, KB_KEY, KB_VALUE, NOTES, NAMESPACE, CATEGORY, TAG_VALUES, REFERENCE, CREATED_ON FROM kbs WHERE KB_ID = ?"
	queryAKBByKeySQL           = "SELECT INTERNAL_ID, KB_ID, KB_KEY, KB_VALUE, NOTES, NAMESPACE, CATEGORY, TAG_VALUES, REFERENCE, CREATED_ON FROM kbs WHERE KB_KEY = ?"
	queryKBsByFilterSQL        = "SELECT k.KB_ID, k.KB_KEY, k.CATEGORY, k.NAMESPACE, k.TAG_VALUES FROM kbs k %s;"
	queryKBsByFilterAndTagsSQL = "SELECT k.KB_ID, k.KB_KEY, k.CATEGORY, k.NAMESPACE, k.TAG_VALUES FROM kbs k JOIN tags_idx t ON (t.rowid = k.INTERNAL_ID) %s;"
	countKBsByCategorySQL      = "SELECT COUNT(k.KB_ID) FROM kbs k WHERE k.CATEGORY = ?"
	countKBsByFilterSQL        = "SELECT COUNT(k.KB_ID) FROM kbs k %s;"
	countKBsByFilterAndTagsSQL = "SELECT COUNT(k.KB_ID) FROM kbs k JOIN tags_idx t ON (t.rowid = k.INTERNAL_ID) %s;"

	countKBsSQL = "SELECT COUNT(k.KB_ID) FROM kbs k %s;"
	queryKBsSQL = `SELECT k.INTERNAL_ID, k.KB_ID, k.KB_KEY, k.KB_VALUE, k.NOTES, k.NAMESPACE, k.CATEGORY, k.TAG_VALUES, k.REFERENCE, k.CREATED_ON 
FROM kbs k %s`

	createKBTableSQL = `DROP TABLE IF EXISTS kbs;
CREATE TABLE IF NOT EXISTS kbs (
	INTERNAL_ID INTEGER PRIMARY KEY AUTOINCREMENT,
	KB_ID VARCHAR(36) UNIQUE,
	KB_KEY VARCHAR(64) NOT NULL UNIQUE,
	KB_VALUE TEXT NOT NULL,
	NOTES TEXT NOT NULL,
	CATEGORY VARCHAR(64) NOT NULL,
	NAMESPACE VARCHAR(64) NOT NULL,
	TAG_VALUES VARCHAR(256) NOT NULL,
	REFERENCE VARCHAR(64),
	CREATED_ON DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

	createFTSKBTableSQL = `DROP TABLE IF EXISTS tags_idx;
CREATE VIRTUAL TABLE IF NOT EXISTS tags_idx 
USING fts5(
	tag_values,
	content='kbs',
	content_rowid='INTERNAL_ID'
);
`

	createTriggerAfterInsertSQL = `DROP TRIGGER IF EXISTS kbs_ai;
CREATE TRIGGER kbs_ai AFTER INSERT ON kbs BEGIN
INSERT INTO tags_idx(rowid, tag_values) VALUES (new.INTERNAL_ID, new.TAG_VALUES);
END;
`

	createTriggerAfterDeleteSQL = `DROP TRIGGER IF EXISTS kbs_ad;
CREATE TRIGGER kbs_ad AFTER DELETE ON kbs BEGIN
DELETE FROM tags_idx WHERE rowid = old.INTERNAL_ID;
END;
`

	createTriggerAfterUpdateSQL = `DROP TRIGGER IF EXISTS kbs_au;
CREATE TRIGGER kbs_au AFTER UPDATE ON kbs BEGIN
INSERT INTO tags_idx(tags_idx, rowid, tag_values) VALUES('delete', old.INTERNAL_ID, old.TAG_VALUES);
INSERT INTO tags_idx(rowid, tag_values) VALUES (new.INTERNAL_ID, new.TAG_VALUES);
END;`
)

func NewSQLite(setup *SQLiteSetup) *SQLite {
	newSQLite := SQLite{
		db: setup.DB,
	}

	return &newSQLite
}

// CreateSQLiteConnection creates a new sqlite connection with the given file path.
func CreateSQLiteConnection(dbPath string) (*sql.DB, error) {
	db, err := sql.Open(sqliteVersion, fmt.Sprintf("file:%s", dbPath))
	if err != nil {
		return nil, fmt.Errorf("unable to create cx to sqlite file: %w", err)
	}

	return db, nil
}

func (s *SQLite) Version() (string, error) {
	var version string

	err := s.db.QueryRow(`SELECT sqlite_version()`).Scan(&version)
	if err != nil {
		return "", fmt.Errorf("unable to read version: %w", err)
	}

	return version, nil
}

func (s *SQLite) InitializeDB(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, createKBTableSQL)
	if err != nil {
		return fmt.Errorf("unable to initialize db with kb table: %w", err)
	}

	_, err = s.db.ExecContext(ctx, createFTSKBTableSQL)
	if err != nil {
		log.Println("unable to initialize db with fts table", err)
		return err
	}

	// -- Triggers to keep the FTS index up to date.
	_, err = s.db.ExecContext(ctx, createTriggerAfterInsertSQL)
	if err != nil {
		return fmt.Errorf("unable to initialize db with trigger ai: %w", err)
	}
	_, err = s.db.ExecContext(ctx, createTriggerAfterDeleteSQL)
	if err != nil {
		return fmt.Errorf("unable to initialize db with trigger ad: %w", err)
	}

	_, err = s.db.ExecContext(ctx, createTriggerAfterUpdateSQL)
	if err != nil {
		return fmt.Errorf("unable to initialize db with trigger au: %w", err)
	}

	return nil
}

func (s *SQLite) Create(ctx context.Context, newKB kbs.KB) (string, error) {
	stmt, err := s.db.Prepare(createKBSQL)
	if err != nil {
		return "", fmt.Errorf("unable to create kb: %w", err)
	}

	defer stmt.Close()

	dbKB := toDBKB(&newKB)

	result, err := stmt.ExecContext(ctx,
		dbKB.KeyID, dbKB.Key, dbKB.Value,
		dbKB.Notes, dbKB.Category, dbKB.Tags, dbKB.Reference,
		dbKB.Namespace, dbKB.DateCreated,
	)
	if err != nil {
		return "", fmt.Errorf("unable to create kb: %w", err)
	}

	lastInserID, err := result.LastInsertId()
	if err != nil {
		slog.Error("unable to get id from embedded database", slog.String("error", err.Error()))
	}

	return fmt.Sprintf("%d", lastInserID), nil
}

func (s *SQLite) GetAll(ctx context.Context, filter kbs.KBQueryFilter) (*kbs.GetAllResult, error) {
	result := kbs.GetAllResult{
		Total:  0,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}

	searchFilters := buildSQLFilters(filter, countKBsSQL, queryKBsSQL)

	count, err := s.queryCount(ctx, searchFilters)
	if err != nil {
		slog.Error("running count query to get all kbs",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.countStatement),
			slog.String("error", err.Error()),
		)

		return nil, errUnableToGetAllKBS
	}

	result.Total = count

	kbsFound, err := s.queryKBs(ctx, searchFilters)
	if err != nil {
		slog.Error("querying all kbs",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.query),
			slog.String("error", err.Error()),
		)

		return nil, errUnableToGetAllKBS
	}

	result.KBs = toKBs(kbsFound)

	return &result, nil
}

func (s *SQLite) Search(ctx context.Context, filter kbs.KBQueryFilter) (*kbs.SearchResult, error) {
	result := kbs.SearchResult{
		Total:  0,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}

	querySQL := queryKBsByFilterSQL
	countSQL := countKBsByFilterSQL

	if filter.Keyword != "" {
		querySQL = queryKBsByFilterAndTagsSQL
		countSQL = countKBsByFilterAndTagsSQL
	}

	searchFilters := buildSQLFilters(filter, countSQL, querySQL)

	count, err := s.queryCount(ctx, searchFilters)
	if err != nil {
		slog.Error("running count query of kbs that match given criteria",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.countStatement),
			slog.String("error", err.Error()),
		)

		return nil, errUnableToSearchKBS
	}

	result.Total = count

	kbsFound, err := s.queryKBItems(ctx, searchFilters)
	if err != nil {
		slog.Error("checking if rows results has an error",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.query),
			slog.String("error", err.Error()),
		)

		return nil, errUnableToSearchKBS
	}

	result.Items = toKBItems(kbsFound)

	return &result, nil
}

func (s *SQLite) GetByID(ctx context.Context, id string) (*kbs.KB, error) {
	kb, err := s.getKBRecord(ctx, queryAKBByIDSQL, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get kb by id: %w", err)
	}

	return kb, nil
}

func (s *SQLite) GetByKey(ctx context.Context, key string) (*kbs.KB, error) {
	kb, err := s.getKBRecord(ctx, queryAKBByKeySQL, key)
	if err != nil {
		return nil, fmt.Errorf("unable to get kb by key: %w", err)
	}

	return kb, nil
}

func (s *SQLite) getKBRecord(ctx context.Context, sqlQuery, value string) (*kbs.KB, error) {
	row := s.db.QueryRowContext(ctx, sqlQuery, value)

	var aKB kb

	err := row.Scan(&aKB.InternalID, &aKB.KeyID, &aKB.Key, &aKB.Value, &aKB.Notes, &aKB.Namespace, &aKB.Category, &aKB.Tags, &aKB.Reference, &aKB.DateCreated)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil // it does not exist
	}

	if err != nil {
		return nil, fmt.Errorf("unable to get kb record: %w", err)
	}

	return aKB.toKB(), nil
}

func (s *SQLite) Update(ctx context.Context, kb *kbs.KB) error {
	stmt, err := s.db.Prepare(updateKBSQL)
	if err != nil {
		return fmt.Errorf("unable to update kb: %w", err)
	}

	defer stmt.Close()

	dbKB := toDBKB(kb)

	result, err := stmt.ExecContext(ctx,
		dbKB.Key, dbKB.Value, dbKB.Notes,
		dbKB.Category, dbKB.Tags, dbKB.Reference,
		dbKB.Namespace, dbKB.KeyID,
	)
	if err != nil {
		return fmt.Errorf("unable to update kb: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("rows affected", slog.String("error", err.Error()))

		return nil
	}

	slog.Info("rows affected", slog.Int64("rows", rowsAffected))

	return nil
}

func (s *SQLite) queryCount(ctx context.Context, searchFilters *filterBuilder) (int, error) {
	var count int

	countStmt, err := s.db.Prepare(searchFilters.countStatement)
	if err != nil {
		slog.Error("building count kbs prepared statement",
			slog.Any("filter", searchFilters),
			slog.String("error", err.Error()),
		)

		return -1, fmt.Errorf("unable to build query to count kbs: %w", err)
	}

	defer countStmt.Close()

	row := countStmt.QueryRowContext(ctx, searchFilters.countArgs...)

	err = row.Scan(&count)
	if err != nil {
		slog.Error("scanning count of kbs found",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.countStatement),
			slog.String("error", err.Error()),
		)

		return -1, fmt.Errorf("unable to scanning count of kbs found: %w", err)
	}

	if row.Err() != nil {
		return -1, fmt.Errorf("unable to count kbs: %w", row.Err())
	}

	return count, nil
}

func (s *SQLite) queryKBItems(ctx context.Context, searchFilters *filterBuilder) ([]kbItem, error) {
	rows, err := s.db.QueryContext(ctx, searchFilters.query, searchFilters.queryArgs...)
	if err != nil {
		slog.Error("running query to find kbs with given criteria",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.query),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("unable to query kbs: %w", err)
	}

	defer rows.Close()

	kbsFound := make([]kbItem, 0)

	for rows.Next() {
		kb := new(kbItem)
		// id, firstname, lastname, nickname, country
		rowErr := rows.Scan(&kb.ID, &kb.Key, &kb.Category, &kb.Namespace, &kb.Tags)
		if rowErr != nil {
			slog.Error("scanning rows for searching kbs with search criteria",
				slog.Any("filter", searchFilters),
				slog.String("query", searchFilters.query),
				slog.String("error", rowErr.Error()),
			)

			return nil, fmt.Errorf("unable to scan kb rows: %w", rowErr)
		}

		kbsFound = append(kbsFound, *kb)
	}

	if err := rows.Err(); err != nil {
		slog.Error("checking if rows results has an error",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.query),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("kb search query had some errors: %w", err)
	}

	return kbsFound, nil
}

func (s *SQLite) queryKBs(ctx context.Context, searchFilters *filterBuilder) ([]kb, error) {
	rows, err := s.db.QueryContext(ctx, searchFilters.query, searchFilters.queryArgs...)
	if err != nil {
		slog.Error("running query to get all kbs",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.query),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("unable to query all kbs: %w", err)
	}

	defer rows.Close()

	kbsFound := make([]kb, 0)

	for rows.Next() {
		kb := new(kb)
		rowErr := rows.Scan(&kb.InternalID, &kb.KeyID, &kb.Key, &kb.Value, &kb.Notes, &kb.Namespace, &kb.Category, &kb.Tags, &kb.Reference, &kb.DateCreated)
		if rowErr != nil {
			slog.Error("scanning rows to get all kbs",
				slog.Any("filter", searchFilters),
				slog.String("query", searchFilters.query),
				slog.String("error", rowErr.Error()),
			)

			return nil, fmt.Errorf("unable to scan kb rows: %w", rowErr)
		}

		kbsFound = append(kbsFound, *kb)
	}

	if err := rows.Err(); err != nil {
		slog.Error("checking if rows results has an error",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.query),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("get all kbs query had some errors: %w", err)
	}

	return kbsFound, nil
}

func (s *SQLite) CountByCategory(ctx context.Context, category string) (int64, error) {
	countStmt, err := s.db.Prepare(countKBsByCategorySQL)
	if err != nil {
		slog.Error("building count kbs by category prepared statement",
			slog.Any("category", category),
			slog.Any("query", countKBsByCategorySQL),
			slog.String("error", err.Error()),
		)

		return -1, fmt.Errorf("unable to build query to count kbs by category: %w", err)
	}

	defer countStmt.Close()

	row := countStmt.QueryRowContext(ctx, category)

	var count int64

	err = row.Scan(&count)
	if err != nil {
		slog.Error("scanning count of kbs by category found",
			slog.Any("category", category),
			slog.String("query", countKBsByCategorySQL),
			slog.String("error", err.Error()),
		)

		return -1, fmt.Errorf("unable to scanning count of kbs found: %w", err)
	}

	if row.Err() != nil {
		return -1, fmt.Errorf("unable to count kbs: %w", row.Err())
	}

	return count, nil
}

func buildSQLFilters(filters kbs.KBQueryFilter, countSQL, querySQL string) *filterBuilder {
	newFilterBuilder := &filterBuilder{
		filters:   make([]string, 0),
		countArgs: make([]interface{}, 0),
		queryArgs: make([]interface{}, 0),
	}

	if filters.Keyword != "" {
		newFilterBuilder.addCondition(tagValuesVirtualColumn, matchOperator, filters.Keyword)
	}

	if filters.Key != "" {
		newFilterBuilder.addCondition(keyColumn, likeOperator, fmt.Sprintf("%%%s%%", filters.Key))
	}

	if filters.Category != "" {
		newFilterBuilder.addCondition(categoryColumn, equalsOperator, filters.Category)
	}

	if filters.Namespace != "" {
		newFilterBuilder.addCondition(namespaceColumn, equalsOperator, filters.Namespace)
	}

	var countWhereClause string
	for _, v := range newFilterBuilder.filters {
		countWhereClause += v
	}

	countStatement := fmt.Sprintf(countSQL, countWhereClause)
	newFilterBuilder.countStatement = countStatement

	newFilterBuilder.addFilter(fmt.Sprintf(" %s", limitOperator), filters.Limit, true)
	newFilterBuilder.addFilter(fmt.Sprintf(" %s", offsetOperator), filters.Offset, true)

	var whereClause string
	for _, v := range newFilterBuilder.filters {
		whereClause += v
	}

	queryStatement := fmt.Sprintf(querySQL, whereClause)
	newFilterBuilder.query = queryStatement

	return newFilterBuilder
}

func (s *SQLite) Close() {
	if s == nil || s.db == nil {
		return // just for initializing app
	}

	err := s.db.Close()
	if err != nil {
		slog.Error("unable to close db connection", slog.String("error", err.Error()))
	}
}
