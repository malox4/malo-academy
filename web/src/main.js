import { createApp } from "vue";
import { createPinia } from "pinia";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import "./style.css";
import { useAuth } from "./stores/auth";
import { bootTelegram, inTelegram } from "./lib/telegram";

const routes = [
  { path: "/", component: () => import("./pages/LandingPage.vue"), meta: { public: true } },
  { path: "/pricing", component: () => import("./pages/PricingPage.vue"), meta: { public: true } },
  { path: "/login", component: () => import("./pages/AuthPage.vue"), meta: { public: true, mode: "login" } },
  { path: "/register", component: () => import("./pages/AuthPage.vue"), meta: { public: true, mode: "register" } },
  { path: "/hall", component: () => import("./pages/HallPage.vue") },
  { path: "/level/:id", component: () => import("./pages/LevelPage.vue") },
  { path: "/lesson/:id", component: () => import("./pages/LessonPage.vue") },
  { path: "/boards", component: () => import("./pages/BoardsPage.vue") },
  { path: "/boards/:id", component: () => import("./pages/BoardsPage.vue") },
  { path: "/pet", component: () => import("./pages/PetPage.vue") },
  { path: "/api-lab", redirect: "/pet?tab=api" },
  { path: "/practice", component: () => import("./pages/PracticePage.vue") },
  { path: "/materials", component: () => import("./pages/MaterialsPage.vue"), meta: { public: true } },
  { path: "/materials/:id", component: () => import("./pages/TopicPage.vue"), meta: { public: true } },
  { path: "/interview", component: () => import("./pages/InterviewPage.vue") },
  { path: "/live", component: () => import("./pages/LiveInterviewPage.vue"), meta: { public: true } },
  { path: "/interview/live", redirect: "/live" },
  { path: "/sobes", redirect: "/interview" },
  { path: "/profile", component: () => import("./pages/ProfilePage.vue") },
  { path: "/admin", component: () => import("./pages/AdminPage.vue"), meta: { admin: true } },
  { path: "/:pathMatch(.*)*", redirect: "/" },
];

const router = createRouter({ history: createWebHistory(), routes });

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);
app.use(router);

const auth = useAuth();
bootTelegram(router);
router.beforeEach(async (to) => {
  if (!auth.ready) await auth.hydrate();
  if (to.path === "/" && auth.user) return "/hall";
  if (inTelegram() && to.path === "/" && !auth.user) return "/login";
  if (to.meta.public) return true;
  if (!auth.user) return { path: "/login", query: { next: to.fullPath } };
  if (to.meta.admin && auth.user.role !== "admin") return "/hall";
  return true;
});

function placeTitle(to) {
  const p = to.path;
  if (p === "/hall") return "Зал";
  if (p.startsWith("/lesson/")) return "";
  if (p.startsWith("/level/")) return "Этаж";
  if (p.startsWith("/pet")) return "Пет";
  if (p.startsWith("/practice")) return "Практика";
  if (p.startsWith("/interview")) return "Собес";
  if (p.startsWith("/live")) return "Live собес";
  if (p.startsWith("/materials")) return "Материалы";
  if (p.startsWith("/boards")) return "Доска";
  if (p === "/profile") return "Профиль";
  if (p === "/admin") return "Журнал";
  return to.name || p;
}

function pingHere() {
  if (!auth.user) return;
  const to = router.currentRoute.value;
  if (to.path === "/login" || to.path === "/register") return;
  fetch("/api/me/here", {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ path: to.fullPath, title: placeTitle(to) }),
  }).catch(() => {});
}

router.afterEach(() => {
  pingHere();
});
setInterval(pingHere, 45000);

app.mount("#app");
