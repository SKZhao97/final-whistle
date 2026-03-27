"use client";

import { useLocale } from "@/components/i18n/LocaleProvider";

export default function MatchDetailError() {
  const { t } = useLocale();

  return (
    <div className="py-10">
      <h1 className="text-2xl font-semibold">{t("matchDetail.error.title")}</h1>
      <p className="mt-2 text-sm text-neutral-600">
        {t("matchDetail.error.body")}
      </p>
      <p className="mt-2 text-sm text-neutral-600">
        {t("matchDetail.error.restart")}
      </p>
    </div>
  );
}
