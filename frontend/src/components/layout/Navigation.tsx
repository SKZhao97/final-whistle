"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { useLocale } from "@/components/i18n/LocaleProvider";

export default function Navigation() {
  const { t } = useLocale();
  const pathname = usePathname();
  const navItems = [
    { href: "/matches", label: t("nav.matches") },
    { href: "/me", label: t("nav.profile") },
  ];

  return (
    <nav className="flex items-center gap-2">
      {navItems.map((item) => (
        <Link
          key={item.href}
          href={item.href}
          aria-current={pathname === item.href ? "page" : undefined}
          className={`flex items-center rounded-full px-4 py-2 text-sm font-medium transition-colors ${
            pathname === item.href
              ? "bg-[var(--fw-nav-active)] text-[var(--fw-nav-active-ink)] shadow-[0_12px_28px_rgba(16,31,24,0.14)]"
              : "text-[var(--fw-muted)] hover:bg-[var(--fw-field-100)] hover:text-[var(--fw-field-900)]"
          }`}
        >
          <span>{item.label}</span>
        </Link>
      ))}
    </nav>
  );
}
