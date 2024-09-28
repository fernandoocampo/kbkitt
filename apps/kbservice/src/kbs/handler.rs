use std::collections::HashMap;

use crate::errors::error::Error;
use crate::kbs::service::Service;
use crate::types::categories::{Category, SaveCategorySuccess};
use crate::types::kbs::{KnowledgeBase, NewKnowledgeBase, SaveKBSuccess, KBID};

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

pub async fn update_kb(
    kb: KnowledgeBase,
    service: Service<impl storage::Storer>,
) -> Result<impl Reply, Rejection> {
    match service.update_kb(kb.clone()).await {
        Ok(_) => {
            debug!("kb was updated");

            Ok(warp::reply())
        }
        Err(e) => {
            error!("updating kb {:?}: {:?}", kb, e);
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
            "KB not found".to_string(),
            StatusCode::NOT_FOUND,
        ))
    } else if let Some(Error::GetKBError) = r.find() {
        Ok(warp::reply::with_status(
            "Unable to get KB".to_string(),
            StatusCode::INTERNAL_SERVER_ERROR,
        ))
    } else if let Some(Error::CreateKBError) = r.find() {
        Ok(warp::reply::with_status(
            "Unable to create KB".to_string(),
            StatusCode::INTERNAL_SERVER_ERROR,
        ))
    } else if let Some(Error::CreateCategoryError) = r.find() {
        Ok(warp::reply::with_status(
            "Unable to create category".to_string(),
            StatusCode::INTERNAL_SERVER_ERROR,
        ))
    } else if let Some(Error::SearchError) = r.find() {
        Ok(warp::reply::with_status(
            "Unable to search KBs".to_string(),
            StatusCode::INTERNAL_SERVER_ERROR,
        ))
    } else if let Some(Error::DuplicateKBError) = r.find() {
        Ok(warp::reply::with_status(
            "KB already exists".to_string(),
            StatusCode::CONFLICT,
        ))
    } else if let Some(Error::DuplicateCategory) = r.find() {
        Ok(warp::reply::with_status(
            "Category already exists".to_string(),
            StatusCode::CONFLICT,
        ))
    } else if let Some(Error::DatabaseQueryError) = r.find() {
        Ok(warp::reply::with_status(
            "Service Database is not available".to_string(),
            StatusCode::INTERNAL_SERVER_ERROR,
        ))
    } else {
        Ok(warp::reply::with_status(
            "Unexpected error".to_string(),
            StatusCode::INTERNAL_SERVER_ERROR,
        ))
    }
}
