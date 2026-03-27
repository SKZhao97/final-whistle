import type { Locale } from "./config";
import { translate } from "./core";

export function formatDateTime(value: string | Date, locale: Locale) {
  const date = value instanceof Date ? value : new Date(value);
  return new Intl.DateTimeFormat(locale === "zh" ? "zh-CN" : "en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

export function formatNumber(value: number | null | undefined, locale: Locale) {
  if (value === null || value === undefined || Number.isNaN(value)) {
    return translate(locale, "common.noSamples");
  }
  return new Intl.NumberFormat(locale === "zh" ? "zh-CN" : "en-US", {
    maximumFractionDigits: 1,
  }).format(value);
}

export function getWatchedTypeLabel(value: "FULL" | "PARTIAL" | "HIGHLIGHTS", locale: Locale) {
  switch (value) {
    case "FULL":
      return translate(locale, "enum.watchedType.full");
    case "PARTIAL":
      return translate(locale, "enum.watchedType.partial");
    case "HIGHLIGHTS":
      return translate(locale, "enum.watchedType.highlights");
  }
}

export function getSupporterSideLabel(
  value: "HOME" | "AWAY" | "NEUTRAL",
  locale: Locale,
) {
  switch (value) {
    case "HOME":
      return translate(locale, "enum.supporterSide.home");
    case "AWAY":
      return translate(locale, "enum.supporterSide.away");
    case "NEUTRAL":
      return translate(locale, "enum.supporterSide.neutral");
  }
}
