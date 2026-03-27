"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { SectionShell } from "@/components/experience/FootballPrimitives";
import { useAuth } from "@/components/auth/AuthProvider";
import { useLocale } from "@/components/i18n/LocaleProvider";
import { ApiError } from "@/lib/api/client";

export default function LoginPage() {
  const router = useRouter();
  const { login, status } = useAuth();
  const { t, locale } = useLocale();
  const [email, setEmail] = useState("demo@final-whistle.test");
  const [name, setName] = useState("Demo User");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      await login({ email, name });
      router.push("/me");
      router.refresh();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError(
          locale === "zh"
            ? "登录失败，请确认后端正在运行且数据库连接正确。"
            : "Login failed. Make sure the backend is running and connected to the database.",
        );
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="mx-auto max-w-3xl py-10">
      <SectionShell
        eyebrow={t("auth.sessionAuth")}
        title={t("auth.loginTitle")}
        description={t("auth.loginDescription")}
        accent="field"
      >
        <form onSubmit={handleSubmit} className="mt-6 grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
          <div className="space-y-4 rounded-[1.5rem] border border-[var(--fw-line)] bg-[var(--fw-surface)]/88 p-6">
            <label className="block text-sm">
              <span className="mb-2 block font-medium text-[var(--fw-ink-soft)]">{t("auth.email")}</span>
              <input
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                className="w-full rounded-[1rem] border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] px-4 py-3 text-[var(--fw-ink)]"
                type="email"
                required
              />
            </label>

            <label className="block text-sm">
              <span className="mb-2 block font-medium text-[var(--fw-ink-soft)]">{t("auth.displayName")}</span>
              <input
                value={name}
                onChange={(event) => setName(event.target.value)}
                className="w-full rounded-[1rem] border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] px-4 py-3 text-[var(--fw-ink)]"
                type="text"
                required
              />
            </label>

            {error ? <p className="text-sm text-red-700">{error}</p> : null}
            {status === "authenticated" ? (
              <p className="text-sm text-[var(--fw-field-700)]">{t("auth.alreadySignedIn")}</p>
            ) : null}

            <button type="submit" disabled={submitting} className="fw-button fw-button--primary w-full disabled:opacity-50">
              {submitting ? t("auth.signingIn") : t("auth.signIn")}
            </button>
          </div>

          <div className="rounded-[1.5rem] border border-[var(--fw-line)] bg-[var(--fw-paper-strong)] p-6">
            <p className="match-eyebrow">{locale === "zh" ? "开发提示" : "Dev Note"}</p>
            <div className="mt-4 space-y-3 text-sm leading-6 text-[var(--fw-ink-soft)]">
              <p>
                {locale === "zh"
                  ? "当前登录仍然依赖本地后端和数据库。若登录失败，优先检查 API 是否在 8080 端口运行。"
                  : "Login still depends on the local backend and database. If sign-in fails, first confirm the API is running on port 8080."}
              </p>
              <p>
                {locale === "zh"
                  ? "推荐使用 demo@final-whistle.test 体验完整链路。"
                  : "Use demo@final-whistle.test for the most reliable smoke path."}
              </p>
            </div>
          </div>
        </form>
      </SectionShell>
    </div>
  );
}
