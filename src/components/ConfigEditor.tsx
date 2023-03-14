import React, { ChangeEvent } from 'react';
import { InlineField, Input, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const onPathChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      test_db: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  const onDBURIKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        db_uri: event.target.value,
      },
    });
  };

  const onResetDBURI = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        db_uri: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        db_uri: '',
      },
    });
  };

  const { jsonData, secureJsonFields } = options;
  const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

  return (
    <div className="gf-form-group">
      <InlineField label="Test Database" labelWidth={12}>
        <Input
          onChange={onPathChange}
          value={jsonData.test_db || ''}
          placeholder="Test Database"
          width={40}
        />
      </InlineField>
      <InlineField label="Database URI" labelWidth={12}>
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.db_uri) as boolean}
          value={secureJsonData.db_uri || ''}
          placeholder="database URI"
          width={40}
          onReset={onResetDBURI}
          onChange={onDBURIKeyChange}
        />
      </InlineField>
    </div>
  );
}
