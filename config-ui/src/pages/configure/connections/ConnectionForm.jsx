import React, { useEffect, useState } from 'react'
import {
  Button, Colors,
  FormGroup, InputGroup, Label,
  Card,
  Icon,
  Tag,
  Elevation,
  Popover,
  Position,
  Intent
} from '@blueprintjs/core'

import GenerateTokenForm from '@/pages/configure/connections/GenerateTokenForm'

import '@/styles/integration.scss'
import '@/styles/connections.scss'
import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function ConnectionForm (props) {
  const {
    isLocked = false,
    activeProvider,
    name,
    endpointUrl,
    token,
    username,
    password,
    isSaving,
    isTesting,
    showError,
    errors,
    testStatus,
    onSave = () => {},
    onCancel = () => {},
    onTest = () => {},
    onNameChange = () => {},
    onEndpointChange = () => {},
    onTokenChange = () => {},
    onUsernameChange = () => {},
    onPasswordChange = () => {},
    authType = 'token',
    sourceLimits = { jenkins: 1, gitlab: 1 },
    showLimitWarning = true
  } = props

  const [allowedAuthTypes, setAllowedAuthTypes] = useState(['token', 'plain'])
  const [showTokenCreator, setShowTokenCreator] = useState(false)

  const getConnectionStatusIcon = () => {
    let statusIcon = <Icon icon='full-circle' size='10' color={Colors.RED3} />
    switch (testStatus) {
      case 1:
        statusIcon = <Icon icon='full-circle' size='10' color={Colors.GREEN3} />
        break
      case 2:
        statusIcon = <Icon icon='full-circle' size='10' color={Colors.RED3} />
        break
      case 0:
      default:
        statusIcon = <Icon icon='full-circle' size='10' color={Colors.GRAY3} />
        break
    }
    return statusIcon
  }

  const handleTokenInteraction = (isOpen) => {
    setShowTokenCreator(isOpen)
  }

  useEffect(() => {
    if (!allowedAuthTypes.includes(authType)) {
      console.log('INVALID AUTH TYPE!')
    }
  }, [authType, allowedAuthTypes])

  useEffect(() => {
    setAllowedAuthTypes(['token', 'plain'])
  }, [])

  return (
    <>
      <form className='form form-add-connection'>
        <div className='headlineContainer'>
          <h2 className='headline' style={{ marginTop: 0, textDecoration: isLocked ? 'line-through' : 'none' }}>Configure Connection</h2>
          <p className='description'>Instance Account & Authentication settings</p>
          {activeProvider && activeProvider.id && sourceLimits[activeProvider.id] && showLimitWarning && (
            <Card interactive={false} elevation={Elevation.TWO} style={{ width: '50%', marginBottom: '20px', backgroundColor: '#f0f0f0' }}>
              <p className='warning-message' intent={Intent.WARNING}>
                <Icon icon='warning-sign' size='16' color={Colors.GRAY1} style={{ marginRight: '5px' }} />
                <strong>CONNECTION SOURCES LIMITED</strong><br />
                You may only add <Tag intent={Intent.PRIMARY}>{sourceLimits[activeProvider.id]}</Tag> instance(s) at this time,
                multiple sources will be supported in a future release.
              </p>
            </Card>
          )}
        </div>

        {showError && (
          <div className='bp3-callout bp3-intent-danger' style={{ margin: '20px 0', maxWidth: '50%' }}>
            <h4 className='bp3-heading'>Operation Failed</h4>
            Your connection could not be saved.
            {errors.length > 0 && (
              <ul>
                {errors.map((errorMessage, idx) => (
                  <li key={`save-error-message-${idx}`}>{errorMessage}</li>
                ))}
              </ul>)}
          </div>)}

        <div className='formContainer'>
          <FormGroup
            disabled={isTesting || isSaving || isLocked}
            readOnly={['gitlab', 'jenkins'].includes(activeProvider.id)}
            label=''
            inline={true}
            labelFor='connection-name'
            helperText='NAME'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label style={{ display: 'inline' }}>
              Connection&nbsp;Name <span className='requiredStar'>*</span>
            </Label>
            <InputGroup
              id='connection-name'
              disabled={isTesting || isSaving || isLocked}
              readOnly={['gitlab', 'jenkins'].includes(activeProvider.id)}
              placeholder='Enter Instance Name eg. ISSUES-AWS-US-EAST'
              value={name}
              onChange={(e) => onNameChange(e.target.value)}
              className='input connection-name-input'
              leftIcon={['gitlab', 'jenkins'].includes(activeProvider.id) ? 'lock' : null}
              fill
            />
          </FormGroup>
        </div>

        <div className='formContainer'>
          <FormGroup
            disabled={isTesting || isSaving || isLocked}
            label=''
            inline={true}
            labelFor='connection-endpoint'
            helperText='ENDPOINT_URL'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label>
              Endpoint&nbsp;URL <span className='requiredStar'>*</span>
            </Label>
            <InputGroup
              id='connection-endpoint'
              disabled={isTesting || isSaving || isLocked}
              placeholder='Enter Endpoint URL eg. https://merico.atlassian.net/rest'
              value={endpointUrl}
              onChange={(e) => onEndpointChange(e.target.value)}
              className='input'
              fill
            />
            {/* <a href='#' style={{ margin: '5px 0 5px 5px' }}><Icon icon='info-sign' size='16' /></a> */}
          </FormGroup>
        </div>

        {authType === 'token' && (
          <div className='formContainer'>
            <FormGroup
              disabled={isTesting || isSaving || isLocked}
              label=''
              inline={true}
              labelFor='connection-token'
              helperText='TOKEN'
              className='formGroup'
              contentClassName='formGroupContent'
            >
              <Label>
                Basic&nbsp;Auth&nbsp;Token <span className='requiredStar'>*</span>
              </Label>
              <InputGroup
                id='connection-token'
                disabled={isTesting || isSaving || isLocked}
                placeholder='Enter Auth Token eg. EJrLG8DNeXADQcGOaaaX4B47'
                value={token}
                onChange={(e) => onTokenChange(e.target.value)}
                className='input'
                fill
                required
              />
              <Popover
                className='popover-generate-token'
                position={Position.RIGHT}
                autoFocus={false}
                enforceFocus={false}
                isOpen={showTokenCreator}
                onInteraction={handleTokenInteraction}
                onClosed={() => setShowTokenCreator(false)}
                usePortal={false}
              >
                <Button
                  disabled={isTesting || isSaving || isLocked}
                  type='button' icon='key' intent={Intent.PRIMARY} style={{ marginLeft: '5px' }}
                />
                <>
                  <div style={{ padding: '15px 20px 15px 15px' }}>
                    <GenerateTokenForm
                      isTesting={isTesting}
                      isSaving={isSaving}
                      isLocked={isLocked}
                      onTokenChange={onTokenChange}
                      setShowTokenCreator={setShowTokenCreator}
                    />
                  </div>
                </>
              </Popover>
              {/* <a href='#' style={{ margin: '5px 0 5px 5px' }}><Icon icon='info-sign' size='16' /></a> */}
            </FormGroup>
          </div>
        )}
        {authType === 'plain' && (
          <>
            <div style={{ marginTop: '20px', marginBottom: '20px' }}>
              <h3 style={{ margin: 0 }}>Username & Password</h3>
              <span className='description' style={{ margin: 0, color: Colors.GRAY2 }}>
                If this connection uses login credentials to generate a token or uses PLAIN Auth, specify it here.
              </span>
            </div>
            <div className='formContainer'>
              <FormGroup
                label=''
                disabled={isTesting || isSaving || isLocked}
                inline={true}
                labelFor='connection-username'
                helperText='USERNAME'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <Label style={{ display: 'inline' }}>
                  Username <span className='requiredStar'>*</span>
                </Label>
                <InputGroup
                  id='connection-username'
                  disabled={isTesting || isSaving || isLocked}
                  placeholder='Enter Username'
                  defaultValue={username}
                  onChange={(e) => onUsernameChange(e.target.value)}
                  className='input'
                  style={{ maxWidth: '300px' }}
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isTesting || isSaving || isLocked}
                label=''
                inline={true}
                labelFor='connection-password'
                helperText='PASSWORD'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <Label style={{ display: 'inline' }}>
                  Password <span className='requiredStar'>*</span>
                </Label>
                <InputGroup
                  id='connection-password'
                  type='password'
                  disabled={isTesting || isSaving || isLocked}
                  placeholder='Enter Password'
                  defaultValue={password}
                  onChange={(e) => onPasswordChange(e.target.value)}
                  className='input'
                  style={{ maxWidth: '300px' }}
                />
              </FormGroup>
            </div>
          </>
        )}
        <div
          className='form-actions-block'
          style={{ display: 'flex', marginTop: '30px', justifyContent: 'space-between' }}
        >
          <div style={{ display: 'flex' }}>
            {/* <Button
              className='btn-test-connection'
              icon={getConnectionStatusIcon()}
              text='Test Connection'
              onClick={onTest}
              loading={isTesting}
              disabled={isTesting || isSaving || isLocked}
            /> */}
          </div>
          <div style={{ display: 'flex' }}>
            <Button className='btn-cancel' icon='remove' text='Cancel' onClick={onCancel} disabled={isSaving || isTesting} />
            <Button
              className='btn-save'
              icon='cloud-upload' intent='primary' text='Save Connection'
              loading={isSaving}
              disabled={isSaving || isTesting || isLocked}
              onClick={onSave}
            />
          </div>
        </div>
      </form>
    </>
  )
}
