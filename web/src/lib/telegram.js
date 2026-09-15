export function tg() {
  return typeof window !== "undefined" ? window.Telegram?.WebApp || null : null;
}

export function inTelegram() {
  return Boolean(tg()?.initData);
}

export function bootTelegram(router) {
  const app = tg();
  if (!app?.initData) return;
  document.documentElement.classList.add("tg-app");
  try {
    app.ready();
    app.expand();
    app.setHeaderColor?.("#ffffff");
    app.setBackgroundColor?.("#f1f3f6");
    app.disableVerticalSwipes?.();
  } catch {
    /* older clients */
  }

  const syncBack = () => {
    if (!app.BackButton) return;
    const p = router.currentRoute.value.path;
    if (p === "/" || p === "/hall" || p === "/login") app.BackButton.hide();
    else app.BackButton.show();
  };

  app.BackButton?.onClick?.(() => {
    if (window.history.length > 1) router.back();
    else router.replace("/hall");
  });
  router.afterEach(syncBack);
  syncBack();
}
