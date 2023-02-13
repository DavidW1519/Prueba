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

import React from 'react';
import { FormGroup, Switch } from '@blueprintjs/core';

import * as S from './styled';

interface Props {
  values: any;
  setValues: any;
}

export const Graphql = ({ values, setValues }: Props) => {
  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    const enableGraphql = (e.target as HTMLInputElement).checked;
    setValues({
      ...values,
      enableGraphql,
    });
  };

  return (
    <FormGroup
      label={<S.Label>Use GraphQL APIs</S.Label>}
      subLabel={
        <S.LabelDescription>
          GraphQL APIs are 10+ times faster than REST APIs, but they may not be supported in GitHub Server.
        </S.LabelDescription>
      }
    >
      <Switch checked={values.enableGraphql ?? false} onChange={handleChange} />
    </FormGroup>
  );
};
