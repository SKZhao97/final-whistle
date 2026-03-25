import Link from "next/link";

export default function Header() {
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
              Matches
            </Link>
            <Link
              href="/teams"
              className="transition-colors hover:text-foreground/80 text-foreground/60"
            >
              Teams
            </Link>
            <Link
              href="/players"
              className="transition-colors hover:text-foreground/80 text-foreground/60"
            >
              Players
            </Link>
          </nav>
        </div>
        <div className="flex flex-1 items-center justify-end space-x-4">
          <nav className="flex items-center space-x-2">
            <Link
              href="/me"
              className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
            >
              My Profile
            </Link>
          </nav>
        </div>
      </div>
    </header>
  );
}
