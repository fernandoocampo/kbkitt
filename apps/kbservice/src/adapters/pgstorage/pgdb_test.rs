#[cfg(test)]
mod pgstorage_test {
    use log::info;
    use std::env;

    use crate::adapters::pgstorage::pgdb::DBData;
    use crate::types::categories::{Category, CategoryFilter};
    use crate::types::kbs::{KBItem, KBQueryFilter, KnowledgeBase, SearchResult, KBID};
    use crate::{adapters::pgstorage::pgdb, kbs::storage::Storer};

    use tokio::runtime::Runtime;

    const INTEGRATION_TEST: &str = "INTEGRATION_TEST";

    #[test]
    fn test_build_db_url() {
        // Given
        let db_data = DBData {
            user: "kibi".to_string(),
            pwd: "kibipwd".to_string(),
            host: "localhost".to_string(),
            db_name: "kbs".to_string(),
            port: 5432,
        };

        let want = "postgres://kibi:kibipwd@localhost:5432/kbs".to_string();

        // When
        let got = db_data.build_url_connection();

        // Then
        assert_eq!(want, got);
    }

    #[test]
    fn test_get_kb_by_id() {
        // Given
        if !is_integration_test() {
            info!("==== skipping test");
            assert_eq!(true, true);
            return;
        }
        info!("==== running integration test");

        let id = KBID(String::from("681cca89-890b-4667-8ca0-e328546e268c"));
        /*
        INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS)
        VALUES ('681cca89-890b-4667-8ca0-e328546e268c', 'red', 'remember this color', 'one color', 'concepts', 'color concepts')
        RETURNING KB_ID
        */
        let want = KnowledgeBase {
            id: KBID(String::from("681cca89-890b-4667-8ca0-e328546e268c")),
            key: String::from("red"),
            kind: String::from("concepts"),
            notes: String::from("one color"),
            value: String::from("remember this color"),
            reference: None,
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

    #[test]
    fn test_get_kb_by_key() {
        // Given
        if !is_integration_test() {
            info!("==== skipping test");
            assert_eq!(true, true);
            return;
        }
        info!("==== running integration test");

        let key = String::from("red");
        /*
        INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, REFERENCE, TAGS)
        VALUES ('681cca89-890b-4667-8ca0-e328546e268c', 'red', 'remember this color', 'one color', 'concepts', 'color concepts')
        RETURNING KB_ID
        */
        let want = KnowledgeBase {
            id: KBID(String::from("681cca89-890b-4667-8ca0-e328546e268c")),
            key: String::from("red"),
            kind: String::from("concepts"),
            notes: String::from("one color"),
            value: String::from("remember this color"),
            reference: None,
            tags: vec![String::from("color"), String::from("concepts")],
        };
        let runtime = Runtime::new().expect("Unable to create a runtime");
        let store = runtime.block_on(new_db_storage());

        // When
        let result = runtime.block_on(store.get_kb_by_key(key));

        // Then
        match result {
            Ok(got) => assert_eq!(want, got),
            Err(err) => panic!("unexpected error: {:?}", err),
        }

        runtime.block_on(store.close());
    }

    #[test]
    fn test_search_by_key() {
        // Given
        if !is_integration_test() {
            info!("==== skipping test");
            assert_eq!(true, true);
            return;
        }
        info!("==== running integration test");

        let query = KBQueryFilter {
            key: String::from("rick"),
            keyword: Default::default(),
            limit: Some(10),
            offset: 0,
        };
        /*
        INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('6411a28b-640a-43d9-b901-1c4b15d91568', 'frederick', 'long name', 'multiple names', 'names', 'name names');
        INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('5a2579f7-83b9-4891-8dbc-e0024b5f3505', 'rick', 'short name', 'just one name', 'names', 'name names');
        INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('22cfc4fb-f9b6-4f6e-9158-9982347ad2a7', 'patrick', 'a saint', 'names', 'words', 'over words');
        SELECT KB_ID, KB_KEY, KIND, TAGS::TEXT AS TAGS FROM kbs WHERE KB_KEY LIKE '%rick%' LIMIT 10 OFFSET 0;
        */
        let want: SearchResult = SearchResult {
            items: vec![
                KBItem {
                    id: KBID(String::from("6411a28b-640a-43d9-b901-1c4b15d91568")),
                    key: String::from("frederick"),
                    kind: String::from("names"),
                    tags: vec![String::from("name"), String::from("names")],
                },
                KBItem {
                    id: KBID(String::from("22cfc4fb-f9b6-4f6e-9158-9982347ad2a7")),
                    key: String::from("patrick"),
                    kind: String::from("words"),
                    tags: vec![String::from("over"), String::from("words")],
                },
                KBItem {
                    id: KBID(String::from("5a2579f7-83b9-4891-8dbc-e0024b5f3505")),
                    key: String::from("rick"),
                    kind: String::from("names"),
                    tags: vec![String::from("name"), String::from("names")],
                },
            ],
            limit: 10,
            offset: 0,
            total: 3,
        };
        let runtime = Runtime::new().expect("Unable to create a runtime");
        let store = runtime.block_on(new_db_storage());

        // When
        let result = runtime.block_on(store.search_by_key(query));

        // Then
        match result {
            Ok(got) => assert_eq!(want, got),
            Err(err) => panic!("unexpected error: {:?}", err),
        }

        runtime.block_on(store.close());
    }

    #[test]
    fn test_search() {
        // Given
        if !is_integration_test() {
            info!("==== skipping test");
            assert_eq!(true, true);
            return;
        }
        info!("==== running integration test");

        let query = KBQueryFilter {
            keyword: String::from("names"),
            key: Default::default(),
            limit: Some(10),
            offset: 0,
        };
        /*
        INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('6411a28b-640a-43d9-b901-1c4b15d91568', 'frederick', 'long name', 'multiple names', 'names', 'name names');
        INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('5a2579f7-83b9-4891-8dbc-e0024b5f3505', 'rick', 'short name', 'just one name', 'names', 'name names');
        INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('22cfc4fb-f9b6-4f6e-9158-9982347ad2a7', 'patrick', 'a saint', 'names', 'words', 'over words');
        SELECT KB_ID, KB_KEY, KIND, TAGS::TEXT AS TAGS FROM kbs WHERE TAGS @@ to_tsquery('names') ORDER BY KB_KEY LIMIT 10 OFFSET 0;
        */
        let want = SearchResult {
            items: vec![
                KBItem {
                    id: KBID(String::from("6411a28b-640a-43d9-b901-1c4b15d91568")),
                    key: String::from("frederick"),
                    kind: String::from("names"),
                    tags: vec![String::from("name"), String::from("names")],
                },
                KBItem {
                    id: KBID(String::from("5a2579f7-83b9-4891-8dbc-e0024b5f3505")),
                    key: String::from("rick"),
                    kind: String::from("names"),
                    tags: vec![String::from("name"), String::from("names")],
                },
            ],
            limit: 10,
            offset: 0,
            total: 3,
        };
        let runtime = Runtime::new().expect("Unable to create a runtime");
        let store = runtime.block_on(new_db_storage());

        // When
        let result = runtime.block_on(store.search(query));

        // Then
        match result {
            Ok(got) => assert_eq!(want, got),
            Err(err) => panic!("unexpected error: {:?}", err),
        }

        runtime.block_on(store.close());
    }

    #[test]
    fn test_save_kb() {
        // Given
        if !is_integration_test() {
            info!("==== skipping test");
            assert_eq!(true, true);
            return;
        }
        info!("==== running integration test");

        let mut newkb = KnowledgeBase::new(String::from("new_kb_key"));
        newkb.value = String::from("a new kb item");
        newkb.notes = String::from("this is a test");
        newkb.kind = String::from("tests");
        newkb.tags = vec![String::from("test"), String::from("save")];

        let wanted_id = newkb.id.clone();

        let runtime = Runtime::new().expect("Unable to create a runtime");
        let store = runtime.block_on(new_db_storage());

        // When
        let result = runtime.block_on(store.save_kb(newkb));

        // Then
        match result {
            Ok(got) => assert_eq!(wanted_id, got),
            Err(err) => panic!("unexpected error: {:?}", err),
        }

        runtime.block_on(store.close());
    }

    #[test]
    fn test_save_category() {
        // Given
        if !is_integration_test() {
            info!("==== skipping test");
            assert_eq!(true, true);
            return;
        }
        info!("==== running integration test");

        let new_category = Category {
            name: String::from("names"),
            description: String::from("set of well known names"),
        };

        let want = String::from("names");

        let runtime = Runtime::new().expect("unable to create a runtime");
        let storage = runtime.block_on(new_db_storage());

        // When
        let result = runtime.block_on(storage.save_category(new_category));

        // Then
        match result {
            Ok(got) => assert_eq!(want, got),
            Err(err) => panic!("unexpected error: {:?}", err),
        }
    }

    #[test]
    fn test_list_categories() {
        // Given
        if !is_integration_test() {
            info!("==== skipping test");
            assert_eq!(true, true);
            return;
        }
        info!("==== running integration test");

        let filter = CategoryFilter::default();
        /*
        INSERT INTO categories (CATEGORY_NAME, CATEGORY_DESC) VALUES ('concepts', 'a principle or idea');
        INSERT INTO categories (CATEGORY_NAME, CATEGORY_DESC) VALUES ('software', 'any topic related to software development');
        INSERT INTO categories (CATEGORY_NAME, CATEGORY_DESC) VALUES ('personal', 'any personal data');
        SELECT CATEGORY_NAME, CATEGORY_DESC FROM categories ORDER BY CATEGORY_NAME LIMIT 5 OFFSET 0;
        */

        let want = vec![
            Category {
                name: "concepts".to_string(),
                description: "a principle or idea".to_string(),
            },
            Category {
                name: "names".to_string(),
                description: "set of well known names".to_string(),
            },
            Category {
                name: "personal".to_string(),
                description: "any personal data".to_string(),
            },
            Category {
                name: "software".to_string(),
                description: "any topic related to software development".to_string(),
            },
        ];

        let runtime = Runtime::new().expect("unable to create runtime");
        let storage = runtime.block_on(new_db_storage());

        // When
        let result = runtime.block_on(storage.list_categories(filter));

        // Then
        match result {
            Ok(got) => assert_eq!(want, got),
            Err(err) => panic!("unexpected error: {:?}", err),
        }
    }

    #[test]
    fn test_list_categories_with_keyword() {
        // Given
        if !is_integration_test() {
            info!("==== skipping test");
            assert_eq!(true, true);
            return;
        }
        info!("==== running integration test");

        let filter = CategoryFilter {
            keyword: Some(String::from("ftw")),
            limit: Some(5),
            offset: 0,
        };
        /*
        INSERT INTO categories (CATEGORY_NAME, CATEGORY_DESC) VALUES ('concepts', 'a principle or idea');
        INSERT INTO categories (CATEGORY_NAME, CATEGORY_DESC) VALUES ('software', 'any topic related to software development');
        INSERT INTO categories (CATEGORY_NAME, CATEGORY_DESC) VALUES ('personal', 'any personal data');
        SELECT CATEGORY_NAME, CATEGORY_DESC FROM categories WHERE CATEGORY_NAME LIKE '%ftw%' ORDER BY CATEGORY_NAME LIMIT 5 OFFSET 0;
        */

        let want = vec![Category {
            name: "software".to_string(),
            description: "any topic related to software development".to_string(),
        }];

        let runtime = Runtime::new().expect("unable to create runtime");
        let storage = runtime.block_on(new_db_storage());

        // When
        let result = runtime.block_on(storage.list_categories(filter));

        // Then
        match result {
            Ok(got) => assert_eq!(want, got),
            Err(err) => panic!("unexpected error: {:?}", err),
        }
    }

    fn is_integration_test() -> bool {
        match env::var(INTEGRATION_TEST) {
            Ok(val) => {
                matches!(val.as_str(), "true" | "t" | "1")
            }
            Err(_) => false,
        }
    }

    async fn new_db_storage() -> pgdb::Store {
        let db_data = DBData {
            db_name: "kbdb".to_string(),
            host: "localhost".to_string(),
            port: 5432,
            pwd: "kbpwd".to_string(),
            user: "kbdb".to_string(),
        };

        pgdb::Store::new(db_data).await
    }
}
