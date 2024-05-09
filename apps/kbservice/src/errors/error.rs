use std::fmt::{Display, Formatter, Result};

#[derive(Debug, PartialEq)]
pub enum Error {
    KBNotFound,
    GetKBError,
    CreateKBError,
    CreateCategoryError,
    SearchError,
    DuplicateKBError,
    DuplicateCategory,
    DatabaseQueryError,
}

impl Display for Error {
    fn fmt(&self, f: &mut Formatter) -> Result {
        match *self {
            Error::KBNotFound => write!(f, "KB not found"),
            Error::GetKBError => write!(f, "Unable to get this KB"),
            Error::CreateKBError => write!(f, "Unable to create KB"),
            Error::DuplicateKBError => write!(f, "KB already exists"),
            Error::CreateCategoryError => write!(f, "Unable to create Category"),
            Error::DuplicateCategory => write!(f, "Category already exists"),
            Error::SearchError => write!(f, "Unable to search knowledge base"),
            Error::DatabaseQueryError => write!(f, "Unable to query repository"),
        }
    }
}
