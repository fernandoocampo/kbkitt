use crate::errors::error::Error;
use crate::types::categories::Category;
use crate::types::kbs::{KBItem, KBQueryFilter, KnowledgeBase, KBID};
use async_trait::async_trait;
use std::fmt::{Debug, Error as FmtError, Formatter};

#[async_trait]
pub trait Storer {
    /// get a Knowledge base with the given id.
    async fn get_kb_by_id(&self, id: KBID) -> Result<KnowledgeBase, Error>;
    /// get a Knowledge base with the given key.
    async fn get_kb_by_key(&self, key: String) -> Result<KnowledgeBase, Error>;
    /// get a list of knowledge base entries where their keys contain the given keywords.
    async fn search_by_key(&self, filter: KBQueryFilter) -> Result<Vec<KBItem>, Error>;
    /// get a list of knowledge base entries where their keys contain the given keywords.
    async fn search(&self, filter: KBQueryFilter) -> Result<Vec<KBItem>, Error>;
    /// save given knowledge base in the repository.
    async fn save_kb(&self, kb: KnowledgeBase) -> Result<KBID, Error>;
    /// save given category in the repository.
    async fn save_category(&self, category: Category) -> Result<String, Error>;
}

impl Debug for dyn Storer {
    fn fmt(&self, f: &mut Formatter<'_>) -> Result<(), FmtError> {
        f.debug_struct("Storer").finish()
    }
}
