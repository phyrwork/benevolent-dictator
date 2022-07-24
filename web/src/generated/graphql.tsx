import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type LikesUpdate = {
  __typename?: 'LikesUpdate';
  added: Array<Scalars['Int']>;
  removed: Array<Scalars['Int']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  createRule: Rule;
  createUser: User;
  deleteRule?: Maybe<Scalars['ID']>;
  like?: Maybe<LikesUpdate>;
  login: UserToken;
  updateUser: User;
};


export type MutationCreateRuleArgs = {
  detail?: InputMaybe<Scalars['String']>;
  summary: Scalars['String'];
};


export type MutationCreateUserArgs = {
  email: Scalars['String'];
  name: Scalars['String'];
  password: Scalars['String'];
};


export type MutationDeleteRuleArgs = {
  id: Scalars['ID'];
};


export type MutationLikeArgs = {
  add?: InputMaybe<Array<Scalars['ID']>>;
  remove?: InputMaybe<Array<Scalars['ID']>>;
};


export type MutationLoginArgs = {
  email: Scalars['String'];
  password: Scalars['String'];
};


export type MutationUpdateUserArgs = {
  name?: InputMaybe<Scalars['String']>;
};

export type PageInfo = {
  __typename?: 'PageInfo';
  endCursor?: Maybe<Scalars['ID']>;
  hasNextPage: Scalars['Boolean'];
  hasPreviousPage: Scalars['Boolean'];
  startCursor?: Maybe<Scalars['ID']>;
};

export type Query = {
  __typename?: 'Query';
  rules: RulePage;
  users: UserPage;
};


export type QueryRulesArgs = {
  after?: Scalars['Int'];
  limit?: Scalars['Int'];
  userId?: InputMaybe<Scalars['ID']>;
};


export type QueryUsersArgs = {
  after?: Scalars['Int'];
  limit?: Scalars['Int'];
  name?: InputMaybe<Scalars['String']>;
};

export type Rule = {
  __typename?: 'Rule';
  created: Scalars['String'];
  detail?: Maybe<Scalars['String']>;
  id: Scalars['ID'];
  likes: UserPage;
  summary: Scalars['String'];
  user: User;
};


export type RuleLikesArgs = {
  after?: Scalars['Int'];
  limit?: Scalars['Int'];
};

export type RulePage = {
  __typename?: 'RulePage';
  pageInfo: PageInfo;
  rules: Array<Rule>;
};

export type User = {
  __typename?: 'User';
  id: Scalars['ID'];
  likes: RulePage;
  name: Scalars['String'];
  rules: RulePage;
};


export type UserLikesArgs = {
  after?: Scalars['Int'];
  limit?: Scalars['Int'];
};


export type UserRulesArgs = {
  after?: Scalars['Int'];
  limit?: Scalars['Int'];
};

export type UserPage = {
  __typename?: 'UserPage';
  pageInfo: PageInfo;
  users: Array<User>;
};

export type UserToken = {
  __typename?: 'UserToken';
  expiresAt: Scalars['Int'];
  token: Scalars['String'];
};

export type RulesListQueryVariables = Exact<{
  limit: Scalars['Int'];
}>;


export type RulesListQuery = { __typename?: 'Query', rules: { __typename?: 'RulePage', rules: Array<{ __typename?: 'Rule', id: string, summary: string, user: { __typename?: 'User', name: string }, likes: { __typename?: 'UserPage', users: Array<{ __typename?: 'User', id: string }> } }> } };


export const RulesListDocument = gql`
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

/**
 * __useRulesListQuery__
 *
 * To run a query within a React component, call `useRulesListQuery` and pass it any options that fit your needs.
 * When your component renders, `useRulesListQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useRulesListQuery({
 *   variables: {
 *      limit: // value for 'limit'
 *   },
 * });
 */
export function useRulesListQuery(baseOptions: Apollo.QueryHookOptions<RulesListQuery, RulesListQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<RulesListQuery, RulesListQueryVariables>(RulesListDocument, options);
      }
export function useRulesListLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<RulesListQuery, RulesListQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<RulesListQuery, RulesListQueryVariables>(RulesListDocument, options);
        }
export type RulesListQueryHookResult = ReturnType<typeof useRulesListQuery>;
export type RulesListLazyQueryHookResult = ReturnType<typeof useRulesListLazyQuery>;
export type RulesListQueryResult = Apollo.QueryResult<RulesListQuery, RulesListQueryVariables>;