use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, PartialOrd, Ord)]
pub struct Category {
    pub name: String,
    pub description: String,
}

impl Category {
    pub fn new(name: String, description: String) -> Self {
        Category { name, description }
    }
}
