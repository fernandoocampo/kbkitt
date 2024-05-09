use crate::errors::error::Error;
use crate::kbs::storage::Storer as kb_storage;
use crate::types::categories::Category;
use crate::types::kbs::{KBID, KBItem, KBQueryFilter, KnowledgeBase};

use tracing::debug;
use async_trait::async_trait;

use sqlx::postgres::{PgPool, PgPoolOptions, PgRow};
use sqlx::Row;

#[derive(Debug, Clone)]
pub struct Store {
    pub connection: PgPool,
}

impl Store {
    pub async fn new(db_url: &str) -> Self {
        let db_pool = match PgPoolOptions::new()
            .max_connections(5)
            .connect(db_url)
            .await
        {
            Ok(pool) => pool,
            Err(e) => panic!("couldn't establish DB connection: {e}")
        };

        Store {
            connection: db_pool,
        }
    }
}

#[async_trait]
impl kb_storage for Store {
    /// get a Knowledge base with the given id.
    async fn get_kb_by_id(&self, id: KBID) -> Result<KnowledgeBase, Error> {
        match sqlx::query("SELECT * FROM kbs WHERE KB_ID = $1")
            .bind(id.to_string())
            .map(|row: PgRow| KnowledgeBase {
                id: KBID(row.get("KB_ID")),
                key: row.get("KB_KEY"),
                value: row.get("KB_VALUE"),
                notes: row.get("NOTES"),
                kind: row.get("KIND"),
                tags: row.get::<String, _>("TAGS")
                    .split(" ")
                    .map(|s| s.to_string())
                    .collect::<Vec<String>>(),
            })
            .fetch_one(&self.connection)
            .await
        {
            Ok(kb) => Ok(kb),
            Err(sqlx::Error::RowNotFound) => Ok(KnowledgeBase::default()),
            Err(e) => {
                tracing::event!(tracing::Level::ERROR, "querying kb by id {:?} : {:?}", id, e);
                Err(Error::DatabaseQueryError)
            }
        }
    }
    /// get a Knowledge base with the given key.
    async fn get_kb_by_key(&self, key: String) -> Result<KnowledgeBase, Error> {
        match sqlx::query("SELECT * FROM kbs WHERE KB_KEY = $1")
            .bind(key)
            .map(|row: PgRow| KnowledgeBase {
                id: KBID(row.get("KB_ID")),
                key: row.get("KB_KEY"),
                value: row.get("KB_VALUE"),
                notes: row.get("NOTES"),
                kind: row.get("KIND"),
                tags: row.get::<String, _>("TAGS")
                    .split(" ")
                    .map(|s| s.to_string())
                    .collect::<Vec<String>>(),
            })
            .fetch_one(&self.connection)
            .await
        {
            Ok(kb) => Ok(kb),
            Err(sqlx::Error::RowNotFound) => Ok(KnowledgeBase::default()),
            Err(e) => {
                tracing::event!(tracing::Level::ERROR, "querying kb by key {:?} : {:?}", key, e);
                Err(Error::DatabaseQueryError)
            }
        }
    }

    /// get a list of knowledge base entries where their keys contain the given keywords.
    async fn search_by_key(&self, filter: KBQueryFilter) -> Result<Vec<KBItem>, Error> {
        match sqlx::query("SELECT KB_ID, KB_KEY, KIND, TAGS FROM kbs WHERE KB_KEY LIKE $1 LIMIT $2 OFFSET $3")
            .bind(format!("%{}%", filter.keyword))
            .bind(filter.limit)
            .bind(filter.offset)
            .map(|row: PgRow| KBItem {
                id: KBID(row.get("KB_ID")),
                key: row.get("KB_KEY"),
                kind: row.get("KIND"),
                tags: row.get::<String, _>("TAGS")
                    .split(" ")
                    .map(|s| s.to_string())
                    .collect::<Vec<String>>(),
            })
            .fetch_all(&self.connection)
            .await {
            Ok(kbs) => {
                debug!("found some kbs: {:?}", kbs);
                Ok(kbs)
            }
            Err(e) => {
                tracing::event!(tracing::Level::ERROR, "{:?}", e);
                Err(Error::DatabaseQueryError)
            }
        }
    }

    /// get a list of knowledge base entries where their keys contain the given keywords.
    async fn search(&self, filter: KBQueryFilter) -> Result<Vec<KBItem>, Error> {
        match sqlx::query("SELECT KB_ID, KB_KEY, KIND, TAGS FROM kbs WHERE TAGS @@ to_tsquery($1) LIMIT $2 OFFSET $3")
            .bind(format!("'{}'", filter.keyword))
            .bind(filter.limit)
            .bind(filter.offset)
            .map(|row: PgRow| KBItem {
                id: KBID(row.get("KB_ID")),
                key: row.get("KB_KEY"),
                kind: row.get("KIND"),
                tags: row.get::<String, _>("TAGS")
                    .split(" ")
                    .map(|s| s.to_string())
                    .collect::<Vec<String>>(),
            })
            .fetch_all(&self.connection)
            .await {
            Ok(kbs) => {
                debug!("found some kbs: {:?}", kbs);
                Ok(kbs)
            }
            Err(e) => {
                tracing::event!(tracing::Level::ERROR, "{:?}", e);
                Err(Error::DatabaseQueryError)
            }
        }
    }

    /// save given knowledge base in the repository.
    async fn save_kb(&self, kb: KnowledgeBase) -> Result<KBID, Error> {
        debug!("adding kb to postgresql db: {:?}", kb);

        match sqlx::query("INSERT INTO kb (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ($1, $2, $3, $4, $5, $6) RETURNING KB_ID")
            .bind(kb.id.to_string())
            .bind(kb.key)
            .bind(kb.value)
            .bind(kb.notes)
            .bind(kb.kind)
            .bind(kb.tags)
            .map(|row: PgRow| KBID(row.get("KB_ID")))
            .fetch_one(&self.connection)
            .await
        {
            Ok(kb_id) => {
                debug!("kb was added to postgres database: {:?}", kb_id);
                Ok(kb_id)
            }
            Err(e) => {
                tracing::event!(tracing::Level::ERROR, "{:?}", e);
                Err(Error::DatabaseQueryError)
            }
        }
    }
    /// save given category in the repository.
    async fn save_category(&self, kb: Category) -> Result<String, Error> {}
}