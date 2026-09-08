<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import SiteNav from "../components/SiteNav.vue";
import { useAuth } from "../stores/auth";
import { TRACKS, trackById } from "../content/materials";

const auth = useAuth();
const route = useRoute();
const router = useRouter();

const hashId = computed(() => String(route.hash || "").replace("#", ""));
const active = ref(TRACKS[0].id);

function syncFromHash() {
  const id = hashId.value;
  if (TRACKS.some((t) => t.id === id)) active.value = id;
}

onMounted(syncFromHash);
watch(hashId, syncFromHash);

function select(id) {
  active.value = id;
  router.replace({ path: "/materials", hash: "#" + id });
}

const track = computed(() => trackById(active.value));
</script>

<template>
  <div>
    <SiteNav v-if="!auth.user" />
    <div :class="auth.user ? '' : 'mx-auto max-w-6xl px-4 pb-20 pt-8 md:px-6'">
      <p class="kicker">самообучение</p>
      <h1 class="font-display mt-3 text-4xl md:text-6xl">Материалы зала</h1>
      <span class="section-rule" aria-hidden="true" />
      <p class="mt-3 max-w-2xl text-[17px] leading-7 text-mute">
        Треки: требования, REST, SQL, интеграции, архитектура, процессы, качество, платежи, софты, собес.
        Карточка открывает статью из шести частей: что это, для чего, как, пример, ошибки, закрепление. Рядом схема.
      </p>

      <div class="mt-8 flex flex-wrap gap-2">
        <button
          v-for="t in TRACKS"
          :key="t.id"
          type="button"
          class="stamp"
          :class="active === t.id ? 'stamp-solid' : 'stamp-ink'"
          @click="select(t.id)"
        >
          {{ t.title }}
        </button>
      </div>

      <section class="mt-10">
        <p class="kicker">{{ track.kicker }}</p>
        <h2 class="font-display mt-3 text-3xl">{{ track.title }}</h2>
        <span class="section-rule" aria-hidden="true" />
        <p class="mt-2 max-w-xl text-[15px] leading-7 text-mute">{{ track.blurb }}</p>
        <div class="catalog-grid mt-8 md:grid-cols-2 lg:grid-cols-3">
          <RouterLink
            v-for="c in track.cards"
            :key="c.id"
            :to="c.to"
            class="card card-hover p-5"
          >
            <div class="flex items-center justify-between">
              <span class="stamp stamp-ink">{{ c.grade }}</span>
              <span v-if="c.grade === 'intern'" class="stamp stamp-ok">открыто</span>
              <span v-else class="stamp stamp-ink">урок · PRO</span>
            </div>
            <h3 class="mt-3 font-display text-2xl leading-none">{{ c.title }}</h3>
            <p class="mt-2 text-sm leading-6 text-mute">{{ c.teaser }}</p>
            <p class="mt-4 font-mono text-[11px] text-accent">Читать статью →</p>
          </RouterLink>
        </div>
      </section>
    </div>
  </div>
</template>
