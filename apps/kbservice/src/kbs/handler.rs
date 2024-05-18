use std::collections::HashMap;

use crate::errors::error::Error;
use crate::kbs::service::Service;
use crate::types::categories::{Category, SaveCategorySuccess};
use crate::types::kbs::{NewKnowledgeBase, SaveKBSuccess, KBID};

use tracing::{debug, error};
use warp::{http::StatusCode, Rejection, Reply};

use super::storage;

pub async fn get_kb_with_id(
    id: String,
    service: Service<impl storage::Storer>,
) -> Result<impl Reply, Rejection> {
    let result = match service.get_kb_with_id(KBID(id)).await {
        Ok(kb) => kb,
        Err(e) => return Err(warp::reject::custom(e)),
    };

    Ok(warp::reply::json(&result))
}

pub async fn get_kb_with_key(
    key: String,
    service: Service<impl storage::Storer>,
) -> Result<impl Reply, Rejection> {
    let result = match service.get_kb_with_key(key).await {
        Ok(kb) => kb,
        Err(e) => return Err(warp::reject::custom(e)),
    };

    Ok(warp::reply::json(&result))
}

pub async fn search(
    params: HashMap<String, String>,
    service: Service<impl storage::Storer>,
) -> Result<impl Reply, Rejection> {
    debug!("start searching");

    let filter = crate::types::kbs::extract_filter_params(params);

    let result = match service.search(filter).await {
        Ok(res) => res,
        Err(e) => return Err(warp::reject::custom(e)),
    };

    Ok(warp::reply::json(&result))
}

pub async fn add_kb(
    new_kb: NewKnowledgeBase,
    service: Service<impl storage::Storer>,
) -> Result<impl Reply, Rejection> {
    match service.add_kb(new_kb.clone()).await {
        Ok(kb) => {
            debug!("new kb was saved {:?}", kb);

            let result = SaveKBSuccess::new(kb.id);

            Ok(warp::reply::json(&result))
        }
        Err(e) => {
            error!("adding kb {:?}", new_kb);
            Err(warp::reject::custom(e))
        }
    }
}

pub async fn add_category(
    new_category: Category,
    service: Service<impl storage::Storer>,
) -> Result<impl Reply, Rejection> {
    match service.add_category(new_category.clone()).await {
        Ok(result) => {
            debug!("new category was saved: {:?}", result);

            let result = SaveCategorySuccess::new(result);

            Ok(warp::reply::json(&result))
        }
        Err(e) => {
            error!("adding category {:?}", new_category);
            Err(warp::reject::custom(e))
        }
    }
}

pub async fn return_error(r: Rejection) -> Result<impl Reply, Rejection> {
    if let Some(Error::KBNotFound) = r.find() {
        Ok(warp::reply::with_status(
            "Route not found".to_string(),
            StatusCode::NOT_FOUND,
        ))
    } else if let Some(Error::GetKBError) = r.find() {
        Ok(warp::reply::with_status(
            "Route not found".to_string(),
            StatusCode::NOT_FOUND,
        ))
    } else if let Some(Error::CreateKBError) = r.find() {
        Ok(warp::reply::with_status(
            "Route not found".to_string(),
            StatusCode::NOT_FOUND,
        ))
    } else if let Some(Error::CreateCategoryError) = r.find() {
        Ok(warp::reply::with_status(
            "Route not found".to_string(),
            StatusCode::NOT_FOUND,
        ))
    } else if let Some(Error::SearchError) = r.find() {
        Ok(warp::reply::with_status(
            "Route not found".to_string(),
            StatusCode::NOT_FOUND,
        ))
    } else if let Some(Error::DuplicateKBError) = r.find() {
        Ok(warp::reply::with_status(
            "Route not found".to_string(),
            StatusCode::NOT_FOUND,
        ))
    } else if let Some(Error::DuplicateCategory) = r.find() {
        Ok(warp::reply::with_status(
            "Route not found".to_string(),
            StatusCode::NOT_FOUND,
        ))
    } else if let Some(Error::DatabaseQueryError) = r.find() {
        Ok(warp::reply::with_status(
            "Route not found".to_string(),
            StatusCode::NOT_FOUND,
        ))
    } else {
        Ok(warp::reply::with_status(
            "Route not found".to_string(),
            StatusCode::NOT_FOUND,
        ))
    }
}
