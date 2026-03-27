import { cookies } from "next/headers";

import { normalizeTheme, THEME_COOKIE_NAME } from "./config";

export async function getServerTheme() {
  const cookieStore = await cookies();
  return normalizeTheme(cookieStore.get(THEME_COOKIE_NAME)?.value);
}
