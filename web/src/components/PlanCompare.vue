<script setup>
import { RouterLink } from "vue-router";

defineProps({
  cta: { type: Boolean, default: true },
});

const rows = [
  { name: "Уроки Intern", free: "да", pro: "да" },
  { name: "Уроки Junior, Middle, Senior", free: "нет", pro: "да" },
  { name: "Учебный банк — чтение", free: "да", pro: "да" },
  { name: "Запись в банк POST /api/v1", free: "нет", pro: "холд, P2P, IBAN" },
  { name: "SQL", free: "витрины intern_*", pro: "полная лаба и эталон" },
  { name: "Лабы", free: "intern", pro: "все 78" },
  { name: "Собес", free: "intern", pro: "middle+ · 57 вопросов" },
  { name: "Live /live", free: "да", pro: "да" },
  { name: "Материалы с инфографикой", free: "intern-статьи", pro: "80+" },
  { name: "Срок", free: "без срока", pro: "пока оплачен" },
];
</script>

<template>
  <section id="compare" class="scroll-mt-24">
    <p id="pay" class="kicker scroll-mt-24">зачем платить</p>
        <h2 class="font-display mt-3 text-3xl md:text-5xl">Intern смотрит. PRO пишет.</h2>
        <span class="section-rule" aria-hidden="true" />
        <p class="mt-3 max-w-2xl text-[16px] leading-7 text-mute">
          Бесплатно не урезают. Платите, когда сами вызываете POST в учебный банк, открываете junior+ и полный SQL.
          Мост: 3 дня правого столбца после регистрации.
        </p>

    <div class="mt-8 grid gap-6 lg:grid-cols-2">
      <article class="card p-6 md:p-8">
        <p class="stamp">бесплатно</p>
        <h3 class="font-display mt-2 text-3xl">Intern</h3>
        <ul class="mt-4 space-y-2 text-[15px] leading-7">
          <li>Уроки Intern, intern-лабы, intern-собес.</li>
          <li>Чтение учебного банка: кошельки, холды, журнал.</li>
          <li>SQL только по витринам intern_*.</li>
          <li>Срок не кончается. Junior+ закрыт страницей PRO.</li>
        </ul>
      </article>
      <article class="card p-6 md:p-8">
        <p class="stamp stamp-solid">PRO</p>
        <h3 class="font-display mt-2 text-3xl">Пишете в банк</h3>
        <ul class="mt-4 space-y-2 text-[15px] leading-7">
          <li>Junior, Middle, Senior и все 78 лаб.</li>
          <li>POST /api/v1: холд, P2P, входящий IBAN.</li>
          <li>Полный SQL к тем же таблицам.</li>
          <li>Собес middle+ (57 вопросов). Live — отдельно, и так открыт.</li>
        </ul>
      </article>
    </div>

    <p class="mt-6 max-w-2xl text-[15px] leading-7 text-mute">
      Триал: 72 часа правого столбца. Потом снова Intern, прогресс на месте.
      Сертификата нет: это зал практики, не диплом.
    </p>

    <div class="card mt-6 overflow-x-auto">
      <table class="pay-table">
        <thead>
          <tr>
            <th>Что есть</th>
            <th>Бесплатно</th>
            <th class="pay-pro">PRO</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in rows" :key="r.name">
            <td>{{ r.name }}</td>
            <td class="text-mute">{{ r.free }}</td>
            <td class="pay-pro">{{ r.pro }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="cta" class="mt-6 flex flex-wrap gap-3">
      <RouterLink to="/register" class="btn btn-accent">Начать бесплатно</RouterLink>
      <RouterLink to="/register?trial=1" class="btn btn-ghost">3 дня PRO</RouterLink>
    </div>
  </section>
</template>
