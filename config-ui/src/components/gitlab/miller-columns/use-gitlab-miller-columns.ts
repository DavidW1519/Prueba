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

import { useMemo, useCallback } from 'react'

import type { ItemType } from '@/components/miller-columns'
import { useLoadItems, ItemTypeEnum } from '@/components/miller-columns'
import request from '@/components/utils/request'

import { getGitLabProxyApiPrefix } from '../config'

export interface UseGitLabMillerColumnsProps {
  connectionId: string
}

export const useGitLabMillerColumns = <T>({
  connectionId
}: UseGitLabMillerColumnsProps) => {
  const prefix = useMemo(
    () => getGitLabProxyApiPrefix(connectionId),
    [connectionId]
  )

  const upadateGroups = (arr: any): Array<ItemType> =>
    arr.map((it: any) => ({
      id: it.id,
      title: it.name,
      type: ItemTypeEnum.BRANCH,
      items: []
    }))

  const updateProjects = (arr: any): Array<ItemType> =>
    arr.map((it: any) => ({
      id: it.id,
      title: it.name,
      type: ItemTypeEnum.LEAF,
      items: [],
      nameWithNameSpace: it.name_with_namespace
    }))

  const getInitItems = useCallback(async () => {
    const user = await request(`${prefix}/user`)
    const [groups, projects] = await Promise.all([
      request(`${prefix}/groups`, {
        data: { top_level_only: 1, per_page: 100 }
      }),
      request(`${prefix}/users/${user.id}/projects`, {
        data: { per_page: 100 }
      })
    ])
    return [...upadateGroups(groups), ...updateProjects(projects)]
  }, [prefix])

  const loadMoreItems = useCallback(
    async (item: ItemType) => {
      const [groups, projects] = await Promise.all([
        request(`${prefix}/groups/${item.id}/subgroups`, {
          data: { per_page: 100 }
        }),
        request(`${prefix}/groups/${item.id}/projects`, {
          data: { per_page: 100 }
        })
      ])
      return [...upadateGroups(groups), ...updateProjects(projects)]
    },
    [prefix]
  )

  const { items, itemTree, loadItems } = useLoadItems<T>({
    getInitItems,
    loadMoreItems
  })

  return {
    items,
    itemTree,
    onExpandItem: loadItems
  }
}
