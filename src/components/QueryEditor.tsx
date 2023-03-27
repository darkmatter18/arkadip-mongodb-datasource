import React, { ChangeEvent } from 'react';
// import { InlineField, Input, InlineLabel,TextArea } from '@grafana/ui';
import { InlineLabel,TextArea } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
   const onQueryChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    onChange({ ...query, q: event.target.value });
    // executes the query
    // onRunQuery();
  };
  
  const { q } = query;

  return (
    <div className="gf-form">
      <InlineLabel width="auto" tooltip="Database Query">
        Query
      </InlineLabel>
      <TextArea onChange={onQueryChange} value={q} rows={10} />
    </div>
  );
}
