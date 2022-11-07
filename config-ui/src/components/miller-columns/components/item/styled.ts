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

import styled from '@emotion/styled'

export const Wrapper = styled.div<{ selected: boolean }>`
  display: flex;
  align-items: center;
  padding: 4px 12px;
  cursor: pointer;

  ${({ selected }) => (selected ? 'background-color: #f5f5f7;' : '')}

  &:hover {
    background-color: #f5f5f7;
  }

  & > span.name {
    font-size: 14px;
  }

  & > span.count {
    font-size: 14px;
  }
`
