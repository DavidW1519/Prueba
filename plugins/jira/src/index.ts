import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import JiraPlugin from './JiraPlugin';
import { ScheduleModule } from '@nestjs/schedule';
import IssueCollector from './runners/IssueCollector';
import { ConfigModule, ConfigService } from '@nestjs/config';

@Module({
  imports: [
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      useFactory: (config: ConfigService) => ({
        type: config.get<'postgres' | 'mysql'>('DB_TYPE', 'mysql'),
        url: config.get<string>('DB_URL'),
        name: 'jiraModuleDb',
        entityPrefix: 'plugin_jira_',
        entities: [],
      }),
      inject: [ConfigService],
    }),
    ScheduleModule.forRoot(),
  ],
  providers: [JiraPlugin, IssueCollector],
  exports: [JiraPlugin],
})
class JiraModule {}

export default JiraModule;
