import * as React from 'react';
import { useRulesListQuery } from '../../generated/graphql';
import RulesList from './RulesList';

const RulesListContainer = () => {
    const { data, error, loading } = useRulesListQuery({ variables: { limit: 20 } });

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error || !data) {
        return <div>ERROR</div>;
    }

    return <RulesList data={data} />;
};

export default RulesListContainer;