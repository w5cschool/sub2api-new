export default {
  batchImageGuide: {
    title: '图片批量生成',
    description: '一次提交多条提示词，任务完成后可统一下载图片结果'
  },
  // Home Page
  home: {
    viewOnGithub: '在 GitHub 上查看',
    viewDocs: '查看文档',
    docs: '文档',
    switchToLight: '切换到浅色模式',
    switchToDark: '切换到深色模式',
    dashboard: '控制台',
    login: '登录',
    getStarted: '立即开始',
    goToDashboard: '进入控制台',
    // 新增：面向用户的价值主张
    heroSubtitle: '一个密钥，畅用多个 AI 模型',
    heroDescription: '无需管理多个订阅账号，一站式接入 Claude、GPT、Gemini 等主流 AI 服务',
    tags: {
      subscriptionToApi: '订阅转 API',
      stickySession: '会话保持',
      realtimeBilling: '按量计费'
    },
    // 用户痛点区块
    painPoints: {
      title: '你是否也遇到这些问题？',
      items: {
        expensive: {
          title: '订阅费用高',
          desc: '每个 AI 服务都要单独订阅，每月支出越来越多'
        },
        complex: {
          title: '多账号难管理',
          desc: '不同平台的账号、密钥分散各处，管理起来很麻烦'
        },
        unstable: {
          title: '服务不稳定',
          desc: '单一账号容易触发限制，影响正常使用'
        },
        noControl: {
          title: '用量无法控制',
          desc: '不知道钱花在哪了，也无法限制团队成员的使用'
        }
      }
    },
    // 解决方案区块
    solutions: {
      title: '我们帮你解决',
      subtitle: '简单三步，开始省心使用 AI'
    },
    features: {
      unifiedGateway: '一键接入',
      unifiedGatewayDesc: '获取一个 API 密钥，即可调用所有已接入的 AI 模型，无需分别申请。',
      multiAccount: '稳定可靠',
      multiAccountDesc: '智能调度多个上游账号，自动切换和负载均衡，告别频繁报错。',
      balanceQuota: '用多少付多少',
      balanceQuotaDesc: '按实际使用量计费，支持设置配额上限，团队用量一目了然。'
    },
    // 优势对比
    comparison: {
      title: '为什么选择我们？',
      headers: {
        feature: '对比项',
        official: '官方订阅',
        us: '本平台'
      },
      items: {
        pricing: {
          feature: '付费方式',
          official: '固定月费，用不完也付',
          us: '按量付费，用多少付多少'
        },
        models: {
          feature: '模型选择',
          official: '单一服务商',
          us: '多模型随意切换'
        },
        management: {
          feature: '账号管理',
          official: '每个服务单独管理',
          us: '统一密钥，一站管理'
        },
        stability: {
          feature: '服务稳定性',
          official: '单账号易触发限制',
          us: '多账号池，自动切换'
        },
        control: {
          feature: '用量控制',
          official: '无法限制',
          us: '可设配额、查明细'
        }
      }
    },
    providers: {
      title: '已支持的 AI 模型',
      description: '一个 API，多种选择',
      supported: '已支持',
      soon: '即将推出',
      claude: 'Claude',
      gemini: 'Gemini',
      antigravity: 'Antigravity',
      more: '更多'
    },
    // CTA 区块
    cta: {
      title: '准备好开始了吗？',
      description: '注册即可获得免费试用额度，体验一站式 AI 服务',
      button: '免费注册'
    },
    footer: {
      allRightsReserved: '保留所有权利。'
    },
    onlyCodex: {
      hero: {
        eyebrow: '专注 Codex · GPT-5.6 / GPT-5.5 / GPT-5.4 稳定中转服务',
        titlePrefix: '让团队更顺畅地接入',
        description: 'OnlyCodex，提供高速、稳定的 Codex 中转服务代理能力，让团队更顺畅地使用 GPT-5.5 与 GPT-5.4 完成开发工作。',
        note: '无需魔法，无需官方账号，无需国外信用卡/手机号。'
      },
      actions: { signUpAndTry: '免费注册，立即体验', viewSetupGuide: '查看配置教程', contactNow: '立即咨询', signUpNow: '立即注册', signUpFree: '免费注册', copy: '复制', close: '关闭' },
      terminal: { connected: '✓ 已连接至 OnlyCodex 中转节点', optimizePrompt: '帮我优化这段代码的性能...', refactorPrompt: '帮我重构这个模块的...' },
      advantages: {
        eyebrow: '核心优势', titlePrefix: '为什么开发者', titleHighlight: '选择我们', titleSuffix: '？',
        global: { title: '全球布局', description: '部署线路服务器，自动负载均衡确保快速响应。', tag1: '全球用户快速响应', tag2: '自动负载均衡' },
        speed: { title: '极致速度', description: 'CN2 GIA 高质量链路，跨境 RTT 最低 60 ms，带宽 200 Mbps 弹性扩容，高峰依然流畅。', tag: '60ms 低延迟' },
        reliability: { title: '稳定如磐', description: '亚太两地互备，5s 内自动切换，SLA 99.99%，全年宕机 < 52 分钟。', tag: '双地互备' },
        security: { title: '数据安全', description: '传输采用 HTTPS 加密。数据平面与管理平面分离，不截留、不分析用户数据。关键数据 AES 算法加密。', tag1: 'HTTPS 加密', tag2: 'AES 加密' }
      },
      usagePricing: {
        eyebrow: '按量计费 · 价格对比', titlePrefix: '官方价格，OnlyCodex 仅需 ', titleHighlight: '1.25 折', description: '每 1M Tokens；上方为官方按量价格，下方为 OnlyCodex 对应价格，平均节省 87.5%。',
        columns: { model: '模型', input: '输入', cachedInput: '缓存输入', cacheWrite: '缓存写入', output: '输出' }, official: '官方', note: '说明：所示官方价为你提供的参考单价；OnlyCodex 按对应官方价格的 12.5% 计算。', unit: '价格单位：USD / 1M Tokens'
      },
      subscription: { eyebrow: '订阅制 · 透明定价', title: '按月订阅，稳定使用', description: '适合长期高频使用 Codex 的个人与团队。', official: '官方订阅', perMonth: '/ 月' },
      plans: {
        usageExperience: '使用体验：', dailyLimit: '每日限制：', weeklyLimit: '每周限制：', monthlyLimit: '每月限制：', gpt56Available: 'GPT-5.6 可用', gpt55Available: 'GPT-5.5 可用', gpt54Available: 'GPT-5.4 可用', dedicatedKey: '独享 API Key & 高速通道',
        pro5x: { badge: '推荐', experience: '每天 8 小时高强度使用', dailyLimit: '60 美元/日', weeklyLimit: '360 美元/周', monthlyLimit: '1440 美元/月' },
        pro10x: { badge: '进阶', experience: '每天 16 小时高强度使用', dailyLimit: '120 美元/日', weeklyLimit: '720 美元/周', monthlyLimit: '2880 美元/月' },
        pro20x: { badge: '顶级', experience: '肆无忌惮的高强度使用', dailyLimit: '240 美元/日', weeklyLimit: '1440 美元/周', monthlyLimit: '5760 美元/月' }
      },
      quickStart: { eyebrow: '快速上手', title: '三步开始使用', description: '不需要懂技术，按步骤做就行' },
      steps: {
        register: { badge: '步骤 01', title: '注册账号', description: '访问 OnlyCodex 首页，填写邮箱和密码即可完成注册，全程不超过 1 分钟。', check1: '无需信用卡，免费注册', check2: '支持邮箱一键登录' },
        key: { badge: '步骤 02', title: '获取 API Key', description: '登录控制台，在 API Keys 页面点击新建，几秒钟就能创建好一个专属密钥。', check1: '一键创建，即刻生效', check2: '可创建多个 Key 分项目', check3: '随时禁用或删除' },
        launch: { badge: '步骤 03', title: '启动 Codex', description: '设置环境变量，然后运行 codex，即刻连接到 GPT-5.5 或 GPT-5.4，开始你的编码之旅。', check1: '支持 Mac / Windows / Linux', check2: '配置一次，长期可用' }
      },
      mock: { email: '邮箱', password: '密码', createKey: '+ 新建密钥', firstKey: '我的第一个 Key', projectB: '项目 B' },
      metrics: {
        eyebrow: '核心优势', title: '数据说话', description: '只做 Codex 中转，每一项指标都认真对待',
        uptime: { title: '稳定在线率', description: '专为 Codex 优化的中转线路，7×24 小时监控，低延迟、不掉线，告别连接焦虑。', tag1: '多节点冗余', tag2: '自动故障切换' },
        latency: { title: '平均响应延迟', description: '国内优化线路，请求直达，响应快速。你专注写代码，网络延迟的事交给我们。', tag1: '国内优化', tag2: '低延迟线路' },
        onboarding: { title: '新手上手时间', description: '两个环境变量搞定一切，提供图文教程和一键配置脚本，完全不懂技术也能轻松完成。', tag1: '图文教程', tag2: '一键配置' }
      },
      contact: { eyebrow: '联系我们', title: '加入社区，一起交流', description: '遇到问题、想反馈建议，或只是想和其他开发者聊聊 AI 写代码。', method: '联系方式（微信）', hint: '如果后台已经配置了联系信息，这里会直接读取并展示；没有配置时会显示默认支持邮箱。', modalTitle: '联系咨询', modalDescription: '请加微信：', copyContact: '复制微信号联系咨询' },
      cta: { eyebrow: '准备好了吗', title: '现在就开始一键接入 GPT-5.6 与 GPT-5.5 Codex 模型', description: '注册账号并配置好环境变量，就能直接开始使用 OnlyCodex。' }
    }
  },

  // Key Usage Query Page
  keyUsage: {
    title: 'API Key 用量查询',
    subtitle: '输入您的 API Key 以查看实时消费金额与使用状态',
    placeholder: 'sk-ant-mirror-xxxxxxxxxxxx',
    query: '查询',
    querying: '查询中...',
    privacyNote: '您的 Key 仅在浏览器本地处理，不会被存储',
    dateRange: '统计范围:',
    dateRangeToday: '今日',
    dateRange7d: '7 天',
    dateRange30d: '30 天',
    dateRange90d: '90 天',
    dateRangeCustom: '自定义',
    apply: '应用',
    used: '已使用',
    detailInfo: '详细信息',
    tokenStats: 'Token 统计',
    dailyDetail: '按日明细',
    modelStats: '模型用量统计',
    // Table headers
    date: '日期',
    model: '模型',
    requests: '请求数',
    inputTokens: '输入 Tokens',
    outputTokens: '输出 Tokens',
    cacheCreationTokens: '缓存创建',
    cacheReadTokens: '缓存读取',
    cacheWriteTokens: '缓存写入',
    totalTokens: '总 Tokens',
    cost: '费用',
    // Status
    quotaMode: 'Key 限额模式',
    walletBalance: '钱包余额',
    // Ring card titles
    totalQuota: '总额度',
    limit5h: '5 小时限额',
    limitDaily: '日限额',
    limit7d: '7 天限额',
    limitWeekly: '周限额',
    limitMonthly: '月限额',
    // Detail rows
    remainingQuota: '剩余额度',
    expiresAt: '过期时间',
    todayExpires: '(今日到期)',
    daysLeft: '({days} 天)',
    usedQuota: '已用额度',
    resetNow: '即将重置',
    subscriptionType: '订阅类型',
    subscriptionExpires: '订阅到期',
    // Usage stat cells
    todayRequests: '今日请求',
    todayInputTokens: '今日输入',
    todayOutputTokens: '今日输出',
    todayTokens: '今日 Tokens',
    todayCacheCreation: '今日缓存创建',
    todayCacheRead: '今日缓存读取',
    todayCost: '今日费用',
    rpmTpm: 'RPM / TPM',
    totalRequests: '累计请求',
    totalInputTokens: '累计输入',
    totalOutputTokens: '累计输出',
    totalTokensLabel: '累计 Tokens',
    totalCacheCreation: '累计缓存创建',
    totalCacheRead: '累计缓存读取',
    totalCost: '累计费用',
    avgDuration: '平均耗时',
    // Messages
    enterApiKey: '请输入 API Key',
    querySuccess: '查询成功',
    queryFailed: '查询失败',
    queryFailedRetry: '查询失败，请稍后重试',
    noDailyUsage: '暂无按日用量数据',
  },

  // Setup Wizard
  setup: {
    title: 'Sub2API 安装向导',
    description: '配置您的 Sub2API 实例',
    database: {
      title: '数据库配置',
      description: '连接到您的 PostgreSQL 数据库',
      host: '主机',
      port: '端口',
      username: '用户名',
      password: '密码',
      databaseName: '数据库名称',
      sslMode: 'SSL 模式',
      passwordPlaceholder: '密码',
      ssl: {
        disable: '禁用',
        require: '要求',
        verifyCa: '验证 CA',
        verifyFull: '完全验证'
      }
    },
    redis: {
      title: 'Redis 配置',
      description: '连接到您的 Redis 服务器',
      host: '主机',
      port: '端口',
      password: '密码（可选）',
      database: '数据库',
      passwordPlaceholder: '密码',
      enableTls: '启用 TLS',
      enableTlsHint: '连接 Redis 时使用 TLS（公共 CA 证书）'
    },
    admin: {
      title: '管理员账户',
      description: '创建您的管理员账户',
      email: '邮箱',
      password: '密码',
      confirmPassword: '确认密码',
      passwordPlaceholder: '至少 8 个字符',
      confirmPasswordPlaceholder: '确认密码',
      passwordMismatch: '密码不匹配'
    },
    ready: {
      title: '准备安装',
      description: '检查您的配置并完成安装',
      database: '数据库',
      redis: 'Redis',
      adminEmail: '管理员邮箱'
    },
    status: {
      testing: '测试中...',
      success: '连接成功',
      testConnection: '测试连接',
      installing: '安装中...',
      completeInstallation: '完成安装',
      completed: '安装完成！',
      redirecting: '正在跳转到登录页面...',
      restarting: '服务正在重启，请稍候...',
      timeout: '服务重启时间超出预期，请手动刷新页面。'
    }
  },

  // Common
}
