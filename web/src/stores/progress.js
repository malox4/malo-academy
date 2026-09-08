import { defineStore } from "pinia";
import { ref } from "vue";

export const useProgress = defineStore("progress", () => {
  const payload = ref({ lessons: {} });
  const catalog = ref(null);

  async function load() {
    const [p, c] = await Promise.all([
      fetch("/api/me/progress", { credentials: "include" }).then((r) => r.json()),
      fetch("/api/catalog", { credentials: "include" }).then((r) => r.json()),
    ]);
    payload.value = p.payload && typeof p.payload === "object" ? p.payload : { lessons: {} };
    if (!payload.value.lessons) payload.value.lessons = {};
    catalog.value = c;
  }

  async function save() {
    await fetch("/api/me/progress", {
      method: "PUT",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ payload: payload.value }),
    });
  }

  function isDone(id) {
    return Boolean(payload.value.lessons?.[id]?.done);
  }

  function sawInterview(id) {
    return Boolean(payload.value.interview?.[id]);
  }

  async function markInterview(id) {
    payload.value.interview = { ...(payload.value.interview || {}), [id]: { at: Date.now() } };
    await save();
  }

  async function complete(id) {
    payload.value.lessons = { ...payload.value.lessons, [id]: { ...(payload.value.lessons?.[id] || {}), done: true, at: Date.now() } };
    await save();
  }

  return { payload, catalog, load, save, isDone, complete, sawInterview, markInterview };
});
