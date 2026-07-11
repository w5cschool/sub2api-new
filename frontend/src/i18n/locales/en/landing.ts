export default {
  batchImageGuide: {
    title: 'Batch Image Generation',
    description: 'Submit multiple prompts in one job and download the generated images when complete'
  },
  // Home Page
  home: {
    viewOnGithub: 'View on GitHub',
    viewDocs: 'View Documentation',
    docs: 'Docs',
    switchToLight: 'Switch to Light Mode',
    switchToDark: 'Switch to Dark Mode',
    dashboard: 'Dashboard',
    login: 'Login',
    getStarted: 'Get Started',
    goToDashboard: 'Go to Dashboard',
    // User-focused value proposition
    heroSubtitle: 'One Key, All AI Models',
    heroDescription: 'No need to manage multiple subscriptions. Access Claude, GPT, Gemini and more with a single API key',
    tags: {
      subscriptionToApi: 'Subscription to API',
      stickySession: 'Session Persistence',
      realtimeBilling: 'Pay As You Go'
    },
    // Pain points section
    painPoints: {
      title: 'Sound Familiar?',
      items: {
        expensive: {
          title: 'High Subscription Costs',
          desc: 'Paying for multiple AI subscriptions that add up every month'
        },
        complex: {
          title: 'Account Chaos',
          desc: 'Managing scattered accounts and API keys across different platforms'
        },
        unstable: {
          title: 'Service Interruptions',
          desc: 'Single accounts hitting rate limits and disrupting your workflow'
        },
        noControl: {
          title: 'No Usage Control',
          desc: "Can't track where your money goes or limit team member usage"
        }
      }
    },
    // Solutions section
    solutions: {
      title: 'We Solve These Problems',
      subtitle: 'Three simple steps to stress-free AI access'
    },
    features: {
      unifiedGateway: 'One-Click Access',
      unifiedGatewayDesc: 'Get a single API key to call all connected AI models. No separate applications needed.',
      multiAccount: 'Always Reliable',
      multiAccountDesc: 'Smart routing across multiple upstream accounts with automatic failover. Say goodbye to errors.',
      balanceQuota: 'Pay What You Use',
      balanceQuotaDesc: 'Usage-based billing with quota limits. Full visibility into team consumption.'
    },
    // Comparison section
    comparison: {
      title: 'Why Choose Us?',
      headers: {
        feature: 'Comparison',
        official: 'Official Subscriptions',
        us: 'Our Platform'
      },
      items: {
        pricing: {
          feature: 'Pricing',
          official: 'Fixed monthly fee, pay even if unused',
          us: 'Pay only for what you use'
        },
        models: {
          feature: 'Model Selection',
          official: 'Single provider only',
          us: 'Switch between models freely'
        },
        management: {
          feature: 'Account Management',
          official: 'Manage each service separately',
          us: 'Unified key, one dashboard'
        },
        stability: {
          feature: 'Stability',
          official: 'Single account rate limits',
          us: 'Multi-account pool, auto-failover'
        },
        control: {
          feature: 'Usage Control',
          official: 'Not available',
          us: 'Quotas & detailed analytics'
        }
      }
    },
    providers: {
      title: 'Supported AI Models',
      description: 'One API, Multiple Choices',
      supported: 'Supported',
      soon: 'Soon',
      claude: 'Claude',
      gemini: 'Gemini',
      antigravity: 'Antigravity',
      more: 'More'
    },
    // CTA section
    cta: {
      title: 'Ready to Get Started?',
      description: 'Sign up now and get free trial credits to experience seamless AI access',
      button: 'Sign Up Free'
    },
    footer: {
      allRightsReserved: 'All rights reserved.'
    },
    onlyCodex: {
      hero: {
        eyebrow: 'Built for Codex · Reliable GPT-5.6 / GPT-5.5 / GPT-5.4 Gateway',
        titlePrefix: 'Give your team a smoother path to',
        description: 'OnlyCodex provides a fast, reliable Codex gateway so teams can use GPT-5.5 and GPT-5.4 more smoothly for development.',
        note: 'We charge just 12.5% of the official price, saving you 87.5%.'
      },
      actions: { signUpAndTry: 'Sign up free and try it now', viewSetupGuide: 'View setup guide', contactNow: 'Contact us', signUpNow: 'Sign up now', signUpFree: 'Sign up free', copy: 'Copy', close: 'Close' },
      terminal: { connected: '✓ Connected to an OnlyCodex gateway node', optimizePrompt: 'Help me optimize the performance of this code...', refactorPrompt: 'Help me refactor this module...' },
      advantages: {
        eyebrow: 'CORE ADVANTAGES', titlePrefix: 'Why developers ', titleHighlight: 'choose us', titleSuffix: '',
        global: { title: 'Global Coverage', description: 'Globally deployed servers with automatic load balancing for fast responses.', tag1: 'Fast global responses', tag2: 'Automatic load balancing' },
        speed: { title: 'Exceptional Speed', description: 'Premium CN2 GIA routes, cross-border RTT as low as 60 ms, and elastic 200 Mbps bandwidth keep things smooth at peak.', tag: '60 ms low latency' },
        reliability: { title: 'Rock-Solid Reliability', description: 'Active-active redundancy across Asia Pacific, automatic failover within 5 seconds, 99.99% SLA, and under 52 minutes of annual downtime.', tag: 'Dual-site redundancy' },
        security: { title: 'Data Security', description: 'HTTPS encryption in transit, separated data and management planes, no retention or analysis of user data, and AES encryption for critical data.', tag1: 'HTTPS encryption', tag2: 'AES encryption' }
      },
      usagePricing: {
        eyebrow: 'USAGE-BASED BILLING · PRICE COMPARISON', titlePrefix: 'Official pricing at just ', titleHighlight: '12.5%', description: 'Per 1M tokens. Official usage pricing is shown above, with corresponding OnlyCodex pricing below—an average savings of 87.5%.',
        columns: { model: 'Model', input: 'Input', cachedInput: 'Cached Input', cacheWrite: 'Cache Write', output: 'Output' }, official: 'Official', note: 'Note: Official prices are the reference unit rates provided; OnlyCodex is calculated at 12.5% of each official rate.', unit: 'Pricing unit: USD / 1M tokens'
      },
      subscription: { eyebrow: 'SUBSCRIPTIONS · TRANSPARENT PRICING', title: 'Subscribe monthly, use reliably', description: 'For individuals and teams who use Codex heavily over the long term.', official: 'Official subscription', perMonth: '/ month' },
      plans: {
        usageExperience: 'Usage profile:', dailyLimit: 'Daily limit:', weeklyLimit: 'Weekly limit:', monthlyLimit: 'Monthly limit:', gpt56Available: 'GPT-5.6 included', gpt55Available: 'GPT-5.5 included', gpt54Available: 'GPT-5.4 included', dedicatedKey: 'Dedicated API key & fast lane',
        pro5x: { badge: 'Recommended', experience: 'Up to 8 hours of intensive use daily', dailyLimit: '$60/day', weeklyLimit: '$360/week', monthlyLimit: '$1,440/month' },
        pro10x: { badge: 'Advanced', experience: 'Up to 16 hours of intensive use daily', dailyLimit: '$120/day', weeklyLimit: '$720/week', monthlyLimit: '$2,880/month' },
        pro20x: { badge: 'Ultimate', experience: 'Unrestricted intensive use', dailyLimit: '$240/day', weeklyLimit: '$1,440/week', monthlyLimit: '$5,760/month' }
      },
      quickStart: { eyebrow: 'QUICK START', title: 'Start in three steps', description: 'No technical expertise needed—just follow the steps.' },
      steps: {
        register: { badge: 'STEP 01', title: 'Create an account', description: 'Open the OnlyCodex home page and register with your email and password in under one minute.', check1: 'No card required; sign up free', check2: 'One-click email sign-in' },
        key: { badge: 'STEP 02', title: 'Get an API key', description: 'Sign in to the dashboard and create your dedicated key from the API Keys page in seconds.', check1: 'Create instantly; active right away', check2: 'Create multiple keys for projects', check3: 'Disable or delete at any time' },
        launch: { badge: 'STEP 03', title: 'Launch Codex', description: 'Set your environment variables, then run codex to connect to GPT-5.5 or GPT-5.4 and start coding.', check1: 'Supports Mac / Windows / Linux', check2: 'Set it up once, use it long term' }
      },
      mock: { email: 'Email', password: 'Password', createKey: '+ New key', firstKey: 'My first key', projectB: 'Project B' },
      metrics: {
        eyebrow: 'CORE ADVANTAGES', title: 'The numbers speak', description: 'A Codex-only gateway where every metric is taken seriously.',
        uptime: { title: 'Uptime', description: 'Codex-optimized gateway routes with 24/7 monitoring, low latency, and reliable connectivity.', tag1: 'Multi-node redundancy', tag2: 'Automatic failover' },
        latency: { title: 'Average Response Latency', description: 'Optimized routes deliver requests directly and respond quickly, so you can focus on writing code.', tag1: 'Optimized routes', tag2: 'Low-latency network' },
        onboarding: { title: 'Time to Get Started', description: 'Two environment variables are all you need, with illustrated guides and one-click setup scripts for an easy start.', tag1: 'Illustrated guide', tag2: 'One-click setup' }
      },
      contact: { eyebrow: 'CONTACT US', title: 'Join the community and connect', description: 'Get help, share feedback, or simply talk with other developers about coding with AI.', method: 'Contact (WeChat)', hint: 'If contact information is configured in the admin panel, it is shown here; otherwise the default support email is displayed.', modalTitle: 'Get in touch', modalDescription: 'Add us on WeChat:', copyContact: 'Copy WeChat ID to get in touch' },
      cta: { eyebrow: 'READY?', title: 'Connect to GPT-5.6 and GPT-5.5 Codex models in one step', description: 'Create an account, set your environment variables, and start using OnlyCodex.' }
    }
  },

  // Key Usage Query Page
  keyUsage: {
    title: 'API Key Usage',
    subtitle: 'Enter your API Key to view real-time spending and usage status',
    placeholder: 'sk-ant-mirror-xxxxxxxxxxxx',
    query: 'Query',
    querying: 'Querying...',
    privacyNote: 'Your Key is processed locally in the browser and will not be stored',
    dateRange: 'Date Range:',
    dateRangeToday: 'Today',
    dateRange7d: '7 Days',
    dateRange30d: '30 Days',
    dateRange90d: '90 Days',
    dateRangeCustom: 'Custom',
    apply: 'Apply',
    used: 'Used',
    detailInfo: 'Detail Information',
    tokenStats: 'Token Statistics',
    dailyDetail: 'Daily Detail',
    modelStats: 'Model Usage Statistics',
    // Table headers
    date: 'Date',
    model: 'Model',
    requests: 'Requests',
    inputTokens: 'Input Tokens',
    outputTokens: 'Output Tokens',
    cacheCreationTokens: 'Cache Creation',
    cacheReadTokens: 'Cache Read',
    cacheWriteTokens: 'Cache Write',
    totalTokens: 'Total Tokens',
    cost: 'Cost',
    // Status
    quotaMode: 'Key Quota Mode',
    walletBalance: 'Wallet Balance',
    // Ring card titles
    totalQuota: 'Total Quota',
    limit5h: '5-Hour Limit',
    limitDaily: 'Daily Limit',
    limit7d: '7-Day Limit',
    limitWeekly: 'Weekly Limit',
    limitMonthly: 'Monthly Limit',
    // Detail rows
    remainingQuota: 'Remaining Quota',
    expiresAt: 'Expires At',
    todayExpires: '(expires today)',
    daysLeft: '({days} days)',
    usedQuota: 'Used Quota',
    resetNow: 'Resetting soon',
    subscriptionType: 'Subscription Type',
    subscriptionExpires: 'Subscription Expires',
    // Usage stat cells
    todayRequests: 'Today Requests',
    todayInputTokens: 'Today Input',
    todayOutputTokens: 'Today Output',
    todayTokens: 'Today Tokens',
    todayCacheCreation: 'Today Cache Creation',
    todayCacheRead: 'Today Cache Read',
    todayCost: 'Today Cost',
    rpmTpm: 'RPM / TPM',
    totalRequests: 'Total Requests',
    totalInputTokens: 'Total Input',
    totalOutputTokens: 'Total Output',
    totalTokensLabel: 'Total Tokens',
    totalCacheCreation: 'Total Cache Creation',
    totalCacheRead: 'Total Cache Read',
    totalCost: 'Total Cost',
    avgDuration: 'Avg Duration',
    // Messages
    enterApiKey: 'Please enter an API Key',
    querySuccess: 'Query successful',
    queryFailed: 'Query failed',
    queryFailedRetry: 'Query failed, please try again later',
    noDailyUsage: 'No daily usage data',
  },

  // Setup Wizard
  setup: {
    title: 'Sub2API Setup',
    description: 'Configure your Sub2API instance',
    database: {
      title: 'Database Configuration',
      description: 'Connect to your PostgreSQL database',
      host: 'Host',
      port: 'Port',
      username: 'Username',
      password: 'Password',
      databaseName: 'Database Name',
      sslMode: 'SSL Mode',
      passwordPlaceholder: 'Password',
      ssl: {
        disable: 'Disable',
        require: 'Require',
        verifyCa: 'Verify CA',
        verifyFull: 'Verify Full'
      }
    },
    redis: {
      title: 'Redis Configuration',
      description: 'Connect to your Redis server',
      host: 'Host',
      port: 'Port',
      password: 'Password (optional)',
      database: 'Database',
      passwordPlaceholder: 'Password',
      enableTls: 'Enable TLS',
      enableTlsHint: 'Use TLS when connecting to Redis (public CA certs)'
    },
    admin: {
      title: 'Admin Account',
      description: 'Create your administrator account',
      email: 'Email',
      password: 'Password',
      confirmPassword: 'Confirm Password',
      passwordPlaceholder: 'Min 8 characters',
      confirmPasswordPlaceholder: 'Confirm password',
      passwordMismatch: 'Passwords do not match'
    },
    ready: {
      title: 'Ready to Install',
      description: 'Review your configuration and complete setup',
      database: 'Database',
      redis: 'Redis',
      adminEmail: 'Admin Email'
    },
    status: {
      testing: 'Testing...',
      success: 'Connection Successful',
      testConnection: 'Test Connection',
      installing: 'Installing...',
      completeInstallation: 'Complete Installation',
      completed: 'Installation completed!',
      redirecting: 'Redirecting to login page...',
      restarting: 'Service is restarting, please wait...',
      timeout: 'Service restart is taking longer than expected. Please refresh the page manually.'
    }
  },

  // Common
}
