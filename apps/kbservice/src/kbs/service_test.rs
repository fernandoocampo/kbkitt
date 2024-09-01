use crate::errors::error::Error;
use crate::kbs::service::Service;
use crate::kbs::storage::Storer;
use crate::types::categories::{Category, CategoryFilter};
use crate::types::kbs::{
    KBItem, KBQueryFilter, KnowledgeBase, NewKnowledgeBase, SearchResult, KBID,
};
use tokio::runtime::Runtime;

use async_trait::async_trait;

#[test]
fn test_get_kb_with_id() {
    // Given
    let kb_id = KBID("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8".to_string());
    let want = KnowledgeBase {
        id: KBID("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8".to_string()),
        key: "red".to_string(),
        value: "of the colour of fresh blood".to_string(),
        kind: "concepts".to_string(),
        notes: String::from("to know about color red"),
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
    let got = runtime.block_on(service.get_kb_with_id(kb_id));

    // Then
    match got {
        Ok(kb_got) => assert_eq!(want, kb_got),
        Err(err) => panic!("unexpected error: {:?}", err),
    }
}

#[test]
fn test_get_kb_with_id_error() {
    // Given
    let kb_id = KBID("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8".to_string());
    let want = Error::GetKBError;
    let store = KBStore::new_with_get_kb(None, true);
    let service = Service::new(store);
    let runtime =
        Runtime::new().expect("unable to create runtime to test get kb with id with error");
    // When
    let got = runtime.block_on(service.get_kb_with_id(kb_id));

    // Then
    match got {
        Ok(kb_got) => panic!("unexpected result: {:?}", kb_got),
        Err(err) => assert_eq!(err, want),
    }
}

#[test]
fn test_get_kb_with_key() {
    // Given
    let key = String::from("red");
    let want = KnowledgeBase {
        id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8")),
        key: String::from("red"),
        value: String::from("of the colour of fresh blood"),
        kind: String::from("concepts"),
        notes: String::from("to know about color red"),
        tags: vec![
            String::from("concept"),
            String::from("color"),
            String::from("paint"),
        ],
    };
    let store_kb = KnowledgeBase {
        id: KBID(String::from("dcb8fac0-0756-4c8a-b625-a9a4d1c871c8")),
        key: String::from("red"),
        value: String::from("of the colour of fresh blood"),
        kind: String::from("concepts"),
        notes: String::from("to know about color red"),
        tags: vec![
            String::from("concept"),
            String::from("color"),
            String::from("paint"),
        ],
    };
    let store = KBStore::new_with_get_kb(Some(store_kb), false);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test get kb with key");
    // When
    let got = runtime.block_on(service.get_kb_with_key(key));
    // Then
    match got {
        Ok(kb_got) => assert_eq!(want, kb_got),
        Err(err) => panic!("unexpected error: {:?}", err),
    }
}

#[test]
fn test_search_by_key() {
    // Given
    let query_params = KBQueryFilter {
        key: String::from("red"),
        keyword: Default::default(),
        limit: Some(0),
        offset: 2,
    };
    let wanted_items = vec![
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
    ];
    let want = SearchResult {
        items: wanted_items,
        limit: 0,
        offset: 2,
        total: 10,
    };
    let stored_kbs = vec![
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
    ];
    let result_kbs = SearchResult {
        items: stored_kbs,
        limit: 0,
        offset: 2,
        total: 10,
    };
    let store = KBStore::new_with_search(result_kbs, false);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test get kb with key");
    // When
    let got = runtime.block_on(service.search(query_params));
    // Then
    match got {
        Ok(kb_got) => assert_eq!(want, kb_got),
        Err(err) => panic!("unexpected error: {:?}", err),
    }
}

#[test]
fn test_search_by_key_but_error() {
    // Given
    let want = Error::SearchError;
    let query_params = KBQueryFilter {
        key: "red".to_string(),
        keyword: Default::default(),
        limit: Some(2),
        offset: 0,
    };
    let store = KBStore::new_with_search(SearchResult::default(), true);
    let service = Service::new(store);
    let runtime = Runtime::new()
        .expect("unable to create runtime to test search by key but there is an error");

    // When
    let got = runtime.block_on(service.search(query_params));
    // Then
    match got {
        Ok(kb_got) => panic!("unexpected result: {:?}", kb_got),
        Err(err) => assert_eq!(want, err),
    }
}

#[test]
fn test_search() {
    // Given
    let query_params = KBQueryFilter {
        keyword: String::from("red"),
        key: Default::default(),
        limit: Some(0),
        offset: 2,
    };
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
        limit: 0,
        offset: 2,
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
        limit: 0,
        offset: 2,
        total: 10,
    };
    let store = KBStore::new_with_search(stored_kbs, false);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test get kb with key");
    // When
    let got = runtime.block_on(service.search(query_params));
    // Then
    match got {
        Ok(kb_got) => assert_eq!(want, kb_got),
        Err(err) => panic!("unexpected error: {:?}", err),
    }
}

#[test]
fn test_search_but_error() {
    // Given
    let want = Error::SearchError;
    let query_params = KBQueryFilter {
        keyword: "red".to_string(),
        key: Default::default(),
        limit: Some(2),
        offset: 0,
    };
    let store = KBStore::new_with_search(SearchResult::default(), true);
    let service = Service::new(store);
    let runtime = Runtime::new()
        .expect("unable to create runtime to test search by key but there is an error");

    // When
    let got = runtime.block_on(service.search(query_params));
    // Then
    match got {
        Ok(kb_got) => panic!("unexpected result: {:?}", kb_got),
        Err(err) => assert_eq!(want, err),
    }
}

#[test]
fn test_get_kb_with_key_error() {
    // Given
    let key = "red".to_string();
    let want = Error::GetKBError;
    let store = KBStore::new_with_get_kb(None, true);
    let service = Service::new(store);
    let runtime =
        Runtime::new().expect("unable to create runtime to test get kb with id with error");
    // When
    let got = runtime.block_on(service.get_kb_with_key(key));

    // Then
    match got {
        Ok(kb_got) => panic!("unexpected result: {:?}", kb_got),
        Err(err) => assert_eq!(err, want),
    }
}

#[test]
fn test_add_kb() {
    // Given
    let new_kb = NewKnowledgeBase {
        key: String::from("red"),
        value: String::from("of the colour of fresh blood"),
        kind: String::from("concepts"),
        notes: String::from("to know about color red"),
        tags: vec![
            String::from("concept"),
            String::from("color"),
            String::from("paint"),
        ],
    };
    let mut want = KnowledgeBase {
        id: KBID(String::from("")),
        key: String::from("red"),
        value: String::from("of the colour of fresh blood"),
        kind: String::from("concepts"),
        notes: String::from("to know about color red"),
        tags: vec![
            String::from("concept"),
            String::from("color"),
            String::from("paint"),
        ],
    };
    let non_existing_kb = Some(KnowledgeBase::default());
    let store = KBStore::new_with_add_kb(non_existing_kb, false);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test add kb");
    // When
    let got = runtime.block_on(service.add_kb(new_kb));
    // Then
    match got {
        Ok(kb_got) => {
            want.id = kb_got.id.clone();
            assert_eq!(want, kb_got.clone());
        }
        Err(err) => panic!("unexpected error: {:?}", err),
    }
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
    let want = true;
    let store = KBStore::new_with_add_category(false);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test add category");
    // When
    let got = runtime.block_on(service.add_category(new_category));
    // Then
    match got {
        Ok(kb_got) => {
            assert_eq!(want, kb_got);
        }
        Err(err) => panic!("unexpected error: {:?}", err),
    }
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
    let got = runtime.block_on(service.add_category(new_category));
    // Then
    match got {
        Ok(cat_got) => panic!("unexpected result: {:?}", cat_got),
        Err(err) => assert_eq!(err, want),
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
        tags: vec![
            String::from("concept"),
            String::from("color"),
            String::from("paint"),
        ],
    };
    let non_existing_kb = Some(KnowledgeBase::default());
    let want = Error::CreateKBError;
    let store = KBStore::new_with_add_kb(non_existing_kb, true);
    let service = Service::new(store);
    let runtime = Runtime::new().expect("unable to create runtime to test add kb with error");
    // When
    let got = runtime.block_on(service.add_kb(new_kb));
    // Then
    match got {
        Ok(kb_got) => panic!("unexpected result: {:?}", kb_got),
        Err(err) => assert_eq!(err, want),
    }
}

#[derive(Debug, Clone)]
struct KBStore {
    get_kb_value: Option<KnowledgeBase>,
    search_value: SearchResult,
    search_error: Option<bool>,
    get_kb_error: Option<bool>,
    save_kb_error: Option<bool>,
    save_category_error: Option<bool>,
}

impl Default for KBStore {
    fn default() -> Self {
        KBStore {
            get_kb_error: Default::default(),
            save_kb_error: Default::default(),
            search_error: Default::default(),
            save_category_error: Default::default(),
            get_kb_value: Default::default(),
            search_value: Default::default(),
        }
    }
}

impl KBStore {
    fn new_with_get_kb(kb: Option<KnowledgeBase>, is_error: bool) -> Self {
        let mut dummy_store = KBStore::default();

        dummy_store.get_kb_value = kb;
        dummy_store.get_kb_error = Some(is_error);

        dummy_store
    }
    fn new_with_search(kbs: SearchResult, is_error: bool) -> Self {
        let mut dummy_store = KBStore::default();

        dummy_store.search_value = kbs;
        dummy_store.search_error = Some(is_error);

        dummy_store
    }
    fn new_with_add_kb(kb: Option<KnowledgeBase>, is_error: bool) -> Self {
        let mut dummy_store = KBStore::default();

        dummy_store.get_kb_value = kb;
        dummy_store.get_kb_error = Some(false);
        dummy_store.save_kb_error = Some(is_error);

        dummy_store
    }
    fn new_with_add_category(is_error: bool) -> Self {
        let mut dummy_store = KBStore::default();
        dummy_store.save_category_error = Some(is_error);

        dummy_store
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
