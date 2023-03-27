import React, { ChangeEvent } from 'react';
// import { InlineField, Input, InlineLabel,TextArea } from '@grafana/ui';
import { InlineLabel,TextArea } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  // const onTimeField = (event: ChangeEvent<HTMLInputElement>) => {
  //   onChange({ ...query, time_field: event.target.value });
  // };

  // const onDatabaseChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   onChange({ ...query, database: event.target.value });
  //   // executes the query
  //   // onRunQuery();
  // };

  // const onCollectionChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   onChange({ ...query, collection: event.target.value });
  //   // executes the query
  //   // onRunQuery();
  // };


   const onQueryChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    onChange({ ...query, q: event.target.value });
    // executes the query
    // onRunQuery();
  };
  

  // const { time_field, database, collection } = query;
  const { q } = query;

  return (
    <div className="gf-form">
      <InlineLabel width="auto" tooltip="Database Query">
        Query
      </InlineLabel>
      <TextArea onChange={onQueryChange} value={q} />

      {/* <InlineField label="Database" tooltip="Name of the Database">
        <Input onChange={onDatabaseChange} value={database} width={16} type="string" />
      </InlineField>
      <InlineField label="Collection" tooltip="Name of the collection">
        <Input onChange={onCollectionChange} value={collection} width={20} type="string" />
      </InlineField>
      <InlineField label="Time field" tooltip="The time field">
        <Input onChange={onTimeField} value={time_field || ''} />
      </InlineField> */}
    </div>
  );
}
