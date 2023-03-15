import React, { ChangeEvent } from 'react';
import { InlineField, Input } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onTimeField = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, time_field: event.target.value });
  };

  const onDatabaseChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, database: event.target.value });
    // executes the query
    // onRunQuery();
  };

  const onCollectionChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, collection: event.target.value });
    // executes the query
    // onRunQuery();
  };

  const { time_field, database, collection } = query;

  return (
    <div className="gf-form">
      <InlineField label="Database" tooltip="Name of the Database">
        <Input onChange={onDatabaseChange} value={database} width={16} type="string" />
      </InlineField>
      <InlineField label="Collection" tooltip="Name of the collection">
        <Input onChange={onCollectionChange} value={collection} width={20} type="string" />
      </InlineField>
      <InlineField label="Time field" tooltip="The time field">
        <Input onChange={onTimeField} value={time_field || ''} />
      </InlineField>
    </div>
  );
}
