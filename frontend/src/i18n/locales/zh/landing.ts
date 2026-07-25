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
    heroSubtitle: '管理已获授权的 API 服务',
    heroDescription: '为开发团队提供清晰的 API 接入、密钥管理和用量审计能力',
    tags: {
      subscriptionToApi: '访问控制',
      stickySession: '用量审计',
      realtimeBilling: '配置管理'
    },
    // 用户痛点区块
    painPoints: {
      title: '清晰的管理边界',
      items: {
        expensive: {
          title: '授权范围不清晰',
          desc: '团队需要确认每项服务都具有适当的访问权限与使用依据'
        },
        complex: {
          title: '接入信息分散',
          desc: '不同项目的密钥和配置需要统一管理与审计'
        },
        unstable: {
          title: '运行状态难追踪',
          desc: '团队需要查看请求记录并定位配置问题'
        },
        noControl: {
          title: '责任边界不明确',
          desc: '使用方、项目和服务范围需要保持可追溯'
        }
      }
    },
    // 解决方案区块
    solutions: {
      title: '为已授权接入提供管理能力',
      subtitle: '先确认权限与文档，再进行配置'
    },
    features: {
      unifiedGateway: '接入管理',
      unifiedGatewayDesc: '集中管理已获授权的服务配置与访问密钥。',
      multiAccount: '运行可见',
      multiAccountDesc: '通过用量和请求记录协助定位运行问题。',
      balanceQuota: '使用控制',
      balanceQuotaDesc: '为团队配置配额和访问边界，保持使用过程可追溯。'
    },
    // 优势对比
    comparison: {
      title: '服务原则',
      headers: {
        feature: '对比项',
        official: '使用方责任',
        us: '本站服务'
      },
      items: {
        pricing: {
          feature: '授权与条款',
          official: '确认对外部服务的使用权与责任',
          us: '提供配置和管理能力，不替代外部服务条款'
        },
        models: {
          feature: '服务范围',
          official: '由使用方确认可用服务',
          us: '仅为已授权的配置提供管理入口'
        },
        management: {
          feature: '账号管理',
          official: '保护第三方账户凭据',
          us: '不在公开页面索取第三方凭据'
        },
        stability: {
          feature: '运行记录',
          official: '按外部服务规则使用',
          us: '展示已配置服务的请求和用量记录'
        },
        control: {
          feature: '数据处理',
          official: '遵守适用法律与合同',
          us: '以公开政策说明处理范围与支持渠道'
        }
      }
    },
    providers: {
      title: '可配置的 API 服务',
      description: '仅连接已获授权的服务',
      supported: '已支持',
      soon: '即将推出',
      claude: '服务 A',
      gemini: '服务 B',
      antigravity: '服务 C',
      more: '其他已授权服务'
    },
    // CTA 区块
    cta: {
      title: '开始前，请先确认授权范围',
      description: '阅读公开文档与服务政策后，再通过控制台配置已获授权的服务。',
      button: '进入控制台'
    },
    footer: {
      allRightsReserved: '保留所有权利。'
    }
  },

  // Key Usage Query Page
  keyUsage: {
    title: '用量查询',
    subtitle: '用量信息仅在登录后的控制台中向账户所有者展示',
    placeholder: '',
    query: '',
    querying: '',
    privacyNote: '',
    accessTitle: '用量查询需要登录。',
    accessDescription: '为了避免在公开页面收集、传输或暴露 API 凭据，用量信息仅在完成身份验证后的控制台中展示给对应账户。',
    accessLogin: '前往登录',
    accessDashboard: '进入控制台',
    securityTitle: '凭据安全提示',
    securityDescription: '不要在公开页面输入任何 API Key、第三方账户密码、验证码、支付信息或私钥。本站不会通过此页面要求这些信息。',
    securityReport: '如发现可疑页面、链接或凭据请求，请停止操作并通过首页公开的运营方联系方式报告。',
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
      username: '用户名（可选）',
      password: '密码（可选）',
      database: '数据库',
      usernamePlaceholder: '默认用户留空',
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
