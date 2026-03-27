"use client";

import { useLocale } from "@/components/i18n/LocaleProvider";

export default function LanguageSwitcher() {
  const { locale, setLocale, t } = useLocale();

  return (
    <div
      className="inline-flex items-center rounded-full border border-neutral-200 bg-white/80 p-0.5 text-xs shadow-sm"
      aria-label={t("locale.label")}
      title={t("locale.label")}
    >
      <button
        type="button"
        onClick={() => setLocale("en")}
        className={`rounded-full px-2.5 py-1 transition-colors ${
          locale === "en" ? "bg-neutral-900 text-white" : "text-neutral-500 hover:bg-neutral-100 hover:text-neutral-800"
        }`}
        aria-pressed={locale === "en"}
      >
        EN
      </button>
      <button
        type="button"
        onClick={() => setLocale("zh")}
        className={`rounded-full px-2.5 py-1 transition-colors ${
          locale === "zh" ? "bg-neutral-900 text-white" : "text-neutral-500 hover:bg-neutral-100 hover:text-neutral-800"
        }`}
        aria-pressed={locale === "zh"}
      >
        中
      </button>
    </div>
  );
}
