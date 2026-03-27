"use client";

import Link from "next/link";

import { useAuth } from "@/components/auth/AuthProvider";
import { useLocale } from "@/components/i18n/LocaleProvider";

export default function AuthStatus() {
  const { status, user, logout } = useAuth();
  const { t } = useLocale();

  if (status === "loading") {
    return <span className="px-2 text-sm text-[var(--fw-muted)]">{t("nav.checkingSession")}</span>;
  }

  if (status === "authenticated" && user) {
    return (
      <div className="flex items-center gap-2">
        <span className="px-2 text-sm text-[var(--fw-muted)]">{user.name}</span>
        <button
          type="button"
          onClick={() => void logout()}
          className="fw-button fw-button--secondary px-4 py-2 text-sm"
        >
          {t("nav.logout")}
        </button>
      </div>
    );
  }

  return (
    <Link
      href="/login"
      className="fw-button fw-button--primary px-4 py-2 text-sm"
    >
      {t("nav.devLogin")}
    </Link>
  );
}
