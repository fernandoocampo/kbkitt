use log::LevelFilter;
use log4rs::{
    append::console::ConsoleAppender,
    config::{Appender, Logger, Root},
    encode::json::JsonEncoder,
    Config,
};

use crate::application::transport;
use crate::{
    adapters::pgstorage::pgdb::{DBData, Store},
    kbs,
    types::configuration::{ApplicationSetup, RepositorySetup},
};

pub async fn run() {
    let settings = initialize_configuration();
    initialize_logger();

    log::info!("🗿\tStarting database connection...");
    let store = new_db_storage(settings.repository).await;

    log::info!("🗒️\tInitializing kb service...");
    let service = new_kb_service(store).await;

    log::info!("🚦\tEstablishing API routes...");
    let app_routes = transport::make_create_routes(service);

    log::info!("🍏\tStarting server at :{}", settings.web_server_port);
    transport::start_web_server(settings.web_server_port, app_routes).await;
}

fn initialize_configuration() -> ApplicationSetup {
    log::info!("⏲️\tLoading configuration");

    ApplicationSetup::new()
}

fn initialize_logger() {
    log::info!("🪵\tInitializing logger...");
    log::info!("🪵\tUsing log4rs");
    let app_stdout = ConsoleAppender::builder().build();
    let kb_stdout = ConsoleAppender::builder()
        .encoder(Box::new(JsonEncoder::new()))
        .build();

    let config = Config::builder()
        .appender(Appender::builder().build("app_stdout", Box::new(app_stdout)))
        .appender(Appender::builder().build("kb_stdout", Box::new(kb_stdout)))
        // .logger(Logger::builder().appender("app_stdout").build("people::application::app", LevelFilter::Debug))
        .logger(
            Logger::builder()
                .appender("kb_stdout")
                .build("warp::*", LevelFilter::Debug),
        )
        .logger(
            Logger::builder()
                .appender("kb_stdout")
                .build("hyper::*", LevelFilter::Debug),
        )
        .logger(
            Logger::builder()
                .appender("kb_stdout")
                .build("people::people::handler", LevelFilter::Debug),
        )
        .build(
            Root::builder()
                .appender("app_stdout")
                .build(LevelFilter::Debug),
        )
        .unwrap();

    log4rs::init_config(config).unwrap();
}

async fn new_db_storage(settings: RepositorySetup) -> Store {
    let db_data = DBData {
        db_name: settings.db_name,
        host: settings.host,
        pwd: settings.db_pwd,
        user: settings.db_user,
        port: settings.port,
    };

    Store::new(db_data).await
}

async fn new_kb_service<T: kbs::storage::Storer>(store: T) -> kbs::service::Service<T> {
    kbs::service::Service::new(store)
}
