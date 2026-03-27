"use client";

import { useLocale } from "@/components/i18n/LocaleProvider";

export default function PlayerDetailError() {
  const { t } = useLocale();

  return (
    <div className="py-10">
      <h1 className="text-2xl font-semibold">{t("player.error.title")}</h1>
      <p className="mt-2 text-sm text-neutral-600">{t("player.error.body")}</p>
    </div>
  );
}
