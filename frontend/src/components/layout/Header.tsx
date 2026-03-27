import Link from "next/link";

import AuthStatus from "@/components/auth/AuthStatus";
import LanguageSwitcher from "@/components/i18n/LanguageSwitcher";
import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";

export default async function Header() {
  const locale = await getServerLocale();

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 max-w-screen-2xl items-center">
        <div className="mr-4 flex">
          <Link href="/" className="mr-6 flex items-center space-x-2">
            <span className="font-bold text-xl">Final Whistle</span>
          </Link>
          <nav className="flex items-center space-x-6 text-sm font-medium">
            <Link
              href="/matches"
              className="transition-colors hover:text-foreground/80 text-foreground/60"
            >
              {translate(locale, "nav.matches")}
            </Link>
          </nav>
        </div>
        <div className="flex flex-1 items-center justify-end space-x-4">
          <nav className="flex items-center space-x-2">
            <LanguageSwitcher />
            <AuthStatus />
          </nav>
        </div>
      </div>
    </header>
  );
}
