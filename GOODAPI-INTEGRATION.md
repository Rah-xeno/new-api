# GoodAPI 前端接入说明

## 交付范围

最初以 `Rah-xeno/new-api` 的 main 源码快照为基础，现已整合 dev 分支的既有功能，迁入用户提供的
`goodapinew-main.zip` 中的界面，接入位置为 `web/default`。

- 首页：球体视频、静态海报、移动端导航、模型展示与页面过渡。
- 登录和注册：GoodAPI 的分栏布局与样式，继续使用 new-api 的登录、注册、OAuth、2FA 和 Cookie 会话流程。
- 控制台：顶栏、侧栏、表格、弹窗、配色、概览和余额用量卡片。
- 模型价格页：模型卡片、分组筛选、搜索、表格模式、价格详情与性能展示。
- 保留 new-api 的渠道、令牌、用户、钱包、订阅、日志及管理功能。
- 遵循 dev 分支仅运行 default 新版前端的设计，保留 dev 的代理分销、邀请计划、管理统计与资源加载恢复功能。

这次是前端界面移植与接口适配，没有替换 new-api 的后端业务系统。
GoodAPI 独有的 JWT 刷新令牌、安全验证、细粒度管理员权限、插件市场、
周奖励等后端扩展没有迁入，也没有添加会调用这些缺失接口的菜单入口。
这些模块若要一并提供，需要另做后端功能合并。

## 兼容处理

1. 沿用 new-api 的认证存储、请求封装和业务路由；不调用 GoodAPI 的 `/api/user/auth/refresh`。
2. 移动端账号安全入口仍归入 new-api 的个人资料页。
3. 性能卡片读取 new-api 的 `recent_success_rates` 数组；没有数据时保留空状态，不伪造时间戳或性能数据。
4. 后台文档和导航继续使用站点配置，不强制跳转 GoodAPI 的文档站。
5. 保留 workspace、dev 的默认主题 Docker 构建以及版权和许可证；首页同时显示上游归属。
6. 在前端 workspace 固定 `date-fns` 2.x 的根依赖，解决 classic 的 `date-fns-tz` 1.x 被错误提升到 4.x 的兼容问题。
7. 修正导入概览页和宿主布局重复嵌套的问题，避免重复标题与内部滚动容器。

## 已有站点切换

dev 分支已统一使用 default 新版前端，本次继续采用该行为。

后台如果设置了自定义首页内容或首页 URL，会继续优先显示该内容。
若要显示这次的球体首页，请清空自定义首页配置。

## 构建

推荐使用 Bun 1.4.2（本次验证版本），Go 1.26.1（本次验证版本）。

```sh
cd web
bun install --frozen-lockfile
cd default
bun run typecheck
bun run test
bun run build
cd ../..
go build -o new-api .
```

Git 仓库不跟踪构建产物；必须先构建 web/default/dist，再构建 Go 程序。此前交付的 ZIP 基于 main，不能替代整合了 dev 功能的当前分支。
修改前端后必须重新构建前端，再重新构建 Go；静态资源被嵌入二进制。

原 `docker-compose.yml` 使用上游公开镜像，不包含本次修改。
新增 `docker-compose.goodapi.yml` 用于从当前源码构建集成版本：

```sh
# 只构建镜像
docker compose -f docker-compose.yml -f docker-compose.goodapi.yml build new-api

# 决定发布后，再使用相同的两个配置文件启动服务
docker compose -f docker-compose.yml -f docker-compose.goodapi.yml up -d
```

已有部署请沿用自己的数据库、卷、域名和环境配置，把应用镜像换为本源码构建的版本。
本次没有执行部署或修改线上数据库。

## 验证

- 新版前端 TypeScript 类型检查通过。
- 新版和 classic 生产构建通过。
- 7 个前端测试文件、54 项测试通过，包含首页导航、动效降级、登录布局和原有交互测试。
- 本次修改的前端代码 lint 通过，并格式化；仓库其他未修改文件仍有上游已有的 lint/格式问题。
- Go 程序构建通过，`go test ./common ./setting/system_setting` 通过（后者无测试文件）。
- 真实本地 Go + SQLite：初始化、登录、用户信息、余额用量、API 密钥列表、模型价格与详情均已检查。
- 浏览器 1440px 桌面和 390px 手机视口检查；首页、价格页手机宽度无横向溢出，检查页面无控制台 error。
- 上游模型调用、真实支付、第三方 OAuth 和完整 Docker 镜像运行未测试。

测试使用独立临时 SQLite 和无效测试渠道密钥；这些数据库和凭据不在交付包中。

## 来源

- 基础仓库：https://github.com/Rah-xeno/new-api
- 界面来源：https://github.com/ljyoukong-cpu/goodapinew ，使用用户提供的 ZIP。
- 上游：https://github.com/QuantumNous/new-api
- 基础源码 ZIP SHA-256：`3AB6D94E5C7688D8548B6DAEFEE3854BEF4DB17821B4547D6CF29B0C639DDD89`
- 用户 ZIP SHA-256：`BD68793C53C12C74B19297B995425D773C50D2ADE365D398A1D1CFF2216B3A91`

许可证和第三方声明见仓库原有的 LICENSE、NOTICE 和 THIRD-PARTY-LICENSES.md。

## dev 合并说明

保留 dev 原有的后端、邀请计划、代理分销、管理统计、注册接口、账号字段、资源加载恢复和新版专用构建。翻译按键合并；双方都修改的既有键优先保留 dev 业务文案，GoodAPI 新增界面键保留。验证结果以 PR 的最新合并说明为准。
