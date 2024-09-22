use crate::errors::error::Error;
use crate::kbs::storage::Storer as kb_storage;
use crate::types::categories::{Category, CategoryFilter};
use crate::types::kbs::{KBItem, KBQueryFilter, KnowledgeBase, SearchResult, KBID};

use async_trait::async_trait;
use log::error;
use tracing::debug;

use sqlx::postgres::{PgPool, PgPoolOptions, PgRow};
use sqlx::Row;

#[derive(Debug, Clone)]
pub struct DBData {
    pub user: String,
    pub pwd: String,
    pub host: String,
    pub db_name: String,
    pub port: u16,
}

#[derive(Debug, Clone)]
pub struct Store {
    pub connection: PgPool,
}

impl DBData {
    pub fn build_url_connection(&self) -> String {
        format!(
            "postgres://{}:{}@{}:{}/{}",
            self.user, self.pwd, self.host, self.port, self.db_name,
        )
    }
}

impl Store {
    pub async fn new(db_data: DBData) -> Self {
        let db_pool = match PgPoolOptions::new()
            .max_connections(5)
            .connect(db_data.build_url_connection().as_str())
            .await
        {
            Ok(pool) => pool,
            Err(e) => panic!("couldn't establish DB connection: {e}"),
        };

        Store {
            connection: db_pool,
        }
    }

    pub async fn close(&self) {
        (self).connection.close().await
    }
}

#[async_trait]
impl kb_storage for Store {
    /// get a Knowledge base with the given id.
    async fn get_kb_by_id(&self, id: KBID) -> Result<KnowledgeBase, Error> {
        match sqlx::query("SELECT KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, REFERENCE, TAG_VALUES AS TAGS FROM kbs WHERE KB_ID = $1")
            .bind(id.to_string())
            .map(|row: PgRow| KnowledgeBase {
                id: KBID(row.get("kb_id")),
                key: row.get("kb_key"),
                value: row.get("kb_value"),
                notes: row.get("notes"),
                kind: row.get("kind"),
                reference: row.get("reference"),
                tags: row.get::<String, _>("tags")
                    .split(' ')
                    .map(|s| s.to_string())
                    .map(|s| s.replace('\'', ""))
                    .collect::<Vec<String>>()
            })
            .fetch_one(&self.connection)
            .await
        {
            Ok(kb) => {
                debug!("found kb: {:?}", kb);
                Ok(kb)
            }
            Err(sqlx::Error::RowNotFound) => Ok(KnowledgeBase::default()),
            Err(e) => {
                tracing::event!(tracing::Level::ERROR, "querying kb by id {:?} : {:?}", id, e);
                Err(Error::DatabaseQueryError)
            }
        }
    }
    /// get a Knowledge base with the given key.
    async fn get_kb_by_key(&self, key: String) -> Result<KnowledgeBase, Error> {
        match sqlx::query("SELECT KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, REFERENCE, TAG_VALUES AS TAGS FROM kbs WHERE KB_KEY = $1")
            .bind(key.clone())
            .map(|row: PgRow| KnowledgeBase {
                id: KBID(row.get("kb_id")),
                key: row.get("kb_key"),
                value: row.get("kb_value"),
                notes: row.get("notes"),
                kind: row.get("kind"),
                reference: row.get("reference"),
                tags: row.get::<String, _>("tags")
                    .split(' ')
                    .map(|s| s.to_string())
                    .map(|s| s.replace('\'', ""))
                    .collect::<Vec<String>>()
            })
            .fetch_one(&self.connection)
            .await
        {
            Ok(kb) => Ok(kb),
            Err(sqlx::Error::RowNotFound) => Ok(KnowledgeBase::default()),
            Err(e) => {
                error!("getting kb by key {:?}: {:?}", key, e);
                tracing::event!(
                    tracing::Level::ERROR,
                    "querying kb by key {:?} : {:?}",
                    key,
                    e
                );
                Err(Error::DatabaseQueryError)
            }
        }
    }

    /// get a list of knowledge base entries where their keys contain the given keywords.
    async fn search_by_key(&self, filter: KBQueryFilter) -> Result<SearchResult, Error> {
        // let's count first
        // let mut count: i64 = 0;
        let count: i64 =
            match sqlx::query_scalar("SELECT COUNT(*) as TOTAL FROM kbs WHERE KB_KEY LIKE $1")
                .bind(format!("%{}%", filter.key))
                .fetch_one(&self.connection)
                .await
            {
                Ok(total) => {
                    debug!("total kbs found: {:?}", total);
                    total
                }
                Err(e) => {
                    error!("searching kb by key {:?}: {:?}", filter.keyword, e);
                    tracing::event!(tracing::Level::ERROR, "{:?}", e);
                    return Err(Error::DatabaseQueryError);
                }
            };
        // Now let's query the data
        match sqlx::query(
            "SELECT KB_ID, KB_KEY, KIND, TAG_VALUES AS TAGS FROM kbs WHERE KB_KEY LIKE $1 ORDER BY KB_KEY LIMIT $2 OFFSET $3",
        )
        .bind(format!("%{}%", filter.key))
        .bind(i32::from(filter.limit.unwrap_or(5)))
        .bind(i32::from(filter.offset))
        .map(|row: PgRow| KBItem {
            id: KBID(row.get("kb_id")),
            key: row.get("kb_key"),
            kind: row.get("kind"),
            tags: row.get::<String, _>("tags")
                .split(' ')
                .map(|s| s.to_string())
                .map(|s| s.replace('\'', ""))
                .collect::<Vec<String>>()
        })
        .fetch_all(&self.connection)
        .await
        {
            Ok(kbs) => {
                debug!("found some kbs: {:?}", kbs);

                Ok(SearchResult {
                    items: kbs,
                    offset: filter.offset,
                    total: count,
                    limit: filter.limit.unwrap_or(0),
                })
            }
            Err(e) => {
                error!("searching kb by key {:?}: {:?}", filter.keyword, e);
                tracing::event!(tracing::Level::ERROR, "{:?}", e);
                Err(Error::DatabaseQueryError)
            }
        }
    }

    /// get a list of knowledge base entries where their keys contain the given keywords.
    async fn search(&self, filter: KBQueryFilter) -> Result<SearchResult, Error> {
        // let's count first
        let count: i64 =
            match sqlx::query_scalar("SELECT COUNT(*) FROM kbs WHERE TAGS @@ to_tsquery($1)")
                .bind(format!("'{}'", filter.keyword))
                .fetch_one(&self.connection)
                .await
            {
                Ok(total) => {
                    debug!("total kbs found: {:?}", total);
                    total
                }
                Err(e) => {
                    error!("searching kb by key {:?}: {:?}", filter.keyword, e);
                    tracing::event!(tracing::Level::ERROR, "{:?}", e);
                    return Err(Error::DatabaseQueryError);
                }
            };
        // Now let's query the data
        match sqlx::query("SELECT KB_ID, KB_KEY, KIND, TAG_VALUES AS TAGS FROM kbs WHERE TAGS @@ to_tsquery($1) ORDER BY KB_KEY LIMIT $2 OFFSET $3")
            .bind(format!("'{}'", filter.keyword))
            .bind(i32::from(filter.limit.unwrap_or(5)))
            .bind(i32::from(filter.offset))
            .map(|row: PgRow| KBItem {
                id: KBID(row.get("kb_id")),
                key: row.get("kb_key"),
                kind: row.get("kind"),
                tags: row.get::<String, _>("tags")
                    .split(' ')
                    .map(|s| s.to_string())
                    .map(|s| s.replace('\'', ""))
                    .collect::<Vec<String>>()
            })
            .fetch_all(&self.connection)
            .await {
            Ok(kbs) => {
                debug!("found some kbs: {:?}", kbs);

                Ok(SearchResult {
                    items: kbs,
                    offset: filter.offset,
                    total: count,
                    limit: filter.limit.unwrap_or(0),
                })
            }
            Err(e) => {
                error!("searching kbs with keywords: {:?}", e);
                tracing::event!(tracing::Level::ERROR, "{:?}", e);
                Err(Error::DatabaseQueryError)
            }
        }
    }

    /// save given knowledge base in the repository.
    async fn save_kb(&self, kb: KnowledgeBase) -> Result<KBID, Error> {
        debug!("adding kb to postgresql db: {:?}", kb);

        match sqlx::query("INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, REFERENCE, TAGS, TAG_VALUES) VALUES ($1, $2, $3, $4, $5, $6, to_tsvector($7), $8) RETURNING KB_ID")
            .bind(kb.id.to_string())
            .bind(kb.key)
            .bind(kb.value)
            .bind(kb.notes)
            .bind(kb.kind)
            .bind(kb.reference)
            .bind(kb.tags.join(" "))
            .bind(kb.tags.join(" "))
            .map(|row: PgRow| KBID(row.get("kb_id")))
            .fetch_one(&self.connection)
            .await
        {
            Ok(kb_id) => {
                debug!("kb was added to postgres database: {:?}", kb_id);
                Ok(kb_id)
            }
            Err(e) => {
                error!("inserting new kb: {:?}", e);
                tracing::event!(tracing::Level::ERROR, "{:?}", e);
                Err(Error::DatabaseQueryError)
            }
        }
    }
    /// save given category in the repository.
    async fn save_category(&self, category: Category) -> Result<String, Error> {
        debug!("adding new category to postgresql db: {:?}", category);

        match sqlx::query("INSERT INTO CATEGORIES (CATEGORY_NAME, CATEGORY_DESC) VALUES ($1, $2) RETURNING CATEGORY_NAME")
            .bind(category.name)
            .bind(category.description)
            .map(|row: PgRow| row.get("category_name"))
            .fetch_one(&self.connection)
            .await
        {
            Ok(category_name) => {
                debug!("category was added to postgres database: {:?}", category_name);
                Ok(category_name)
            }
            Err(e) => {
                error!("inserting new category: {:?}", e);
                tracing::event!(tracing::Level::ERROR, "{:?}", e);
                Err(Error::DatabaseQueryError)
            }
        }
    }

    /// get a list of categories.
    async fn list_categories(&self, filter: CategoryFilter) -> Result<Vec<Category>, Error> {
        let query = match filter.keyword {
            Some(name) => {
                sqlx::query("SELECT CATEGORY_NAME, CATEGORY_DESC FROM categories WHERE CATEGORY_NAME LIKE $1 ORDER BY CATEGORY_NAME LIMIT $2 OFFSET $3")
                .bind(format!("%{}%", name))
                .bind(i32::from(filter.limit.unwrap_or(5)))
                .bind(i32::from(filter.offset))
            },
            None => {
                sqlx::query("SELECT CATEGORY_NAME, CATEGORY_DESC FROM categories ORDER BY CATEGORY_NAME LIMIT $1 OFFSET $2")
                .bind(i32::from(filter.limit.unwrap_or(5)))
                .bind(i32::from(filter.offset))
            }
        };

        match query
            .map(|row: PgRow| Category {
                name: row.get("category_name"),
                description: row.get("category_desc"),
            })
            .fetch_all(&self.connection)
            .await
        {
            Ok(kbs) => {
                debug!("found some kbs: {:?}", kbs);
                Ok(kbs)
            }
            Err(e) => {
                println!("searching kbs with keywords: {:?}", e);
                error!("searching kbs with keywords: {:?}", e);
                tracing::event!(tracing::Level::ERROR, "{:?}", e);
                Err(Error::DatabaseQueryError)
            }
        }
    }
}
