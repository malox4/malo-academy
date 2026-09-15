<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { RouterLink, useRouter, useRoute } from "vue-router";
import { useAuth } from "../stores/auth";
import BrandMark from "../components/BrandMark.vue";
import { inTelegram, tg } from "../lib/telegram";

const props = defineProps({ mode: String });
const route = useRoute();
const router = useRouter();
const auth = useAuth();
const isReg = computed(() => (props.mode || route.meta.mode) === "register");
const wantTrial = computed(() => isReg.value && String(route.query.trial || "") !== "");
const tgOn = computed(() => inTelegram() && Boolean(tg()?.initData));
const tgBusy = ref(false);

const form = reactive({ name: "", email: "", password: "" });

async function go() {
  const path = isReg.value ? "/api/auth/register" : "/api/auth/login";
  const body = isReg.value
    ? { name: form.name, email: form.email, password: form.password }
    : { email: form.email, password: form.password };
  const ok = await auth.submit(path, body);
  if (ok) router.replace(route.query.next || "/hall");
}

async function goTg() {
  tgBusy.value = true;
  const ok = await auth.submit("/api/auth/telegram", { initData: tg().initData });
  tgBusy.value = false;
  if (ok) router.replace(route.query.next || "/hall");
}

onMounted(() => {
  if (auth.user) router.replace(route.query.next || "/hall");
});
</script>

<template>
  <div class="flex min-h-dvh flex-col">
    <header class="mx-auto flex h-12 w-full max-w-md items-center px-4 sm:h-14 sm:px-6">
      <BrandMark to="/" compact />
    </header>
    <div class="mx-auto flex w-full max-w-md flex-1 flex-col justify-center px-4 pb-16 sm:px-6">
    <section class="card px-5 py-8 sm:px-6 sm:py-10">
      <p class="kicker">{{ wantTrial ? "триал PRO" : isReg ? "регистрация" : "вход" }}</p>
      <h1 class="font-display mt-3 text-3xl sm:text-4xl">{{ wantTrial ? "3 дня как PRO" : isReg ? "Создать аккаунт" : "Войти" }}</h1>
      <span class="section-rule" aria-hidden="true" />
      <p class="mt-3 text-[16px] leading-7 text-mute">
        {{
          wantTrial
            ? "После регистрации 72 часа: junior+, POST в учебный банк, полный SQL, middle-собес. Потом снова Intern."
            : isReg
              ? "Intern сразу: уроки, практика, пет на чтение, собес intern. Тем же аккаунтом — 3 дня PRO."
              : tgOn
                ? "В Mini App можно войти тем же Telegram. Email — если зал открыт в браузере."
                : "Email и пароль. После входа — зал, пет, практика, собес."
        }}
      </p>
      <button v-if="tgOn && !isReg" type="button" class="btn btn-accent mt-8 w-full" :disabled="tgBusy" @click="goTg">
        Войти через Telegram
      </button>
      <form class="mt-8 space-y-5" @submit.prevent="go">
        <label v-if="isReg" class="block text-sm font-medium">
          Имя
          <input v-model="form.name" class="mt-1 w-full border-b border-line bg-transparent py-2 outline-none" required minlength="2" />
        </label>
        <label class="block text-sm font-medium">
          Email
          <input v-model="form.email" type="email" class="mt-1 w-full border-b border-line bg-transparent py-2 outline-none" required />
        </label>
        <label class="block text-sm font-medium">
          Пароль
          <input v-model="form.password" type="password" class="mt-1 w-full border-b border-line bg-transparent py-2 outline-none" required minlength="8" />
        </label>
        <p v-if="auth.error" class="text-sm text-bad">{{ auth.error }}</p>
        <button type="submit" class="btn w-full sm:w-auto">
          {{ wantTrial ? "Открыть 3 дня PRO" : isReg ? "Начать бесплатно" : "Войти" }}
        </button>
      </form>
      <p class="mt-8 text-sm text-mute">
        <RouterLink v-if="!isReg" to="/register" class="font-medium text-accent">Нет аккаунта — регистрация</RouterLink>
        <RouterLink v-else to="/login" class="font-medium text-accent">Уже есть — вход</RouterLink>
      </p>
    </section>
    </div>
  </div>
</template>
