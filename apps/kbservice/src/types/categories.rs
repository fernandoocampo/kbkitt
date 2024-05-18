use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub struct Category {
    pub name: String,
    pub description: String,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub struct CategoryFilter {
    pub keyword: Option<String>,
    /// determines the number of rows.
    pub limit: Option<u16>,
    /// skips the offset rows before beginning to return the rows.
    pub offset: u16,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub struct SaveCategorySuccess {
    pub ok: bool,
}

impl Category {
    pub fn new(name: String, description: String) -> Self {
        Category { name, description }
    }
}

impl Default for CategoryFilter {
    fn default() -> Self {
        CategoryFilter {
            keyword: None,
            limit: Some(5),
            offset: 0,
        }
    }
}

impl SaveCategorySuccess {
    pub fn new(is_ok: bool) -> Self {
        SaveCategorySuccess { ok: is_ok }
    }
}
