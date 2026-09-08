<script setup>
import { computed, ref, watch } from "vue";

const props = defineProps({ desk: { type: Object, required: true } });
const emit = defineEmits(["done", "play"]);

const place = ref({});
const hold = ref(null);
const marks = ref({});
const step = ref(0);
const pick = ref(null);
const probePick = ref(null);
const probeStep = ref(0);
const last = ref(null);
const orderSeq = ref([]);

watch(
  () => props.desk,
  () => {
    place.value = {};
    hold.value = null;
    marks.value = {};
    step.value = 0;
    pick.value = null;
    probePick.value = null;
    probeStep.value = 0;
    last.value = null;
    orderSeq.value = [];
  },
);

const kind = computed(() => props.desk.kind);
const probes = computed(() => props.desk.probes || []);
const probeNow = computed(() => probes.value[probeStep.value]);

function report(ev) {
  last.value = ev;
  emit("play", ev);
}

function put(id, col) {
  place.value = { ...place.value, [id]: col };
  const card = (props.desk.cards || []).find((c) => c.id === id);
  if (card) {
    const good = card.column === col;
    report({
      good,
      why: good ? card.why || "На месте. Смотрите, какой узел загорелся." : card.miss || "Не та куча — узел на доске красный.",
      hit: card.hit || col,
    });
  }
  maybeDoneDesk();
}

function maybeDoneDesk() {
  const cards = props.desk.cards || [];
  if (cards.length && cards.every((c) => place.value[c.id])) emit("done");
}

function drop(col) {
  if (!hold.value) return;
  put(hold.value, col);
  hold.value = null;
}

function toggleMark(id) {
  marks.value = { ...marks.value, [id]: !marks.value[id] };
  const line = (props.desk.lines || []).find((l) => l.id === id);
  const on = !!marks.value[id];
  if (line) {
    const good = line.bad ? on : !on;
    report({
      good,
      why: line.why || (good ? "Брак пойман." : "Это живая строка, не трогайте."),
      hit: id,
      lineId: id,
      tickets: { [id]: line.bad && on ? "bad" : on ? "ok" : undefined },
    });
  }
  emit("done");
}

function choose(i) {
  pick.value = i;
  const o = chatStep.value?.options?.[i];
  if (!o) return;
  report({ good: !!o.good, why: o.why, hit: o.hit });
  emit("done");
}

function nextChat() {
  if (step.value < (props.desk.steps?.length || 1) - 1) {
    step.value += 1;
    pick.value = null;
  }
}

const chatStep = computed(() => props.desk.steps?.[step.value]);

function ledgerPick(id) {
  pick.value = id;
  const leg = (props.desk.legs || []).find((l) => l.id === id);
  const good = !!(leg?.hole || id === props.desk.missing);
  report({
    good,
    why: good ? "Этой ноги быть не должно — дыра контракта." : "Не эта. Ищите лишнюю или недостающую — на доске появится призрак.",
    hit: id,
    extraDebit: good ? [] : [{ id: "fx-ghost", text: "DR лишняя · ход мимо дыры", ok: false }],
  });
  emit("done");
}

function mapPut(pinId, slot) {
  place.value = { ...place.value, [pinId]: slot };
  const pin = (props.desk.pins || []).find((p) => p.id === pinId);
  if (pin) {
    const good = pin.slot === slot;
    report({
      good,
      why: good ? pin.why || "Человек на своей зоне." : pin.miss || "Не та зона — узел роли красный.",
      hit: pin.hit || pin.slot || pin.id,
    });
  }
  const pins = props.desk.pins || [];
  if (pins.every((p) => place.value[p.id])) emit("done");
}

function dropPin(slot) {
  if (!hold.value) return;
  mapPut(hold.value, slot);
  hold.value = null;
}

function orderClick(item) {
  const next = orderSeq.value.length + 1;
  const good = item.pos === next;
  orderSeq.value = [...orderSeq.value, item.id];
  report({
    good,
    why: good ? "Шаг на месте — процесс на доске держится." : "Не этот шаг. Сначала предыдущий — иначе шов рвётся.",
    hit: item.id,
  });
  if (orderSeq.value.length === (props.desk.items || []).length) emit("done");
}

function orderReset() {
  orderSeq.value = [];
  last.value = null;
}

function chooseProbe(i) {
  probePick.value = i;
  const o = probeNow.value?.options?.[i];
  if (!o) return;
  report({ good: !!o.good, why: o.why, hit: o.hit });
  emit("done");
}

function nextProbe() {
  if (probeStep.value < probes.value.length - 1) {
    probeStep.value += 1;
    probePick.value = null;
  }
}

function optClass(i, o, current) {
  if (current === null) return "border-line bg-white hover:border-accent";
  if (current === i && o.good) return "border-ok bg-ok/10";
  if (current === i && !o.good) return "border-bad bg-bad/10";
  return "border-line bg-white/70 opacity-70";
}
</script>

<template>
  <section class="card rounded-[28px] p-5 md:p-6">
    <div class="flex items-start justify-between gap-3">
      <div>
        <div class="font-mono text-[10px] uppercase tracking-[0.2em] text-accent">Стол · ход двигает доску</div>
        <h3 class="mt-1 text-xl font-semibold">{{ desk.title }}</h3>
        <p class="mt-2 max-w-xl text-[17px] leading-7 text-mute">{{ desk.prompt || desk.setting }}</p>
      </div>
      <span v-if="last" class="shrink-0 rounded-full px-2.5 py-1 font-mono text-[10px]" :class="last.good ? 'bg-ok/15 text-ok' : 'bg-bad/15 text-bad'">
        {{ last.good ? "шов ок" : "доска сломалась" }}
      </span>
    </div>

    <div v-if="kind === 'desk'" class="mt-5">
      <div class="mb-4 flex flex-wrap gap-2">
        <button
          v-for="c in (desk.cards || []).filter((x) => !place[x.id])"
          :key="c.id"
          type="button"
          class="rounded-2xl border border-line bg-white px-3 py-2 text-left text-[15px] transition hover:border-accent"
          :class="hold === c.id ? 'border-accent ring-2 ring-accent/30' : ''"
          @click="hold = c.id"
        >
          {{ c.text }}
        </button>
      </div>
      <div class="grid gap-3" :class="(desk.columns || []).length > 2 ? 'sm:grid-cols-3' : 'md:grid-cols-2'">
        <button
          v-for="col in desk.columns"
          :key="col.id"
          type="button"
          class="min-h-36 rounded-2xl border border-dashed border-line p-3 text-left transition hover:border-accent/50"
          @click="drop(col.id)"
        >
          <div class="text-[11px] uppercase tracking-widest text-mute">{{ col.title }}</div>
          <div class="mt-2 space-y-2">
            <button
              v-for="c in (desk.cards || []).filter((x) => place[x.id] === col.id)"
              :key="c.id"
              type="button"
              class="w-full rounded-xl bg-paper px-3 py-2 text-left text-[15px]"
              :class="place[c.id] === c.column ? 'ring-1 ring-ok/40' : 'ring-1 ring-bad/40'"
              @click="put(c.id, desk.columns.find((x) => x.id !== col.id)?.id)"
            >
              {{ c.text }}
            </button>
          </div>
        </button>
      </div>
    </div>

    <div v-if="kind === 'chat' && chatStep" class="mt-5">
      <div class="rounded-2xl rounded-tl-sm bg-ink px-4 py-3 text-white shadow-sm">
        <div class="font-mono text-[11px] uppercase tracking-widest text-white/45">{{ chatStep.from }}</div>
        <p class="mt-1 text-[17px] leading-7">{{ chatStep.line }}</p>
      </div>
      <div class="mt-3 grid gap-2">
        <button
          v-for="(o, i) in chatStep.options"
          :key="o.text"
          type="button"
          class="rounded-2xl border px-4 py-3 text-left text-[16px] leading-6 transition"
          :class="optClass(i, o, pick)"
          @click="choose(i)"
        >
          {{ o.text }}
          <p v-if="pick === i" class="mt-2 text-sm text-mute">{{ o.why }}</p>
        </button>
      </div>
      <button
        v-if="pick !== null && step < desk.steps.length - 1"
        type="button"
        class="mt-4 rounded-full bg-accent px-4 py-2 text-sm text-white"
        @click="nextChat"
      >
        Дальше в чате
      </button>
      <p v-if="pick !== null" class="mt-3 text-[12px] text-mute">Можно выбрать другой ход — доска перерисуется. Дальше не закрыто.</p>
    </div>

    <div v-if="kind === 'ticket'" class="mt-5 overflow-hidden rounded-2xl border border-line">
      <pre v-if="desk.snippet" class="overflow-auto border-b border-line bg-ink px-4 py-3 font-mono text-[12px] leading-5 text-white/85">{{ desk.snippet }}</pre>
      <button
        v-for="line in desk.lines"
        :key="line.id"
        type="button"
        class="flex w-full gap-3 border-b border-line px-4 py-3 text-left text-[15px] last:border-0 transition"
        :class="marks[line.id] ? (line.bad ? 'bg-bad/12' : 'bg-ok/10') : 'bg-white hover:bg-paper'"
        @click="toggleMark(line.id)"
      >
        <span class="w-6 font-mono text-xs text-mute">{{ line.id }}</span>
        <span class="flex-1">{{ line.text }}</span>
      </button>
      <p class="px-4 py-3 text-sm text-mute">Тап по строке-браку — тикет на доске краснеет. Живую строку лучше не трогать. Дальше открыто.</p>
    </div>

    <div v-if="kind === 'ledger'" class="mt-5 grid gap-2">
      <button
        v-for="leg in desk.legs"
        :key="leg.id"
        type="button"
        class="rounded-2xl border px-4 py-3 text-left transition"
        :class="pick === leg.id ? (leg.hole ? 'border-ok bg-ok/10' : 'border-bad bg-bad/10') : 'border-line bg-white hover:border-accent'"
        @click="ledgerPick(leg.id)"
      >
        {{ leg.text }}
      </button>
      <p v-if="pick" class="text-sm text-mute">
        {{ last?.why }}
      </p>
    </div>

    <div v-if="kind === 'map'" class="mt-5">
      <div class="mb-3 flex flex-wrap gap-2">
        <button
          v-for="p in (desk.pins || []).filter((x) => !place[x.id])"
          :key="p.id"
          type="button"
          class="rounded-full border border-line px-3 py-1 text-sm transition hover:border-accent"
          :class="hold === p.id ? 'border-accent ring-2 ring-accent/30' : ''"
          @click="hold = p.id"
        >
          {{ p.text }}
        </button>
      </div>
      <div class="grid gap-3 sm:grid-cols-3">
        <button
          v-for="s in desk.slots"
          :key="s.id"
          type="button"
          class="min-h-28 rounded-2xl border border-dashed border-line p-3 text-left transition hover:border-accent/50"
          @click="dropPin(s.id)"
        >
          <div class="text-[11px] uppercase tracking-widest text-mute">{{ s.title }}</div>
          <div class="mt-2 space-y-1">
            <div
              v-for="p in (desk.pins || []).filter((x) => place[x.id] === s.id)"
              :key="p.id"
              class="rounded-lg px-2 py-1 text-sm"
              :class="place[p.id] === p.slot ? 'bg-ok/12' : 'bg-bad/12'"
            >
              {{ p.text }}
            </div>
          </div>
        </button>
      </div>
    </div>

    <div v-if="kind === 'order'" class="mt-5">
      <div class="grid gap-2">
        <button
          v-for="item in desk.items"
          :key="item.id"
          type="button"
          class="rounded-2xl border px-4 py-3 text-left text-[15px] leading-6 transition"
          :class="
            orderSeq.includes(item.id)
              ? orderSeq.indexOf(item.id) + 1 === item.pos
                ? 'border-ok bg-ok/10'
                : 'border-bad bg-bad/10'
              : 'border-line bg-white hover:border-accent'
          "
          @click="orderClick(item)"
        >
          <span class="mr-2 font-mono text-[11px] text-mute">{{ orderSeq.includes(item.id) ? orderSeq.indexOf(item.id) + 1 : "·" }}</span>
          {{ item.text }}
        </button>
      </div>
      <button type="button" class="mt-3 text-sm text-accent" @click="orderReset">Сбросить порядок</button>
    </div>

    <div v-if="probes.length && probeNow" class="mt-6 border-t border-line pt-5">
      <div class="font-mono text-[10px] uppercase tracking-[0.2em] text-accent2">Игра · ход {{ probeStep + 1 }}/{{ probes.length }}</div>
      <div class="mt-3 rounded-2xl rounded-tl-sm bg-ink px-4 py-3 text-white">
        <div class="font-mono text-[11px] text-white/45">{{ probeNow.from }}</div>
        <p class="mt-1 text-[16px] leading-7">{{ probeNow.line }}</p>
      </div>
      <div class="mt-3 grid gap-2">
        <button
          v-for="(o, i) in probeNow.options"
          :key="o.text"
          type="button"
          class="rounded-2xl border px-4 py-3 text-left text-[15px] leading-6 transition"
          :class="optClass(i, o, probePick)"
          @click="chooseProbe(i)"
        >
          {{ o.text }}
          <p v-if="probePick === i" class="mt-2 text-sm text-mute">{{ o.why }}</p>
        </button>
      </div>
      <button
        v-if="probePick !== null && probeStep < probes.length - 1"
        type="button"
        class="mt-4 rounded-full bg-ink px-4 py-2 text-sm text-white"
        @click="nextProbe"
      >
        Ещё ход
      </button>
    </div>
  </section>
</template>
