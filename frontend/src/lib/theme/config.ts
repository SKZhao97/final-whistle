export const DEFAULT_THEME = "light" as const;
export const THEME_COOKIE_NAME = "final_whistle_theme";

export type Theme = "light" | "dark";

export function normalizeTheme(value?: string | null): Theme {
  return value === "dark" ? "dark" : DEFAULT_THEME;
}
