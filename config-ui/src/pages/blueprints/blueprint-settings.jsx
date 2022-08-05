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
import React, { useEffect, useState, useCallback, useRef, useMemo } from 'react'
import { useParams, useHistory } from 'react-router-dom'
import dayjs from '@/utils/time'
import {
  API_PROXY_ENDPOINT,
  ISSUE_TYPES_ENDPOINT,
  ISSUE_FIELDS_ENDPOINT,
  BOARDS_ENDPOINT,
} from '@/config/jiraApiProxy'
import request from '@/utils/request'
import {
  Button,
  Elevation,
  Intent,
  Switch,
  Card,
  Tag,
  Tooltip,
  Icon,
  Colors,
  Divider,
  Spinner,
  Classes,
  Position,
  Popover,
  Collapse,
  Dialog,
} from '@blueprintjs/core'

import { integrationsData } from '@/data/integrations'
import { NullBlueprint, BlueprintMode } from '@/data/NullBlueprint'
import { NullPipelineRun } from '@/data/NullPipelineRun'
import { Providers, ProviderLabels, ProviderIcons } from '@/data/Providers'
import {
  StageStatus,
  TaskStatus,
  TaskStatusLabels,
  StatusColors,
  StatusBgColors,
} from '@/data/Task'

import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import Content from '@/components/Content'
import { ToastNotification } from '@/components/Toast'
import BlueprintNameCard from '@/components/blueprints/BlueprintNameCard'
import DataSync from '@/components/blueprints/create-workflow/DataSync'

import { DataEntities, DataEntityTypes } from '@/data/DataEntities'
import { DEFAULT_DATA_ENTITIES } from '@/data/BlueprintWorkflow'

import useBlueprintManager from '@/hooks/useBlueprintManager'
import useConnectionManager from '@/hooks/useConnectionManager'
import useJIRA from '@/hooks/useJIRA'
import useBlueprintValidation from '@/hooks/useBlueprintValidation'
import BlueprintDialog from '@/components/blueprints/BlueprintDialog'
import BlueprintDataScopesDialog from '@/components/blueprints/BlueprintDataScopesDialog'
import BlueprintNavigationLinks from '@/components/blueprints/BlueprintNavigationLinks'

// eslint-disable-next-line no-unused-vars
const TEST_CONNECTIONS = [
  {
    id: 0,
    provider: integrationsData.find((i) => i.id === Providers.GITHUB),
    providerLabel: ProviderLabels[Providers.GITHUB],
    name: 'Merico Github',
    entities: ['Source Code Management', 'Issue Tracking', 'Code Review'],
    projects: [
      'apache/incubator-devlake',
      'merico/devstream',
      'merico/another-project',
    ],
    boards: [],
    transformation: {},
    transformationStates: ['Added', '-', 'Added'],
    editable: true,
  },
  {
    id: 1,
    provider: integrationsData.find((i) => i.id === Providers.JIRA),
    providerLabel: ProviderLabels[Providers.JIRA],
    name: 'Merico JIRA',
    entities: ['Source Code Management', 'Issue Tracking', 'Code Review'],
    projects: [],
    boards: ['Board 1', 'Board 2', 'Board 3', 'Board 4'],
    transformation: {},
    transformationStates: ['Added', 'Added', '-', '-'],
    editable: true,
  },
]

const BlueprintSettings = (props) => {
  // eslint-disable-next-line no-unused-vars
  const history = useHistory()
  const { bId } = useParams()

  const [activeProvider, setActiveProvider] = useState(integrationsData[0])
  const [activeTransformation, setActiveTransformation] = useState()

  const [blueprintId, setBlueprintId] = useState()
  const [activeBlueprint, setActiveBlueprint] = useState(NullBlueprint)
  const [currentRun, setCurrentRun] = useState(NullPipelineRun)
  const [dataEntitiesList, setDataEntitiesList] = useState([
    ...DEFAULT_DATA_ENTITIES,
  ])

  const [connections, setConnections] = useState([])
  const [blueprintConnections, setBlueprintConnections] = useState([])
  const [configuredConnection, setConfiguredConnection] = useState()
  const [configuredProject, setConfiguredProject] = useState(null)
  const [configuredBoard, setConfiguredBoard] = useState(null)
  const [scopeConnection, setScopeConnection] = useState()

  const [boards, setBoards] = useState({ [configuredConnection?.id]: [] })
  const [projects, setProjects] = useState({ [configuredConnection?.id]: [] })

  const [blueprintDialogIsOpen, setBlueprintDialogIsOpen] = useState(false)
  const [blueprintScopesDialogIsOpen, setBlueprintScopesDialogIsOpen] = useState(false)

  const [activeSetting, setActiveSetting] = useState({
    id: null,
    title: '',
    payload: {},
  })
  const [activeSettingComponent, setActiveSettingComponent] = useState(null)

  const {
    // eslint-disable-next-line no-unused-vars
    activeStep,
    blueprint,
    name: blueprintName,
    cronConfig,
    customCronConfig,
    enable,
    tasks: blueprintTasks,
    settings: blueprintSettings,
    mode,
    interval,
    isSaving,
    isFetching: isFetchingBlueprint,
    activateBlueprint,
    deactivateBlueprint,
    getNextRunDate,
    // eslint-disable-next-line no-unused-vars
    fetchBlueprint,
    patchBlueprint,
    setName: setBlueprintName,
    setCronConfig,
    setCustomCronConfig,
    setEnable,
    setMode,
    setInterval,
    setIsManual,
    setSettings,
    createCron,
    getCronPreset,
    getCronPresetByConfig,
    detectCronInterval,
    fetchAllBlueprints,
    saveBlueprint,
    saveComplete,
  } = useBlueprintManager()

  const {
    fetchConnection,
    allProviderConnections,
    connectionsList,
    isFetching: isFetchingConnection,
    fetchAllConnections,
  } = useConnectionManager(
    {
      activeProvider,
      connectionId: configuredConnection?.connectionId,
    },
    configuredConnection && configuredConnection?.id !== null
  )

  const {
    validate: validateBlueprint,
    errors: blueprintValidationErrors,
    isValid: isValidBlueprint,
    fieldHasError,
    getFieldError,
    isValidCronExpression,
    validateBlueprintName,
  } = useBlueprintValidation({
    name: blueprintName,
    boards,
    projects,
    cronConfig,
    customCronConfig,
    enable,
    tasks: blueprintTasks,
    mode,
    // connections: blueprintConnections,
    // entities: dataEntities,
    activeStep,
    // activeProvider: provider,
    // activeConnection: configuredConnection
  })

  const {
    fetchIssueTypes,
    fetchFields,
    fetchBoards,
    boards: jiraApiBoards,
    issueTypes: jiraApiIssueTypes,
    fields: jiraApiFields,
    isFetching: isFetchingJIRA,
    error: jiraProxyError,
  } = useJIRA(
    {
      apiProxyPath: API_PROXY_ENDPOINT,
      issuesEndpoint: ISSUE_TYPES_ENDPOINT,
      fieldsEndpoint: ISSUE_FIELDS_ENDPOINT,
      boardsEndpoint: BOARDS_ENDPOINT,
    },
    configuredConnection
  )

  // const withBlueprintName = useMemo(() => BlueprintNameCard => props => (
  //   <BlueprintNameCard
  //     name={blueprintName}
  //     setBlueprintName={setBlueprintName}
  //     fieBldHasError={fieldHasError}
  //     getFieldError={getFieldError}
  //     elevation={Elevation.ZERO}
  //     enableDivider={false}
  //     cardStyle={{ padding: 0 }}
  //     {...props}
  //   />
  // ), [blueprintName])

  // const NameSettings = useMemo(() => withBlueprintName(BlueprintNameCard), [BlueprintNameCard, withBlueprintName])

  const handleBlueprintActivation = useCallback(
    (blueprint) => {
      if (blueprint.enable) {
        deactivateBlueprint(blueprint)
      } else {
        activateBlueprint(blueprint)
      }
    },
    [activateBlueprint, deactivateBlueprint]
  )

  const handleBlueprintDialogClose = useCallback(() => {
    setBlueprintDialogIsOpen(false)
    setBlueprintName(activeBlueprint?.name)
  }, [activeBlueprint, setBlueprintName])

  const handleBlueprintScopesDialogClose = useCallback(() => {
    setBlueprintScopesDialogIsOpen(false)
  }, [])

  const handleBlueprintSave = useCallback(() => {
    patchBlueprint(activeBlueprint, activeSetting?.payload, () =>
      handleBlueprintDialogClose()
    )
  }, [activeSetting, activeBlueprint, patchBlueprint, handleBlueprintDialogClose])

  const handleConnectionStepChange = useCallback((newStepId, lastStepId, e) => {
    console.log('>>> CONNECTION SETTINGS STEP CHANGED...', newStepId, lastStepId, e)
    switch (newStepId) {
      case 'scopes':
        setActiveSetting((aS) => ({ ...aS, id: 'scopes', title: 'Change Data Scope' }))
        break
      case 'transformations':
        setActiveSetting((aS) => ({ ...aS, id: 'transformations', title: 'Change Transformation Rules' }))
        break
    }
  }, [])

  const viewBlueprintStatus = useCallback(() => {
    history.push(`/blueprints/detail/${blueprintId}`)
  }, [history, blueprintId])

  const viewBlueprintSettings = useCallback(() => {
    history.push(`/blueprints/settings/${blueprintId}`)
  }, [history, blueprintId])

  const modifySetting = useCallback(
    (settingId) => {
      let title = null
      switch (settingId) {
        case 'name':
          title = 'Change Blueprint Name'
          break
        case 'cronConfig':
          title = 'Change Sync Frequency'
          break
        case 'plan':
          title = 'Change Task Configuration'
          break
        default:
          break
      }
      setActiveSetting((aS) => ({ ...aS, id: settingId, title }))
      setBlueprintDialogIsOpen(true)
      fetchBlueprint(blueprintId)
    },
    [blueprintId, fetchBlueprint]
  )

  const modifyConnection = useCallback((connectionIdx, connectionId, provider) => {
    console.log('>>> MODIFYING DATA CONNECTION SCOPE...', provider?.id === Providers.JIRA, connectionIdx, connectionId, provider?.id)
    const connection = connectionsList.find(c => c.connectionId === connectionId && c.provider === provider?.id)
    setActiveSetting((aS) => ({ ...aS, id: 'scopes', title: 'Change Data Scope' }))
    setConfiguredConnection(connection)
    setScopeConnection(connections.find(c => c.connectionId === connectionId && c.provider?.id === provider?.id))
    setActiveProvider(aP => connection ? integrationsData.find(i => i.id === connection?.provider) : aP)
    // @todo: restrict jira proxy calls to jira provider only
    // // if (provider?.id === Providers.JIRA) {
    //   console.log('>>> JIRA PROVIDER DETECTED, FETCHING API PROXY RESOURCES...', provider)
    //   fetchBoards()
    //   fetchIssueTypes()
    //   fetchFields()
    // // }
    setBlueprintScopesDialogIsOpen(true)
  }, [
    connectionsList,
    connections,
    // fetchBoards,
    // fetchFields,
    // fetchIssueTypes
  ])

  const validateActiveSetting = useCallback(() => {
    let isValid = false
    switch (activeSetting?.id) {
      case 'name':
        isValid = validateBlueprintName(blueprintName)
        break
      case 'cronConfig':
        isValid =
          cronConfig === 'custom'
            ? isValidCronExpression(customCronConfig)
            : ['manual', 'custom'].includes(cronConfig) ||
              isValidCronExpression(cronConfig)
        break
    }
    return isValid
  }, [
    activeSetting?.id,
    blueprintName,
    cronConfig,
    customCronConfig,
    validateBlueprintName,
    isValidCronExpression,
  ])

  const getDefaultEntities = useCallback((providerId) => {
    let entities = []
    switch (providerId) {
      case Providers.GITHUB:
      case Providers.GITLAB:
        entities = DEFAULT_DATA_ENTITIES.filter((d) => d.name !== 'ci-cd')
        break
      case Providers.JIRA:
        entities = DEFAULT_DATA_ENTITIES.filter((d) => d.name === 'issue-tracking' || d.name === 'cross-domain')
        break
      case Providers.JENKINS:
        entities = DEFAULT_DATA_ENTITIES.filter((d) => d.name === 'ci-cd')
        break
    }
    return entities
  }, [])

  useEffect(() => {
    console.log('>>> ACTIVE PROVIDER!', activeProvider)
    setDataEntitiesList(deList => activeProvider ? getDefaultEntities(activeProvider?.id) : deList)
    if (activeProvider?.id === Providers.JIRA) {
      console.log('>>> JIRA PROVIDER DETECTED, FETCHING API PROXY RESOURCES...', activeProvider)
      fetchBoards()
      fetchIssueTypes()
      fetchFields()
    }
  }, [activeProvider, getDefaultEntities, setDataEntitiesList, fetchBoards, fetchIssueTypes, fetchFields])

  useEffect(() => {
    setBlueprintId(bId)
    console.log('>>> REQUESTED SETTINGS for BLUEPRINT ID ===', bId)
  }, [bId])

  useEffect(() => {
    if (!isNaN(blueprintId)) {
      console.log('>>>> FETCHING BLUEPRINT ID...', blueprintId)
      fetchBlueprint(blueprintId)
      fetchAllConnections(false, true)
    }
  }, [blueprintId, fetchBlueprint, fetchAllConnections])

  useEffect(() => {
    console.log('>>>> SETTING ACTIVE BLUEPRINT...', blueprint)
    if (blueprint?.id) {
      setActiveBlueprint((b) => ({
        ...b,
        ...blueprint,
      }))
    }
  }, [blueprint])

  useEffect(() => {
    console.log('>>> ACTIVE BLUEPRINT ....', activeBlueprint, jiraApiBoards)
    setConnections(
      activeBlueprint?.settings?.connections.map((c, cIdx) => ({
        id: cIdx,
        connectionId: c.connectionId,
        provider: integrationsData.find((i) => i.id === c.plugin),
        providerLabel: ProviderLabels[c.plugin.toUpperCase()],
        icon: ProviderIcons[c.plugin] ? ProviderIcons[c.plugin](18, 18) : null,
        // name: `Connection ID #${c.connectionId}`,
        name: allProviderConnections.find(pC => pC.connectionId === c.connectionId && pC.provider === c.plugin)?.name || `Connection ID #${c.connectionId}`,
        entities: c.scope[0]?.entities?.map((e) => DEFAULT_DATA_ENTITIES.find(de => de.value === e)?.title),
        entityList: c.scope[0]?.entities?.map((e) => DEFAULT_DATA_ENTITIES.find(de => de.value === e)),
        projects: [Providers.GITLAB, Providers.GITHUB].includes(c.plugin)
          ? c.scope.map((s) => `${s.options.owner}/${s.options?.repo}`)
          : [],
        boards: [Providers.JIRA].includes(c.plugin)
          ? c.scope.map((s) => `Board ${s.options?.boardId}`)
          : [],
        boardsList: [Providers.JIRA].includes(c.plugin)
          ? c.scope.map((s) => jiraApiBoards.find(apiBoard => apiBoard.id === s.options?.boardId))
          : [],
        transformation: c.scope.map((s) => ({ ...s.transformation })),
        transformationStates: c.scope.map((s) =>
          Object.values(s.transformation).some((v) => v?.toString().length > 0)
            ? 'Added'
            : '-'
        ),
        editable: true,
      }))
    )
    setBlueprintName(activeBlueprint?.name)
    setCronConfig(
      [
        getCronPreset('hourly').cronConfig,
        getCronPreset('daily').cronConfig,
        getCronPreset('weekly').cronConfig,
        getCronPreset('monthly').cronConfig,
      ].includes(activeBlueprint?.cronConfig)
        ? activeBlueprint?.cronConfig
        : activeBlueprint?.isManual
          ? 'manual'
          : 'custom'
    )
    setCustomCronConfig(
      !['custom', 'manual'].includes(activeBlueprint?.cronConfig)
        ? activeBlueprint?.cronConfig
        : '0 0 * * *'
    )
    setInterval(detectCronInterval(activeBlueprint?.cronConfig))
    setMode(activeBlueprint?.mode)
    setEnable(activeBlueprint?.enable)
    setIsManual(activeBlueprint?.isManual)
    setSettings(activeBlueprint?.settings)
  }, [
    activeBlueprint,
    setBlueprintName,
    setConnections,
    detectCronInterval,
    getCronPreset,
    setCronConfig,
    setCustomCronConfig,
    setEnable,
    setInterval,
    setIsManual,
    setMode,
    setSettings,
    jiraApiBoards,
    allProviderConnections
  ])

  useEffect(() => {
    console.log('>>> SETTING ACTIVE SETTINGS PAYLOAD....', {
      name: blueprintName,
    })
    const isCustomCron = cronConfig === 'custom'
    const isManualCron = cronConfig === 'manual'

    switch (activeSetting?.id) {
      case 'name':
        setActiveSetting((aS) => ({
          ...aS,
          payload: {
            name: blueprintName,
          },
        }))
        break
      case 'cronConfig':
        setActiveSetting((aS) => ({
          ...aS,
          payload: {
            isManual: !!isManualCron,
            cronConfig: isManualCron
              ? getCronPreset('daily').cronConfig
              : isCustomCron
                ? customCronConfig
                : cronConfig,
          },
        }))
        break
    }
  }, [blueprintName, cronConfig, customCronConfig, activeSetting?.id, getCronPreset])

  useEffect(() => {
    console.log(
      '>>> RECEIVED ACTIVE SETTINGS PAYLOAD....',
      activeSetting?.payload
    )
  }, [activeSetting?.payload])

  useEffect(() => {
    console.log(
      '>>> ACTIVE UI SETTING OBJECT...',
      activeSetting
    )
  }, [activeSetting])

  useEffect(() => {
    validateBlueprint()
  }, [
    blueprintName,
    // @todo: fix dependency warning with validateBlueprint
    // validateBlueprint
  ])

  useEffect(() => {
    console.log('>>> DATA SCOPE CONNECTIONS...', connections)
    setBlueprintConnections(
      connections.map(
        c => connectionsList.find(
          cItem => cItem.connectionId === c.connectionId && cItem.provider === c.provider.id
        )
      )
    )
  }, [connections, connectionsList])

  useEffect(() => {
    console.log('>>> AVAILABLE DATA ENTITIES...', dataEntitiesList)
  }, [dataEntitiesList])

  useEffect(() => {
    console.log('>>> SELECTED BLUEPRINT CONNECTIONS...', blueprintConnections)
  }, [blueprintConnections])

  useEffect(() => {
    console.log('>>> SELECTED / CONFIGURING DATA SCOPE CONNECTION...', scopeConnection)
    switch (scopeConnection?.provider?.id) {
      case Providers.GITHUB:
        setActiveTransformation(scopeConnection?.transformation[scopeConnection?.projects.findIndex(p => p === configuredProject)])
        setProjects({ [configuredConnection?.id]: scopeConnection?.projects })
        break
      case Providers.JIRA:
        setActiveTransformation(scopeConnection?.transformation[scopeConnection?.boards.findIndex(b => b === `Board ${configuredBoard?.id}`)])
        setBoards({ [configuredConnection?.id]: scopeConnection?.boardsList })
        break
    }
  }, [scopeConnection, configuredProject, configuredBoard])

  useEffect(() => {
    console.log('>>> CONFIGURING / MODIFYING CONNECTION', configuredConnection)
  }, [configuredConnection])

  useEffect(() => {
    console.log('>>> ACTIVE TRANSFORMATION RULES...', activeTransformation)
  }, [activeTransformation])

  useEffect(() => {
    console.log('>>> jiraApiBoards...', jiraApiBoards)
  }, [jiraApiBoards])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <div
              className='blueprint-header'
              style={{
                display: 'flex',
                width: '100%',
                justifyContent: 'space-between',
                marginBottom: '10px',
                whiteSpace: 'nowrap',
              }}
            >
              <div className='blueprint-name' style={{}}>
                <h2
                  style={{
                    fontWeight: 'bold',
                    display: 'flex',
                    alignItems: 'center',
                  }}
                >
                  {activeBlueprint?.name}
                  <Tag
                    minimal
                    intent={
                      activeBlueprint.mode === BlueprintMode.ADVANCED
                        ? Intent.DANGER
                        : Intent.PRIMARY
                    }
                    style={{ marginLeft: '10px' }}
                  >
                    {activeBlueprint?.mode?.toString().toUpperCase()}
                  </Tag>
                </h2>
              </div>
              <div
                className='blueprint-info'
                style={{ display: 'flex', alignItems: 'center' }}
              >
                <div className='blueprint-schedule'>
                  {activeBlueprint?.isManual ? (
                    <strong>Manual Mode</strong>
                  ) : (
                    <span
                      className='blueprint-schedule-interval'
                      style={{ textTransform: 'capitalize', padding: '0 10px' }}
                    >
                      {activeBlueprint?.interval} (at{' '}
                      {dayjs(
                        getNextRunDate(activeBlueprint?.cronConfig)
                      ).format(
                        `hh:mm A ${
                          activeBlueprint?.interval !== 'Hourly'
                            ? ' MM/DD/YYYY'
                            : ''
                        }`
                      )}
                      )
                    </span>
                  )}
                  &nbsp;{' '}
                  <span className='blueprint-schedule-nextrun'>
                    {!activeBlueprint?.isManual && (
                      <>
                        Next Run{' '}
                        {dayjs(
                          getNextRunDate(activeBlueprint?.cronConfig)
                        ).fromNow()}
                      </>
                    )}
                  </span>
                </div>
                <div
                  className='blueprint-actions'
                  style={{ padding: '0 10px' }}
                >
                  {/* <Button
                    intent={Intent.PRIMARY}
                    small
                    text='Run Now'
                    onClick={runBlueprint}
                    disabled={!activeBlueprint?.enable || currentRun?.status === TaskStatus.RUNNING}
                  /> */}
                </div>
                <div className='blueprint-enabled'>
                  <Switch
                    id='blueprint-enable'
                    name='blueprint-enable'
                    checked={activeBlueprint?.enable}
                    label={
                      activeBlueprint?.enable
                        ? 'Blueprint Enabled'
                        : 'Blueprint Disabled'
                    }
                    onChange={() => handleBlueprintActivation(activeBlueprint)}
                    style={{
                      marginBottom: 0,
                      marginTop: 0,
                      color: !activeBlueprint?.enable
                        ? Colors.GRAY3
                        : 'inherit',
                    }}
                    disabled={currentRun?.status === TaskStatus.RUNNING}
                  />
                </div>
                <div style={{ padding: '0 10px' }}>
                  <Button
                    intent={Intent.PRIMARY}
                    icon='trash'
                    small
                    minimal
                    disabled
                  />
                </div>
              </div>
            </div>

            <BlueprintNavigationLinks
              blueprint={activeBlueprint}
            />
            {/* <div
              className='blueprint-navigation'
              style={{
                alignSelf: 'center',
                display: 'flex',
                margin: '20px auto',
              }}
            >
              <div style={{ marginRight: '10px' }}>
                <a
                  href='#'
                  className='blueprint-navigation-link'
                  onClick={viewBlueprintStatus}
                >
                  Status
                </a>
              </div>
              <div style={{ marginLeft: '10px' }}>
                <a
                  href='#'
                  className='blueprint-navigation-link active'
                  onClick={viewBlueprintSettings}
                >
                  Settings
                </a>
              </div>
            </div> */}

            <div
              className='blueprint-main-settings'
              style={{ display: 'flex', alignSelf: 'flex-start' }}
            >
              <div className='configure-settings-name'>
                <h3>Name</h3>
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <div className='blueprint-name'>{activeBlueprint?.name}</div>
                  <Button
                    icon='annotation'
                    intent={Intent.PRIMARY}
                    size={12}
                    small
                    minimal
                    onClick={() => modifySetting('name')}
                  />
                </div>
              </div>
              <div
                className='configure-settings-frequency'
                style={{ marginLeft: '40px' }}
              >
                <h3>Sync Frequency</h3>
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <div className='blueprint-frequency'>
                    {activeBlueprint?.isManual ? (
                      'Manual'
                    ) : (
                      <span>
                        {activeBlueprint?.interval} (at{' '}
                        {dayjs(
                          getNextRunDate(activeBlueprint?.cronConfig)
                        ).format('hh:mm A')}
                        )
                      </span>
                    )}
                  </div>
                  <Button
                    icon='annotation'
                    intent={Intent.PRIMARY}
                    size={12}
                    small
                    minimal
                    onClick={() => modifySetting('cronConfig')}
                  />
                </div>
              </div>
            </div>

            {
              activeBlueprint?.mode === BlueprintMode.NORMAL && (
                <div
                  className='data-scopes-grid'
                  style={{
                    width: '100%',
                    marginTop: '40px',
                    alignSelf: 'flex-start',
                  }}
                >
                  <h2 style={{ fontWeight: 'bold' }}>
                    Data Scope and Transformation
                  </h2>

                  <Card
                    elevation={Elevation.TWO}
                    style={{ padding: 0, minWidth: '878px' }}
                  >
                    <div
                      className='simplegrid'
                      style={{
                        display: 'flex',
                        flex: 1,
                        width: '100%',
                        flexDirection: 'column',
                      }}
                    >
                      <div
                        className='simplegrid-header'
                        style={{
                          display: 'flex',
                          flex: 1,
                          width: '100%',
                          minHeight: '48px',
                          lineHeight: 'auto',
                          padding: '16px 20px',
                          fontWeight: 'bold',
                          borderBottom: '1px solid #BDCEFB',
                          justfiyContent: 'space-evenly',
                        }}
                      >
                        <div
                          className='cell-header connections'
                          style={{ flex: 1 }}
                        >
                          Data Connections
                        </div>
                        <div
                          className='cell-header entities'
                          style={{ flex: 1 }}
                        >
                          Data Entities
                        </div>
                        <div className='cell-header scope' style={{ flex: 1 }}>
                          Data Scope
                        </div>
                        <div
                          className='cell-header transformation'
                          style={{ flex: 1 }}
                        >
                          Transformation
                        </div>
                        <div
                          className='cell-header actions'
                          style={{ minWidth: '100px' }}
                        >
                          &nbsp;
                        </div>
                      </div>

                      {connections.map((c, cIdx) => (
                        <div
                          key={`connection-row-key-${cIdx}`}
                          className='simplegrid-row'
                          style={{
                            display: 'flex',
                            flex: 1,
                            width: '100%',
                            minHeight: '48px',
                            lineHeight: 'auto',
                            padding: '10px 20px',
                            borderBottom: '1px solid #BDCEFB',
                            justfiyContent: 'space-evenly',
                          }}
                        >
                          <div className='cell connections' style={{ display: 'flex', flex: 1, alignItems: 'center' }}>
                            <span style={{ marginBottom: '-5px', marginRight: '10px' }}>
                              {c.icon}
                            </span>
                            <span>{c.name}</span>
                          </div>
                          <div className='cell entities' style={{ display: 'flex', flex: 1, alignItems: 'center' }}>
                            <ul
                              style={{
                                listStyle: 'none',
                                margin: 0,
                                padding: 0,
                              }}
                            >
                              {c.entities.map((entityLabel, eIdx) => (
                                <li key={`list-item-key-${eIdx}`}>
                                  {entityLabel}
                                </li>
                              ))}
                            </ul>
                          </div>
                          <div className='cell scope' style={{ display: 'flex', flex: 1, alignItems: 'center' }}>
                            {[Providers.GITLAB, Providers.GITHUB].includes(
                              c.provider?.id
                            ) && (
                              <ul
                                style={{
                                  listStyle: 'none',
                                  margin: 0,
                                  padding: 0,
                                }}
                              >
                                {c.projects.map((project, pIdx) => (
                                  <li key={`list-item-key-${pIdx}`}>
                                    {project}
                                  </li>
                                ))}
                              </ul>
                            )}
                            {[Providers.JIRA].includes(c.provider?.id) && (
                              <ul
                                style={{
                                  listStyle: 'none',
                                  margin: 0,
                                  padding: 0,
                                }}
                              >
                                {c.boards.map((board, bIdx) => (
                                  <li key={`list-item-key-${bIdx}`}>{board}</li>
                                ))}
                              </ul>
                            )}
                          </div>
                          <div
                            className='cell transformation'
                            style={{ display: 'flex', flex: 1, alignItems: 'center' }}
                          >
                            <ul
                              style={{
                                listStyle: 'none',
                                margin: 0,
                                padding: 0,
                              }}
                            >
                              {c.transformationStates.map((state, sIdx) => (
                                <li key={`list-item-key-${sIdx}`}>{state}</li>
                              ))}
                            </ul>
                          </div>
                          <div
                            className='cell actions'
                            style={{
                              display: 'flex',
                              minWidth: '100px',
                              textAlign: 'right',
                              alignItems: 'center',
                              justifyContent: 'flex-end',
                            }}
                          >
                            <Button
                              icon='annotation'
                              intent={Intent.PRIMARY}
                              size={12}
                              small
                              minimal
                              onClick={() => modifyConnection(cIdx, c.connectionId, c.provider)}
                            />
                          </div>
                        </div>
                      ))}
                    </div>
                  </Card>
                </div>
              )
}

            {mode === BlueprintMode.ADVANCED && (
              <div
                className='data-advanced'
                style={{
                  width: '100%',
                  maxWidth: '540px',
                  marginTop: '40px',
                  alignSelf: 'flex-start',
                }}
              >
                <h2 style={{ fontWeight: 'bold' }}>JSON Tasks Configuration</h2>

                <Card className='workflow-card' elevation={Elevation.TWO}>
                  <h4>
                    <Button
                      icon='annotation'
                      intent={Intent.PRIMARY}
                      size={12}
                      small
                      minimal
                      onClick={() => modifySetting('plan')}
                      style={{ float: 'right' }}
                    />
                    Task Editor
                  </h4>
                  <p>Modify JSON Tasks or preload from a template</p>
                  <code>
                    <pre
                      style={{
                        padding: '10px',
                        backgroundColor: '#f0f0f0',
                        borderRadius: '4px',
                      }}
                    >
                      {JSON.stringify(activeBlueprint?.plan, null, ' ')}
                    </pre>
                  </code>
                </Card>
              </div>
            )}
          </main>
        </Content>
      </div>

      <BlueprintDialog
        isOpen={blueprintDialogIsOpen}
        title={activeSetting?.title}
        blueprint={activeBlueprint}
        onSave={handleBlueprintSave}
        isSaving={isSaving}
        isValid={validateActiveSetting()}
        onClose={handleBlueprintDialogClose}
        onCancel={handleBlueprintDialogClose}
        content={(() => {
          let Settings = null
          switch (activeSetting?.id) {
            case 'name':
              Settings = (
                <BlueprintNameCard
                  name={blueprintName}
                  setBlueprintName={setBlueprintName}
                  fieBldHasError={fieldHasError}
                  getFieldError={getFieldError}
                  elevation={Elevation.ZERO}
                  enableDivider={false}
                  cardStyle={{ padding: 0 }}
                  isSaving={isSaving}
                />
              )
              break
            case 'cronConfig':
              Settings = (
                <DataSync
                  cronConfig={cronConfig}
                  customCronConfig={customCronConfig}
                  createCron={createCron}
                  setCronConfig={setCronConfig}
                  getCronPreset={getCronPreset}
                  fieldHasError={fieldHasError}
                  getFieldError={getFieldError}
                  setCustomCronConfig={setCustomCronConfig}
                  getCronPresetByConfig={getCronPresetByConfig}
                  elevation={Elevation.ZERO}
                  enableHeader={false}
                  cardStyle={{ padding: 0 }}
                />
              )
              break
          }
          return Settings
        })()}
      />

      <BlueprintDataScopesDialog
        isOpen={blueprintScopesDialogIsOpen}
        title={activeSetting?.title}
        dataEntitiesList={dataEntitiesList}
        blueprint={activeBlueprint}
        blueprintConnections={blueprintConnections}
        configuredConnection={configuredConnection}
        configuredProject={configuredProject}
        configuredBoard={configuredBoard}
        scopeConnection={scopeConnection}
        activeTransformation={activeTransformation}
        provider={activeProvider}
        boards={boards}
        boardsList={jiraApiBoards}
        projects={projects}
        issueTypesList={jiraApiIssueTypes}
        fieldsList={jiraApiFields}
        isFetching={isFetchingBlueprint}
        isFetchingJIRA={isFetchingJIRA}
        setConfiguredProject={setConfiguredProject}
        setConfiguredBoard={setConfiguredBoard}
        onSave={handleBlueprintSave}
        isSaving={isSaving}
        // @todo: validation status
        isValid={true}
        onClose={handleBlueprintScopesDialogClose}
        onCancel={handleBlueprintScopesDialogClose}
        onStepChange={handleConnectionStepChange}
        fieldHasError={fieldHasError}
        getFieldError={getFieldError}
        jiraProxyError={jiraProxyError}
      />
    </>
  )
}

export default BlueprintSettings
