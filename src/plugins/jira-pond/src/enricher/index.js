require('module-alias/register')

const issueCollector = require('../collector/issues')
const changelogCollector = require('../collector/changelogs')

const closedStatuses = ['Done', 'Closed', '关闭', '完全的', '已关闭']

module.exports = {
  async enrich(rawDb, enrichedDb, projectId) {
    console.log('Jira Enrichment', projectId)
    await module.exports.enrichLeadTimeOnIssues(
      rawDb,
      enrichedDb,
      projectId
    )
    console.log('Done enriching issues')
  },

  async enrichLeadTimeOnIssues(rawDb, enrichedDb, projectId) {
    const {
      JiraIssue
    } = enrichedDb

    const issues = await issueCollector.findIssues({
      'fields.project.id': `${projectId}`
    }, rawDb)

    const creationPromises = []
    const leadTimePromises = []
    const issuesToCreate = []
    issues.forEach(async issue => {
      leadTimePromises.push(module.exports.calculateLeadTime(issue, rawDb))
      issuesToCreate.push({
        id: issue.id,
        url: issue.self,
        title: issue.fields.summary,
        projectId: issue.fields.project.id
        // description: issue.fields.description
      })
    })

    const leadTimes = await Promise.all(leadTimePromises)

    leadTimes.forEach((leadTime, index) => {
      console.log('JON >>> leadTime', leadTime)
      let issue = issuesToCreate[index]
      issue = {
        leadTime,
        ...issue
      }
      creationPromises.push(JiraIssue.findOrCreate({
        where: {
          id: issue.id
        },
        defaults: issue
      }))
    })

    await Promise.all(creationPromises)
  },

  async calculateLeadTime(issue, db) {
    const changelogs = await changelogCollector.findChangelogs({
      issueId: `${issue.id}`
    }, db)

    let leadTime = 0
    let lastTime = new Date(issue.fields.created).getTime()
    let isDone = false

    for (const change of changelogs) {
      for (const item of change.items) {
        if (item.field === 'status') {
          console.log('JON >>> item.field', item.field)
          const changeTime = new Date(change.created).getTime()

          if (!closedStatuses.includes(item.fromString)) {
            const elapsedTime = changeTime - lastTime

            leadTime += elapsedTime
          }

          lastTime = changeTime
          isDone = closedStatuses.includes(item.toString)
        }
      }
    }

    return isDone ?
      Math.round(leadTime / 1000) :
      0
  }
}