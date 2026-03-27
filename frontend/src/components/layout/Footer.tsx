import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";

export default async function Footer() {
  const locale = await getServerLocale();

  return (
    <footer className="border-t border-[var(--fw-line)] py-8">
      <div className="container flex flex-col items-center justify-between gap-4 md:flex-row">
        <div className="flex flex-col items-center gap-3 px-4 md:flex-row md:px-0">
          <p className="text-center text-sm leading-loose text-[var(--fw-muted)] md:text-left">
            &copy; {new Date().getFullYear()} Final Whistle. {translate(locale, "footer.tagline")}
          </p>
        </div>
        <p className="text-sm text-[var(--fw-muted)]">
          {translate(locale, "footer.foundation")}
        </p>
      </div>
    </footer>
  );
}
