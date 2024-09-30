use warp::{http::Method, Filter};

use crate::kbs::{self, handler, service, storage::Storer};

// https://github.com/firegloves/Hexagonal-Sandwich-Recipes/blob/main/chapter-6/Cargo.toml
// Map<impl Filter<Extract = (), Error = Infallible> + Copy, impl Fn() -> Service<Store>>
// https://stackoverflow.com/questions/68227799/how-to-create-function-to-return-routes-with-warp
// https://stackoverflow.com/questions/67413744/warp-asks-for-absurdly-long-and-complex-explicit-type-annotations-is-there-anot
pub fn make_create_routes(
    // service: service::Service<impl Storer + Clone + Filter<Extract = (), Error = Infallible> + Send + 'static + Sync>,
    service: service::Service<impl Storer + Clone + Send + 'static + Sync>,
    // service_filter: impl Clone + Filter<Extract = (HashMap<String, String>, Service<impl Storer>), Error = Infallible> + Copy + Send + 'static + Sync,
    // service_filter: impl Clone + Filter<Extract = (), Error = Infallible> + Copy + Send + 'static + Sync,
    // service_filter: impl Filter<Extract = (impl Filter<Extract = (), Error = Infallible> + Copy, impl Fn() -> Service<Store>), Error = Infallible> + Clone,
    // service_filter: impl Filter<Extract = (), Error = Infallible> + Clone,
) -> impl Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    let service_filter = warp::any().map(move || service.clone());

    let cors = warp::cors()
        .allow_any_origin()
        .allow_header("content-type")
        .allow_methods(&[Method::PUT, Method::DELETE, Method::GET, Method::POST]);

    log::info!("📚\tCreating search kbs endpoint: GET /kbs");
    let search_kbs = warp::get()
        .and(warp::path("kbs"))
        .and(warp::path::end())
        .and(warp::query())
        .and(service_filter.clone())
        .and_then(kbs::handler::search)
        .with(warp::trace(|info| {
            tracing::info_span!(
                "get_people request",
                method = %info.method(),
                path = %info.path(),
                id = %uuid::Uuid::new_v4(),
            )
        }));

    log::info!("📖\tCreating get kb by id endpoint: GET /kbs/{{id}}");
    let get_kb_by_id = warp::get()
        .and(warp::path("kbs"))
        .and(warp::path::param::<String>())
        .and(warp::path::end())
        .and(service_filter.clone())
        .and_then(kbs::handler::get_kb_with_id);

    log::info!("📖\tCreating add kb endpoint: POST /kbs");
    let add_kb = warp::post()
        .and(warp::path("kbs"))
        .and(warp::path::end())
        .and(warp::body::json())
        .and(service_filter.clone())
        .and_then(kbs::handler::add_kb);

    log::info!("📖\tCreating update kb endpoint: PATCH /kbs");
    let update_kb = warp::patch()
        .and(warp::path("kbs"))
        .and(warp::path::end())
        .and(warp::body::json())
        .and(service_filter.clone())
        .and_then(kbs::handler::update_kb);

    log::info!("📗\tCreating add category endpoint: POST /categories");
    let add_category = warp::post()
        .and(warp::path("categories"))
        .and(warp::path::end())
        .and(warp::body::json())
        .and(service_filter.clone())
        .and_then(kbs::handler::add_category);

    search_kbs
        .or(get_kb_by_id)
        .or(add_kb)
        .or(update_kb)
        .or(add_category)
        .with(cors)
        .with(warp::trace::request())
        .recover(handler::return_error)
}

pub async fn start_web_server(
    port: u16,
    routes: impl Filter<Extract = impl warp::Reply, Error = warp::Rejection>
        + Clone
        + Send
        + Sync
        + 'static,
) {
    warp::serve(routes).run(([0, 0, 0, 0], port)).await;
}
