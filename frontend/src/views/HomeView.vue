<template>
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <div
    v-else
    class="relative min-h-screen overflow-hidden bg-[#effcf5] text-slate-900 dark:bg-[#07110c] dark:text-emerald-50"
  >
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div class="absolute left-[-10rem] top-[3rem] h-[24rem] w-[24rem] rounded-full bg-emerald-200/45 blur-3xl dark:bg-emerald-500/16"></div>
      <div class="absolute right-[-6rem] top-[7rem] h-[18rem] w-[18rem] rounded-full bg-lime-100/55 blur-3xl dark:bg-lime-300/10"></div>
      <div class="absolute bottom-[12%] left-[35%] h-[16rem] w-[16rem] rounded-full bg-teal-200/25 blur-3xl dark:bg-teal-300/10"></div>
      <div class="home-grid-mask absolute inset-0"></div>
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.45),transparent_45%)] dark:bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.04),transparent_45%)]"></div>
    </div>

    <header
      :class="[
        'sticky top-0 z-30 transition-all duration-300',
        isNavScrolled
          ? 'border-b border-emerald-200/70 bg-[#effcf5]/86 shadow-[0_10px_40px_rgba(6,78,59,0.07)] backdrop-blur-xl dark:border-emerald-900/70 dark:bg-[#07110c]/84'
          : 'bg-transparent'
      ]"
    >
      <nav class="mx-auto flex h-16 w-full max-w-6xl items-center justify-between px-6">
        <div class="flex items-center gap-3">
          <div
            v-if="siteLogo"
            class="flex h-10 w-10 items-center justify-center overflow-hidden rounded-2xl border border-stone-200/80 bg-white/80 shadow-sm dark:border-stone-800 dark:bg-stone-900/80"
          >
            <img :src="siteLogo" alt="Codex Gateway logo" class="h-full w-full object-contain p-1.5" />
          </div>
          <div class="leading-none">
            <div class="text-[1.35rem] font-black tracking-[-0.08em]">
              <span class="text-emerald-600 dark:text-emerald-300">to</span>
              <span>codex</span>
            </div>
            <div class="mt-1 text-[11px] uppercase tracking-[0.24em] text-stone-500 dark:text-stone-400">
              GPT Codex Gateway
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2 sm:gap-3">
          <LocaleSwitcher />

          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-stone-200/70 bg-white/70 text-stone-600 transition hover:-translate-y-0.5 hover:border-stone-300 hover:text-stone-900 dark:border-stone-800 dark:bg-stone-900/70 dark:text-stone-300 dark:hover:border-stone-700 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5 4.462 5 2 6.462 2 8.266v11.468C2 17.962 4.462 16.5 7.5 16.5c1.746 0 3.332.477 4.5 1.253m0-11.5C13.168 5.477 14.754 5 16.5 5c3.038 0 5.5 1.462 5.5 3.266v11.468C22 17.962 19.538 16.5 16.5 16.5c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
          </a>

          <button
            @click="toggleTheme"
            class="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-stone-200/70 bg-white/70 text-stone-600 transition hover:-translate-y-0.5 hover:border-stone-300 hover:text-stone-900 dark:border-stone-800 dark:bg-stone-900/70 dark:text-stone-300 dark:hover:border-stone-700 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <svg
              v-if="isDark"
              class="h-4 w-4"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.7"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1.5m0 15V21m9-9h-1.5M4.5 12H3m14.864 6.364-1.06-1.06M7.197 7.197l-1.06-1.06m12.728 0-1.06 1.06M7.197 16.803l-1.06 1.06M16.5 12a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0z" />
            </svg>
            <svg
              v-else
              class="h-4 w-4"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.7"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" />
            </svg>
          </button>

          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="inline-flex items-center gap-2 rounded-full bg-emerald-950 px-2.5 py-1.5 text-sm font-medium text-white transition hover:-translate-y-0.5 hover:bg-emerald-800 dark:bg-emerald-300 dark:text-emerald-950 dark:hover:bg-emerald-200"
          >
            <span class="flex h-6 w-6 items-center justify-center rounded-full bg-white/15 text-xs font-semibold dark:bg-stone-950/10">
              {{ userInitial }}
            </span>
            <span>{{ t('home.dashboard') }}</span>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="inline-flex items-center rounded-full bg-emerald-950 px-4 py-2 text-sm font-medium text-white transition hover:-translate-y-0.5 hover:bg-emerald-800 dark:bg-emerald-300 dark:text-emerald-950 dark:hover:bg-emerald-200"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="relative z-10">
      <section class="mx-auto grid max-w-6xl gap-14 px-6 pb-10 pt-12 lg:grid-cols-[1.02fr_0.98fr] lg:items-center lg:pt-20">
        <div class="reveal-section max-w-2xl">
          <div class="hero-wave-chip reveal-item inline-flex items-center gap-2 rounded-full border border-emerald-200/80 bg-white/80 px-4 py-2 text-sm font-semibold text-emerald-700 shadow-sm dark:border-emerald-500/20 dark:bg-slate-950 dark:text-emerald-300">
            <span class="h-2 w-2 rounded-full bg-emerald-500"></span>
            专注 Codex · GPT-5.5 / GPT-5.4 稳定中转服务
          </div>

          <h1 class="mt-6 max-w-3xl text-[clamp(2.4rem,6vw,4.85rem)] font-black leading-[0.92] tracking-[-0.065em] text-stone-950 dark:text-stone-50">
            <span class="reveal-item block text-stone-500 dark:text-stone-400">让团队更顺畅地接入</span>
            <span class="hero-wave-text reveal-item mt-2 inline-block bg-gradient-to-r from-emerald-600 via-lime-400 to-teal-300 bg-clip-text text-transparent dark:from-emerald-200 dark:via-lime-200 dark:to-emerald-500">Codex</span>
          </h1>

          <div class="reveal-item mt-8 flex flex-wrap items-center gap-4">
            <router-link
              :to="heroPrimaryTarget"
              class="inline-flex items-center gap-2 rounded-full bg-emerald-400 px-6 py-3 text-sm font-semibold text-emerald-950 shadow-[0_18px_45px_rgba(52,211,153,0.28)] transition hover:-translate-y-0.5 hover:bg-emerald-300"
            >
              {{ isAuthenticated ? t('home.goToDashboard') : '免费注册，立即体验' }}
              <svg class="h-4 w-4" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 10h12M12 6l4 4-4 4" />
              </svg>
            </router-link>
            <a
              v-if="docUrl"
              :href="docUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center gap-2 rounded-full border border-stone-300/80 bg-white/80 px-5 py-3 text-sm font-medium text-stone-700 transition hover:-translate-y-0.5 hover:border-stone-400 dark:border-stone-700 dark:bg-stone-900/80 dark:text-stone-200 dark:hover:border-stone-600"
            >
              查看配置教程
            </a>
          </div>
        </div>

        <div class="reveal-section">
          <div class="reveal-item relative mx-auto max-w-[34rem]">
            <div class="overflow-hidden rounded-[2rem] border border-stone-200/80 bg-[#18120f] shadow-[0_30px_90px_rgba(28,25,23,0.28)] dark:border-stone-800">
              <div class="flex items-center gap-2 border-b border-white/5 bg-white/5 px-5 py-4">
                <span class="h-3 w-3 rounded-full bg-[#ff5f56]"></span>
                <span class="h-3 w-3 rounded-full bg-[#ffbd2e]"></span>
                <span class="h-3 w-3 rounded-full bg-[#27c93f]"></span>
                <span class="ml-3 text-xs uppercase tracking-[0.24em] text-stone-400">codex</span>
              </div>
              <div class="space-y-4 px-5 py-6 font-mono text-sm text-stone-200">
                <div class="text-stone-500">
                  <span class="mr-2 text-emerald-300">$</span>
                  export OPENAI_BASE_URL="<span class="text-stone-100">{{ gatewayBaseUrl }}</span>"
                </div>
                <div class="text-stone-500">
                  <span class="mr-2 text-emerald-300">$</span>
                  codex
                </div>
                <div class="rounded-2xl border border-white/5 bg-white/5 px-4 py-3 text-stone-300">
                  <div class="text-emerald-300">◆ Codex CLI · GPT-5.5 ready</div>
                </div>
                <div class="text-stone-400">
                  <span class="mr-2 text-emerald-300">›</span>
                  帮我优化这段代码的性能...
                </div>
                <div>
                  <span class="mr-2 text-emerald-300">$</span>
                  <span class="cursor-blink inline-block h-5 w-2 rounded-sm bg-emerald-300/80 align-middle"></span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="reveal-section mx-auto max-w-6xl px-6 py-10">
        <div class="reveal-item text-center">
          <div class="text-xs font-semibold uppercase tracking-[0.3em] text-stone-500 dark:text-stone-400">透明定价</div>
        </div>

        <div class="mx-auto mt-8 grid max-w-6xl gap-5 lg:grid-cols-3">
          <article
            v-for="card in pricingCards"
            :key="card.model"
            class="reveal-item relative overflow-hidden rounded-[1.75rem] border bg-white/80 p-5 shadow-[0_18px_60px_rgba(28,25,23,0.08)] backdrop-blur-sm transition hover:-translate-y-1 dark:bg-stone-900 dark:shadow-[0_18px_60px_rgba(0,0,0,0.3)] md:p-6"
            :class="[card.borderClass, card.featured ? 'ring-2 ring-emerald-300/50 dark:ring-emerald-400/30' : '']"
          >
            <div v-if="card.featured" class="absolute inset-x-8 top-0 h-px bg-gradient-to-r from-transparent via-emerald-300 to-transparent"></div>
            <div class="inline-flex rounded-full px-3 py-1 text-xs font-semibold" :class="card.badgeClass">{{ card.badge }}</div>
            <div class="mt-3 text-3xl font-black tracking-[-0.05em]">{{ card.model }}</div>
            <div class="mt-6 flex items-end justify-center gap-1.5 text-center">
              <span class="text-[2.8rem] font-black tracking-[-0.07em]">{{ card.price }}</span>
              <span class="pb-1 text-sm font-semibold text-stone-500 dark:text-stone-400">/月</span>
            </div>

            <div class="mt-6 border-t border-stone-200/80 pt-5 dark:border-stone-800">
              <div v-for="row in card.facts" :key="`${card.model}-${row.label}`" class="grid grid-cols-[auto_1fr] items-start gap-4 py-2 text-sm">
                <div class="text-stone-500 dark:text-stone-400">{{ row.label }}</div>
                <div class="text-right font-medium text-stone-800 dark:text-stone-100">{{ row.value }}</div>
              </div>
            </div>

            <div class="mt-5 border-t border-stone-200/80 pt-5 dark:border-stone-800">
              <div v-for="feature in card.features" :key="`${card.model}-${feature}`" class="flex items-center gap-3 py-2 text-sm font-medium text-stone-800 dark:text-stone-100">
                <span class="inline-flex h-6 w-6 items-center justify-center rounded-full bg-emerald-500 text-white">
                  <svg class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 10l4 4 8-8" /></svg>
                </span>
                <span>{{ feature }}</span>
              </div>
            </div>

            <button type="button" class="mt-6 inline-flex w-full items-center justify-center rounded-full border border-emerald-300 bg-emerald-400 px-6 py-2.5 text-sm font-semibold text-emerald-950 shadow-[0_18px_40px_rgba(52,211,153,0.25)] transition hover:-translate-y-0.5 hover:bg-emerald-300" @click="showWechatModal = true">
              立即咨询
            </button>
          </article>
        </div>

        <p class="reveal-item mx-auto mt-4 max-w-6xl text-center text-xs leading-5 text-stone-500 dark:text-stone-400 sm:text-right">
          注：每月按 4 周（28 天）计算。
        </p>
      </section>

      <section class="reveal-section mx-auto max-w-6xl px-6 py-10">
        <div class="reveal-item overflow-hidden rounded-[2.4rem] border border-emerald-200/80 bg-emerald-950 text-emerald-50 shadow-[0_24px_80px_rgba(6,78,59,0.18)] dark:border-emerald-500/20 dark:bg-[#0b1912]">
          <div class="grid lg:grid-cols-[0.78fr_1.22fr]">
            <div class="relative overflow-hidden border-b border-emerald-800/80 px-7 py-8 lg:border-b-0 lg:border-r lg:px-9 lg:py-10">
              <div class="pointer-events-none absolute -right-10 -top-14 h-44 w-44 rounded-full border-[28px] border-emerald-400/10"></div>
              <div class="relative">
                <div class="text-xs font-semibold uppercase tracking-[0.3em] text-emerald-300">用户奖励政策</div>
                <h2 class="mt-4 max-w-sm text-3xl font-black tracking-[-0.05em] text-white md:text-4xl">
                  使用越久，<br class="hidden lg:block" />分享越多
                </h2>
                <p class="mt-4 max-w-md text-sm leading-7 text-emerald-100/70">
                  奖励以对应套餐的按量计费额度发放，规则清晰，权益直接到账。
                </p>
              </div>
            </div>

            <div class="divide-y divide-emerald-800/80">
              <article
                v-for="(policy, index) in rewardPolicies"
                :key="policy.title"
                class="group grid gap-5 px-7 py-7 transition-colors hover:bg-emerald-900/55 sm:grid-cols-[auto_1fr_auto] sm:items-center sm:px-9"
              >
                <div class="text-xs font-semibold tracking-[0.22em] text-emerald-400/70">
                  {{ String(index + 1).padStart(2, '0') }}
                </div>
                <div>
                  <h3 class="text-base font-bold text-white sm:text-lg">{{ policy.title }}</h3>
                  <p class="mt-1.5 text-sm leading-6 text-emerald-100/70">{{ policy.trigger }}</p>
                </div>
                <div class="flex items-baseline gap-2 sm:block sm:text-right">
                  <span class="text-4xl font-black tracking-[-0.06em] text-emerald-300">1</span>
                  <span class="text-sm font-semibold text-emerald-100">天额度</span>
                </div>
              </article>
            </div>
          </div>

          <div class="flex items-start gap-3 border-t border-emerald-800/80 bg-emerald-900/35 px-7 py-4 text-xs leading-5 text-emerald-100/65 sm:items-center sm:px-9">
            <svg class="mt-0.5 h-4 w-4 shrink-0 text-emerald-300 sm:mt-0" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.7">
              <circle cx="10" cy="10" r="7.25" />
              <path stroke-linecap="round" d="M10 9v4M10 6.75h.01" />
            </svg>
            <span>赠送额度按所续费或所购买的对应套餐计算。</span>
          </div>
        </div>
      </section>

      <section class="reveal-section mx-auto max-w-6xl px-6 py-10">
        <div class="overflow-hidden rounded-[2.4rem] border border-stone-200/80 bg-white/78 p-8 shadow-[0_24px_80px_rgba(28,25,23,0.1)] backdrop-blur-sm dark:border-stone-800 dark:bg-stone-900 md:p-10">
          <div class="reveal-item text-center">
            <div class="text-xs font-semibold uppercase tracking-[0.3em] text-stone-500 dark:text-stone-400">联系我们</div>
            <h2 class="mt-4 text-3xl font-black tracking-[-0.05em] md:text-4xl">加入社区，一起交流</h2>
            <p class="mt-3 text-base text-stone-600 dark:text-stone-300">
              遇到问题、想反馈建议，或只是想和其他开发者聊聊 AI 写代码。
            </p>
          </div>

          <div class="mt-8 flex justify-center">
            <article class="reveal-item w-full max-w-xl rounded-[1.8rem] border border-stone-200/80 bg-stone-50/90 p-6 dark:border-stone-800 dark:bg-stone-950/60">
              <div class="text-xs font-semibold uppercase tracking-[0.22em] text-stone-400">联系方式（微信）</div>
              <div class="mt-4 flex flex-wrap items-center gap-3">
                <code class="rounded-2xl bg-white px-4 py-3 text-sm font-medium text-stone-700 shadow-sm dark:bg-stone-900 dark:text-stone-200">
                  {{ contactValue }}
                </code>
                <button
                  type="button"
                  class="inline-flex items-center gap-2 rounded-full border border-stone-200 bg-white px-4 py-2 text-sm font-medium text-stone-600 transition hover:-translate-y-0.5 hover:text-stone-900 dark:border-stone-700 dark:bg-stone-900 dark:text-stone-300 dark:hover:text-white"
                  @click="copyContact"
                >
                  <svg
                    v-if="copySucceeded"
                    class="h-4 w-4"
                    viewBox="0 0 20 20"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.8"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4 10l4 4 8-8" />
                  </svg>
                  <svg
                    v-else
                    class="h-4 w-4"
                    viewBox="0 0 20 20"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.8"
                  >
                    <rect x="6" y="6" width="10" height="10" rx="2" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4 14V5a1 1 0 0 1 1-1h9" />
                  </svg>
                  <span>{{ copySucceeded ? t('common.copied') : '复制' }}</span>
                </button>
              </div>
              <p class="mt-4 text-sm leading-7 text-stone-600 dark:text-stone-300">
                如果后台已经配置了联系信息，这里会直接读取并展示；没有配置时会显示默认支持邮箱。
              </p>
            </article>
          </div>
        </div>
      </section>

      <section class="reveal-section mx-auto max-w-6xl px-6 pb-20 pt-10">
        <div class="reveal-item relative overflow-hidden rounded-[2.6rem] border border-emerald-300/60 bg-gradient-to-br from-emerald-300 to-lime-200 px-8 py-12 text-center text-emerald-950 shadow-[0_30px_90px_rgba(52,211,153,0.22)]">
          <div class="absolute left-10 top-8 h-28 w-28 rounded-full bg-white/12 blur-2xl"></div>
          <div class="absolute bottom-0 right-8 h-32 w-32 rounded-full bg-stone-950/10 blur-2xl"></div>
          <div class="relative">
            <div class="text-xs font-semibold uppercase tracking-[0.32em] text-stone-900/70">准备好了吗</div>
            <h2 class="mt-4 text-3xl font-black tracking-[-0.06em] md:text-5xl">现在就开始一键接入 GPT-5.5 与 GPT-5.4 Codex 模型</h2>
            <p class="mx-auto mt-4 max-w-2xl text-base leading-8 text-stone-900/80">
              注册账号并配置好环境变量，就能直接开始使用 tocodex。
            </p>
            <router-link
              :to="heroPrimaryTarget"
              class="mt-8 inline-flex items-center gap-2 rounded-full bg-stone-950 px-6 py-3 text-sm font-semibold text-white transition hover:-translate-y-0.5 hover:bg-stone-800"
            >
              {{ isAuthenticated ? t('home.goToDashboard') : '免费注册' }}
              <svg class="h-4 w-4" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 10h12M12 6l4 4-4 4" />
              </svg>
            </router-link>
          </div>
        </div>
      </section>
    </main>

    <footer class="relative z-10 bg-[#10110f] text-stone-400">
      <div class="mx-auto max-w-6xl px-6 py-12 md:py-14">
        <div class="grid gap-x-12 gap-y-10 sm:grid-cols-2 lg:grid-cols-[0.8fr_1.25fr_0.8fr]">
          <section>
            <h2 class="text-base font-bold text-stone-100">关于 apiclub.cc</h2>
            <nav class="mt-5 flex flex-col items-start gap-3 text-sm">
              <RouterLink to="/legal/terms" class="transition hover:text-white">服务条款</RouterLink>
              <RouterLink to="/legal/privacy-policy" class="transition hover:text-white">隐私政策</RouterLink>
              <RouterLink to="/legal/refund-cancellation-policy" class="transition hover:text-white">退款与取消政策</RouterLink>
              <a
                v-if="docUrl"
                :href="docUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="transition hover:text-white"
              >
                使用文档
              </a>
            </nav>
          </section>

          <section>
            <h2 class="text-base font-bold text-stone-100">运营主体</h2>
            <p class="mt-5 text-sm font-semibold text-stone-200">Pluto Horizon LLC</p>
            <p class="mt-2 max-w-md text-sm leading-7">30 N Gould St #65046, Sheridan, Wyoming 82801, United States</p>
            <p class="mt-3 max-w-md text-sm leading-7">本站独立运营，不代表或隶属于任何第三方服务品牌。</p>
          </section>

          <section class="sm:text-right">
            <h2 class="text-base font-bold text-stone-100">联系我们</h2>
            <div class="mt-5 space-y-2 text-sm leading-7">
              <a :href="`mailto:${supportEmail}`" class="block transition hover:text-white">{{ supportEmail }}</a>
              <span class="block">管理员微信：{{ administratorWechat }}</span>
            </div>
          </section>
        </div>

        <div class="mt-12 border-t border-stone-800 pt-6 text-center text-sm text-stone-500">
          Copyright © {{ copyrightYears }} Pluto Horizon LLC. All rights reserved.
        </div>
      </div>
    </footer>

    <Transition name="modal-fade">
      <div
        v-if="showWechatModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-stone-950/50 backdrop-blur-sm"
        @click.self="showWechatModal = false"
      >
        <div class="relative w-80 rounded-3xl border border-emerald-200/70 bg-white p-8 text-center shadow-[0_30px_80px_rgba(52,211,153,0.18)] dark:border-emerald-500/20 dark:bg-slate-950">
          <div class="mx-auto mb-4 inline-flex h-12 w-12 items-center justify-center rounded-full bg-emerald-500/10">
            <svg class="h-6 w-6 text-emerald-500" viewBox="0 0 24 24" fill="currentColor">
              <path d="M8.5 10a1 1 0 1 0 0-2 1 1 0 0 0 0 2Zm7 0a1 1 0 1 0 0-2 1 1 0 0 0 0 2Z"/>
              <path d="M12 2C6.477 2 2 6.224 2 11.5c0 2.79 1.23 5.3 3.2 7.05L4 22l3.83-1.56A10.64 10.64 0 0 0 12 21c5.523 0 10-4.224 10-9.5S17.523 2 12 2Z"/>
            </svg>
          </div>
          <h3 class="text-lg font-bold text-stone-900 dark:text-stone-50">联系咨询</h3>
          <p class="mt-3 text-sm text-stone-600 dark:text-stone-300">请加微信：</p>
          <div class="mt-4 flex flex-col items-center gap-3">
            <p class="text-xl font-bold text-emerald-500">{{ contactValue }}</p>
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-full border border-emerald-200 bg-emerald-50 px-4 py-2 text-sm font-semibold text-emerald-700 transition hover:-translate-y-0.5 hover:border-emerald-300 hover:bg-emerald-100 dark:border-emerald-500/20 dark:bg-emerald-500/10 dark:text-emerald-200 dark:hover:bg-emerald-500/15"
              @click="copyContact"
            >
              <svg
                v-if="copySucceeded"
                class="h-4 w-4"
                viewBox="0 0 20 20"
                fill="none"
                stroke="currentColor"
                stroke-width="1.8"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 10l4 4 8-8" />
              </svg>
              <svg
                v-else
                class="h-4 w-4"
                viewBox="0 0 20 20"
                fill="none"
                stroke="currentColor"
                stroke-width="1.8"
              >
                <rect x="6" y="6" width="10" height="10" rx="2" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 14V5a1 1 0 0 1 1-1h9" />
              </svg>
              <span>{{ copySucceeded ? t('common.copied') : '复制微信号联系咨询' }}</span>
            </button>
          </div>
          <button
            type="button"
            class="mt-6 inline-flex w-full items-center justify-center rounded-full bg-emerald-400 px-6 py-2.5 text-sm font-semibold text-emerald-950 transition hover:-translate-y-0.5 hover:bg-emerald-300"
            @click="showWechatModal = false"
          >
            关闭
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import { sanitizeUrl } from '@/utils/url'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

type PricingCard = {
  badge: string
  model: string
  price: string
  facts: Array<{ label: string; value: string }>
  features: string[]
  badgeClass: string
  borderClass: string
  featured?: boolean
}

const pricingCards: PricingCard[] = [
  {
    badge: '推荐', model: 'Pro 5X', price: '¥560', featured: true,
    badgeClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300',
    borderClass: 'border-emerald-300/80 dark:border-emerald-500/30',
    facts: [{ label: '每日限制：', value: '60美元/日' }, { label: '每周限制：', value: '360美元/周' }, { label: '每月限制：', value: '1440美元/月' }],
    features: ['GPT-5.6 可用', 'GPT-5.5 可用', 'GPT-5.4 可用']
  },
  {
    badge: '进阶', model: 'Pro 10X', price: '¥900',
    badgeClass: 'bg-lime-100 text-lime-700 dark:bg-lime-500/10 dark:text-lime-300',
    borderClass: 'border-lime-200/80 dark:border-lime-500/20',
    facts: [{ label: '每日限制：', value: '120美元/日' }, { label: '每周限制：', value: '720美元/周' }, { label: '每月限制：', value: '2880美元/月' }],
    features: ['GPT-5.6 可用', 'GPT-5.5 可用', 'GPT-5.4 可用']
  },
  {
    badge: '顶级', model: 'Pro 20X', price: '¥1800',
    badgeClass: 'bg-teal-100 text-teal-700 dark:bg-teal-500/10 dark:text-teal-300',
    borderClass: 'border-teal-200/80 dark:border-teal-500/20',
    facts: [{ label: '每日限制：', value: '240美元/日' }, { label: '每周限制：', value: '1440美元/周' }, { label: '每月限制：', value: '5760美元/月' }],
    features: ['GPT-5.6 可用', 'GPT-5.5 可用', 'GPT-5.4 可用']
  }
]

const rewardPolicies = [
  {
    title: '续费优惠政策',
    trigger: '续费任一套餐，即可获赠对应套餐的按量计费额度。'
  },
  {
    title: '邀请新人购买赠送策略',
    trigger: '成功邀请一位新用户购买套餐，即可获赠对应套餐的按量计费额度。'
  }
]

const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const gatewayBaseUrl = computed(() => appStore.cachedPublicSettings?.api_base_url || appStore.apiBaseUrl || 'https://tocodex.cc')
const contactValue = computed(() => appStore.cachedPublicSettings?.contact_info?.trim() || 'Icanmeetu')

const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const isDark = ref(document.documentElement.classList.contains('dark'))
const isNavScrolled = ref(false)
const copySucceeded = ref(false)
const showWechatModal = ref(false)

let copyTimer: ReturnType<typeof setTimeout> | null = null
let revealObserver: IntersectionObserver | null = null

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const heroPrimaryTarget = computed(() => isAuthenticated.value ? dashboardPath.value : '/login')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return 'A'
  return user.email.charAt(0).toUpperCase()
})
const currentYear = computed(() => new Date().getFullYear())
const copyrightYears = computed(() => currentYear.value > 2026 ? `2026–${currentYear.value}` : '2026')
const supportEmail = 'steveweng.s.w@gmail.com'
const administratorWechat = 'Icanmeetu'

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

function updateNavState() {
  isNavScrolled.value = window.scrollY > 16
}

function setupRevealObserver() {
  revealObserver = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (!entry.isIntersecting) return
      entry.target.classList.add('is-visible')
      entry.target.querySelectorAll<HTMLElement>('.reveal-item').forEach((item, index) => {
        item.style.transitionDelay = `${index * 70}ms`
        item.classList.add('is-visible')
      })
    })
  }, { threshold: 0.08, rootMargin: '0px 0px -72px 0px' })

  document.querySelectorAll<HTMLElement>('.reveal-section').forEach((section) => revealObserver?.observe(section))
}

async function copyContact() {
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(contactValue.value)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = contactValue.value
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }

    copySucceeded.value = true
    if (copyTimer) clearTimeout(copyTimer)
    copyTimer = setTimeout(() => {
      copySucceeded.value = false
    }, 1600)
  } catch {
    copySucceeded.value = false
  }
}

onMounted(() => {
  initTheme()
  updateNavState()
  window.addEventListener('scroll', updateNavState, { passive: true })
  window.setTimeout(setupRevealObserver, 120)

  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', updateNavState)
  revealObserver?.disconnect()
  if (copyTimer) clearTimeout(copyTimer)
})
</script>

<style scoped>
.home-grid-mask {
  background-image:
    linear-gradient(rgba(120, 113, 108, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(120, 113, 108, 0.08) 1px, transparent 1px);
  background-size: 40px 40px;
  mask-image: radial-gradient(circle at center, black 25%, transparent 82%);
}

.hero-wave-chip,
.hero-wave-text {
  position: relative;
  overflow: hidden;
  isolation: isolate;
}

.hero-wave-chip::after,
.hero-wave-text::after {
  content: '';
  position: absolute;
  inset: -28%;
  transform: translateX(-165%) skewX(-24deg);
  background:
    linear-gradient(
      90deg,
      transparent 0%,
      rgba(255, 255, 255, 0) 24%,
      rgba(255, 248, 240, 0.22) 40%,
      rgba(255, 255, 255, 0.65) 48%,
      rgba(255, 242, 214, 0.95) 50%,
      rgba(255, 255, 255, 0.65) 52%,
      rgba(255, 248, 240, 0.22) 60%,
      rgba(255, 255, 255, 0) 76%,
      transparent 100%
    );
  opacity: 0;
  pointer-events: none;
}

.hero-wave-chip::after {
  animation: hero-wave-sheen 5.8s cubic-bezier(0.22, 1, 0.36, 1) infinite;
}

.hero-wave-text::after {
  inset: -18%;
  animation: hero-wave-sheen 4.9s cubic-bezier(0.22, 1, 0.36, 1) 0.5s infinite;
}

.reveal-section,
.reveal-item {
  opacity: 0;
  transform: translateY(24px);
  transition:
    opacity 0.7s ease,
    transform 0.7s cubic-bezier(0.22, 1, 0.36, 1);
}

.reveal-section.is-visible,
.reveal-item.is-visible {
  opacity: 1;
  transform: translateY(0);
}

.cursor-blink {
  animation: cursor-blink 1s steps(2, start) infinite;
}

@keyframes cursor-blink {
  0%,
  45% {
    opacity: 1;
  }

  46%,
  100% {
    opacity: 0;
  }
}

@keyframes hero-wave-sheen {
  0%,
  12% {
    opacity: 0;
    transform: translateX(-165%) skewX(-24deg);
  }

  20%,
  34% {
    opacity: 1;
  }

  44% {
    opacity: 0;
    transform: translateX(165%) skewX(-24deg);
  }

  100% {
    opacity: 0;
    transform: translateX(165%) skewX(-24deg);
  }
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

@media (prefers-reduced-motion: reduce) {
  .reveal-section,
  .reveal-item,
  .modal-fade-enter-active,
  .modal-fade-leave-active {
    transition-duration: 1ms !important;
    transition-delay: 0ms !important;
  }

  .reveal-section,
  .reveal-item {
    transform: none;
  }

  .cursor-blink,
  .hero-wave-chip::after,
  .hero-wave-text::after {
    animation: none !important;
  }

  .hero-wave-chip::after,
  .hero-wave-text::after {
    opacity: 0 !important;
  }
}
</style>
