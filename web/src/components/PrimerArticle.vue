<script setup>
import { computed } from "vue";
import Infographic from "./Infographic.vue";

const props = defineProps({
  primer: { type: Object, required: true },
  /** all | head (term+explanation+examples) | rest (purpose→case) */
  part: { type: String, default: "all" },
});

function hasText(block) {
  if (!block) return false;
  if (block.paragraphs?.length) return true;
  if (block.intro || block.heading || block.body) return true;
  if (block.steps?.length) return true;
  if (block.items?.length || block.cases?.length) return true;
  if (block.infographic) return true;
  if (block.name) return true;
  return false;
}

const beats = computed(() => {
  const p = props.primer || {};
  const list = [];
  if (hasTerm.value) list.push({ id: "primer-term", n: "01", label: "термин" });
  if (hasAbout.value) list.push({ id: "primer-about", n: "02", label: "объяснение" });
  if (hasText(p.example)) list.push({ id: "primer-example", n: "03", label: "примеры" });
  if (hasText(p.purpose)) list.push({ id: "primer-purpose", n: "04", label: "для чего" });
  if (caseItems.value.length) list.push({ id: "primer-case", n: "05", label: "кейс" });
  if (hasText(p.practice)) list.push({ id: "primer-practice", n: "06", label: "закрепление" });
  return list;
});

const term = computed(() => props.primer?.term || {});
const termParas = computed(() => {
  const p = props.primer || {};
  if (p.term?.paragraphs?.length) return p.term.paragraphs;
  const w = p.what?.paragraphs || [];
  return w.slice(0, 1);
});
const aboutParas = computed(() => {
  const p = props.primer || {};
  if (p.about?.paragraphs?.length) return p.about.paragraphs;
  if (p.term?.paragraphs?.length) return [];
  const w = p.what?.paragraphs || [];
  return w.slice(1);
});
const aboutFig = computed(() => props.primer?.about?.infographic || props.primer?.what?.infographic || null);
const aboutSteps = computed(() => props.primer?.about?.steps || props.primer?.how?.steps || []);
const aboutIntro = computed(() => props.primer?.about?.intro || props.primer?.how?.intro || "");

const hasTerm = computed(() => term.value.name || termParas.value.length);
const hasAbout = computed(
  () => aboutParas.value.length || aboutFig.value || aboutSteps.value.length || aboutIntro.value,
);

const caseItems = computed(() => {
  const p = props.primer || {};
  if (p.cases?.items?.length) return p.cases.items;
  return p.mistakes?.cases || [];
});
const showHead = computed(() => props.part === "all" || props.part === "head");
const showRest = computed(() => props.part === "all" || props.part === "rest");
const showRail = computed(() => props.part === "all" || props.part === "head");
</script>

<template>
  <article v-if="primer" class="primer-article" :class="{ 'has-rail': showRail && beats.length }">
    <nav v-if="showRail && beats.length" class="article-rail" aria-label="Содержание статьи">
      <a v-for="b in beats" :key="b.id" :href="'#' + b.id">
        <span class="block">{{ b.n }}</span>
        {{ b.label }}
      </a>
    </nav>
    <div class="primer-col">
      <section v-if="showHead && hasTerm" id="primer-term" class="primer-block primer-term">
        <span class="stamp">термин</span>
        <div class="term-pair">
          <div>
            <p v-if="term.en" class="term-en">{{ term.en }}</p>
            <h2 class="term-word">{{ term.name || term.title || "Термин" }}</h2>
          </div>
          <div v-if="term.name2">
            <p v-if="term.en2" class="term-en">{{ term.en2 }}</p>
            <h2 class="term-word">{{ term.name2 }}</h2>
          </div>
        </div>
        <p v-for="(p, i) in termParas" :key="'t' + i">{{ p }}</p>
      </section>

      <section v-if="showHead && hasAbout" id="primer-about" class="primer-block">
        <span class="stamp">объяснение</span>
        <h2>{{ primer.about?.title || "Объяснение" }}</h2>
        <p v-for="(p, i) in aboutParas" :key="'a' + i">{{ p }}</p>
        <Infographic v-if="aboutFig" :fig="aboutFig" />
        <p v-if="aboutIntro">{{ aboutIntro }}</p>
        <ol v-if="aboutSteps.length">
          <li v-for="(s, i) in aboutSteps" :key="'h' + i">{{ s }}</li>
        </ol>
      </section>

      <section v-if="showHead && primer.example" id="primer-example" class="primer-block">
        <span class="stamp">примеры</span>
        <h2>{{ primer.example.title || "Примеры" }}</h2>
        <h3 v-if="primer.example.heading">{{ primer.example.heading }}</h3>
        <p
          v-for="(p, i) in primer.example.paragraphs || (primer.example.body ? [primer.example.body] : [])"
          :key="'e' + i"
        >
          {{ p }}
        </p>
        <Infographic v-if="primer.example.infographic" :fig="primer.example.infographic" />
        <p v-if="primer.example.note" class="primer-note">{{ primer.example.note }}</p>
        <template v-if="primer.mistakes?.items?.length">
          <h3>Так нельзя</h3>
          <div v-for="(m, i) in primer.mistakes.items" :key="'m' + i" class="primer-item">
            <h3>{{ m.title }}</h3>
            <p>{{ m.body }}</p>
          </div>
        </template>
      </section>

      <section v-if="showRest && primer.purpose" id="primer-purpose" class="primer-block">
        <span class="stamp">для чего</span>
        <h2>{{ primer.purpose.title }}</h2>
        <p v-for="(p, i) in primer.purpose.paragraphs" :key="'p' + i">{{ p }}</p>
        <Infographic v-if="primer.purpose.infographic" :fig="primer.purpose.infographic" />
      </section>

      <section v-if="showRest && caseItems.length" id="primer-case" class="primer-block">
        <span class="stamp stamp-case">кейс</span>
        <h2>{{ primer.cases?.title || "Кейс" }}</h2>
        <div v-for="(c, i) in caseItems" :key="'c' + i" class="primer-case">
          <h3>{{ c.title }}</h3>
          <p><strong>Что сломалось.</strong> {{ c.broke }}</p>
          <p><strong>Почему.</strong> {{ c.why }}</p>
          <p><strong>Что сделать.</strong> {{ c.should }}</p>
        </div>
      </section>

      <section v-if="showRest && primer.practice" id="primer-practice" class="primer-block primer-drill">
        <span class="stamp stamp-solid">закрепление</span>
        <h2>{{ primer.practice.title }}</h2>
        <p
          v-for="(p, i) in primer.practice.paragraphs || (primer.practice.body ? [primer.practice.body] : [])"
          :key="'d' + i"
        >
          {{ p }}
        </p>
        <slot name="drill" />
      </section>
    </div>
  </article>
</template>
