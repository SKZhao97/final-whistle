export const SUPPORTED_LOCALES = ["en", "zh"] as const;

export type Locale = (typeof SUPPORTED_LOCALES)[number];

export const DEFAULT_LOCALE: Locale = "en";
export const LOCALE_COOKIE_NAME = "final_whistle_locale";

export function normalizeLocale(value?: string | null): Locale {
  if (!value) {
    return DEFAULT_LOCALE;
  }

  const normalized = value.toLowerCase().trim();
  if (normalized === "zh" || normalized === "zh-cn" || normalized === "zh-hans") {
    return "zh";
  }
  if (normalized === "en" || normalized === "en-us" || normalized === "en-gb") {
    return "en";
  }

  return DEFAULT_LOCALE;
}
