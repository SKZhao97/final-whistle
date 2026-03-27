import type { Metadata } from "next";
import "./globals.css";
import { AuthProvider } from "@/components/auth/AuthProvider";
import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";
import { LocaleProvider } from "@/components/i18n/LocaleProvider";
import { getServerLocale } from "@/lib/i18n/server";
import { ThemeProvider } from "@/components/theme/ThemeProvider";
import { getServerTheme } from "@/lib/theme/server";

export const metadata: Metadata = {
  title: "Final Whistle",
  description: "A post-match recording product for football viewers.",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const locale = await getServerLocale();
  const theme = await getServerTheme();

  return (
    <html lang={locale} data-theme={theme} className="h-full antialiased">
      <body className="app-shell flex min-h-full flex-col">
        <ThemeProvider initialTheme={theme}>
          <LocaleProvider initialLocale={locale}>
            <AuthProvider>
              <Header />
              <main className="container flex-1 py-8 md:py-10">{children}</main>
              <Footer />
            </AuthProvider>
          </LocaleProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
