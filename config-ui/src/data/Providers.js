import React from 'react'
import { Tooltip, Icon } from '@blueprintjs/core'
import { ReactComponent as GitlabProviderIcon } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProviderIcon } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProviderIcon } from '@/images/integrations/jira.svg'
import { ReactComponent as GitHubProviderIcon } from '@/images/integrations/github.svg'
// import GitExtractorIcon from '@/images/git.png'
// import RefDiffIcon from '@/images/git-diff.png'
import FeishuIcon from '@/images/feishu.png'
// import DBTIcon from '@/images/dbt.png'
// import AEIcon from '@/images/ae.png'

const Providers = {
  NULL: 'null',
  GITLAB: 'gitlab',
  JENKINS: 'jenkins',
  JIRA: 'jira',
  GITHUB: 'github',
  REFDIFF: 'refdiff',
  GITEXTRACTOR: 'gitextractor',
  FEISHU: 'feishu',
  AE: 'ae',
  DBT: 'dbt'
}

const ProviderTypes = {
  PLUGIN: 'plugin',
  INTEGRATION: 'integration',
  PIPELINE: 'pipeline'
}

const ProviderLabels = {
  NULL: 'NullProvider',
  GITLAB: 'GitLab',
  JENKINS: 'Jenkins',
  JIRA: 'JIRA',
  GITHUB: 'GitHub',
  REFDIFF: 'RefDiff',
  GITEXTRACTOR: 'GitExtractor',
  FEISHU: 'Feishu',
  AE: 'Analysis Engine (AE)',
  DBT: 'Data Build Tool (DBT)'
}

const ProviderConnectionLimits = {
  gitlab: 1,
  jenkins: 1,
  // jira: null, // (Multi-connection, no-limit)
  github: 1
}

// NOTE: Not all fields may be referenced/displayed for a provider,
// ie. JIRA prefers $token over $username and $password
const ProviderFormLabels = {
  null: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password'
  },
  gitlab: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password'
  },
  jenkins: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password'
  },
  jira: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    token: 'Basic Auth Token',
    username: 'Username',
    password: 'Password'
  },
  github: {
    name: 'Connection Name',
    endpoint: 'Endpoint URL',
    proxy: 'Proxy URL',
    // token: 'Auth Token(s)',
    token: (
      <>
        <Tooltip
          content={(<span>Due to Github’s rate limit, input more tokens, <br />comma separated, to accelerate data collection.</span>)}
          intent='primary'
        >
          <Icon
            icon='info-sign'
            size={12}
            style={{
              float: 'left',
              display: 'inline-block',
              alignContent: 'center',
              marginBottom: '4px',
              marginRight: '8px',
              color: '#999'
            }}
          />
        </Tooltip>
        Auth Token(s)
      </>),
    username: 'Username',
    password: 'Password'
  },
}

const ProviderFormPlaceholders = {
  null: {
    name: 'eg. Enter Instance Name',
    endpoint: 'eg. https://null-api.localhost',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 3f5cda2a23ff410792e0',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  },
  gitlab: {
    name: 'eg. GitLab',
    endpoint: 'eg. https://gitlab.com/api/v4',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. ff9d1ad0e5c04f1f98fa',
    username: 'Enter Username / E-mail',
    password: 'Enter Password'
  },
  jenkins: {
    name: 'eg. Jenkins',
    endpoint: 'URL eg. https://api.jenkins.io',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 6b057ffe68464c93a057',
    username: 'eg. admin',
    password: 'eg. ************'
  },
  jira: {
    name: 'eg. JIRA',
    endpoint: 'eg. https://your-domain.atlassian.net/rest/',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 8c06a7cc50b746bfab30',
    username: 'eg. admin',
    password: 'eg. ************'
  },
  github: {
    name: 'eg. GitHub',
    endpoint: 'eg. https://api.github.com',
    proxy: 'eg. http://proxy.localhost:8080',
    token: 'eg. 4c5cbdb62c165e2b3d18, 40008ebccff9837bb8d2',
    username: 'eg. admin',
    password: 'eg. ************'
  }
}

const ProviderIcons = {
  [Providers.GITLAB]: (w, h) => <GitlabProviderIcon width={w || 24} height={h || 24} />,
  [Providers.JENKINS]: (w, h) => <JenkinsProviderIcon width={w || 24} height={h || 24} />,
  [Providers.JIRA]: (w, h) => <JiraProviderIcon width={w || 24} height={h || 24} />,
  [Providers.GITHUB]: (w, h) => <GitHubProviderIcon width={w || 24} height={h || 24} />,
  [Providers.REFDIFF]: (w, h) => <Icon icon='box' size={w || 24} />,
  [Providers.GITEXTRACTOR]: (w, h) => <Icon icon='box' size={w || 24} />,
  [Providers.FEISHU]: (w, h) => <img src={FeishuIcon} width={w || 24} height={h || 24} />,
  [Providers.AE]: (w, h) => <Icon icon='box' size={w || 24} />,
  [Providers.DBT]: (w, h) => <Icon icon='box' size={w || 24} />,
}

export {
  Providers,
  ProviderTypes,
  ProviderIcons,
  ProviderLabels,
  ProviderConnectionLimits,
  ProviderFormLabels,
  ProviderFormPlaceholders
}
