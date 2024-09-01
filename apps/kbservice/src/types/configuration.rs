use serde::{Deserialize, Serialize};
use std::env;

const KB_WEB_SERVER_PORT_ENV_VAR: &str = "KB_WEB_SERVER_PORT";
const WEB_SERVER_PORT_DEFAULT: u16 = 3030;
const KB_DB_USER_ENV_VAR: &str = "KB_DB_USER";
const KB_DB_PWD_ENV_VAR: &str = "KB_DB_PWD";
const KB_DB_HOST_ENV_VAR: &str = "KB_DB_HOST";
const KB_DB_NAME_ENV_VAR: &str = "KB_DB_NAME";
const KB_DB_PORT_ENV_VAR: &str = "KB_DB_PORT";
const DB_PORT_DEFAULT: u16 = 5432;

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct ApplicationSetup {
    pub web_server_port: u16,
    pub repository: RepositorySetup,
}

#[derive(Debug, Clone, Serialize, Deserialize, Eq, Hash, PartialEq, PartialOrd, Ord)]
pub struct RepositorySetup {
    pub db_user: String,
    pub db_pwd: String,
    pub host: String,
    pub db_name: String,
    pub port: u16,
}

impl ApplicationSetup {
    pub fn new() -> Self {
        ApplicationSetup {
            web_server_port: load_env_var_int(KB_WEB_SERVER_PORT_ENV_VAR, WEB_SERVER_PORT_DEFAULT),
            repository: RepositorySetup::new(),
        }
    }
}

impl RepositorySetup {
    pub fn new() -> Self {
        RepositorySetup {
            db_name: load_env_var(KB_DB_NAME_ENV_VAR),
            db_pwd: load_env_var(KB_DB_PWD_ENV_VAR),
            db_user: load_env_var(KB_DB_USER_ENV_VAR),
            host: load_env_var(KB_DB_HOST_ENV_VAR),
            port: load_env_var_int(KB_DB_PORT_ENV_VAR, DB_PORT_DEFAULT),
        }
    }
}

impl Default for ApplicationSetup {
    fn default() -> Self {
        Self::new()
    }
}

impl Default for RepositorySetup {
    fn default() -> Self {
        Self::new()
    }
}

fn load_env_var(key: &str) -> String {
    match env::var(key) {
        Ok(val) => val,
        Err(_) => {
            println!("{} env var was not set using default", key);
            "".to_string()
        }
    }
}

fn load_env_var_int(key: &str, default: u16) -> u16 {
    match env::var(key) {
        Ok(val) => match val.parse::<u16>() {
            Ok(num) => num,
            Err(_) => {
                println!("{} env var is not a valid number using default", key);
                default
            }
        },
        Err(_) => {
            println!("{} env var was not set using default", key);
            default
        }
    }
}
