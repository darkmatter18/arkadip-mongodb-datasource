import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  database: string;
  collection: string;
  time_field: string
}

export const DEFAULT_QUERY: Partial<MyQuery> = {
  database: 'admin'
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  test_db: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  db_uri: string;
}
