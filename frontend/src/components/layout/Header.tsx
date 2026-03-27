import Link from "next/link";

import AuthStatus from "@/components/auth/AuthStatus";
import Navigation from "@/components/layout/Navigation";
import { BrandMark, Wordmark } from "@/components/experience/FootballPrimitives";
import LanguageSwitcher from "@/components/i18n/LanguageSwitcher";
import ThemeSwitcher from "@/components/theme/ThemeSwitcher";
import { getServerLocale } from "@/lib/i18n/server";

export default async function Header() {
  await getServerLocale();

  return (
    <header className="sticky top-0 z-50 w-full border-b border-[var(--fw-line)] bg-[var(--fw-header)] backdrop-blur supports-[backdrop-filter]:bg-[var(--fw-header)]/90">
      <div className="container flex h-18 max-w-screen-2xl items-center justify-between gap-6 py-3">
        <div className="flex items-center gap-8">
          <Link href="/" className="flex items-center gap-3 rounded-full pr-3 transition-opacity hover:opacity-90">
            <BrandMark />
            <Wordmark />
          </Link>
          <Navigation />
        </div>
        <div className="flex items-center gap-3">
          <div className="header-control rounded-full p-1">
            <ThemeSwitcher />
          </div>
          <div className="header-control rounded-full p-1">
            <LanguageSwitcher />
          </div>
          <div className="header-control rounded-full px-2 py-1">
            <AuthStatus />
          </div>
        </div>
      </div>
    </header>
  );
}
