use core::fmt;

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct KBID(pub String);

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct KnowledgeBase {
    pub id: KBID,
    pub key: String,
    pub value: String,
    pub notes: String,
    pub kind: String,
    pub tags: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct NewKnowledgeBase {
    pub key: String,
    pub value: String,
    pub notes: String,
    pub kind: String,
    pub tags: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct KBItem {
    pub id: KBID,
    pub key: String,
    pub kind: String,
    pub tags: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct KBQueryFilter {
    pub keyword: String,
    /// determines the number of rows.
    pub limit: Option<i32>,
    /// skips the offset rows before beginning to return the rows.
    pub offset: i32,
}

impl Default for KnowledgeBase {
    fn default() -> Self {
        KnowledgeBase {
            id: KBID("".to_string()),
            key: Default::default(),
            value: Default::default(),
            notes: Default::default(),
            kind: Default::default(),
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
            tags: self.tags.clone(),
        }
    }
}

impl fmt::Display for KBID {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}
