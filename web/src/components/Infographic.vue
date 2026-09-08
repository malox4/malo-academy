<script setup>
import { computed } from "vue";

const props = defineProps({
  fig: { type: Object, default: null },
  kicker: { type: String, default: "схема" },
  plate: { type: String, default: "" },
  title: { type: String, default: "" },
  caption: { type: String, default: "" },
});

const data = computed(() => props.fig || {});
const kind = computed(() => data.value.kind || "");
const headTitle = computed(() => data.value.title || props.title);
const headCaption = computed(() => data.value.caption || props.caption);
const kickerLabel = computed(() => {
  if (!props.fig) return props.kicker;
  const map = {
    flow: "схема",
    compare: "сравнение",
    http: "контракт",
    er: "сущности",
    layers: "слои",
    break: "сходится / рвётся",
    seq: "последовательность",
    sequence: "последовательность",
    code: "бланк",
    c4: "контур",
    hold: "три числа",
    stack: "слои",
    codes: "коды",
  };
  return map[kind.value] || "схема";
});

function money(n) {
  const x = Number(n);
  if (!Number.isFinite(x)) return n ?? "—";
  return x.toLocaleString("ru-RU");
}

function toneClass(tone) {
  if (tone === "ok") return "ig-tone-ok";
  if (tone === "bad") return "ig-tone-bad";
  if (tone === "accent") return "ig-tone-accent";
  return "";
}

const seq = computed(() => {
  const actors = data.value.actors || [];
  const steps = (data.value.steps || data.value.messages || []).map((st) => ({
    ...st,
    label: st.label || st.text || "",
  }));
  const col = 120;
  const top = 28;
  const rowH = 44;
  const width = Math.max(actors.length * col + 24, 280);
  const height = top + 28 + steps.length * rowH + 16;
  const x = (i) => 28 + i * col + col / 2;
  const actorIndex = (name) => {
    const i = actors.indexOf(name);
    return i >= 0 ? i : 0;
  };
  return { actors, steps, col, top, rowH, width, height, x, actorIndex };
});

function seqArrow(st, i) {
  const s = seq.value;
  const y = s.top + 18 + i * s.rowH;
  const x2 = s.x(s.actorIndex(st.to));
  const x1 = s.x(s.actorIndex(st.from));
  const dir = x2 >= x1 ? 1 : -1;
  return `${x2},${y} ${x2 - 8 * dir},${y - 4} ${x2 - 8 * dir},${y + 4}`;
}
</script>

<template>
  <figure class="infographic">
    <div class="infographic-head">
      <span class="stamp">{{ kickerLabel }}</span>
      <span v-if="plate" class="plate-n">{{ plate }}</span>
      <span v-if="headTitle" class="infographic-title">{{ headTitle }}</span>
    </div>
    <div class="infographic-frame">
      <slot />

      <div v-if="kind === 'flow'" class="ig-flow">
        <template v-for="(n, i) in data.nodes || []" :key="n.id || i">
          <div v-if="i" class="ig-join" aria-hidden="true" />
          <article class="ig-plate" :class="toneClass(n.tone)">
            <span class="plate-n">0{{ i + 1 }}</span>
            <h4>{{ n.label }}</h4>
            <p v-if="n.sub" class="ig-sub">{{ n.sub }}</p>
            <p v-if="n.detail">{{ n.detail }}</p>
          </article>
        </template>
        <p v-if="data.notes?.length" class="ig-notes">
          <span v-for="(note, i) in data.notes" :key="i" class="stamp stamp-ink">{{ note }}</span>
        </p>
      </div>

      <div v-else-if="kind === 'compare'" class="ig-compare">
        <article v-for="(col, i) in data.columns || []" :key="col.title || i" class="ig-plate" :class="toneClass(col.tone)">
          <span class="plate-n">0{{ i + 1 }}</span>
          <h4>{{ col.title }}</h4>
          <p v-if="col.lead" class="ig-sub">{{ col.lead }}</p>
          <ul>
            <li v-for="(row, ri) in col.rows || []" :key="ri">{{ row }}</li>
          </ul>
        </article>
      </div>

      <div v-else-if="kind === 'http'" class="ig-http">
        <article class="ig-plate">
          <span class="stamp stamp-ink">клиент</span>
          <h4>{{ data.client }}</h4>
          <p class="ig-sub">{{ data.clientSub }}</p>
        </article>
        <div class="ig-http-wire">
          <p class="font-mono">{{ data.req }}</p>
          <p v-if="data.reqNote" class="ig-sub">{{ data.reqNote }}</p>
          <div class="ig-join ig-join-wide" aria-hidden="true" />
          <p class="font-mono">{{ data.res }}</p>
          <p v-if="data.resNote" class="ig-sub">{{ data.resNote }}</p>
        </div>
        <article class="ig-plate ig-tone-accent">
          <span class="stamp">сервер</span>
          <h4>{{ data.server }}</h4>
          <p class="ig-sub">{{ data.serverSub }}</p>
        </article>
      </div>

      <div v-else-if="kind === 'er'" class="ig-er">
        <article v-for="(ent, i) in data.entities || []" :key="ent.name || i" class="ig-plate">
          <span class="plate-n">0{{ i + 1 }}</span>
          <h4>{{ ent.name }}</h4>
          <ul class="ig-fields">
            <li v-for="f in ent.fields || []" :key="f.name" :class="{ 'ig-key': f.key }">
              {{ f.name }}
            </li>
          </ul>
        </article>
        <p v-if="data.rel" class="ig-rel font-mono">{{ data.rel }}</p>
      </div>

      <div v-else-if="kind === 'layers'" class="ig-layers">
        <template v-for="(n, i) in data.layers || []" :key="n.id || i">
          <div v-if="i" class="ig-join ig-join-down" aria-hidden="true" />
          <article class="ig-plate" :class="toneClass(n.tone)">
            <span class="plate-n">0{{ i + 1 }}</span>
            <h4>{{ n.label }}</h4>
            <p v-if="n.sub" class="ig-sub">{{ n.sub }}</p>
            <p v-if="n.detail">{{ n.detail }}</p>
          </article>
        </template>
        <p v-if="data.notes?.length" class="ig-notes">
          <span v-for="(note, i) in data.notes" :key="i" class="stamp stamp-ink">{{ note }}</span>
        </p>
      </div>

      <div v-else-if="kind === 'break'" class="ig-break">
        <div class="ig-break-head">
          <span class="stamp stamp-ok">{{ data.okHead || "живёт" }}</span>
          <span class="stamp stamp-case">{{ data.badHead || "рвётся" }}</span>
        </div>
        <article v-for="(row, i) in data.rows || []" :key="row.label || i" class="ig-break-row">
          <span class="plate-n">0{{ i + 1 }}</span>
          <h4>{{ row.label }}</h4>
          <p class="ig-break-ok">{{ row.ok }}</p>
          <p class="ig-break-bad">{{ row.bad }}</p>
        </article>
      </div>

      <div v-else-if="kind === 'stack'" class="ig-stack">
        <article v-for="(p, i) in data.parts || []" :key="p.title || i" class="ig-plate" :class="toneClass(p.tone)">
          <span class="plate-n">0{{ i + 1 }}</span>
          <h4>{{ p.title }}</h4>
          <p v-if="p.body">{{ p.body }}</p>
        </article>
      </div>

      <div v-else-if="kind === 'codes'" class="ig-codes">
        <article v-for="(g, i) in data.groups || []" :key="g.code || i" class="ig-plate" :class="toneClass(g.tone)">
          <span class="plate-n">0{{ i + 1 }}</span>
          <h4>{{ g.code }}</h4>
          <p class="ig-sub">{{ g.mean }}</p>
        </article>
      </div>

      <div v-else-if="kind === 'seq' || kind === 'sequence'" class="ig-seq">
        <svg
          class="ig-seq-svg"
          :viewBox="'0 0 ' + seq.width + ' ' + seq.height"
          role="img"
          :aria-label="headTitle || 'последовательность'"
        >
          <line
            v-for="(a, i) in seq.actors"
            :key="'lifeline-' + i"
            :x1="seq.x(i)"
            :y1="seq.top"
            :x2="seq.x(i)"
            :y2="seq.height - 8"
            class="ig-seq-life"
          />
          <g v-for="(a, i) in seq.actors" :key="'act-' + i">
            <rect
              :x="seq.x(i) - 46"
              y="4"
              width="92"
              height="22"
              class="ig-seq-box"
            />
            <text :x="seq.x(i)" y="19" class="ig-seq-actor">{{ a }}</text>
          </g>
          <g v-for="(st, i) in seq.steps" :key="'st-' + i">
            <line
              :x1="seq.x(seq.actorIndex(st.from))"
              :y1="seq.top + 18 + i * seq.rowH"
              :x2="seq.x(seq.actorIndex(st.to))"
              :y2="seq.top + 18 + i * seq.rowH"
              class="ig-seq-msg"
              :class="st.tone === 'bad' ? 'ig-seq-msg-bad' : ''"
            />
            <polygon
              :points="seqArrow(st, i)"
              class="ig-seq-head"
              :class="st.tone === 'bad' ? 'ig-seq-msg-bad' : ''"
            />
            <text
              :x="(seq.x(seq.actorIndex(st.from)) + seq.x(seq.actorIndex(st.to))) / 2"
              :y="seq.top + 10 + i * seq.rowH"
              class="ig-seq-label"
            >{{ st.label }}</text>
          </g>
        </svg>
      </div>

      <div v-else-if="kind === 'code'" class="ig-code">
        <p v-if="data.lead" class="ig-sub">{{ data.lead }}</p>
        <pre class="font-mono">{{ data.body }}</pre>
      </div>

      <div v-else-if="kind === 'c4'" class="ig-c4">
        <article class="ig-plate ig-tone-accent ig-c4-system">
          <span class="stamp stamp-ink">{{ data.contextLabel || "контур" }}</span>
          <h4>{{ data.system }}</h4>
          <p v-if="data.actors" class="ig-sub">{{ data.actors }}</p>
        </article>
        <div class="ig-flow">
          <article v-for="(box, i) in data.boxes || []" :key="box.label || i" class="ig-plate" :class="toneClass(box.tone)">
            <span class="plate-n">0{{ i + 1 }}</span>
            <h4>{{ box.label }}</h4>
            <p v-if="box.detail">{{ box.detail }}</p>
          </article>
        </div>
      </div>

      <div v-else-if="kind === 'hold'" class="ig-hold">
        <article class="ig-plate">
          <span class="stamp">книга</span>
          <h4>ledger</h4>
          <p class="ig-sub">{{ data.ledgerLabel || "остаток в журнале" }}</p>
          <p class="ig-hold-n">{{ money(data.ledger) }}</p>
        </article>
        <span class="ig-hold-op" aria-hidden="true">−</span>
        <article class="ig-plate ig-tone-accent">
          <span class="stamp">холд</span>
          <h4>hold OPEN</h4>
          <p class="ig-sub">{{ data.holdLabel || "резерв, проводок нет" }}</p>
          <p class="ig-hold-n">{{ money(data.holdAmt) }}</p>
        </article>
        <span class="ig-hold-op" aria-hidden="true">=</span>
        <article class="ig-plate ig-tone-ok">
          <span class="stamp stamp-ok">витрина</span>
          <h4>available</h4>
          <p class="ig-sub">{{ data.availNote || "то, что клиент может потратить" }}</p>
          <p class="ig-hold-n">{{ money(Number(data.ledger || 0) - Number(data.holdAmt || 0)) }}</p>
        </article>
      </div>
    </div>
    <figcaption v-if="headCaption || $slots.caption" class="infographic-caption">
      <slot name="caption">{{ headCaption }}</slot>
    </figcaption>
  </figure>
</template>
