#[cfg(test)]
mod kbs_tests {
    use crate::types::kbs::{KnowledgeBase, NewKnowledgeBase, KBID};
    use uuid::Uuid;

    #[test]
    fn test_to_knowledge_base() {
        // Given
        let new_knowledge_base = NewKnowledgeBase {
            key: String::from("red"),
            value: String::from("of the colour of fresh blood"),
            category: String::from("concepts"),
            notes: String::from("to know about color red"),
            reference: Some(String::from("Some Author")),
            tags: vec![
                "concept".to_string(),
                "color".to_string(),
                "paint".to_string(),
            ],
        };

        let mut want = KnowledgeBase {
            id: KBID("".to_string()),
            key: "red".to_string(),
            value: "of the colour of fresh blood".to_string(),
            category: "concepts".to_string(),
            notes: String::from("to know about color red"),
            reference: Some(String::from("Some Author")),
            tags: vec![
                "concept".to_string(),
                "color".to_string(),
                "paint".to_string(),
            ],
        };

        // When
        let got = new_knowledge_base.to_knowledge_base();
        // Then
        let id_result = Uuid::parse_str(got.id.to_string().as_str());
        match id_result {
            Ok(_) => {}
            Err(e) => {
                panic!("unexpected error: {}", e);
            }
        }
        want.id = got.id.clone();
        assert_eq!(want, got);
    }
}
