use std::fmt::{Display, Formatter, Result};
use std::num::ParseIntError;

use warp::reject::Reject;

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
    MissingParameters,
    UpdateKBError,
    KBWasNotUpdatedError,
    ParseError(ParseIntError),
}

impl Display for Error {
    fn fmt(&self, f: &mut Formatter) -> Result {
        match *self {
            Error::KBNotFound => write!(f, "KB not found"),
            Error::GetKBError => write!(f, "Unable to get this KB"),
            Error::CreateKBError => write!(f, "Unable to create KB"),
            Error::UpdateKBError => write!(f, "Unable to update KB"),
            Error::DuplicateKBError => write!(f, "KB already exists"),
            Error::CreateCategoryError => write!(f, "Unable to create Category"),
            Error::DuplicateCategory => write!(f, "Category already exists"),
            Error::SearchError => write!(f, "Unable to search knowledge base"),
            Error::DatabaseQueryError => write!(f, "Unable to query repository"),
            Error::MissingParameters => write!(f, "Missing parameter"),
            Error::KBWasNotUpdatedError => write!(f, "KB was not updated"),
            Error::ParseError(ref err) => write!(f, "Cannot parse parameter: {err}"),
        }
    }
}

impl Reject for Error {}
