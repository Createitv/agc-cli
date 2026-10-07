## 问题与结果 / Problem and outcome

<!-- 说明触发场景、之前的行为、改后行为。关联 Issue：Closes #... -->
<!-- Describe the trigger, previous behavior, and resulting behavior. Link an issue if applicable. -->

## 改动 / Changes

<!-- 列出帮助审查的关键改动；说明是否影响命令、参数、JSON 输出或默认行为。 -->
<!-- Summarize review-relevant changes. Note changes to commands, flags, JSON output, or defaults. -->

## 兼容性 / Compatibility

<!-- 无破坏性变更时写“无 / None”；否则给出旧命令 → 新命令、配置迁移或输出变化。 -->
<!-- For breaking changes, show old → new usage, configuration migration, or output changes. -->

- CLI / REST：
- 配置与凭据 / Configuration and credentials：
- 输出格式 / Output schema：

## 验证 / Validation

<!-- 写实际运行的命令和结果。区分本地模拟、真实 API 请求与远程状态确认。 -->
<!-- List commands actually run and their results. Separate mocks, live API calls, and remote state checks. -->

- 检查命令及结果 / Checks and results：
- 尚未验证或不适用 / Not verified or not applicable：

## 文档与完成条件 / Documentation and completion

<!-- 仅勾选已完成项；不适用的项目写明原因。纯文档 PR 可运行 make check-docs。 -->
<!-- Check completed items only; explain items that do not apply. Docs-only PRs can run make check-docs. -->

- [ ] 改动范围明确，无无关文件 / Focused scope, no unrelated files
- [ ] 行为变更有相关测试；纯文档改动注明不适用 / Behavior changes have relevant tests; mark docs-only as N/A
- [ ] 共享能力已检查 CLI 与 REST / CLI and REST checked for shared behavior
- [ ] 参数、JSON 输出和默认行为的兼容性已说明 / Flag, JSON, and default compatibility documented
- [ ] 受影响的中文与英文文档同步 / Affected Chinese and English docs updated
- [ ] 用户可见变化已写入 CHANGELOG / User-visible changes added to CHANGELOG
- [ ] 已运行适用检查，未验证项已说明 / Relevant checks run, unverified areas disclosed
- [ ] 示例和日志已脱敏 / Examples and logs contain no secrets
