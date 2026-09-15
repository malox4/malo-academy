<script setup>
import { computed, ref, watch } from "vue";

const props = defineProps({
  board: { type: Object, required: true },
  fx: { type: Object, default: () => ({}) },
});
const emit = defineEmits(["select"]);
const selected = ref(props.board.nodes?.[0]?.id || props.board.debit?.[0]?.id || props.board.lines?.[0]?.id || null);

watch(
  () => props.board,
  (b) => {
    selected.value = b.nodes?.[0]?.id || b.debit?.[0]?.id || b.lines?.[0]?.id || null;
  },
);

watch(
  () => props.fx?.pulse,
  (id) => {
    if (id) selected.value = id;
  },
);

const nodes = computed(() => props.board.nodes || []);
const debit = computed(() => [...(props.board.debit || []), ...(props.fx?.extraDebit || [])]);
const credit = computed(() => props.board.credit || []);

const active = computed(() => {
  const b = props.board;
  if (b.kind === "flow") return nodes.value.find((n) => n.id === selected.value);
  if (b.kind === "taccount") return [...debit.value, ...credit.value].find((n) => n.id === selected.value);
  if (b.kind === "ticket") return (b.lines || []).find((n) => n.id === selected.value);
  return null;
});

function pick(id) {
  selected.value = id;
  emit("select", id);
}

function light(id) {
  return props.fx?.lights?.[id];
}

function isBroken(id) {
  return (props.fx?.broken || []).includes(id);
}

function plateClass(n) {
  const L = light(n.id);
  const on = selected.value === n.id;
  if (isBroken(n.id) || L === "bad") return "ig-tone-bad";
  if (L === "ok") return "ig-tone-ok";
  if (L === "hold" || n.mood === "hold") return "ig-tone-accent";
  if (on) return "ig-tone-accent";
  return "";
}

function ticketHeat(line) {
  const t = props.fx?.tickets?.[line.id];
  if (t === "bad" || line.bad) return "bg-bad/15 text-ink";
  if (t === "ok") return "bg-ok/12";
  if (selected.value === line.id) return "bg-accent/10";
  return "hover:bg-paper";
}
</script>

<template>
  <section class="card overflow-hidden" :class="fx.flash === 'ok' ? 'fx-ok' : fx.flash === 'bad' ? 'fx-bad' : ''">
    <div class="flex items-start justify-between gap-4 border-b-2 border-ink px-5 py-4 md:px-6">
      <div>
        <span class="stamp">доска</span>
        <h2 class="font-display mt-2 text-2xl leading-none md:text-3xl">{{ board.title }}</h2>
        <p class="mt-2 max-w-2xl text-[15px] leading-6 text-mute">{{ board.caption }}</p>
      </div>
    </div>

    <div v-if="board.kind === 'flow'" class="bg-paper px-4 py-5 md:px-6">
      <div class="ig-flow">
        <template v-for="(n, i) in nodes" :key="n.id">
          <div v-if="i" class="ig-join" aria-hidden="true" />
          <button
            type="button"
            class="ig-plate text-left"
            :class="[plateClass(n), isBroken(n.id) ? 'node-shake' : '']"
            @click="pick(n.id)"
          >
            <span class="plate-n">0{{ i + 1 }}</span>
            <h4>{{ n.label }}</h4>
            <p v-if="n.sub" class="ig-sub">{{ n.sub }}</p>
          </button>
        </template>
      </div>
      <p v-if="board.edges?.length" class="ig-notes mt-3">
        <span v-for="(e, i) in board.edges" :key="'e' + i" class="stamp stamp-ink">{{ e.note }}</span>
      </p>
      <p v-if="active" class="mt-4 border-t border-line pt-3 text-[15px] leading-6">
        <span class="stamp">{{ active.label }}</span>
        <span class="mt-2 block font-medium">{{ active.sub }}</span>
        <span class="mt-1 block text-mute">{{ active.detail }}</span>
      </p>
    </div>

    <div v-else-if="board.kind === 'taccount'" class="grid md:grid-cols-2">
      <div class="border-b border-line p-5 md:border-r md:border-b-0">
        <span class="stamp stamp-case">дебет</span>
        <button
          v-for="row in debit"
          :key="row.id"
          type="button"
          class="mt-3 block w-full border px-4 py-3 text-left font-medium transition"
          :class="[
            row.ok === false || isBroken(row.id) || light(row.id) === 'bad' ? 'border-bad bg-bad/8' : '',
            selected === row.id ? 'border-accent' : 'border-line',
          ]"
          style="border-radius: 2px"
          @click="pick(row.id)"
        >
          {{ row.text }}
        </button>
      </div>
      <div class="p-5">
        <span class="stamp stamp-ok">кредит</span>
        <button
          v-for="row in credit"
          :key="row.id"
          type="button"
          class="mt-3 block w-full border px-4 py-3 text-left font-medium transition"
          :class="selected === row.id ? 'border-accent' : 'border-line'"
          style="border-radius: 2px"
          @click="pick(row.id)"
        >
          {{ row.text }}
        </button>
      </div>
      <p class="border-t border-line px-5 py-4 text-[15px] leading-6 md:col-span-2">{{ board.hole }}</p>
    </div>

    <div v-else-if="board.kind === 'ticket'" class="divide-y divide-line">
      <button
        v-for="line in board.lines"
        :key="line.id"
        type="button"
        class="flex w-full gap-3 px-5 py-3.5 text-left text-[15px] leading-6 transition"
        :class="ticketHeat(line)"
        @click="pick(line.id)"
      >
        <span class="w-6 font-mono text-xs text-mute">{{ line.id }}</span>
        <span>{{ line.text }}</span>
      </button>
      <p v-if="active?.why" class="px-5 py-4 text-[15px] leading-6">{{ active.why }}</p>
    </div>

    <div v-else-if="board.kind === 'compare'" class="overflow-x-auto">
      <table class="w-full text-left text-[14px]">
        <thead class="font-mono text-[10px] uppercase tracking-[0.16em] text-mute">
          <tr>
            <th class="px-5 py-3"></th>
            <th class="px-5 py-3 text-accent">{{ board.leftTitle }}</th>
            <th class="px-5 py-3">{{ board.rightTitle }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in board.rows" :key="row.label" class="border-t border-line">
            <td class="px-5 py-3 font-medium">{{ row.label }}</td>
            <td class="px-5 py-3">{{ row.left }}</td>
            <td class="px-5 py-3 text-mute">{{ row.right }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div
      v-if="fx.caption"
      class="border-t px-5 py-3 text-[14px] leading-6"
      :class="fx.flash === 'bad' ? 'border-bad/20 bg-bad/8 text-ink' : 'border-ok/20 bg-ok/8'"
    >
      <span class="font-mono text-[10px] uppercase tracking-[0.18em]" :class="fx.flash === 'bad' ? 'text-bad' : 'text-ok'">
        {{ fx.flash === "bad" ? "Ошибка хода" : "Верно" }}
      </span>
      <p class="mt-1">{{ fx.caption }}</p>
    </div>
  </section>
</template>
