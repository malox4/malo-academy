import { createApp } from "vue";
import { createPinia } from "pinia";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import "./style.css";
import { useAuth } from "./stores/auth";

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
  { path: "/admin", component: () => import("./pages/AdminPage.vue") },
  { path: "/:pathMatch(.*)*", redirect: "/" },
];

const router = createRouter({ history: createWebHistory(), routes });

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);
app.use(router);

const auth = useAuth();
router.beforeEach(async (to) => {
  if (!auth.ready) await auth.hydrate();
  if (to.path === "/" && auth.user) return "/hall";
  if (to.meta.public) return true;
  if (!auth.user) return { path: "/login", query: { next: to.fullPath } };
  return true;
});

app.mount("#app");
