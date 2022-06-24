import React from 'react'
import { Icon } from '@blueprintjs/core'
import { Providers, ProviderLabels } from '@/data/Providers'

const NullProvider = {
  id: Providers.NULL, // Unique ID, for a Provider (alphanumeric, lowercase)
  enabled: false, // Enabled Flag
  multiConnection: false, // If Provider is Multi-connection
  name: ProviderLabels.NULL, // Display Name of Data Provider
  // eslint-disable-next-line max-len
  icon: <Icon icon='box' size={30} />, // Provider Icon
  iconDashboard: <Icon icon='box' size={42} />, // Provider Icon on INTEGRATIONS Dashboard
  settings: ({ activeProvider, activeConnection, isSaving, setSettings }) => (<></>) // REACT Settings Component for Render
}

export {
  NullProvider
}
