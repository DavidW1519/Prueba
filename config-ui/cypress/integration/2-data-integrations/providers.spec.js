/// <reference types="cypress" />

context('Data Integration Providers', () => {
  beforeEach(() => {
    cy.visit('/integrations')
  })

  describe('JIRA Data Provider', () => {
    it('provides access to jira integration', () => {
      cy.visit('/integrations/jira')
      cy.get('.headlineContainer')
        .find('h1')
        .contains(/jira integration/i)
    })

    it('displays connection sources data table', () => {
      cy.visit('/integrations/jira')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('thead')
        .find('th')
        .should('contain', 'ID')
        .should('contain', 'Connection Name')
        .should('contain', 'Endpoint')
        .should('contain', 'Status')
    })

    it('displays add connection button', () => {
      cy.visit('/integrations/jira')
      cy.get('button.bp3-button').contains('Add Connection')
        .should('be.visible')
    })

    it('displays refresh connections button', () => {
      cy.visit('/integrations/jira')
      cy.get('button.bp3-button').contains('Refresh Connections')
        .should('be.visible')
    })

    it('can create a new jira connection source', () => {
      cy.fixture('new-jira-connection').as('JIRAConnectionSourceJSON')
      cy.intercept('POST', '/api/plugins/jira/sources', { statusCode: 201, body: '@JIRAConnectionSourceJSON' }).as('createJIRASource')
      cy.visit('/integrations/jira')
      cy.get('button#btn-add-new-connection').click()
      cy.get('button#btn-save')
        .should('be.visible')
        .should('be.disabled')

      cy.get('input#connection-name').type('TEST JIRA INSTANCE')
      cy.get('input#connection-endpoint').type('https://test-46f2c29a-2955-4fa6-8de8-fffeaf8cf8e0.atlassian.net/rest/')
      cy.get('input#connection-token').type('xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx')
      cy.get('input#connection-proxy').type('http://proxy.localhost:8800')

      cy.get('button#btn-save')
        .should('be.visible')
        .should('be.enabled')
        .click()
      
      cy.wait('@createJIRASource').its('response.statusCode').should('eq', 201)
      cy.url().should('include', '/integrations/jira')
    })
  })

  describe('GitLab Data Provider', () => {
    it('provides access to gitlab integration', () => {
      cy.visit('/integrations/gitlab')
      cy.get('.headlineContainer')
        .find('h1')
        .contains(/gitlab integration/i)
    })
    it('displays connection sources data table', () => {
      cy.visit('/integrations/gitlab')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('thead')
        .find('th')
        .should('contain', 'Connection Name')
        .should('contain', 'Endpoint')
        .should('contain', 'Status')
    })
    it('limited to one (1) connection source', () => {
      cy.visit('/integrations/gitlab')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('tbody').should('have.length', 1)
    })
    it('cannot add a new connection', () => {
      cy.visit('/integrations/gitlab')
      cy.get('button.bp3-button').contains('Add Connection')
        .parent()
        .should('have.class', 'bp3-disabled')
        .should('have.attr', 'disabled')
    })
  })

  describe('GitHub Data Provider', () => {
    it('provides access to github integration', () => {
      cy.visit('/integrations/github')
      cy.get('.headlineContainer')
        .find('h1')
        .contains(/github integration/i)
    })
    it('displays connection sources data table', () => {
      cy.visit('/integrations/github')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('thead')
        .find('th')
        .should('contain', 'Connection Name')
        .should('contain', 'Endpoint')
        .should('contain', 'Status')
    })
    it('limited to one (1) connection source', () => {
      cy.visit('/integrations/github')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('tbody').should('have.length', 1)
    })
    it('cannot add a new connection', () => {
      cy.visit('/integrations/github')
      cy.get('button.bp3-button').contains('Add Connection')
        .parent()
        .should('have.class', 'bp3-disabled')
        .should('have.attr', 'disabled')
    })
  })

  describe('Jenkins Data Provider', () => {
    it('provides access to jenkins integration', () => {
      cy.visit('/integrations/jenkins')
      cy.get('.headlineContainer')
        .find('h1')
        .contains(/jenkins integration/i)
    })
    it('displays connection sources data table', () => {
      cy.visit('/integrations/jenkins')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('thead')
        .find('th')
        .should('contain', 'Connection Name')
        .should('contain', 'Endpoint')
        .should('contain', 'Status')
    })
    it('limited to one (1) connection source', () => {
      cy.visit('/integrations/jenkins')
      cy.get('.connections-table')
        .should('have.class', 'bp3-html-table')
        .should('be.visible')
        .find('tbody').should('have.length', 1)
    })
    it('cannot add a new connection', () => {
      cy.visit('/integrations/jenkins')
      cy.get('button.bp3-button').contains('Add Connection')
        .parent()
        .should('have.class', 'bp3-disabled')
        .should('have.attr', 'disabled')
    })
  })
})