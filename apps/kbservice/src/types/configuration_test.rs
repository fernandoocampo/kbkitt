#[cfg(test)]
mod configuration_tests {
    use crate::types::configuration::{ApplicationSetup, RepositorySetup};
    use std::env;

    #[test]
    fn test_new_repository() {
        // Given
        env::set_var("KB_DB_USER", "testuser");
        env::set_var("KB_DB_PWD", "testpwd");
        env::set_var("KB_DB_HOST", "testhost");
        env::set_var("KB_DB_NAME", "testdb");
        env::set_var("KB_DB_PORT", "5433");

        let want = RepositorySetup {
            db_name: String::from("testdb"),
            db_pwd: String::from("testpwd"),
            db_user: String::from("testuser"),
            host: String::from("testhost"),
            port: 5433,
        };

        // When
        let got = RepositorySetup::new();
        // Then
        assert_eq!(want, got)
    }

    #[test]
    fn test_new_repository_no_env_vars() {
        // Given
        env::remove_var("KB_DB_USER");
        env::remove_var("KB_DB_PWD");
        env::remove_var("KB_DB_HOST");
        env::remove_var("KB_DB_NAME");
        env::remove_var("KB_DB_PORT");

        let want = RepositorySetup {
            db_name: String::from(""),
            db_pwd: String::from(""),
            db_user: String::from(""),
            host: String::from(""),
            port: 5432,
        };

        // When
        let got = RepositorySetup::new();
        // Then
        assert_eq!(want, got)
    }

    #[test]
    fn test_new_application() {
        // Given
        env::set_var("KB_DB_USER", "testuser");
        env::set_var("KB_DB_PWD", "testpwd");
        env::set_var("KB_DB_HOST", "testhost");
        env::set_var("KB_DB_NAME", "testdb");
        env::set_var("KB_DB_PORT", "5433");
        env::set_var("KB_WEB_SERVER_PORT", "3035");

        let want = ApplicationSetup {
            web_server_port: 3035,
            repository: RepositorySetup {
                db_name: String::from("testdb"),
                db_pwd: String::from("testpwd"),
                db_user: String::from("testuser"),
                host: String::from("testhost"),
                port: 5433,
            },
        };

        // When
        let got = ApplicationSetup::new();
        // Then
        assert_eq!(want, got)
    }
}
