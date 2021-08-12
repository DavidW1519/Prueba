export type PrimaryValues = Record<string, string | number | boolean>;

interface Collector {
  /**
   * 当前插件的名字，如果未继承默认为去掉后缀的类名
   * name of plugin
   * If you do not inherit name(),  it defaults to the class name removing the suffix
   */
  name?(): string;

  /**
   * 返回改依赖项，collector只能依赖collector
   * return dependencies
   * dependencies could only include collector
   * 返回格式(return format):
   * {
   *   # depend on collector in this plugin
   *   [collectorName: string]: primaryKeyObject
   *   # or depend on collector in other plugin
   *   'pluginName' + '/' + 'collectorName': primaryKeyObject
   * }
   */
  dependencies?(primaryKeys: PrimaryValues): Record<string, PrimaryValues>;

  /**
   * 返回对应主键的数据是否已经准备好，需要快速的返回结果
   * get if data have prepared for these primaryKeys
   * it should return result as soon as.
   * @param primaryKeys
   */
  isDataPrepared(primaryKeys: PrimaryValues): Promise<boolean>;

  /**
   * 使用队列开始导入数据
   * start collect data within queue
   * @param primaryKeys
   * @param consumerModule
   */
  collectData(primaryKeys: PrimaryValues, consumerModule): Promise<void>;

  /**
   * 清除指定主键对应的数据
   * clean data for these primaryKeys
   * @param primaryKeys
   */
  cleanData(primaryKeys?: PrimaryValues): Promise<boolean>;
}

export default Collector;
