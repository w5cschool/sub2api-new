<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white">
    <header class="border-b border-gray-200 bg-white/95 dark:border-dark-800 dark:bg-dark-900/95">
      <div class="mx-auto flex max-w-5xl items-center justify-between gap-4 px-4 py-4 sm:px-6">
        <RouterLink to="/home" class="flex min-w-0 items-center gap-3">
          <span class="flex h-10 w-10 flex-shrink-0 items-center justify-center overflow-hidden rounded-xl bg-white shadow-sm ring-1 ring-gray-200 dark:bg-dark-800 dark:ring-dark-700">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </span>
          <span class="truncate text-base font-semibold text-gray-950 dark:text-white">
            {{ siteName }}
          </span>
        </RouterLink>
        <RouterLink
          to="/login"
          class="inline-flex flex-shrink-0 items-center justify-center rounded-lg bg-primary-600 px-4 py-2 text-sm font-semibold text-white shadow-sm shadow-primary-600/20 transition hover:bg-primary-700"
        >
          {{ t('home.login') }}
        </RouterLink>
      </div>
    </header>

    <main class="mx-auto max-w-4xl px-4 py-8 sm:px-6 lg:py-10">
      <div v-if="loading" class="flex min-h-[320px] items-center justify-center">
        <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
      </div>

      <section
        v-else-if="loadError"
        class="rounded-lg border border-red-200 bg-red-50 p-6 text-red-700 dark:border-red-500/30 dark:bg-red-500/10 dark:text-red-200"
      >
        <h1 class="text-lg font-semibold">{{ t('legal.loadFailed') }}</h1>
        <p class="mt-2 text-sm">{{ t('legal.retryLater') }}</p>
      </section>

      <section
        v-else-if="!currentDocument"
        class="rounded-lg border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-900"
      >
        <div class="flex items-start gap-3">
          <span class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-md bg-gray-100 text-gray-600 dark:bg-dark-800 dark:text-dark-300">
            <Icon name="document" size="sm" />
          </span>
          <div>
            <h1 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('legal.notFound') }}</h1>
            <p class="mt-2 text-sm leading-6 text-gray-600 dark:text-dark-300">
              {{ t('legal.notFoundDescription') }}
            </p>
          </div>
        </div>
      </section>

      <article v-else>
        <div class="mb-8 border-b border-gray-200 pb-6 dark:border-dark-700">
          <div class="flex items-start gap-4">
            <span class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-md bg-primary-50 text-primary-700 dark:bg-primary-500/10 dark:text-primary-300">
              <Icon :name="documentIcon" size="md" />
            </span>
            <div class="min-w-0">
              <p class="text-sm font-medium text-primary-700 dark:text-primary-300">{{ documentTypeLabel }}</p>
              <h1 class="mt-2 break-words text-2xl font-bold tracking-normal text-gray-950 dark:text-white sm:text-3xl">
                {{ currentDocument.title }}
              </h1>
              <p v-if="updatedAt" class="mt-3 text-sm text-gray-500 dark:text-dark-400">
                {{ t('legal.updatedAt', { date: updatedAt }) }}
              </p>
            </div>
          </div>
        </div>

        <div
          v-if="hasContent"
          class="legal-document-content"
          v-html="renderedHtml"
        ></div>
        <div
          v-else
          class="rounded-lg border border-dashed border-gray-300 bg-white px-6 py-14 text-center text-sm text-gray-500 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-400"
        >
          {{ t('legal.empty') }}
        </div>
      </article>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { getPublicSettings } from '@/api/auth'
import { getLocale } from '@/i18n'
import { sanitizeUrl } from '@/utils/url'
import type { LoginAgreementDocument, PublicSettings } from '@/types'
import zhAdminCompliance from '../../../../docs/legal/admin-compliance.zh.md?raw'
import enAdminCompliance from '../../../../docs/legal/admin-compliance.en.md?raw'

type LegalDocumentIcon = 'document' | 'shield' | 'globe' | 'cog'

const operatorName = 'Pluto Horizon LLC'
const operatorAddress = '30 N Gould St #65046, Sheridan, Wyoming 82801, United States'

const publicLegalDocuments: LoginAgreementDocument[] = [
  {
    id: 'terms',
    title: '服务条款',
    content_md: `# 服务条款

**运营主体：** ${operatorName}

**注册地址：** ${operatorAddress}

本网站及其控制台由 ${operatorName} 独立运营，用于管理已获授权的 API 接入。本站不代表、隶属于或获任何第三方模型、支付、身份或账户服务品牌授权，除非双方以书面协议明确约定。

## 使用条件

您仅可将本站用于您已获得授权、且符合适用法律、合同和第三方服务条款的用途。您不得利用本站冒充第三方、规避访问限制、收集他人凭据，或实施任何违法、欺诈或侵权行为。

## 账户与凭据

请妥善保管本站账户与 API 凭据。本站不会在公开页面要求您提交第三方账户密码、验证码、银行卡信息或私钥。发现可疑请求时，请立即停止操作并通过首页公开联系方式报告。

## 政策与联系

使用本站前，请一并阅读《隐私政策》及《退款与取消政策》。如您不同意这些条款，请停止使用本站服务。`,
  },
  {
    id: 'privacy-policy',
    title: '隐私政策',
    content_md: `# 隐私政策

**运营主体：** ${operatorName}

**注册地址：** ${operatorAddress}

本政策说明 ${operatorName} 在提供本站服务时处理信息的基本原则。

## 可能处理的信息

为提供账户、身份验证、访问控制、故障排查和用量管理功能，本站可能处理您提供的联系信息、登录与安全事件、API 请求所需的认证信息、使用记录以及为完成请求而提交的内容。请不要在公开页面提交任何第三方账户密码、验证码、银行卡信息或私钥。

## 使用与披露

信息仅用于提供、维护、保护和改进本站服务，以及履行适用法律义务。我们可能使用受授权的基础设施或服务提供商处理必要数据；除非法律要求、取得您的同意或为提供服务所必需，不会将信息出售给无关第三方。

## 安全与您的选择

我们采取合理的技术和组织措施保护信息，但任何网络传输或存储方式都不能保证绝对安全。您可通过首页公开的运营方联系方式咨询个人信息、账户或安全问题；在适用法律要求的范围内，我们会处理相关请求。`,
  },
  {
    id: 'refund-cancellation-policy',
    title: '退款与取消政策',
    content_md: `# 退款与取消政策

**运营主体：** ${operatorName}

**注册地址：** ${operatorAddress}

本政策适用于由 ${operatorName} 直接提供且明确标示为可购买的服务。

## 取消

如服务提供取消入口，您可以按照购买页面、订单说明或与 ${operatorName} 的书面协议办理取消。取消生效时间及已产生服务费用以购买时展示的条款或书面协议为准。

## 退款申请

仅当因本站原因导致无法继续提供服务时，您可以通过首页公开的运营方联系方式提交订单编号、购买账户和申请原因，申请退款。符合退款条件的订单按以下规则计算：

1. 先扣除订单实付金额的 10% 作为服务费；
2. 再按照订阅剩余天数占订阅总天数的比例，计算可退金额；
3. 除上述情形外，其他情况均不支持退款。

退款计算以订单实际支付金额、订阅周期和服务停止时点为准。

## 法定权利

本政策不排除或限制适用法律赋予您的任何不可放弃的消费者权利。`,
  },
]

const route = useRoute()
const { t } = useI18n()
const settings = ref<PublicSettings | null>(null)
const loading = ref(true)
const loadError = ref(false)

marked.setOptions({
  breaks: true,
  gfm: true,
})

const documentId = computed(() => String(route.params.documentId || ''))
const isAdminComplianceDocument = computed(() => documentId.value === 'admin-compliance')
const documents = computed(() => settings.value?.login_agreement_documents ?? [])
const siteName = computed(() => settings.value?.site_name || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(settings.value?.site_logo || '', {
  allowRelative: true,
  allowDataUrl: true,
}))
const updatedAt = computed(() =>
  isAdminComplianceDocument.value ? '' : settings.value?.login_agreement_updated_at || ''
)
const documentTypeLabel = computed(() =>
  isAdminComplianceDocument.value ? t('legal.adminCompliance') : t('legal.loginAgreement')
)

const currentDocument = computed<LoginAgreementDocument | null>(() => {
  if (isAdminComplianceDocument.value) {
    return {
      id: 'admin-compliance',
      title: t('adminCompliance.title'),
      content_md: getLocale() === 'zh' ? zhAdminCompliance : enAdminCompliance
    }
  }
  const id = documentId.value
  if (!id) {
    return null
  }
  return documents.value.find((doc) => doc.id === id && doc.content_md.trim())
    ?? publicLegalDocuments.find((doc) => doc.id === id)
    ?? null
})

const hasContent = computed(() => Boolean(currentDocument.value?.content_md?.trim()))

const renderedHtml = computed(() => {
  const content = currentDocument.value?.content_md?.trim() || ''
  if (!content) {
    return ''
  }
  const html = marked.parse(content) as string
  return DOMPurify.sanitize(html)
})

const documentIcon = computed<LegalDocumentIcon>(() => {
  const title = currentDocument.value?.title || ''
  if (title.includes('政策') || title.includes('隐私')) {
    return 'shield'
  }
  if (title.includes('国家') || title.includes('地区')) {
    return 'globe'
  }
  if (title.includes('特定')) {
    return 'cog'
  }
  return 'document'
})

onMounted(async () => {
  loading.value = true
  loadError.value = false
  try {
    settings.value = await getPublicSettings()
  } catch {
    loadError.value = true
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.legal-document-content {
  line-height: 1.75;
  overflow-wrap: anywhere;
  color: inherit;
}

.legal-document-content :deep(h1) {
  @apply mb-4 mt-8 border-b border-gray-200 pb-3 text-3xl font-bold dark:border-dark-700;
}

.legal-document-content :deep(h2) {
  @apply mb-3 mt-7 text-2xl font-bold;
}

.legal-document-content :deep(h3) {
  @apply mb-2 mt-6 text-xl font-semibold;
}

.legal-document-content :deep(h4) {
  @apply mb-2 mt-5 text-lg font-semibold;
}

.legal-document-content :deep(p) {
  @apply mb-4 text-gray-700 dark:text-dark-200;
}

.legal-document-content :deep(a) {
  @apply text-primary-600 underline underline-offset-4 hover:text-primary-700 dark:text-primary-300 dark:hover:text-primary-200;
}

.legal-document-content :deep(ul) {
  @apply mb-4 list-disc pl-6;
}

.legal-document-content :deep(ol) {
  @apply mb-4 list-decimal pl-6;
}

.legal-document-content :deep(li) {
  @apply mb-1 text-gray-700 dark:text-dark-200;
}

.legal-document-content :deep(blockquote) {
  @apply my-5 border-l-4 border-gray-300 pl-4 text-gray-600 dark:border-dark-600 dark:text-dark-300;
}

.legal-document-content :deep(code) {
  @apply rounded bg-gray-100 px-1.5 py-0.5 font-mono text-sm dark:bg-dark-800;
}

.legal-document-content :deep(pre) {
  @apply my-5 overflow-x-auto rounded-lg bg-gray-950 p-4 text-gray-100;
}

.legal-document-content :deep(pre code) {
  @apply bg-transparent p-0 text-inherit;
}

.legal-document-content :deep(table) {
  @apply my-5 block w-full overflow-x-auto border-collapse;
}

.legal-document-content :deep(th) {
  @apply border border-gray-300 bg-gray-50 px-3 py-2 text-left font-semibold dark:border-dark-600 dark:bg-dark-800;
}

.legal-document-content :deep(td) {
  @apply border border-gray-300 px-3 py-2 dark:border-dark-600;
}

.legal-document-content :deep(img) {
  @apply my-5 h-auto max-w-full rounded-lg;
}

.legal-document-content :deep(hr) {
  @apply my-7 border-gray-200 dark:border-dark-700;
}
</style>
