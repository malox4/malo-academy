<script setup>
import { computed } from "vue";
import { RouterLink } from "vue-router";
import { useAuth } from "../stores/auth";

defineProps({
  title: { type: String, default: "Этот раздел в PRO" },
});

const auth = useAuth();
const trialEnded = computed(() => auth.trialEnded);
</script>

<template>
  <div class="card mx-auto max-w-lg p-8">
    <p class="stamp">{{ trialEnded ? "пробный период закончился" : "PRO" }}</p>
    <h2 class="font-display mt-2 text-3xl">{{ title }}</h2>
    <p v-if="trialEnded" class="mt-3 text-[16px] leading-7 text-mute">
      Три дня полного доступа прошли. Intern по-прежнему открыт: уроки, практика intern, пет на чтение, собес intern.
      Junior+, запись в пет и senior-собес — в платном PRO.
    </p>
    <p v-else class="mt-3 text-[16px] leading-7 text-mute">
      Intern открыт бесплатно. Этот грейд входит в PRO. После регистрации даём 3 дня, чтобы пройти его целиком.
    </p>
    <div class="mt-6 flex flex-wrap gap-4 text-sm font-medium">
      <RouterLink to="/hall" class="text-accent">← В зал</RouterLink>
      <RouterLink to="/pricing" class="text-ink">Сравнение тарифов</RouterLink>
    </div>
  </div>
</template>
