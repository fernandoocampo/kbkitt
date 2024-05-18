use crate::errors::error::Error;
use crate::kbs::storage::Storer;
use crate::types::categories::Category;
use crate::types::kbs::{KBItem, KBQueryFilter, KnowledgeBase, NewKnowledgeBase, KBID};
use log::error;
use tracing::debug;

#[derive(Debug, Clone)]
pub struct Service<T: Storer> {
    store: T,
}

impl<T: Storer> Service<T> {
    pub fn new(a_store: T) -> Self {
        Service { store: a_store }
    }

    pub async fn get_kb_with_id(&self, kb_id: KBID) -> Result<KnowledgeBase, Error> {
        debug!("start getting kb {}", kb_id);

        match self.store.get_kb_by_id(kb_id).await {
            Ok(kb) => Ok(kb),
            Err(e) => {
                error!("getting kb from repository with id: {:?}", e);
                Err(Error::GetKBError)
            }
        }
    }

    pub async fn get_kb_with_key(&self, key: String) -> Result<KnowledgeBase, Error> {
        debug!("start getting kb '{}'", key);

        match self.store.get_kb_by_key(key).await {
            Ok(kb) => Ok(kb),
            Err(e) => {
                error!("getting kb from repository with key: {:?}", e);
                Err(Error::GetKBError)
            }
        }
    }

    pub async fn add_kb(&self, new_kb: NewKnowledgeBase) -> Result<KnowledgeBase, Error> {
        debug!("start adding kb: {:?}", new_kb);

        let kb_to_save = new_kb.to_knowledge_base();

        match self.store.save_kb(kb_to_save.clone()).await {
            Ok(_) => Ok(kb_to_save),
            Err(e) => {
                error!("saving kb: {:?}", e);
                Err(Error::CreateKBError)
            }
        }
    }

    pub async fn add_category(&self, new_category: Category) -> Result<bool, Error> {
        debug!("start adding category: {:?}", new_category);

        match self.store.save_category(new_category).await {
            Ok(_) => Ok(true),
            Err(e) => {
                error!("saving category: {:?}", e);
                Err(Error::CreateCategoryError)
            }
        }
    }

    pub async fn search(&self, query_params: KBQueryFilter) -> Result<Vec<KBItem>, Error> {
        let mut result: Vec<KBItem> = vec![];

        if query_params.key.is_empty() && query_params.keyword.is_empty() {
            return Ok(result);
        }
        if !query_params.key.is_empty() {
            result = match self.search_by_key(query_params).await {
                Ok(res) => res,
                Err(e) => return Err(e),
            };
        } else if !query_params.keyword.is_empty() {
            result = match self.search_by_keyword(query_params).await {
                Ok(res) => res,
                Err(e) => return Err(e),
            };
        }

        Ok(result)
    }

    async fn search_by_key(&self, query_params: KBQueryFilter) -> Result<Vec<KBItem>, Error> {
        debug!("start searching by key: {:?}", query_params);

        match self.store.search_by_key(query_params).await {
            Ok(kbs) => Ok(kbs),
            Err(e) => {
                error!("searching by key: {:?}", e);

                Err(Error::SearchError)
            }
        }
    }

    async fn search_by_keyword(&self, query_params: KBQueryFilter) -> Result<Vec<KBItem>, Error> {
        debug!("start search: {:?}", query_params);

        match self.store.search(query_params).await {
            Ok(kbs) => Ok(kbs),
            Err(e) => {
                error!("searching by key: {:?}", e);

                Err(Error::SearchError)
            }
        }
    }
}
