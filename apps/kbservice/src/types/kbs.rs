use core::fmt;
use std::collections::HashMap;

use serde::{Deserialize, Serialize};

const LIMIT_KEY: &str = "limit";
const OFFSET_KEY: &str = "offset";
const KEY_KEY: &str = "key";
const KEYWORD_KEY: &str = "keyword";
const DEFAULT_LIMIT: u16 = 5;
const DEFAULT_OFFSET: u16 = 0;

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct KBID(pub String);

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct KnowledgeBase {
    pub id: KBID,
    pub key: String,
    pub value: String,
    pub notes: String,
    pub kind: String,
    pub reference: Option<String>,
    pub tags: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct NewKnowledgeBase {
    pub key: String,
    pub value: String,
    pub notes: String,
    pub kind: String,
    pub reference: Option<String>,
    pub tags: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord, Default)]
pub struct SearchResult {
    pub items: Vec<KBItem>,
    pub total: i64,
    pub limit: u16,
    pub offset: u16,
}

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct KBItem {
    pub id: KBID,
    pub key: String,
    pub kind: String,
    pub tags: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord, Default)]
pub struct KBQueryFilter {
    /// value used to search kbs by tags.
    pub keyword: String,
    /// value used to search kbs by keys.
    pub key: String,
    /// determines the number of rows.
    pub limit: Option<u16>,
    /// skips the offset rows before beginning to return the rows.
    pub offset: u16,
}

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct SaveKBSuccess {
    pub id: String,
}

impl Default for KnowledgeBase {
    fn default() -> Self {
        KnowledgeBase {
            id: KBID("".to_string()),
            key: Default::default(),
            value: Default::default(),
            notes: Default::default(),
            kind: Default::default(),
            reference: None,
            tags: Default::default(),
        }
    }
}

/// a knowledge base entry in your knowledge base repository.
/// this could be anything from a concept, an image to a block of code.
impl KnowledgeBase {
    pub fn new(kb_key: String) -> Self {
        KnowledgeBase {
            id: KBID(uuid::Uuid::new_v4().to_string()),
            key: kb_key,
            value: Default::default(),
            notes: Default::default(),
            kind: Default::default(),
            reference: None,
            tags: Default::default(),
        }
    }
}

/// A new knowledge base entry is represented here.
impl NewKnowledgeBase {
    pub fn to_knowledge_base(&self) -> KnowledgeBase {
        KnowledgeBase {
            id: KBID(uuid::Uuid::new_v4().to_string()),
            key: self.key.clone(),
            value: self.value.clone(),
            notes: self.notes.clone(),
            kind: self.kind.clone(),
            reference: self.reference.clone(),
            tags: self.tags.clone(),
        }
    }
}

impl fmt::Display for KBID {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

pub fn extract_filter_params(params: HashMap<String, String>) -> KBQueryFilter {
    let mut filter = KBQueryFilter::default();
    if params.contains_key(LIMIT_KEY) {
        match params.get(LIMIT_KEY) {
            Some(item) => match item.parse::<u16>() {
                Ok(value) => filter.limit = Some(value),
                Err(_) => filter.limit = Some(DEFAULT_LIMIT),
            },
            _ => filter.limit = Some(DEFAULT_LIMIT),
        };
    }
    if params.contains_key(OFFSET_KEY) {
        match params.get(OFFSET_KEY) {
            Some(item) => match item.parse::<u16>() {
                Ok(value) => filter.offset = value,
                Err(_) => filter.offset = DEFAULT_OFFSET,
            },
            _ => filter.offset = DEFAULT_OFFSET,
        };
    }

    if params.contains_key(KEY_KEY) {
        match params.get(KEY_KEY) {
            Some(value) => filter.key = value.into(),
            _ => filter.key = "".to_string(),
        }
    }

    if params.contains_key(KEYWORD_KEY) {
        match params.get(KEYWORD_KEY) {
            Some(value) => filter.keyword = value.into(),
            _ => filter.keyword = "".to_string(),
        }
    }

    filter
}

impl SaveKBSuccess {
    pub fn new(kb_id: KBID) -> Self {
        SaveKBSuccess {
            id: kb_id.to_string(),
        }
    }
}
