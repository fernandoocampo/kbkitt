#[cfg(test)]
mod pgstorage_test {
    use log::info;
    use std::env;

    use crate::types::kbs::{KBID, KnowledgeBase};
    use crate::{adapters::pgstorage::pgdb, kbs::storage::Storer};

    use tokio::runtime::Runtime;

    const INTEGRATION_TEST: &str = "INTEGRATION_TEST";

    #[test]
    fn test_get_kb_by_id() {
        // Given
        if !is_integration_test() {
            info!("==== skipping test");
            assert_eq!(true, true);
            return
        }
        info!("==== running integration test");

        let id = KBID(String::from("681cca89-890b-4667-8ca0-e328546e268c"));
        let want = KnowledgeBase {
            id: KBID(String::from("681cca89-890b-4667-8ca0-e328546e268c")),
            key: String::from("red"),
            kind: String::from("concepts"),
            notes: String::from("one color"),
            value: String::from("remember this color"),
            tags: vec![String::from("color"), String::from("concepts")],
        };
        let runtime = Runtime::new().expect("Unable to create a runtime");
        let store = runtime.block_on(new_db_storage());

        // When
        let result = runtime.block_on(store.get_kb_by_id(id));

        // Then
        match result {
            Ok(got) => assert_eq!(want, got),
            Err(err) => panic!("unexpected error: {:?}", err),
        }

        runtime.block_on(store.close());
    }

    fn is_integration_test() -> bool {
        match env::var(INTEGRATION_TEST) {
            Ok(val) => {
                matches!(val.as_str(), "true" | "t" | "1")
            },
            Err(_) => {
                false
            },
        }
    }

    async fn new_db_storage() -> pgdb::Store {
        let db_url = "postgres://kb_user:kb_pwd@localhost:5432/knowledgebase";
    
        pgdb::Store::new(db_url).await
    }
}
