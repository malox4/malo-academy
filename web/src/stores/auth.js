import { defineStore } from "pinia";
import { ref, computed } from "vue";

async function api(path, opts = {}) {
  const res = await fetch(path, {
    credentials: "include",
    headers: { "Content-Type": "application/json", ...(opts.headers || {}) },
    ...opts,
  });
  const data = await res.json().catch(() => ({}));
  return { ok: res.ok, status: res.status, data };
}

function mergeUser(data) {
  if (!data?.user) return null;
  return {
    ...data.user,
    plan: data.plan ?? data.user.plan,
    trialUntil: data.trialUntil ?? data.user.trialUntil,
    trialActive: data.trialActive ?? data.user.trialActive,
    trialRemaining: data.trialRemaining ?? data.user.trialRemaining,
    isPro: data.isPro ?? data.user.isPro,
  };
}

export const useAuth = defineStore("auth", () => {
  const user = ref(null);
  const entitlements = ref({ plan: "guest", grades: [] });
  const ready = ref(false);
  const error = ref("");

  const isPro = computed(() => Boolean(user.value?.isPro));
  const trialActive = computed(() => Boolean(user.value?.trialActive));
  const trialEnded = computed(() => Boolean(user.value?.trialUntil) && !isPro.value);
  const trialLabel = computed(() => {
    const sec = Number(user.value?.trialRemaining || 0);
    if (sec <= 0) return "";
    const days = Math.floor(sec / 86400);
    if (days >= 1) return `${days}д`;
    const hours = Math.max(1, Math.ceil(sec / 3600));
    return `${hours}ч`;
  });

  function canGrade(id) {
    if (isPro.value) return true;
    const g = entitlements.value.grades;
    return Array.isArray(g) && g.includes(id);
  }

  function canLab(gradeId) {
    if (isPro.value) return true;
    return gradeId === "intern";
  }

  function apply(data) {
    user.value = mergeUser(data);
    entitlements.value = data.entitlements || { plan: "guest", grades: [] };
  }

  async function hydrate() {
    const { data } = await api("/api/auth/me");
    apply(data);
    ready.value = true;
  }

  async function submit(path, body) {
    error.value = "";
    const { ok, data } = await api(path, { method: "POST", body: JSON.stringify(body) });
    if (!ok) {
      error.value = data.error?.message || "Не вышло.";
      return false;
    }
    apply(data);
    return true;
  }

  async function logout() {
    await api("/api/auth/logout", { method: "POST", body: "{}" });
    user.value = null;
    entitlements.value = { plan: "guest", grades: [] };
  }

  return {
    user,
    entitlements,
    ready,
    error,
    isPro,
    trialActive,
    trialEnded,
    trialLabel,
    canGrade,
    canLab,
    hydrate,
    submit,
    logout,
  };
});
