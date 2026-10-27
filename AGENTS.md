## graphify

本项目知识图谱存放于 `graphify-out/`，包含聚合节点、模块群落结构以及跨文件依赖关系。

当用户输入 `/graphify --code-only`，执行任何操作前优先启用 graphify 相关能力。

规则：
- 询问代码库相关问题时，若存在 `graphify-out/graph.json`，优先执行 `graphify query "<question>" --code-only`；
  查询依赖链路使用 `graphify path "<A>" "<B>" --code-only`；定向概念解析使用 `graphify explain "<concept>" --code-only`。
  以上指令会返回范围收敛的子图谱，相比完整 GRAPH_REPORT.md 或全局文本检索体量更小。
- 钩子触发增量更新后，`graphify-out/` 内出现临时脏文件属于正常现象，不能以此为由跳过 graphify；
  只有任务目标是排查图谱过时、图谱数据错误，或是用户明确禁止使用图谱时，才可跳过。
- 如果存在 `graphify-out/wiki/index.md`，优先依靠该文档进行整体导航，而非直接浏览源码。
- 仅在全局架构评审、query/path/explain 无法获取足够信息时，再读取 `graphify-out/GRAPH_REPORT.md`。
- 修改代码完成后，执行 `graphify update . --code-only` 同步更新图谱（仅解析AST，无额外接口开销）。
- Monorepo Vue 项目中，出现节点数量为0的JSON文件、Vue节点去重相关警告属于正常现象，**不要视作构建错误**，正常继续使用图谱。
