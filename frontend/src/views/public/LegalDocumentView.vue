<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white">
    <header class="border-b border-gray-200 bg-white/95 dark:border-dark-800 dark:bg-dark-900/95">
      <div class="mx-auto flex max-w-5xl items-center justify-between gap-4 px-4 py-4 sm:px-6">
        <RouterLink to="/home" class="flex min-w-0 items-center gap-3">
          <template v-if="settings">
            <span class="flex h-10 w-10 flex-shrink-0 items-center justify-center overflow-hidden rounded-xl bg-white shadow-sm ring-1 ring-gray-200 dark:bg-dark-800 dark:ring-dark-700">
              <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
            </span>
            <span class="truncate text-base font-semibold text-gray-950 dark:text-white">
              {{ siteName }}
            </span>
          </template>
          <template v-else>
            <span class="h-10 w-10 flex-shrink-0 animate-pulse rounded-xl bg-gray-200 dark:bg-dark-700" aria-hidden="true"></span>
            <span class="h-5 w-28 animate-pulse rounded bg-gray-200 dark:bg-dark-700" aria-hidden="true"></span>
          </template>
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
import { getLocale } from '@/i18n'
import { sanitizeUrl } from '@/utils/url'
import { useAppStore } from '@/stores/app'
import type { LoginAgreementDocument } from '@/types'
import zhAdminCompliance from '../../../../docs/legal/admin-compliance.zh.md?raw'
import enAdminCompliance from '../../../../docs/legal/admin-compliance.en.md?raw'

type LegalDocumentIcon = 'document' | 'shield' | 'globe' | 'cog'

const operatorName = 'Pluto Horizon LLC'
const operatorAddress = '30 N Gould St #65046, Sheridan, Wyoming 82801, United States'
const supportEmail = 'steveweng.s.w@gmail.com'
const administratorWechat = 'Icanmeetu'

const publicLegalDocuments: LoginAgreementDocument[] = [
  {
    id: 'terms',
    title: '服务条款',
    content_md: `# 使用条款

本网站及服务由 **${operatorName}** 提供。注册地址：${operatorAddress}。管理员微信：${administratorWechat}。客户服务邮箱：${supportEmail}。

## 第一部分：定义与限制

### 1.1 定义

“争议”是指因本条款引起或与本条款以任何方式相关的任何索赔、冲突、争论、分歧或涉嫌违约行为。

“重大违约”是指发生该违约后，处于非违约方立场的合理人士会希望立即终止本条款的任何违约行为。

“条款”和/或“协议”是指本条款与条件。

“我们”“本公司”和/或“我们的”是指 ${operatorName} 及 apiclub.cc。

“网站”是指域名为 apiclub.cc 的网站，包括所有子域名。

“您”和/或“用户”是指网站和/或服务的用户。

### 1.2 同意受约束；授权许可

以下条款与条件，以及本网站上列明的相关信息，包括任何可用功能、服务及相关服务，均受以下条款与条件约束。请您仔细阅读。任何使用本网站和/或服务的行为，均表示用户同意受本条款约束，无需另行确认。

本条款与条件受隐私政策约束，隐私政策同样适用于您对本网站和/或服务的使用。您确认并同意，我们的每一位供应商均为本条款与条件的第三方受益人，并可直接执行任何赋予其利益的条款。

## 第二部分：一般规定

### 2.1 信息的准确性、完整性与及时性

如果网站和/或服务上提供的信息不准确、不完整或不是最新的，我们对此不承担责任。网站和服务仅提供一般信息，不应作为您作出决定的唯一依据。在作出决定前，您应咨询更准确、更完整或更及时的信息来源。

### 2.2 网站和/或服务中的错误

我们不保证网站和/或服务中的任何错误都会被纠正。

### 2.3 条款与条件的修改和变更

我们保留自行决定修改或替换本条款与条件的权利，并会通过在网站上发布更新后的条款来进行修改。如果修改构成对条款与条件的重大变更，我们将通过电子邮件通知您。

### 2.4 网站和/或服务的修改和变更

我们保留自行决定随时修改或替换网站和/或服务的权利。如果修订属于重大变更，我们将在新条款生效前至少提前 30 天通知您。

### 2.5 网站和/或服务的访问

虽然我们会尽力使网站和服务每天 24 小时、每周 7 天可用，计划维护除外，但我们不保证网站和服务始终可用。我们不保证您的设备能够访问和/或支持网站或所有服务。

### 2.6 拒绝、限制、中止及终止权利

我们保留随时自行决定因任何理由拒绝向您提供网站和/或服务访问权限的权利。我们可自行决定因任何原因限制或取消用户账户。如果我们更改或取消账户，可能会尝试通知您；但未能通知并不构成任何责任，包括因账户终止导致数据删除而产生的损失责任。

### 2.7 网站和服务的禁止用途

您不得将网站和/或服务用于任何非法目的；不得引诱他人实施非法行为；不得违反法律或条例；不得侵犯知识产权；不得骚扰、辱骂、诽谤、恐吓或歧视他人；不得提交虚假信息；不得上传恶意代码；不得收集或追踪他人的个人信息；不得发送垃圾信息、钓鱼、域欺骗、冒充、使用蜘蛛程序、抓取或爬取数据；不得用于任何淫秽或不道德目的；不得干扰或规避网站和/或服务的安全功能。

此外，您不得以第三方服务使用政策所禁止的任何目的或方式使用 apiclub.cc。我们保留因您违反任何禁止用途或其他任何原因而自行决定终止您使用网站和/或服务的权利。

## 第三部分：账户与服务使用

### 3.1 在线账户

用户可能有机会通过在线注册表创建用户账户，以参与网站和服务的某些功能。注册时，您声明并保证您提供的信息据您所知是最新、完整且准确的，并同意维护并及时更新这些信息。

在注册过程中，您可能需要选择密码。您确认并同意，我们可以依赖该密码来识别您的身份。您应对您账户的所有使用行为负责，无论该访问或使用是否由您授权，并应确保您账户的所有使用行为均符合本条款与条件。

### 3.2 禁止多个账户

您同意您不得拥有超过一个账户，也不得出售、交易或转让该账户给任何其他个人或实体。

### 3.3 账户准则

网站和/或服务可能包含用户与第三方可以互动交流的功能。参与此类互动区域时，您同意不得上传非法、诽谤、淫秽、色情、辱骂或其他违法材料；不得威胁或口头辱骂其他用户；不得使用歧视性或仇恨性语言；不得对其他用户进行人身攻击；不得发布招揽信息、垃圾信息或未经授权的商业通信；不得上传违反法律、鼓吹非法活动或包含病毒及破坏性组件的内容；不得冒充任何个人或实体；不得收集可识别个人身份的信息；不得从事非法多层次营销；不得干扰网站、服务或相关网络；不得协助或鼓励任何违反本条款与条件或我们政策的行为。

用户进一步同意并确认，其创建的任何个人资料可能会被我们编辑、删除、修改、发布、传输和展示，并且用户放弃因相关材料被以其不同意的方式更改或修改而可能享有的任何权利。

### 3.4 提交内容的权利

如果您向互动区域提交、展示、发布或以其他方式上传任何内容，您即授予我们以及我们的合作伙伴和关联方一项有限、非独占、可再许可、全球范围、已全额支付、免版税的许可，用于托管、索引、缓存、分发、标记、营销以及所有其他合法目的，包括使用、修改、公开表演、公开展示、复制和分发该等提交内容。

您声明并保证您拥有或已获得有效许可，可以使用所有提交内容，并且有权授予本条款所述许可。

### 3.5 监控权

我们有权自行决定对您的账户进行监控。

## 第四部分：费用与支付

### 4.1 费用

使用网站本身是免费的；但是，服务的某些功能可能需要支付费用。相关费用会在网站上和/或购买时列明。

### 4.2 支付

付款应根据您注册时选择的会员等级条款到期并支付。

### 4.3 退款

除《退款与取消政策》另有规定外，任何情况下均不予退款，包括预付服务未被完全使用的情况。退款处理最长可能需要 30 个自然日。

## 第五部分：知识产权

### 5.1 知识产权未被放弃

本协议仅为您访问和使用网站和/或服务的协议，本条款与条件并未授予您任何软件许可。网站和服务受美国以及适用情况下的国际知识产权法律保护。网站和服务属于我们，且为我们或我们的许可方所有。我们保留网站和服务的全部所有权。

本网站和/或服务上展示或传输的所有材料，包括文字、照片、图片、插图、视频片段、音频片段和图形，均由我们拥有，并受版权、商标、服务标识及其他专有权利保护。除本条款允许的情况外，您不得复制、再制作、发布、传输、转让、出售、出租、修改、基于其创作衍生作品、分发、重新发布、表演、展示或商业利用网站上的材料，也不得删除或更改材料上的任何版权、商标、其他专有声明、视觉标识或标志。

您可以打印一份我们在本网站和/或服务上提供的材料，仅供个人、非商业用途使用，前提是您不得删除材料上的任何版权、商标、其他专有声明、视觉标识或标志。如需获得存档、保留或重新发布任何材料部分的许可，可通过 ${supportEmail} 提交申请。

### 5.2 反馈

如果您提供评论、建议、想法或反馈，您即授予我们对该等内容的唯一所有权，包括无需向您支付补偿即可使用、利用、复制、传输、发布、分发、公开展示、公开表演、创作衍生作品、托管、索引、缓存、标记、编码、修改和改编该等反馈的权利。所有此类反馈均视为非保密信息。您确认，您的想法可能与我们从第三方收到的想法或我们自行开发的想法相似，我们使用任何类似想法均不对您承担义务。

## 第六部分：第三方广告、推广与链接

### 6.1 第三方广告和推广

我们可能不时在网站和/或服务上投放来自第三方的广告和推广活动。您与非本公司的广告商之间的交易、通信或参与推广活动，均仅发生在您与该第三方之间。我们不对因此类交易产生的任何损失或损害承担责任。

### 6.2 第三方工具的使用

我们可能向您提供第三方工具的访问权限，而我们既不监控，也无法控制或参与这些工具。您确认并同意，我们以“现状”和“可用”为基础提供此类工具访问权限，不作任何形式的保证、声明或条件，也不代表我们对其进行认可。您通过网站使用任何可选第三方工具，完全由您自行承担风险并自行决定。

### 6.3 第三方链接

通过我们网站和/或服务提供的某些内容、产品和服务可能包含来自第三方的材料。第三方链接可能会将您导向与我们无关联的网站和/或服务。我们不负责审查或评估任何第三方材料或网站的内容或准确性，也不对与任何第三方相关的商品、服务、资源、内容或其他交易的购买或使用所造成的任何损害承担责任。

## 第七部分：免责声明、责任限制与赔偿

### 7.1 保证免责声明；责任限制

您同意，使用网站和服务的风险完全由您自行承担。网站、服务以及任何可下载软件、产品或其他材料均按“现状”提供，不作任何明示或默示保证，但适用法律规定不得排除、限制或修改的默示保证除外。

在任何情况下，我们、我们的员工、子公司、母公司、代理、合作伙伴、第三方内容提供商、供应商，以及我们或其各自的董事、高级职员和成员，均不对您或任何其他人承担任何损失或损害责任，包括直接、间接、特殊、后果性、附带性、惩罚性、信赖性或示范性损害、利润损失、人身伤害、死亡、财产损害、声誉损害或信息、数据损失，无论该等损害是否因使用或无法使用网站和/或服务而产生或与之相关。

对于任何未经授权访问或使用您的个人身份信息的行为，我们不承担任何形式的责任。如果您对网站和/或服务不满意，您唯一且排他的权利和救济方式是停止访问或使用网站和/或服务。

### 7.2 赔偿

您同意为我们以及我们的关联方、供应商及其各自的董事、高级职员、用户和代理进行抗辩、赔偿并使其免受损害，赔偿范围包括因以下事项引起或与之相关的所有索赔、诉讼和费用，包括律师费：您使用网站和/或服务；您未遵守或违反本协议；您使用第三方服务、产品、链接、广告和/或工具；您侵犯任何第三方权利；或任何其他人使用您的信息未经授权使用网站和/或服务。

## 第八部分：适用法律与仲裁

### 8.1 适用法律

本条款应受美国法律管辖并据其解释，不考虑其法律冲突规则。

### 8.2 仲裁

双方同意，任何涉及、关联或指向本条款和/或网站和/或服务的争议，均应根据美国实体法律通过具有约束力的仲裁方式专属解决。

## 第九部分：其他条款

### 9.1 客户服务

如果您对网站和/或服务有任何问题、意见或疑虑，可随时通过 ${supportEmail} 联系客户服务。我们会尽力在 72 个工作小时内回复所有客户服务咨询。

### 9.2 联盟披露

我们可能与第三方及关联方存在商业关系，并通过服务链接和推广其产品和/或服务。由于该关系，当用户从第三方联盟购买产品时，我们可能会获得佣金。

### 9.3 服务器位置；国际传输

我们在全球范围内运营，因此有必要在国际范围内传输您的信息。特别是，您的信息很可能会被传输至美国的服务器并在该地处理。使用网站和/或服务即表示您同意您的信息按照隐私政策所述方式被收集、使用和传输。

### 9.4 授权

双方各自向对方声明并保证，其拥有签署本协议的全部权力和授权，且本协议对该方具有约束力，并可根据其条款执行。

### 9.5 弃权

对本条款项下任何权利的放弃，仅在书面同意或书面声明时有效。延迟行使权利或未行使权利不应被视为放弃，也不妨碍任何一方日后行使该权利。

### 9.6 不可抗力

如果因不可抗力或任何超出我们控制范围的事件导致我们无法履行义务，我们不受该义务约束。此类事件包括自然灾害、战争、内乱、恐怖活动、紧急状态、制裁、禁运、国有化、罢工以及公共设施故障等。

### 9.7 转让

我们有权在按照本条款通知您后，将本条款以及我们在本条款项下的权利和义务转让给任何第三方。未经我们事先书面同意，您不得转让或转移您的权利，也不得委托他人履行您在本条款项下的任何义务。

### 9.8 第三方权利

除本条款明确规定外，本条款不赋予任何第三方任何权利。

### 9.9 可分割性

如果本协议的任何部分根据适用法律被认定为无效或不可执行，则该无效或不可执行条款将被视为由最接近原条款意图的有效且可执行条款替代，本协议其余部分继续有效。

### 9.10 更新与生效日期

本条款的生效日期为 **2026 年 7 月 10 日**。我们可能会不时更新本条款与条件，因此建议您经常查看。`,
  },
  {
    id: 'privacy-policy',
    title: '隐私政策',
    content_md: `# 隐私政策

**运营主体：** ${operatorName}

**注册地址：** ${operatorAddress}

**支持邮箱：** ${supportEmail}

**管理员微信：** ${administratorWechat}

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

**支持邮箱：** ${supportEmail}

**管理员微信：** ${administratorWechat}

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
const appStore = useAppStore()
const settings = computed(() => appStore.cachedPublicSettings)
const loading = ref(!settings.value)
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
  loadError.value = false
  const loadedSettings = await appStore.fetchPublicSettings()
  if (!loadedSettings) {
    loadError.value = true
  }
  loading.value = false
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
