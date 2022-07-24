import * as React from 'react';
import { RulesListQuery } from '../../generated/graphql';
import './styles.css';

interface Props {
    data: RulesListQuery;
}

const className = 'RulesList';

const RulesList: React.FC<Props> = ({ data }) => (
    <div className={className}>
        <h3>Rules</h3>
        <table>
            <tr>
                <td>#</td>
                <td>Summary</td>
                <td>Created by</td>
                <td># Likes</td>
            </tr>
                {!!data.rules && data.rules.rules.map(
                    (rule, i) =>
                        !!rule && (
                            <tr key={i}>
                                <td>{rule.id}</td>
                                <td>{rule.summary}</td>
                                <td>{rule.user.name}</td>
                                <td>{rule.likes.users.length}</td>
                            </tr>
                        ),
                )}
        </table>
    </div>
);

export default RulesList;
