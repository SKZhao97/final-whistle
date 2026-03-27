import { getServerLocale } from "@/lib/i18n/server";
import { translate } from "@/lib/i18n/core";

export default async function Footer() {
  const locale = await getServerLocale();

  return (
    <footer className="border-t py-6 md:py-0">
      <div className="container flex flex-col items-center justify-between gap-4 md:h-24 md:flex-row">
        <div className="flex flex-col items-center gap-4 px-8 md:flex-row md:gap-2 md:px-0">
          <p className="text-center text-sm leading-loose text-muted-foreground md:text-left">
            &copy; {new Date().getFullYear()} Final Whistle. {translate(locale, "footer.tagline")}
          </p>
        </div>
        <p className="text-sm text-muted-foreground">
          {translate(locale, "footer.foundation")}
        </p>
      </div>
    </footer>
  );
}
