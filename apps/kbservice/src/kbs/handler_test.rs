use crate::errors::error;
use crate::errors::error::Error;
use crate::kbs::handler;
use crate::kbs::service::Service;
use crate::kbs::storage::Storer;
use crate::types::categories::{Category, CategoryFilter, SaveCategorySuccess};
use crate::types::kbs::{
    KBItem, KBQueryFilter, KnowledgeBase, NewKnowledgeBase, SaveKBSuccess, SearchResult, KBID,
};
use async_trait::async_trait;
use hyper::StatusCode;
use std::collections::HashMap;
use tokio::runtime::Runtime;
use warp::Reply;

#[test]
fn test_get_kb_with_id() {
    // Given
    let kb_id = "dcb8fac0-0756-4c8a-b625-a9a4d1c871c8".to_string();

    let want = KnowledgeBase {
        id: KBID("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8".to_string()),
        key: "red".to_string(),
        value: "of the colour of fresh blood".to_string(),
        kind: "concepts".to_string(),
        notes: String::from("to know about color red"),
        reference: Some(String::from("Some Author")),
        tags: vec![
            "concept".to_string(),
            "color".to_string(),
            "paint".to_string(),
        ],
    };
    let stored_kb = KnowledgeBase {
        id: KBID("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8".to_string()),
        key: "red".to_string(),
        value: "of the colour of fresh blood".to_string(),
        kind: "concepts".to_string(),
        notes: String::from("to know about color red"),
        reference: Some(String::from("Some Author")),
        tags: vec![
            "concept".to_string(),
            "color".to_string(),
            "paint".to_string(),
        ],
    };
    let store = KBStore::new_with_get_kb(Some(stored_kb), false);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test get kb with id");
    // When
    let response = runtime.block_on(handler::get_kb_with_id(kb_id, service));

    // Then
    let response_body = response.unwrap().into_response().into_body();
    let body_bytes = runtime
        .block_on(hyper::body::to_bytes(response_body))
        .unwrap();

    let got: KnowledgeBase = serde_json::from_slice(&body_bytes).unwrap();

    assert_eq!(want, got)
}

#[test]
fn test_get_kb_with_id_error() {
    // Given
    let kb_id = "dcb8fac0-0756-4c8a-b625-a9a4d1c871c8".to_string();
    let want = Error::GetKBError;
    let store = KBStore::new_with_get_kb(None, true);
    let service = Service::new(store);
    let runtime =
        Runtime::new().expect("unable to create runtime to test get kb with id with error");
    // When
    let got = runtime.block_on(handler::get_kb_with_id(kb_id, service));

    // Then
    assert_eq!(true, got.is_err());
    let got_error = match got {
        Ok(value) => panic!("unexpected result {:?}", value.into_response()),
        Err(err) => err,
    };

    if let Some(e) = got_error.find::<error::Error>() {
        assert_eq!(want, *e);
        return;
    }
}

#[test]
fn test_search_by_key() {
    // Given
    let mut params: HashMap<String, String> = HashMap::new();
    params.insert(String::from("offset"), String::from("0"));
    params.insert(String::from("limit"), String::from("2"));
    params.insert(String::from("key"), String::from("red"));

    let want = SearchResult {
        items: vec![
            KBItem {
                id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8")),
                key: String::from("red"),
                kind: String::from("concepts"),
                tags: vec![
                    String::from("concept"),
                    String::from("color"),
                    String::from("paint"),
                ],
            },
            KBItem {
                id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c9")),
                key: String::from("redemption"),
                kind: String::from("concepts"),
                tags: vec![
                    String::from("concept"),
                    String::from("saving"),
                    String::from("absolution"),
                ],
            },
        ],
        limit: 2,
        offset: 0,
        total: 10,
    };
    let stored_kbs = SearchResult {
        items: vec![
            KBItem {
                id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8")),
                key: String::from("red"),
                kind: String::from("concepts"),
                tags: vec![
                    String::from("concept"),
                    String::from("color"),
                    String::from("paint"),
                ],
            },
            KBItem {
                id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c9")),
                key: String::from("redemption"),
                kind: String::from("concepts"),
                tags: vec![
                    String::from("concept"),
                    String::from("saving"),
                    String::from("absolution"),
                ],
            },
        ],
        limit: 2,
        offset: 0,
        total: 10,
    };
    let store = KBStore::new_with_search(stored_kbs, false);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test get kb with key");
    // When
    let response = runtime.block_on(handler::search(params, service));
    // Then
    let response_body = response.unwrap().into_response().into_body();
    let body_bytes = runtime
        .block_on(hyper::body::to_bytes(response_body))
        .unwrap();

    let got: SearchResult = serde_json::from_slice(&body_bytes).unwrap();

    assert_eq!(want, got)
}

#[test]
fn test_search_by_key_but_error() {
    // Given
    let mut params: HashMap<String, String> = HashMap::new();
    params.insert(String::from("offset"), String::from("0"));
    params.insert(String::from("limit"), String::from("2"));
    params.insert(String::from("key"), String::from("red"));

    let want = Error::SearchError;
    let store = KBStore::new_with_search(SearchResult::default(), true);
    let service = Service::new(store);
    let runtime = Runtime::new()
        .expect("unable to create runtime to test search by key but there is an error");

    // When
    let response = runtime.block_on(handler::search(params, service));
    // Then
    assert_eq!(true, response.is_err());
    let got_error = match response {
        Ok(value) => panic!("unexpected result {:?}", value.into_response()),
        Err(err) => err,
    };

    if let Some(e) = got_error.find::<error::Error>() {
        assert_eq!(want, *e);
        return;
    }
}

#[test]
fn test_search() {
    // Given
    let mut params: HashMap<String, String> = HashMap::new();
    params.insert(String::from("offset"), String::from("0"));
    params.insert(String::from("limit"), String::from("2"));
    params.insert(String::from("keyword"), String::from("red"));

    let want = SearchResult {
        items: vec![
            KBItem {
                id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8")),
                key: String::from("red"),
                kind: String::from("concepts"),
                tags: vec![
                    String::from("concept"),
                    String::from("color"),
                    String::from("paint"),
                ],
            },
            KBItem {
                id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c9")),
                key: String::from("redemption"),
                kind: String::from("concepts"),
                tags: vec![
                    String::from("concept"),
                    String::from("saving"),
                    String::from("absolution"),
                ],
            },
        ],
        limit: 2,
        offset: 0,
        total: 10,
    };
    let stored_kbs = SearchResult {
        items: vec![
            KBItem {
                id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8")),
                key: String::from("red"),
                kind: String::from("concepts"),
                tags: vec![
                    String::from("concept"),
                    String::from("color"),
                    String::from("paint"),
                ],
            },
            KBItem {
                id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c9")),
                key: String::from("redemption"),
                kind: String::from("concepts"),
                tags: vec![
                    String::from("concept"),
                    String::from("saving"),
                    String::from("absolution"),
                ],
            },
        ],
        limit: 2,
        offset: 0,
        total: 10,
    };
    let store = KBStore::new_with_search(stored_kbs, false);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test get kb with key");
    // When
    let response = runtime.block_on(handler::search(params, service));
    // Then
    let response_body = response.unwrap().into_response().into_body();
    let body_bytes = runtime
        .block_on(hyper::body::to_bytes(response_body))
        .unwrap();

    let got: SearchResult = serde_json::from_slice(&body_bytes).unwrap();

    assert_eq!(want, got)
}

#[test]
fn test_add_kb() {
    // Given
    let new_kb = NewKnowledgeBase {
        key: String::from("red"),
        value: String::from("of the colour of fresh blood"),
        kind: String::from("concepts"),
        notes: String::from("to know about color red"),
        reference: Some(String::from("Some Author")),
        tags: vec![
            String::from("concept"),
            String::from("color"),
            String::from("paint"),
        ],
    };
    let mut want = SaveKBSuccess {
        id: String::from(""),
    };

    let non_existing_kb = Some(KnowledgeBase::default());
    let store = KBStore::new_with_add_kb(non_existing_kb, false);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test add kb");
    // When
    let response = runtime.block_on(handler::add_kb(new_kb, service));
    // Then
    let response_body = response.unwrap().into_response().into_body();
    let body_bytes = runtime
        .block_on(hyper::body::to_bytes(response_body))
        .unwrap();

    let got: SaveKBSuccess = serde_json::from_slice(&body_bytes).unwrap();
    want.id = got.clone().id;

    assert_eq!(want, got)
}

#[test]
fn test_update_kb() {
    // Given
    let updated_kb = KnowledgeBase {
        id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c9")),
        key: String::from("btc"),
        kind: String::from("crypto"),
        value: String::from("new currency state"),
        notes: String::from("it is here to stay"),
        reference: Some(String::from("Satoshi Nakamoto")),
        tags: vec![
            String::from("crypto"),
            String::from("btc"),
            String::from("satoshi"),
        ],
    };

    let want = StatusCode::OK;
    let non_existing_kb = Some(KnowledgeBase::default());

    let store = KBStore::new_with_update_kb(false, non_existing_kb, true);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test add kb");
    // When
    let response = runtime.block_on(handler::update_kb(updated_kb, service));
    // Then
    let response = response.unwrap().into_response();

    assert_eq!(want, response.status())
}

#[test]
fn test_add_category() {
    // Given
    let new_category = Category {
        name: String::from("concept"),
        description: String::from(
            "an abstract or general idea that is derived from specific instances or occurrences",
        ),
    };
    let want = SaveCategorySuccess { ok: true };
    let store = KBStore::new_with_add_category(false);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test add category");
    // When
    let response = runtime.block_on(handler::add_category(new_category, service));
    // Then
    let response_body = response.unwrap().into_response().into_body();
    let body_bytes = runtime
        .block_on(hyper::body::to_bytes(response_body))
        .unwrap();

    let got: SaveCategorySuccess = serde_json::from_slice(&body_bytes).unwrap();

    assert_eq!(want, got)
}

#[test]
fn test_add_category_with_error() {
    // Given
    let new_category = Category {
        name: String::from("concept"),
        description: String::from(
            "an abstract or general idea that is derived from specific instances or occurrences",
        ),
    };
    let want = Error::CreateCategoryError;
    let store = KBStore::new_with_add_category(true);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test add category with error");
    // When
    let response = runtime.block_on(handler::add_category(new_category, service));
    // Then
    assert_eq!(true, response.is_err());
    let got_error = match response {
        Ok(value) => panic!("unexpected result {:?}", value.into_response()),
        Err(err) => err,
    };

    if let Some(e) = got_error.find::<error::Error>() {
        assert_eq!(want, *e);
        return;
    }
}

#[test]
fn test_add_kb_with_error() {
    // Given
    let new_kb = NewKnowledgeBase {
        key: String::from("red"),
        value: String::from("of the colour of fresh blood"),
        kind: String::from("concepts"),
        notes: String::from("to know about color red"),
        reference: Some(String::from("Some Author")),
        tags: vec![
            String::from("concept"),
            String::from("color"),
            String::from("paint"),
        ],
    };
    let want = Error::CreateKBError;

    let non_existing_kb = Some(KnowledgeBase::default());
    let store = KBStore::new_with_add_kb(non_existing_kb, true);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test add kb with error");
    // When
    let response = runtime.block_on(handler::add_kb(new_kb, service));
    // Then
    assert_eq!(true, response.is_err());
    let got_error = match response {
        Ok(value) => panic!("unexpected result {:?}", value.into_response()),
        Err(err) => err,
    };

    if let Some(e) = got_error.find::<error::Error>() {
        assert_eq!(want, *e);
        return;
    }
}

#[derive(Debug, Clone, Default)]
struct KBStore {
    get_kb_value: Option<KnowledgeBase>,
    search_value: SearchResult,
    search_error: Option<bool>,
    get_kb_error: Option<bool>,
    save_kb_error: Option<bool>,
    save_category_error: Option<bool>,
    update_kb_error: Option<bool>,
    update_kb_result: bool,
}

impl KBStore {
    fn new_with_get_kb(kb: Option<KnowledgeBase>, is_error: bool) -> Self {
        KBStore {
            get_kb_value: kb,
            get_kb_error: Some(is_error),
            ..Default::default()
        }
    }
    fn new_with_search(kbs: SearchResult, is_error: bool) -> Self {
        KBStore {
            search_value: kbs,
            search_error: Some(is_error),
            ..Default::default()
        }
    }
    fn new_with_add_kb(kb: Option<KnowledgeBase>, is_error: bool) -> Self {
        KBStore {
            get_kb_value: kb,
            get_kb_error: Some(false),
            save_kb_error: Some(is_error),
            ..Default::default()
        }
    }
    fn new_with_update_kb(is_error: bool, kb: Option<KnowledgeBase>, result: bool) -> Self {
        KBStore {
            get_kb_value: kb,
            get_kb_error: Some(false),
            update_kb_error: Some(is_error),
            update_kb_result: result,
            ..Default::default()
        }
    }
    fn new_with_add_category(is_error: bool) -> Self {
        KBStore {
            save_category_error: Some(is_error),
            ..Default::default()
        }
    }
}

#[async_trait]
impl Storer for KBStore {
    async fn get_kb_by_id(&self, _: KBID) -> Result<KnowledgeBase, Error> {
        match &self.get_kb_error.unwrap() {
            false => Ok(self.get_kb_value.clone().unwrap()),
            true => Err(Error::GetKBError),
        }
    }

    async fn get_kb_by_key(&self, _: String) -> Result<KnowledgeBase, Error> {
        match &self.get_kb_error.unwrap() {
            false => Ok(self.get_kb_value.clone().unwrap()),
            true => Err(Error::GetKBError),
        }
    }

    async fn search_by_key(&self, _: KBQueryFilter) -> Result<SearchResult, Error> {
        match &self.search_error.unwrap() {
            false => Ok(self.search_value.clone()),
            true => Err(Error::SearchError),
        }
    }

    async fn search(&self, _: KBQueryFilter) -> Result<SearchResult, Error> {
        match &self.search_error.unwrap() {
            false => Ok(self.search_value.clone()),
            true => Err(Error::SearchError),
        }
    }

    async fn save_kb(&self, new_kb: KnowledgeBase) -> Result<KBID, Error> {
        match &self.save_kb_error.unwrap() {
            false => Ok(new_kb.id.clone()),
            true => Err(Error::CreateKBError),
        }
    }

    async fn update_kb(&self, _: KnowledgeBase) -> Result<bool, Error> {
        match &self.update_kb_error.unwrap() {
            false => Ok(self.update_kb_result),
            true => Err(Error::UpdateKBError),
        }
    }

    async fn save_category(&self, category: Category) -> Result<String, Error> {
        match &self.save_category_error.unwrap() {
            false => Ok(category.name.clone()),
            true => Err(Error::CreateCategoryError),
        }
    }

    async fn list_categories(&self, _: CategoryFilter) -> Result<Vec<Category>, Error> {
        Ok(vec![])
    }
}
