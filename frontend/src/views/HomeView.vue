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
    class="onlycodex-home relative min-h-screen overflow-hidden bg-[#f6f8ff] text-slate-900 dark:bg-[#0b1020] dark:text-blue-50"
  >
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div class="absolute left-[-10rem] top-[3rem] h-[24rem] w-[24rem] rounded-full bg-blue-200/45 blur-3xl dark:bg-blue-500/16"></div>
      <div class="absolute right-[-6rem] top-[7rem] h-[18rem] w-[18rem] rounded-full bg-violet-100/55 blur-3xl dark:bg-violet-300/10"></div>
      <div class="absolute bottom-[12%] left-[35%] h-[16rem] w-[16rem] rounded-full bg-sky-200/25 blur-3xl dark:bg-sky-300/10"></div>
      <div class="home-grid-mask absolute inset-0"></div>
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.45),transparent_45%)] dark:bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.04),transparent_45%)]"></div>
    </div>

    <header
      :class="[
        'sticky top-0 z-30 transition-all duration-300',
        isNavScrolled
          ? 'border-b border-blue-200/70 bg-[#f6f8ff]/86 shadow-[0_10px_40px_rgba(66,104,184,0.08)] backdrop-blur-xl dark:border-blue-900/70 dark:bg-[#0b1020]/84'
          : 'bg-transparent'
      ]"
    >
      <nav class="mx-auto flex h-16 w-full max-w-6xl items-center justify-between px-6">
        <div class="flex items-center gap-3">
          <div
            v-if="siteLogo"
            class="flex h-10 w-10 items-center justify-center overflow-hidden rounded-2xl border border-stone-200/80 bg-white/80 shadow-sm dark:border-stone-800 dark:bg-stone-900/80"
          >
            <img :src="siteLogo" alt="OnlyCodex logo" class="h-full w-full object-contain p-1.5" />
          </div>
          <div class="leading-none">
            <div class="text-[1.35rem] font-black tracking-[-0.08em]">
              <span class="text-blue-600 dark:text-blue-300">Only</span>
              <span>Codex</span>
            </div>
            <div class="mt-1 text-[11px] uppercase tracking-[0.24em] text-stone-500 dark:text-stone-400">
              AI Coding Gateway
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
            {{ t('home.onlyCodex.hero.eyebrow') }}
          </div>

          <h1 class="mt-6 max-w-3xl text-[clamp(2.4rem,6vw,4.85rem)] font-black leading-[0.92] tracking-[-0.065em] text-stone-950 dark:text-stone-50">
            <span class="reveal-item block text-stone-500 dark:text-stone-400">{{ t('home.onlyCodex.hero.titlePrefix') }}</span>
            <span class="hero-wave-text reveal-item mt-2 inline-block bg-gradient-to-r from-emerald-600 via-lime-400 to-teal-300 bg-clip-text text-transparent dark:from-emerald-200 dark:via-lime-200 dark:to-emerald-500">Codex</span>
          </h1>

          <div class="mt-6 max-w-xl space-y-4">
            <p class="reveal-item text-base leading-8 text-stone-600 dark:text-stone-300 md:text-lg">
              {{ t('home.onlyCodex.hero.description') }}
            </p>
            <p class="reveal-item max-w-lg text-sm font-semibold uppercase tracking-[0.18em] text-stone-700 dark:text-stone-200 md:text-[0.95rem]">
              {{ t('home.onlyCodex.hero.note') }}
            </p>
          </div>

          <div class="reveal-item mt-8 flex flex-wrap items-center gap-4">
            <router-link
              :to="heroPrimaryTarget"
              class="inline-flex items-center gap-2 rounded-full bg-emerald-400 px-6 py-3 text-sm font-semibold text-emerald-950 shadow-[0_18px_45px_rgba(52,211,153,0.28)] transition hover:-translate-y-0.5 hover:bg-emerald-300"
            >
              {{ isAuthenticated ? t('home.goToDashboard') : t('home.onlyCodex.actions.signUpAndTry') }}
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
              {{ t('home.onlyCodex.actions.viewSetupGuide') }}
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
                  <div class="text-emerald-300">◆ Codex CLI · GPT-5.6 Sol ready</div>
                  <div class="mt-1 text-stone-400">{{ t('home.onlyCodex.terminal.connected') }}</div>
                </div>
                <div class="text-stone-400">
                  <span class="mr-2 text-emerald-300">›</span>
                  {{ t('home.onlyCodex.terminal.optimizePrompt') }}
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
          <div class="text-xs font-semibold uppercase tracking-[0.3em] text-stone-500 dark:text-stone-400">{{ t('home.onlyCodex.advantages.eyebrow') }}</div>
          <h2 class="mt-4 text-3xl font-black tracking-[-0.05em] md:text-4xl">
            {{ t('home.onlyCodex.advantages.titlePrefix') }}<span class="text-emerald-600 dark:text-emerald-300">{{ t('home.onlyCodex.advantages.titleHighlight') }}</span>{{ t('home.onlyCodex.advantages.titleSuffix') }}
          </h2>
        </div>

        <div class="mt-10 grid gap-5 md:grid-cols-2 lg:grid-cols-4">
          <article class="reveal-item relative overflow-hidden rounded-[1.75rem] border border-stone-200/80 bg-white/78 p-6 shadow-[0_18px_60px_rgba(28,25,23,0.08)] backdrop-blur-sm transition hover:-translate-y-1 dark:border-stone-800 dark:bg-stone-900">
            <div class="mb-5 inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-emerald-100 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-300">
              <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 21a9.004 9.004 0 008.716-6.747M12 21a9.004 9.004 0 01-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 017.843 4.582M12 3a8.997 8.997 0 00-7.843 4.582m15.686 0A11.953 11.953 0 0112 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0121 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0112 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 013 12c0-1.605.42-3.113 1.157-4.418" />
              </svg>
            </div>
            <h3 class="text-lg font-black tracking-[-0.03em]">{{ t('home.onlyCodex.advantages.global.title') }}</h3>
            <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">{{ t('home.onlyCodex.advantages.global.description') }}</p>
            <div class="mt-5 flex flex-wrap gap-2">
              <span class="rounded-full border border-stone-200/80 px-3 py-1 text-xs font-medium text-stone-500 dark:border-stone-700 dark:text-stone-400">{{ t('home.onlyCodex.advantages.global.tag1') }}</span>
              <span class="rounded-full border border-stone-200/80 px-3 py-1 text-xs font-medium text-stone-500 dark:border-stone-700 dark:text-stone-400">{{ t('home.onlyCodex.advantages.global.tag2') }}</span>
            </div>
          </article>

          <article class="reveal-item relative overflow-hidden rounded-[1.75rem] border border-emerald-300/60 bg-white/78 p-6 shadow-[0_18px_60px_rgba(52,211,153,0.12)] backdrop-blur-sm transition hover:-translate-y-1 dark:border-emerald-500/30 dark:bg-slate-950">
            <div class="absolute inset-x-8 top-0 h-px bg-gradient-to-r from-transparent via-emerald-300 to-transparent"></div>
            <div class="mb-5 inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-emerald-100 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-300">
              <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
              </svg>
            </div>
            <h3 class="text-lg font-black tracking-[-0.03em]">{{ t('home.onlyCodex.advantages.speed.title') }}</h3>
            <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">{{ t('home.onlyCodex.advantages.speed.description') }}</p>
            <div class="mt-5 flex flex-wrap gap-2">
              <span class="rounded-full border border-emerald-300/80 bg-emerald-50 px-3 py-1 text-xs font-medium text-emerald-700 dark:border-emerald-500/30 dark:bg-emerald-500/10 dark:text-emerald-300">{{ t('home.onlyCodex.advantages.speed.tag') }}</span>
              <span class="rounded-full border border-emerald-300/80 bg-emerald-50 px-3 py-1 text-xs font-medium text-emerald-700 dark:border-emerald-500/30 dark:bg-emerald-500/10 dark:text-emerald-300">200 Mbps</span>
            </div>
          </article>

          <article class="reveal-item relative overflow-hidden rounded-[1.75rem] border border-stone-200/80 bg-white/78 p-6 shadow-[0_18px_60px_rgba(28,25,23,0.08)] backdrop-blur-sm transition hover:-translate-y-1 dark:border-stone-800 dark:bg-stone-900">
            <div class="mb-5 inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-lime-100 text-lime-700 dark:bg-lime-500/15 dark:text-lime-300">
              <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" />
              </svg>
            </div>
            <h3 class="text-lg font-black tracking-[-0.03em]">{{ t('home.onlyCodex.advantages.reliability.title') }}</h3>
            <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">{{ t('home.onlyCodex.advantages.reliability.description') }}</p>
            <div class="mt-5 flex flex-wrap gap-2">
              <span class="rounded-full border border-stone-200/80 px-3 py-1 text-xs font-medium text-stone-500 dark:border-stone-700 dark:text-stone-400">SLA 99.99%</span>
              <span class="rounded-full border border-stone-200/80 px-3 py-1 text-xs font-medium text-stone-500 dark:border-stone-700 dark:text-stone-400">{{ t('home.onlyCodex.advantages.reliability.tag') }}</span>
            </div>
          </article>

          <article class="reveal-item relative overflow-hidden rounded-[1.75rem] border border-stone-200/80 bg-white/78 p-6 shadow-[0_18px_60px_rgba(28,25,23,0.08)] backdrop-blur-sm transition hover:-translate-y-1 dark:border-stone-800 dark:bg-stone-900">
            <div class="mb-5 inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-emerald-100 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-400">
              <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6">
                <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z" />
              </svg>
            </div>
            <h3 class="text-lg font-black tracking-[-0.03em]">{{ t('home.onlyCodex.advantages.security.title') }}</h3>
            <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">{{ t('home.onlyCodex.advantages.security.description') }}</p>
            <div class="mt-5 flex flex-wrap gap-2">
              <span class="rounded-full border border-stone-200/80 px-3 py-1 text-xs font-medium text-stone-500 dark:border-stone-700 dark:text-stone-400">{{ t('home.onlyCodex.advantages.security.tag1') }}</span>
              <span class="rounded-full border border-stone-200/80 px-3 py-1 text-xs font-medium text-stone-500 dark:border-stone-700 dark:text-stone-400">{{ t('home.onlyCodex.advantages.security.tag2') }}</span>
            </div>
          </article>
        </div>
      </section>

      <section class="reveal-section mx-auto max-w-6xl px-6 py-10">
        <div class="reveal-item overflow-hidden rounded-[2rem] border border-blue-100/90 bg-white/78 shadow-[0_20px_70px_rgba(62,99,172,0.10)] backdrop-blur-sm dark:border-blue-400/15 dark:bg-slate-950/70">
          <div class="border-b border-blue-100/80 px-6 py-7 dark:border-white/10 md:px-8">
            <div>
              <div class="text-xs font-semibold uppercase tracking-[0.3em] text-blue-600 dark:text-blue-300">{{ t('home.onlyCodex.usagePricing.eyebrow') }}</div>
              <h2 class="mt-3 text-3xl font-black tracking-[-0.05em] md:text-4xl">{{ t('home.onlyCodex.usagePricing.titlePrefix') }}<span class="text-blue-600 dark:text-blue-300">{{ t('home.onlyCodex.usagePricing.titleHighlight') }}</span></h2>
              <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">{{ t('home.onlyCodex.usagePricing.description') }}</p>
            </div>
          </div>

          <div class="overflow-x-auto">
            <table class="usage-price-table w-full min-w-[850px] border-collapse text-left">
              <thead>
                <tr class="border-b border-blue-100/80 text-xs font-semibold uppercase tracking-[0.12em] text-stone-500 dark:border-white/10 dark:text-stone-400">
                  <th scope="col" class="px-6 py-4 md:px-8">{{ t('home.onlyCodex.usagePricing.columns.model') }}</th>
                  <th scope="col" class="px-4 py-4">{{ t('home.onlyCodex.usagePricing.columns.input') }}</th>
                  <th scope="col" class="px-4 py-4">{{ t('home.onlyCodex.usagePricing.columns.cachedInput') }}</th>
                  <th scope="col" class="px-4 py-4">{{ t('home.onlyCodex.usagePricing.columns.cacheWrite') }}</th>
                  <th scope="col" class="px-4 py-4">{{ t('home.onlyCodex.usagePricing.columns.output') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="rate in usagePricing" :key="rate.model" class="border-b border-stone-100 last:border-b-0 dark:border-white/5">
                  <th scope="row" class="whitespace-nowrap px-6 py-4 text-sm font-black tracking-[-0.025em] text-stone-800 dark:text-stone-100 md:px-8">{{ rate.model }}</th>
                  <td v-for="field in usagePriceFields" :key="`${rate.model}-${field.key}`" class="px-4 py-3">
                    <div v-if="rate.standard[field.key] !== null" class="price-compare-cell">
                      <span class="price-official"><span>{{ t('home.onlyCodex.usagePricing.official') }}</span>{{ formatUsd(rate.standard[field.key]) }}</span>
                      <span class="price-onlycodex"><span>OnlyCodex</span>{{ formatUsd(onlyCodexPrice(rate.standard[field.key])) }}</span>
                    </div>
                    <span v-else class="text-sm text-stone-300 dark:text-stone-600">—</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="flex flex-col gap-2 bg-blue-50/65 px-6 py-4 text-xs text-blue-800 dark:bg-blue-500/10 dark:text-blue-200 md:flex-row md:items-center md:justify-between md:px-8">
            <span>{{ t('home.onlyCodex.usagePricing.note') }}</span>
            <span class="font-semibold">{{ t('home.onlyCodex.usagePricing.unit') }}</span>
          </div>
        </div>
      </section>

      <section class="reveal-section mx-auto max-w-6xl px-6 py-10">
        <div class="reveal-item text-center">
          <div class="text-xs font-semibold uppercase tracking-[0.3em] text-stone-500 dark:text-stone-400">{{ t('home.onlyCodex.subscription.eyebrow') }}</div>
          <h2 class="mt-4 text-3xl font-black tracking-[-0.05em] md:text-4xl">{{ t('home.onlyCodex.subscription.title') }}</h2>
          <p class="mt-3 text-base text-stone-600 dark:text-stone-300">{{ t('home.onlyCodex.subscription.description') }}</p>
        </div>

        <div class="mx-auto mt-8 grid max-w-6xl gap-5 lg:grid-cols-3">
          <article
            v-for="card in pricingCards"
            :key="card.model"
            class="reveal-item relative overflow-hidden rounded-[1.75rem] border bg-white/80 p-5 shadow-[0_18px_60px_rgba(28,25,23,0.08)] backdrop-blur-sm transition hover:-translate-y-1 dark:bg-stone-900 dark:shadow-[0_18px_60px_rgba(0,0,0,0.3)] md:p-6"
            :class="[card.borderClass, card.featured ? 'ring-2 ring-emerald-300/50 dark:ring-emerald-400/30' : '']"
          >
            <div v-if="card.featured" class="absolute inset-x-8 top-0 h-px bg-gradient-to-r from-transparent via-emerald-300 to-transparent"></div>
            <div class="inline-flex rounded-full px-3 py-1 text-xs font-semibold" :class="card.badgeClass">
              {{ card.badge }}
            </div>
            <div class="mt-3 text-3xl font-black tracking-[-0.05em]">{{ card.model }}</div>
            <div v-if="card.officialPrice" class="mt-5 flex items-center justify-center gap-2 text-sm text-stone-500 dark:text-stone-400">
              <span>{{ t('home.onlyCodex.subscription.official') }}</span>
              <span class="font-semibold line-through decoration-stone-400/80">{{ card.officialPrice }}</span>
              <span>{{ t('home.onlyCodex.subscription.perMonth') }}</span>
            </div>
            <div class="mt-2 flex items-end justify-center gap-1.5 text-center">
              <span class="pb-1 text-xs font-bold uppercase tracking-[0.12em] text-blue-600 dark:text-blue-300">OnlyCodex</span>
              <span class="text-[2.8rem] font-black tracking-[-0.07em]">{{ card.price }}</span>
              <span class="pb-1 text-sm font-semibold text-stone-500 dark:text-stone-400">{{ t('home.onlyCodex.subscription.perMonth') }}</span>
            </div>

            <div class="mt-5 border-t border-stone-200/80 pt-5 dark:border-stone-800">
              <div
                v-for="row in card.facts"
                :key="`${card.model}-${row.label}`"
                class="grid grid-cols-[auto_1fr] items-start gap-4 py-2 text-sm"
              >
                <div class="text-stone-500 dark:text-stone-400">{{ row.label }}</div>
                <div class="text-right font-medium text-stone-800 dark:text-stone-100">{{ row.value }}</div>
              </div>
            </div>

            <div class="mt-5 border-t border-stone-200/80 pt-5 dark:border-stone-800">
              <div
                v-for="feature in card.features"
                :key="`${card.model}-${feature}`"
                class="flex items-center gap-3 py-2 text-sm font-medium text-stone-800 dark:text-stone-100"
              >
                <span class="inline-flex h-6 w-6 items-center justify-center rounded-full bg-emerald-500 text-white">
                  <svg class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4 10l4 4 8-8" />
                  </svg>
                </span>
                <span>{{ feature }}</span>
              </div>
            </div>

            <button
              type="button"
              class="mt-6 inline-flex w-full items-center justify-center rounded-full border border-emerald-300 bg-emerald-400 px-6 py-2.5 text-sm font-semibold text-emerald-950 shadow-[0_18px_40px_rgba(52,211,153,0.25)] transition hover:-translate-y-0.5 hover:bg-emerald-300"
              @click="showWechatModal = true"
            >
              {{ t('home.onlyCodex.actions.contactNow') }}
            </button>
          </article>
        </div>
      </section>

      <section class="reveal-section mx-auto max-w-6xl px-6 py-10">
        <div class="reveal-item text-center">
          <div class="text-xs font-semibold uppercase tracking-[0.3em] text-stone-500 dark:text-stone-400">{{ t('home.onlyCodex.quickStart.eyebrow') }}</div>
          <h2 class="mt-4 text-3xl font-black tracking-[-0.05em] md:text-4xl">{{ t('home.onlyCodex.quickStart.title') }}</h2>
          <p class="mt-3 text-base text-stone-600 dark:text-stone-300">{{ t('home.onlyCodex.quickStart.description') }}</p>
        </div>

        <div class="mt-10 space-y-6">
          <article
            v-for="step in stepItems"
            :key="step.badge"
            class="reveal-item flex flex-col gap-6 rounded-[2rem] border border-stone-200/80 bg-white/78 p-6 shadow-[0_18px_60px_rgba(28,25,23,0.08)] backdrop-blur-sm dark:border-stone-800 dark:bg-stone-900 md:p-8"
            :class="step.reverse ? 'md:flex-row-reverse' : 'md:flex-row'"
          >
            <div class="flex-1">
              <div class="inline-flex rounded-full border border-stone-200/90 bg-stone-50 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-stone-500 dark:border-stone-700 dark:bg-stone-950 dark:text-stone-400">
                {{ step.badge }}
              </div>
              <h3 class="mt-4 text-2xl font-black tracking-[-0.04em]">{{ step.title }}</h3>
              <p class="mt-3 max-w-xl text-sm leading-7 text-stone-600 dark:text-stone-300">{{ step.description }}</p>
              <ul class="mt-5 space-y-3">
                <li
                  v-for="check in step.checks"
                  :key="check"
                  class="flex items-start gap-3 text-sm text-stone-700 dark:text-stone-200"
                >
                  <span class="mt-0.5 inline-flex h-5 w-5 items-center justify-center rounded-full bg-emerald-500/15 text-xs font-bold text-emerald-600 dark:text-emerald-300">✓</span>
                  <span>{{ check }}</span>
                </li>
              </ul>
            </div>

            <div class="flex-1">
              <div
                v-if="step.mockType === 'register'"
                class="overflow-hidden rounded-[1.7rem] border border-stone-200/80 bg-stone-50 shadow-inner dark:border-stone-800 dark:bg-stone-950/60"
              >
                <div class="flex items-center gap-2 border-b border-stone-200/80 px-4 py-3 text-xs text-stone-500 dark:border-stone-800 dark:text-stone-400">
                  <span class="h-2.5 w-2.5 rounded-full bg-stone-300 dark:bg-stone-700"></span>
                  <span class="h-2.5 w-2.5 rounded-full bg-stone-300 dark:bg-stone-700"></span>
                  <span class="h-2.5 w-2.5 rounded-full bg-stone-300 dark:bg-stone-700"></span>
                  <span class="ml-2">onlycodex/register</span>
                </div>
                <div class="space-y-4 p-5">
                  <div class="text-2xl font-black tracking-[-0.06em]">
                    <span class="text-blue-600 dark:text-blue-300">Only</span>Codex
                  </div>
                  <div>
                    <div class="mb-2 text-xs uppercase tracking-[0.18em] text-stone-400">{{ t('home.onlyCodex.mock.email') }}</div>
                    <div class="rounded-2xl border border-stone-200/80 bg-white px-4 py-3 text-sm text-stone-600 dark:border-stone-800 dark:bg-stone-900 dark:text-stone-300">
                      your@email.com
                    </div>
                  </div>
                  <div>
                    <div class="mb-2 text-xs uppercase tracking-[0.18em] text-stone-400">{{ t('home.onlyCodex.mock.password') }}</div>
                    <div class="rounded-2xl border border-stone-200/80 bg-white px-4 py-3 text-sm text-stone-600 dark:border-stone-800 dark:bg-stone-900 dark:text-stone-300">
                      ••••••••••
                    </div>
                  </div>
                  <div class="rounded-2xl bg-emerald-400 px-4 py-3 text-center text-sm font-semibold text-emerald-950">
                    {{ t('home.onlyCodex.actions.signUpNow') }}
                  </div>
                </div>
              </div>

              <div
                v-else-if="step.mockType === 'keys'"
                class="overflow-hidden rounded-[1.7rem] border border-stone-200/80 bg-stone-50 shadow-inner dark:border-stone-800 dark:bg-stone-950/60"
              >
                <div class="flex items-center justify-between border-b border-stone-200/80 px-4 py-3 dark:border-stone-800">
                  <span class="text-sm font-semibold">API Keys</span>
                  <span class="rounded-full bg-emerald-400 px-3 py-1 text-xs font-semibold text-emerald-950">{{ t('home.onlyCodex.mock.createKey') }}</span>
                </div>
                <div class="space-y-3 p-4">
                  <div
                    v-for="key in mockKeys"
                    :key="key.name"
                    class="flex items-center justify-between rounded-2xl border border-stone-200/80 bg-white px-4 py-3 dark:border-stone-800 dark:bg-stone-900"
                  >
                    <div>
                      <div class="text-sm font-semibold">{{ key.name }}</div>
                      <div class="mt-1 text-xs text-stone-500 dark:text-stone-400">{{ key.value }}</div>
                    </div>
                    <span class="rounded-full border border-stone-200 px-3 py-1 text-xs text-stone-500 dark:border-stone-700 dark:text-stone-400">
                      {{ t('home.onlyCodex.actions.copy') }}
                    </span>
                  </div>
                </div>
              </div>

              <div
                v-else
                class="overflow-hidden rounded-[1.7rem] border border-stone-200/80 bg-[#18120f] shadow-[0_20px_60px_rgba(28,25,23,0.22)] dark:border-stone-800"
              >
                <div class="flex items-center gap-2 border-b border-white/5 px-4 py-3 text-xs text-stone-400">
                  <span class="h-2.5 w-2.5 rounded-full bg-[#ff5f56]"></span>
                  <span class="h-2.5 w-2.5 rounded-full bg-[#ffbd2e]"></span>
                  <span class="h-2.5 w-2.5 rounded-full bg-[#27c93f]"></span>
                  <span class="ml-2">terminal</span>
                </div>
                <div class="space-y-3 px-4 py-5 font-mono text-sm text-stone-200">
                  <div class="text-stone-500"><span class="mr-2 text-emerald-300">$</span>export OPENAI_BASE_URL="{{ gatewayBaseUrl }}"</div>
                  <div class="text-stone-500"><span class="mr-2 text-emerald-300">$</span>codex --model gpt-5.5</div>
                  <div class="text-emerald-300">◆ Codex CLI · GPT-5.6 Sol ready</div>
                  <div class="text-stone-400">{{ t('home.onlyCodex.terminal.connected') }}</div>
                  <div class="text-stone-400"><span class="mr-2 text-emerald-300">›</span>{{ t('home.onlyCodex.terminal.refactorPrompt') }}</div>
                  <div><span class="mr-2 text-emerald-300">$</span><span class="cursor-blink inline-block h-5 w-2 rounded-sm bg-emerald-300/80 align-middle"></span></div>
                </div>
              </div>
            </div>
          </article>
        </div>
      </section>

      <section class="reveal-section mx-auto max-w-6xl px-6 py-10">
        <div class="reveal-item text-center">
          <div class="text-xs font-semibold uppercase tracking-[0.3em] text-stone-500 dark:text-stone-400">{{ t('home.onlyCodex.metrics.eyebrow') }}</div>
          <h2 class="mt-4 text-3xl font-black tracking-[-0.05em] md:text-4xl">{{ t('home.onlyCodex.metrics.title') }}</h2>
          <p class="mt-3 text-base text-stone-600 dark:text-stone-300">{{ t('home.onlyCodex.metrics.description') }}</p>
        </div>

        <div class="mt-10 grid gap-6 lg:grid-cols-3">
          <article
            v-for="feature in featureCards"
            :key="feature.title"
            class="reveal-item overflow-hidden rounded-[2rem] border border-stone-200/80 bg-white/78 p-6 shadow-[0_18px_60px_rgba(28,25,23,0.08)] backdrop-blur-sm dark:border-stone-800 dark:bg-stone-900"
          >
            <div class="flex items-start justify-between gap-4">
              <div class="inline-flex h-12 w-12 items-center justify-center rounded-2xl" :class="feature.iconClass">
                <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6">
                  <path stroke-linecap="round" stroke-linejoin="round" :d="feature.iconPath" />
                </svg>
              </div>
              <div class="text-4xl font-black tracking-[-0.06em]">{{ feature.metric }}</div>
            </div>
            <h3 class="mt-6 text-2xl font-black tracking-[-0.04em]">{{ feature.title }}</h3>
            <p class="mt-3 text-sm leading-7 text-stone-600 dark:text-stone-300">{{ feature.description }}</p>
            <div class="mt-5 h-2 rounded-full bg-stone-200/80 dark:bg-stone-800">
              <div class="h-full rounded-full" :class="feature.barClass" :style="{ width: feature.barWidth }"></div>
            </div>
            <div class="mt-5 flex flex-wrap gap-2">
              <span
                v-for="tag in feature.tags"
                :key="tag"
                class="rounded-full border border-stone-200/80 px-3 py-1 text-xs font-medium text-stone-500 dark:border-stone-700 dark:text-stone-400"
              >
                {{ tag }}
              </span>
            </div>
          </article>
        </div>
      </section>

      <section class="reveal-section mx-auto max-w-6xl px-6 py-10">
        <div class="overflow-hidden rounded-[2.4rem] border border-stone-200/80 bg-white/78 p-8 shadow-[0_24px_80px_rgba(28,25,23,0.1)] backdrop-blur-sm dark:border-stone-800 dark:bg-stone-900 md:p-10">
          <div class="reveal-item text-center">
            <div class="text-xs font-semibold uppercase tracking-[0.3em] text-stone-500 dark:text-stone-400">{{ t('home.onlyCodex.contact.eyebrow') }}</div>
            <h2 class="mt-4 text-3xl font-black tracking-[-0.05em] md:text-4xl">{{ t('home.onlyCodex.contact.title') }}</h2>
            <p class="mt-3 text-base text-stone-600 dark:text-stone-300">
              {{ t('home.onlyCodex.contact.description') }}
            </p>
          </div>

          <div class="mt-8 flex justify-center">
            <article class="reveal-item w-full max-w-xl rounded-[1.8rem] border border-stone-200/80 bg-stone-50/90 p-6 dark:border-stone-800 dark:bg-stone-950/60">
              <div class="text-xs font-semibold uppercase tracking-[0.22em] text-stone-400">{{ t('home.onlyCodex.contact.method') }}</div>
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
                  <span>{{ copySucceeded ? t('common.copied') : t('home.onlyCodex.actions.copy') }}</span>
                </button>
              </div>
              <p class="mt-4 text-sm leading-7 text-stone-600 dark:text-stone-300">
                {{ t('home.onlyCodex.contact.hint') }}
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
            <div class="text-xs font-semibold uppercase tracking-[0.32em] text-stone-900/70">{{ t('home.onlyCodex.cta.eyebrow') }}</div>
            <h2 class="mt-4 text-3xl font-black tracking-[-0.06em] md:text-5xl">{{ t('home.onlyCodex.cta.title') }}</h2>
            <p class="mx-auto mt-4 max-w-2xl text-base leading-8 text-stone-900/80">
              {{ t('home.onlyCodex.cta.description') }}
            </p>
            <router-link
              :to="heroPrimaryTarget"
              class="mt-8 inline-flex items-center gap-2 rounded-full bg-stone-950 px-6 py-3 text-sm font-semibold text-white transition hover:-translate-y-0.5 hover:bg-stone-800"
            >
              {{ isAuthenticated ? t('home.goToDashboard') : t('home.onlyCodex.actions.signUpFree') }}
              <svg class="h-4 w-4" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 10h12M12 6l4 4-4 4" />
              </svg>
            </router-link>
          </div>
        </div>
      </section>
    </main>

    <footer class="relative z-10 border-t border-stone-200/70 px-6 py-8 text-center text-sm text-stone-500 dark:border-stone-800 dark:text-stone-400">
      © {{ currentYear }} OnlyCodex. {{ t('home.footer.allRightsReserved') }}
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
          <h3 class="text-lg font-bold text-stone-900 dark:text-stone-50">{{ t('home.onlyCodex.contact.modalTitle') }}</h3>
          <p class="mt-3 text-sm text-stone-600 dark:text-stone-300">{{ t('home.onlyCodex.contact.modalDescription') }}</p>
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
              <span>{{ copySucceeded ? t('common.copied') : t('home.onlyCodex.contact.copyContact') }}</span>
            </button>
          </div>
          <button
            type="button"
            class="mt-6 inline-flex w-full items-center justify-center rounded-full bg-emerald-400 px-6 py-2.5 text-sm font-semibold text-emerald-950 transition hover:-translate-y-0.5 hover:bg-emerald-300"
            @click="showWechatModal = false"
          >
            {{ t('home.onlyCodex.actions.close') }}
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

type PricingFact = {
  label: string
  value: string
}

type PricingCard = {
  badge: string
  model: string
  price: string
  officialPrice?: string
  facts: PricingFact[]
  features: string[]
  badgeClass: string
  borderClass: string
  featured?: boolean
}

type StepItem = {
  badge: string
  title: string
  description: string
  checks: string[]
  mockType: 'register' | 'keys' | 'terminal'
  reverse?: boolean
}

type FeatureCard = {
  metric: string
  title: string
  description: string
  tags: string[]
  iconClass: string
  barClass: string
  barWidth: string
  iconPath: string
}

type UsagePrice = {
  input: number | null
  cachedInput: number | null
  cacheWrite: number | null
  output: number | null
}

type UsagePricing = {
  model: string
  standard: UsagePrice
}

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

const pricingCards = computed<PricingCard[]>(() => [
  {
    badge: t('home.onlyCodex.plans.pro5x.badge'),
    model: 'Pro 5X',
    price: '$56',
    officialPrice: '$100',
    badgeClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300',
    borderClass: 'border-emerald-300/80 dark:border-emerald-500/30',
    featured: true,
    facts: [
      { label: t('home.onlyCodex.plans.usageExperience'), value: t('home.onlyCodex.plans.pro5x.experience') },
      { label: t('home.onlyCodex.plans.dailyLimit'), value: t('home.onlyCodex.plans.pro5x.dailyLimit') },
      { label: t('home.onlyCodex.plans.weeklyLimit'), value: t('home.onlyCodex.plans.pro5x.weeklyLimit') },
      { label: t('home.onlyCodex.plans.monthlyLimit'), value: t('home.onlyCodex.plans.pro5x.monthlyLimit') }
    ],
    features: [
      t('home.onlyCodex.plans.gpt56Available'),
      t('home.onlyCodex.plans.gpt55Available'),
      t('home.onlyCodex.plans.gpt54Available'),
      t('home.onlyCodex.plans.dedicatedKey')
    ]
  },
  {
    badge: t('home.onlyCodex.plans.pro10x.badge'),
    model: 'Pro 10X',
    price: '$90',
    badgeClass: 'bg-lime-100 text-lime-700 dark:bg-lime-500/10 dark:text-lime-300',
    borderClass: 'border-lime-200/80 dark:border-lime-500/20',
    facts: [
      { label: t('home.onlyCodex.plans.usageExperience'), value: t('home.onlyCodex.plans.pro10x.experience') },
      { label: t('home.onlyCodex.plans.dailyLimit'), value: t('home.onlyCodex.plans.pro10x.dailyLimit') },
      { label: t('home.onlyCodex.plans.weeklyLimit'), value: t('home.onlyCodex.plans.pro10x.weeklyLimit') },
      { label: t('home.onlyCodex.plans.monthlyLimit'), value: t('home.onlyCodex.plans.pro10x.monthlyLimit') }
    ],
    features: [
      t('home.onlyCodex.plans.gpt56Available'),
      t('home.onlyCodex.plans.gpt55Available'),
      t('home.onlyCodex.plans.gpt54Available'),
      t('home.onlyCodex.plans.dedicatedKey')
    ]
  },
  {
    badge: t('home.onlyCodex.plans.pro20x.badge'),
    model: 'Pro 20X',
    price: '$180',
    officialPrice: '$200',
    badgeClass: 'bg-teal-100 text-teal-700 dark:bg-teal-500/10 dark:text-teal-300',
    borderClass: 'border-teal-200/80 dark:border-teal-500/20',
    facts: [
      { label: t('home.onlyCodex.plans.usageExperience'), value: t('home.onlyCodex.plans.pro20x.experience') },
      { label: t('home.onlyCodex.plans.dailyLimit'), value: t('home.onlyCodex.plans.pro20x.dailyLimit') },
      { label: t('home.onlyCodex.plans.weeklyLimit'), value: t('home.onlyCodex.plans.pro20x.weeklyLimit') },
      { label: t('home.onlyCodex.plans.monthlyLimit'), value: t('home.onlyCodex.plans.pro20x.monthlyLimit') }
    ],
    features: [
      t('home.onlyCodex.plans.gpt56Available'),
      t('home.onlyCodex.plans.gpt55Available'),
      t('home.onlyCodex.plans.gpt54Available'),
      t('home.onlyCodex.plans.dedicatedKey')
    ]
  }
])

const stepItems = computed<StepItem[]>(() => [
  {
    badge: t('home.onlyCodex.steps.register.badge'),
    title: t('home.onlyCodex.steps.register.title'),
    description: t('home.onlyCodex.steps.register.description'),
    checks: [t('home.onlyCodex.steps.register.check1'), t('home.onlyCodex.steps.register.check2')],
    mockType: 'register'
  },
  {
    badge: t('home.onlyCodex.steps.key.badge'),
    title: t('home.onlyCodex.steps.key.title'),
    description: t('home.onlyCodex.steps.key.description'),
    checks: [t('home.onlyCodex.steps.key.check1'), t('home.onlyCodex.steps.key.check2'), t('home.onlyCodex.steps.key.check3')],
    mockType: 'keys',
    reverse: true
  },
  {
    badge: t('home.onlyCodex.steps.launch.badge'),
    title: t('home.onlyCodex.steps.launch.title'),
    description: t('home.onlyCodex.steps.launch.description'),
    checks: [t('home.onlyCodex.steps.launch.check1'), t('home.onlyCodex.steps.launch.check2')],
    mockType: 'terminal'
  }
])

const featureCards = computed<FeatureCard[]>(() => [
  {
    metric: '99.9%',
    title: t('home.onlyCodex.metrics.uptime.title'),
    description: t('home.onlyCodex.metrics.uptime.description'),
    tags: [t('home.onlyCodex.metrics.uptime.tag1'), t('home.onlyCodex.metrics.uptime.tag2')],
    iconClass: 'bg-lime-500/15 text-lime-700 dark:text-lime-300',
    barClass: 'bg-lime-400',
    barWidth: '100%',
    iconPath: 'M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z'
  },
  {
    metric: '<50ms',
    title: t('home.onlyCodex.metrics.latency.title'),
    description: t('home.onlyCodex.metrics.latency.description'),
    tags: [t('home.onlyCodex.metrics.latency.tag1'), t('home.onlyCodex.metrics.latency.tag2')],
    iconClass: 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-300',
    barClass: 'bg-emerald-500',
    barWidth: '30%',
    iconPath: 'M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z'
  },
  {
    metric: '10min',
    title: t('home.onlyCodex.metrics.onboarding.title'),
    description: t('home.onlyCodex.metrics.onboarding.description'),
    tags: [t('home.onlyCodex.metrics.onboarding.tag1'), t('home.onlyCodex.metrics.onboarding.tag2')],
    iconClass: 'bg-sky-500/15 text-sky-600 dark:text-sky-300',
    barClass: 'bg-sky-500',
    barWidth: '15%',
    iconPath: 'M6.75 7.5l3 2.25-3 2.25m4.5 0h3m-9 8.25h13.5A2.25 2.25 0 0021 18V6a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 6v12a2.25 2.25 0 002.25 2.25z'
  }
])

const usagePriceFields: Array<{ key: keyof UsagePrice }> = [
  { key: 'input' },
  { key: 'cachedInput' },
  { key: 'cacheWrite' },
  { key: 'output' }
]

const usagePricing: UsagePricing[] = [
  { model: 'gpt-5.6-sol', standard: { input: 5, cachedInput: 0.5, cacheWrite: 6.25, output: 30 } },
  { model: 'gpt-5.6-terra', standard: { input: 2.5, cachedInput: 0.25, cacheWrite: 3.125, output: 15 } },
  { model: 'gpt-5.6-luna', standard: { input: 1, cachedInput: 0.1, cacheWrite: 1.25, output: 6 } },
  { model: 'gpt-5.5', standard: { input: 5, cachedInput: 0.5, cacheWrite: null, output: 30 } },
  { model: 'gpt-5.5-pro', standard: { input: 30, cachedInput: null, cacheWrite: null, output: 180 } },
  { model: 'gpt-5.4', standard: { input: 2.5, cachedInput: 0.25, cacheWrite: null, output: 15 } },
  { model: 'gpt-5.4-mini', standard: { input: 0.75, cachedInput: 0.075, cacheWrite: null, output: 4.5 } },
  { model: 'gpt-5.4-nano', standard: { input: 0.2, cachedInput: 0.02, cacheWrite: null, output: 1.25 } },
  { model: 'gpt-5.4-pro', standard: { input: 30, cachedInput: null, cacheWrite: null, output: 180 } }
]

const mockKeys = computed(() => [
  { name: t('home.onlyCodex.mock.firstKey'), value: 'sk-tocodex-Kx9m••••••••••••' },
  { name: t('home.onlyCodex.mock.projectB'), value: 'sk-tocodex-Rp3n••••••••••••' }
])

const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const gatewayBaseUrl = computed(() => appStore.cachedPublicSettings?.api_base_url || appStore.apiBaseUrl || 'https://onlycodex.dev')
const contactValue = computed(() => appStore.cachedPublicSettings?.contact_info?.trim() || 'itwillbe626')

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

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function onlyCodexPrice(value: number | null): number | null {
  return value === null ? null : value * 0.125
}

function formatUsd(value: number | null): string {
  if (value === null) return '—'
  return `$${value.toFixed(4).replace(/\.?0+$/, '')}`
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
.onlycodex-home {
  --onlycodex-blue-50: 239 246 255;
  --onlycodex-blue-100: 219 234 254;
  --onlycodex-blue-200: 191 219 254;
  --onlycodex-blue-300: 147 197 253;
  --onlycodex-blue-400: 96 165 250;
  --onlycodex-blue-500: 79 143 240;
  --onlycodex-blue-600: 55 119 223;
  --onlycodex-blue-700: 40 94 194;
  --onlycodex-blue-950: 23 37 84;
}

/* Keep the existing information architecture, while translating its former green
   accents into the airy blue / violet OnlyCodex visual system. */
.onlycodex-home :deep([class~='bg-emerald-50']),
.onlycodex-home :deep([class~='bg-emerald-100']) { background-color: rgb(var(--onlycodex-blue-50) / var(--tw-bg-opacity, 1)) !important; }
.onlycodex-home :deep([class~='bg-emerald-400']),
.onlycodex-home :deep([class~='bg-emerald-500']) { background-color: rgb(var(--onlycodex-blue-500) / var(--tw-bg-opacity, 1)) !important; }
.onlycodex-home :deep([class~='bg-emerald-500/10']),
.onlycodex-home :deep([class~='bg-emerald-500/15']) { background-color: rgb(var(--onlycodex-blue-500) / .12) !important; }
.onlycodex-home :deep([class~='dark:bg-emerald-300']) { background-color: rgb(var(--onlycodex-blue-300)) !important; }
.onlycodex-home :deep([class~='dark:bg-emerald-500/10']),
.onlycodex-home :deep([class~='dark:bg-emerald-500/15']) { background-color: rgb(var(--onlycodex-blue-500) / .16) !important; }
.onlycodex-home :deep([class~='bg-lime-100']),
.onlycodex-home :deep([class~='bg-lime-500/15']) { background-color: rgb(139 92 246 / .12) !important; }
.onlycodex-home :deep([class~='bg-teal-100']),
.onlycodex-home :deep([class~='bg-teal-500/15']) { background-color: rgb(56 189 248 / .12) !important; }

.onlycodex-home :deep([class~='text-emerald-300']),
.onlycodex-home :deep([class~='text-emerald-400']) { color: rgb(var(--onlycodex-blue-300)) !important; }
.onlycodex-home :deep([class~='text-emerald-500']),
.onlycodex-home :deep([class~='text-emerald-600']),
.onlycodex-home :deep([class~='text-emerald-700']) { color: rgb(var(--onlycodex-blue-600)) !important; }
.onlycodex-home :deep([class~='text-emerald-950']) { color: rgb(var(--onlycodex-blue-950)) !important; }
.onlycodex-home :deep([class~='dark:text-emerald-300']) { color: rgb(var(--onlycodex-blue-300)) !important; }
.onlycodex-home :deep([class~='dark:text-emerald-400']) { color: rgb(var(--onlycodex-blue-400)) !important; }
.onlycodex-home :deep([class~='dark:text-emerald-950']) { color: rgb(var(--onlycodex-blue-950)) !important; }
.onlycodex-home :deep([class~='text-lime-300']) { color: rgb(196 181 253) !important; }
.onlycodex-home :deep([class~='text-lime-700']) { color: rgb(109 40 217) !important; }
.onlycodex-home :deep([class~='text-teal-300']),
.onlycodex-home :deep([class~='text-teal-700']) { color: rgb(14 116 144) !important; }

.onlycodex-home :deep([class~='border-emerald-200']),
.onlycodex-home :deep([class~='border-emerald-300']) { border-color: rgb(var(--onlycodex-blue-200) / var(--tw-border-opacity, 1)) !important; }
.onlycodex-home :deep([class~='border-emerald-200/70']),
.onlycodex-home :deep([class~='border-emerald-200/80']),
.onlycodex-home :deep([class~='border-emerald-300/60']),
.onlycodex-home :deep([class~='border-emerald-300/80']) { border-color: rgb(var(--onlycodex-blue-200) / .8) !important; }
.onlycodex-home :deep([class~='border-emerald-500/20']),
.onlycodex-home :deep([class~='border-emerald-500/30']) { border-color: rgb(var(--onlycodex-blue-500) / .25) !important; }
.onlycodex-home :deep([class~='border-lime-200']) { border-color: rgb(221 214 254 / var(--tw-border-opacity, 1)) !important; }
.onlycodex-home :deep([class~='border-teal-200']) { border-color: rgb(186 230 253 / var(--tw-border-opacity, 1)) !important; }

.onlycodex-home :deep([class~='from-emerald-300']) { --tw-gradient-from: rgb(var(--onlycodex-blue-300)) var(--tw-gradient-from-position) !important; --tw-gradient-to: rgb(var(--onlycodex-blue-300) / 0) var(--tw-gradient-to-position) !important; }
.onlycodex-home :deep([class~='from-emerald-600']) { --tw-gradient-from: rgb(var(--onlycodex-blue-600)) var(--tw-gradient-from-position) !important; --tw-gradient-to: rgb(var(--onlycodex-blue-600) / 0) var(--tw-gradient-to-position) !important; }
.onlycodex-home :deep([class~='via-emerald-300']) { --tw-gradient-stops: var(--tw-gradient-from), rgb(var(--onlycodex-blue-300)) var(--tw-gradient-via-position), var(--tw-gradient-to) !important; }
.onlycodex-home :deep([class~='via-lime-400']) { --tw-gradient-stops: var(--tw-gradient-from), rgb(167 139 250) var(--tw-gradient-via-position), var(--tw-gradient-to) !important; }
.onlycodex-home :deep([class~='to-emerald-500']) { --tw-gradient-to: rgb(var(--onlycodex-blue-500)) var(--tw-gradient-to-position) !important; }
.onlycodex-home :deep([class~='to-teal-300']) { --tw-gradient-to: rgb(125 211 252) var(--tw-gradient-to-position) !important; }
.onlycodex-home :deep([class~='to-lime-200']) { --tw-gradient-to: rgb(221 214 254) var(--tw-gradient-to-position) !important; }
.onlycodex-home :deep([class~='via-lime-200']) { --tw-gradient-stops: var(--tw-gradient-from), rgb(221 214 254) var(--tw-gradient-via-position), var(--tw-gradient-to) !important; }

.price-compare-cell {
  display: grid;
  min-width: 8.3rem;
  gap: .35rem;
  border-left: 2px solid rgb(219 234 254);
  padding-left: .7rem;
}

.price-official,
.price-onlycodex {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: .65rem;
  font-variant-numeric: tabular-nums;
  font-size: .82rem;
  font-weight: 700;
}

.price-official { color: #64748b; }
.price-onlycodex { color: #2563eb; }

.price-official span,
.price-onlycodex span {
  font-size: .64rem;
  font-weight: 700;
  letter-spacing: .03em;
}

.price-official span { color: #94a3b8; }
.price-onlycodex span { color: #5b8fe8; }

:global(.dark) .price-compare-cell { border-color: rgb(96 165 250 / .28); }
:global(.dark) .price-official { color: #94a3b8; }
:global(.dark) .price-onlycodex { color: #93c5fd; }
:global(.dark) .price-official span { color: #64748b; }
:global(.dark) .price-onlycodex span { color: #93c5fd; }

.home-grid-mask {
  background-image:
    linear-gradient(rgba(79, 143, 240, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(79, 143, 240, 0.08) 1px, transparent 1px);
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

@media (max-width: 640px) {
  .usage-price-table th:first-child,
  .usage-price-table td:first-child { position: sticky; left: 0; z-index: 1; background: rgb(255 255 255 / .96); }
  :global(.dark) .usage-price-table th:first-child,
  :global(.dark) .usage-price-table td:first-child { background: rgb(15 23 42 / .96); }
}
</style>
