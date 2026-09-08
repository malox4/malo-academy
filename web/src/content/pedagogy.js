/** Primer: термин → объяснение → примеры → для чего → кейс. */
import { TERMS } from "./glossary.js";

export const H = {
  term: "Термин",
  about: "Объяснение",
  example: "Примеры",
  purpose: "Для чего это нужно",
  cases: "Кейс",
  mistakes: "Ошибки",
  practice: "Закрепление",
};

function paras(v) {
  if (!v) return [];
  return Array.isArray(v) ? v.filter(Boolean) : [String(v)];
}

export function arc(s) {
  const what = paras(s.what);
  const term = paras(s.term).length ? paras(s.term) : what.slice(0, 1);
  const about = paras(s.about).length ? paras(s.about) : what.slice(1);
  return {
    term: {
      title: H.term,
      name: s.name || "",
      en: s.en || "",
      name2: s.name2 || "",
      en2: s.en2 || "",
      paragraphs: term,
    },
    about: {
      title: H.about,
      paragraphs: about,
      intro: s.howIntro || "",
      steps: s.howSteps || s.how || [],
      infographic: s.aboutFig || s.whatFig || s.howFig || null,
    },
    example: {
      title: H.example,
      heading: s.exampleHeading,
      paragraphs: paras(s.example),
      note: s.exampleNote || "",
      infographic: s.exampleFig || null,
    },
    purpose: { title: H.purpose, paragraphs: paras(s.purpose), infographic: s.purposeFig || null },
    mistakes: { title: H.mistakes, intro: s.mistakesIntro, items: s.mistakes || [] },
    cases: { title: H.cases, items: s.cases || [] },
    practice: { title: H.practice, paragraphs: paras(s.practice), to: s.to, cta: s.cta },
  };
}

/** Real term first. Long story goes to explanation. */
export function teach(s) {
  return arc({
    name: s.name,
    en: s.en,
    name2: s.name2,
    en2: s.en2,
    term: s.term,
    about: [...(s.not ? [s.not] : []), ...paras(s.about)],
    whatFig: s.fig,
    purpose: s.purpose,
    purposeFig: s.purposeFig,
    howIntro: s.howIntro,
    howSteps: s.how,
    howFig: s.howFig,
    exampleHeading: s.exampleTitle,
    example: s.example,
    exampleNote: s.exampleNote,
    exampleFig: s.exampleFig,
    mistakesIntro: s.mistakesIntro,
    mistakes: s.mistakes,
    cases: s.cases,
    practice: s.practice,
    to: s.to,
    cta: s.cta,
  });
}

/** Stamp dictionary term onto a built article. Long old term text moves into explanation. */
export function stampTerms(topic) {
  const g = TERMS[topic.id];
  if (!g || !topic.sections?.term) return topic;
  const term = topic.sections.term;
  const about = topic.sections.about;
  term.name = g.name;
  term.en = g.en || "";
  term.name2 = g.name2 || "";
  term.en2 = g.en2 || "";
  const old = term.paragraphs || [];
  if (g.def) {
    term.paragraphs = [g.def];
    const extra = old.filter((p) => p && p !== g.def);
    if (about && extra.length) {
      about.paragraphs = [...extra, ...(about.paragraphs || [])];
    }
  }
  return topic;
}
