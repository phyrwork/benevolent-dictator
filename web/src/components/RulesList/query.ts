import { gql } from '@apollo/client';

export const QUERY_RULES_LIST = gql`
    query RulesList($limit: Int!) {
        rules(limit: $limit) {
            rules {
                id
                user {
                    name
                }
                summary
                likes(limit: 20) {
                    users {
                        id
                    }
                }
            }
        }
    }
`;
