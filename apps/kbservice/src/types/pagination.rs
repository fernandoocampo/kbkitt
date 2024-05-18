use crate::errors::error;
use std::collections::HashMap;

#[derive(Default, Debug, PartialEq)]
pub struct Pagination {
    /// The index of the last item which has to be returned
    pub limit: Option<i32>,
    /// The index of the first item which has to be returned
    pub offset: i32,
}

pub fn extract_pagination(params: HashMap<String, String>) -> Result<Pagination, error::Error> {
    if params.contains_key("limit") && params.contains_key("offset") {
        return Ok(Pagination {
            limit: Some(
                params
                    .get("limit")
                    .unwrap()
                    .parse::<i32>()
                    .map_err(error::Error::ParseError)?,
            ),
            offset: params
                .get("offset")
                .unwrap()
                .parse::<i32>()
                .map_err(error::Error::ParseError)?,
        });
    }
    Err(error::Error::MissingParameters)
}
