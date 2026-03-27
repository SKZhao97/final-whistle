"use client";

import { useLocale } from "@/components/i18n/LocaleProvider";
import { useTheme } from "@/components/theme/ThemeProvider";

export default function ThemeSwitcher() {
  const { theme, setTheme } = useTheme();
  const { t } = useLocale();

  return (
    <div
      className="inline-flex items-center rounded-full bg-transparent p-0.5 text-xs"
      aria-label={t("theme.label")}
      title={t("theme.label")}
    >
      <button
        type="button"
        onClick={() => setTheme("light")}
        className={`rounded-full px-3 py-1.5 transition-colors ${
          theme === "light"
            ? "bg-[var(--fw-control-active)] text-[var(--fw-control-active-ink)]"
            : "text-[var(--fw-muted)] hover:bg-[var(--fw-field-100)] hover:text-[var(--fw-ink)]"
        }`}
        aria-pressed={theme === "light"}
      >
        {t("theme.light")}
      </button>
      <button
        type="button"
        onClick={() => setTheme("dark")}
        className={`rounded-full px-3 py-1.5 transition-colors ${
          theme === "dark"
            ? "bg-[var(--fw-control-active)] text-[var(--fw-control-active-ink)]"
            : "text-[var(--fw-muted)] hover:bg-[var(--fw-field-100)] hover:text-[var(--fw-ink)]"
        }`}
        aria-pressed={theme === "dark"}
      >
        {t("theme.dark")}
      </button>
    </div>
  );
}
