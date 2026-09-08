<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink, useRoute } from "vue-router";
import SiteNav from "../components/SiteNav.vue";
import LandingStage from "../components/LandingStage.vue";
import PlanCompare from "../components/PlanCompare.vue";
import { LANDING_PATH } from "../content/path";
import {
  LANDING_CREDITS,
  LANDING_ROOMS,
  LANDING_STRIP,
  LANDING_PATH_SHOTS,
  LANDING_PATH_ALTS,
} from "../content/landingMedia";

const route = useRoute();
const room = ref("bank");
const rooms = LANDING_ROOMS;
const strip = LANDING_STRIP;
const credits = LANDING_CREDITS;

onMounted(() => {
  const id = route.hash === "#compare" ? "compare" : route.hash === "#pay" ? "pay" : "";
  if (id) document.getElementById(id)?.scrollIntoView({ behavior: "smooth" });
});

const picked = computed(() => rooms.find((x) => x.id === room.value) || rooms[0]);

const path = LANDING_PATH.map((s, i) => ({
  ...s,
  img: LANDING_PATH_SHOTS[i],
  alt: LANDING_PATH_ALTS[i],
}));
</script>

<template>
  <div class="landing min-h-screen">
    <SiteNav />

    <main class="mx-auto max-w-6xl px-4 pb-16 md:px-6">
      <figure class="shot landing-hero mt-4">
        <img
          src="/media/landing/hero-desk.jpg"
          alt="Стол аналитика: ноутбуки и схема процесса на бумаге"
          width="1600"
          height="1068"
        />
        <figcaption>
          <span class="stamp stamp-solid">зал</span>
          Стол аналитика: схема на бумаге, экран рядом. Не Zoom и не вебинар.
        </figcaption>
      </figure>

      <section class="mt-10 grid items-center gap-10 lg:grid-cols-[0.92fr_1.08fr]">
        <div>
          <p class="kicker">зал аналитика</p>
          <h1 class="font-display mt-4 text-4xl leading-[0.92] md:text-6xl">
            Intern читает банк. PRO жмёт холд.
          </h1>
          <span class="section-rule" aria-hidden="true" />
          <p class="mt-5 max-w-xl text-[18px] leading-7">
            Analyst Hall — самообучение BA и SA. Справа — тот же учебный банк, что после входа: метод, полный путь, JSON, ответ.
          </p>
          <ul class="mt-5 space-y-2 text-[15px] leading-6">
            <li><span class="stamp">intern</span> уроки входа, GET кошельки и журнал, intern-лабы и intern-собес. Бесплатно, без срока.</li>
            <li><span class="stamp stamp-solid">pro</span> POST холд / P2P / IBAN, полный SQL, junior+, все 78 лаб, middle-собес.</li>
            <li><span class="stamp">триал</span> 72 часа как PRO с регистрации. Потом снова Intern, прогресс на месте.</li>
          </ul>
          <div class="mt-7 flex flex-wrap gap-3">
            <RouterLink to="/register" class="btn btn-accent">Начать бесплатно</RouterLink>
            <RouterLink to="/register?trial=1" class="btn btn-ghost">3 дня PRO</RouterLink>
            <a href="#compare" class="btn btn-ghost">Зачем платить</a>
          </div>
        </div>
        <LandingStage />
      </section>

      <section class="shot-strip mt-14">
        <figure v-for="s in strip" :key="s.img" class="shot">
          <img :src="s.img" :alt="s.alt" width="1400" height="933" />
          <figcaption>
            <span class="stamp">{{ s.stamp }}</span>
            {{ s.cap }}
          </figcaption>
        </figure>
      </section>

      <section class="mt-16">
        <p class="kicker">продукт</p>
        <h2 class="font-display mt-3 text-3xl md:text-5xl">Четыре комнаты зала</h2>
        <span class="section-rule" aria-hidden="true" />
        <p class="mt-3 max-w-2xl text-[16px] leading-7 text-mute">
          Это не каталог вебинаров. Платите, когда нужно писать в ядро и открыть этажи выше Intern.
        </p>
        <div class="room-board mt-8">
          <button
            v-for="r in rooms"
            :key="r.id"
            type="button"
            class="shot room-tile text-left"
            :class="{ 'room-on': room === r.id }"
            :aria-pressed="room === r.id"
            @click="room = r.id"
          >
            <img :src="r.img" :alt="r.alt" width="1400" height="933" />
            <span class="room-tile-cap">
              <span class="plate-n">{{ r.n }}</span>
              <span class="font-display mt-2 block text-2xl leading-none md:text-3xl">{{ r.title }}</span>
            </span>
          </button>
        </div>
        <article class="shot shot-picked mt-5" :key="room">
          <img :src="picked.img" :alt="picked.alt" width="1400" height="933" />
          <div class="shot-picked-copy">
            <p class="stamp">{{ picked.n }}</p>
            <h3 class="font-display mt-3 text-3xl md:text-4xl">{{ picked.title }}</h3>
            <p class="mt-4 text-[17px] leading-7">{{ picked.sell }}</p>
            <p class="mt-3 font-mono text-[13px] leading-6 text-accent">{{ picked.proof }}</p>
          </div>
        </article>
      </section>

      <section class="mt-16">
        <p class="kicker">с нуля</p>
        <h2 class="font-display mt-3 text-3xl md:text-5xl">Шесть шагов до решения платить</h2>
        <span class="section-rule" aria-hidden="true" />
        <div class="path-board mt-8">
          <article v-for="s in path" :key="s.n" class="shot path-tile">
            <img :src="s.img" :alt="s.alt" width="1400" height="933" />
            <div class="path-tile-cap">
              <div class="plate-n">{{ s.n }}</div>
              <h3 class="font-display mt-3 text-2xl leading-none md:text-3xl">{{ s.t }}</h3>
              <p class="mt-2 text-[15px] leading-6 text-mute">{{ s.d }}</p>
            </div>
          </article>
        </div>
      </section>

      <PlanCompare class="mt-16" />

      <section class="cta-plate cta-photo mt-20">
        <img src="/media/landing/office.jpg" alt="" width="1400" height="935" />
        <div class="cta-inner">
          <p class="stamp stamp-solid">вход</p>
          <h2 class="font-display mt-3 text-3xl md:text-5xl">Intern сразу. 3 дня проверить PRO</h2>
          <p class="mt-3 max-w-xl text-[16px] leading-7 text-white/80">
            Регистрация открывает чтение банка. Тем же аккаунтом 72 часа пишете холд и P2P. Потом решаете, платить ли.
          </p>
          <div class="mt-6 flex flex-wrap gap-3">
            <RouterLink to="/register" class="btn btn-on-ink">Начать бесплатно</RouterLink>
            <RouterLink to="/register?trial=1" class="btn btn-on-ink-ghost">3 дня PRO</RouterLink>
          </div>
        </div>
      </section>
    </main>

    <div class="landing-dock">
      <p class="landing-dock-copy">Intern бесплатно · новым 3 дня PRO</p>
      <div class="flex flex-wrap gap-2">
        <RouterLink to="/register" class="btn btn-on-ink btn-sm">Начать бесплатно</RouterLink>
        <RouterLink to="/register?trial=1" class="btn btn-on-ink-ghost btn-sm">3 дня PRO</RouterLink>
      </div>
    </div>

    <footer class="border-t-2 border-ink bg-white">
      <div class="mx-auto flex max-w-6xl flex-col gap-3 px-4 py-8 text-sm text-mute md:flex-row md:items-center md:justify-between md:px-6">
        <p class="font-display text-ink">Analyst Hall</p>
        <p>зал аналитика · самообучение BA / SA</p>
        <nav class="flex flex-wrap gap-4">
          <RouterLink to="/pricing">Intern и PRO</RouterLink>
          <RouterLink to="/pet">Учебный банк</RouterLink>
          <RouterLink to="/interview">Собес</RouterLink>
        </nav>
      </div>
      <p class="mx-auto max-w-6xl px-4 pb-8 text-[11px] leading-5 text-mute md:px-6">{{ credits }}</p>
    </footer>
  </div>
</template>
