<script setup>
import { computed } from "vue";
import { RouterLink, useRoute } from "vue-router";
import SiteNav from "../components/SiteNav.vue";
import PrimerArticle from "../components/PrimerArticle.vue";
import { useAuth } from "../stores/auth";
import { topicById, TOPICS } from "../content/topics";
import { TRACKS } from "../content/materials";

const auth = useAuth();
const route = useRoute();

const topic = computed(() => topicById(String(route.params.id || "")));
const track = computed(() => TRACKS.find((t) => t.id === topic.value?.trackId));
const siblings = computed(() => TOPICS.filter((t) => t.trackId === topic.value?.trackId));
</script>

<template>
  <div>
    <SiteNav v-if="!auth.user" />
    <div :class="auth.user ? '' : 'mx-auto max-w-5xl px-4 pb-24 pt-8 md:px-6'">
      <div v-if="!topic" class="space-y-4">
        <p class="text-mute">Такой статьи нет.</p>
        <RouterLink to="/materials" class="text-accent">← Все материалы</RouterLink>
      </div>
      <article v-else>
        <RouterLink :to="'/materials#' + topic.trackId" class="text-sm text-mute hover:text-accent">
          ← {{ track?.title || "Материалы" }}
        </RouterLink>
        <div class="mt-6 flex flex-wrap items-center gap-2">
          <span class="stamp stamp-ink">{{ topic.grade }}</span>
          <span class="stamp">{{ topic.minutes }} мин</span>
          <span class="stamp">статья</span>
          <span v-if="topic.grade === 'intern'" class="stamp stamp-ok">открыто</span>
          <span v-else class="stamp stamp-ink">урок · PRO</span>
        </div>
        <h1 class="font-display mt-3 text-4xl leading-[0.95] md:text-6xl">{{ topic.title }}</h1>
        <span class="section-rule" aria-hidden="true" />
        <p class="mt-4 text-[18px] leading-8 text-mute">{{ topic.teaser }}</p>

        <PrimerArticle class="mt-10" :primer="topic.sections">
          <template #drill>
            <RouterLink
              :to="topic.sections.practice.to"
              class="btn btn-accent mt-4"
            >
              {{ topic.sections.practice.cta }} →
            </RouterLink>
            <p class="mt-3 text-sm text-mute">Дальше на уроке и в лабе не запирается. Ошибка хода показывает, почему — не зелёную галочку.</p>
          </template>
        </PrimerArticle>

        <nav v-if="siblings.length > 1" class="mt-14 border-t-2 border-ink pt-6">
          <p class="kicker">ещё в треке</p>
          <div class="mt-4 flex flex-col gap-2">
            <RouterLink
              v-for="s in siblings.filter((x) => x.id !== topic.id)"
              :key="s.id"
              :to="'/materials/' + s.id"
              class="text-[15px] text-accent hover:underline"
            >
              {{ s.title }}
            </RouterLink>
          </div>
        </nav>
      </article>
    </div>
  </div>
</template>
