use crate::errors::error::Error;
use crate::kbs::storage::Storer;
use crate::types::categories::Category;
use crate::types::kbs::{KBQueryFilter, KnowledgeBase, NewKnowledgeBase, SearchResult, KBID};
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
        debug!("checking if kb already exists: {:?}", new_kb.key);

        match self.get_kb_with_key(new_kb.clone().key).await {
            Ok(kb) => {
                if kb.id != KBID("".to_string()) {
                    return Err(Error::DuplicateKBError);
                }
            }
            Err(e) => {
                error!("reading if kb already exists: {:?}", e);
                return Err(Error::CreateKBError);
            }
        }

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

    pub async fn update_kb(&self, kb: KnowledgeBase) -> Result<(), Error> {
        debug!("checking if kb already exists: {:?}", kb.key);

        match self.get_kb_with_key(kb.clone().key).await {
            Ok(existing_kb) => {
                if existing_kb.id != KBID("".to_string()) && existing_kb.id != kb.id {
                    return Err(Error::DuplicateKBError);
                }
            }
            Err(e) => {
                error!("reading if kb already exists: {:?}", e);
                return Err(Error::UpdateKBError);
            }
        }

        debug!("start updating kb: {:?}", kb);

        match self.store.update_kb(kb).await {
            Ok(result) => {
                if !result {
                    return Err(Error::KBWasNotUpdatedError);
                }

                Ok(())
            }
            Err(e) => {
                error!("updating kb: {:?}", e);
                Err(Error::UpdateKBError)
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

    pub async fn search(&self, query_params: KBQueryFilter) -> Result<SearchResult, Error> {
        let mut result: SearchResult = SearchResult::default();

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

    async fn search_by_key(&self, query_params: KBQueryFilter) -> Result<SearchResult, Error> {
        debug!("start searching by key: {:?}", query_params);

        match self.store.search_by_key(query_params).await {
            Ok(kbs) => Ok(kbs),
            Err(e) => {
                error!("searching by key: {:?}", e);

                Err(Error::SearchError)
            }
        }
    }

    async fn search_by_keyword(&self, query_params: KBQueryFilter) -> Result<SearchResult, Error> {
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
