/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import type { TableColumnType } from 'antd';
import { Table } from 'antd';

import type { JenkinsItemType } from './typed';

export interface JenkinsListProps {
  style?: React.CSSProperties;
  extraColumn?: TableColumnType<any>[];
  loading: boolean;
  data: Array<JenkinsItemType>;
}

export const JenkinsConnectionList = ({ extraColumn, loading, data, ...props }: JenkinsListProps) => {
  const columns: TableColumnType<any>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Endpoint',
      dataIndex: 'endpoint',
      key: 'endpoint',
    },
    ...(extraColumn ?? []),
  ];

  return <Table {...props} rowKey="id" loading={loading} columns={columns} dataSource={data} pagination={false} />;
};
