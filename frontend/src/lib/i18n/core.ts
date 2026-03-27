import { DEFAULT_LOCALE, type Locale } from "./config.js";
import { messages } from "./messages.js";

type MessageParams = Record<string, string | number>;

export type TranslationKey = keyof (typeof messages)[typeof DEFAULT_LOCALE];

export function translate(
  locale: Locale,
  key: TranslationKey,
  params?: MessageParams,
) {
  const entry = messages[locale][key] ?? messages[DEFAULT_LOCALE][key];
  if (typeof entry === "function") {
    return entry(params);
  }
  return entry;
}
